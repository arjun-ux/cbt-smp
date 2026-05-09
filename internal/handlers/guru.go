package handlers

import (
	"cbt-smp/internal/database"
	"cbt-smp/internal/models"
	"cbt-smp/pkg/utils"
	encoding_csv "encoding/csv"
	"strings"
	"sync"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type GuruInput struct {
	NIP      string `json:"nip"`
	NamaGuru string `json:"nama_guru"`
	Password string `json:"password"` // Opsional saat update
}

// GetGurus mengambil semua data guru
func GetGurus(c *fiber.Ctx) error {
	var gurus []models.MasterGuru
	// Kita join dengan tabel user agar bisa lihat is_active
	database.DB.Preload("User").Find(&gurus)

	return c.JSON(fiber.Map{
		"data": gurus,
	})
}

// CreateGuru menambah guru baru (sekaligus membuatkan akun login)
func CreateGuru(c *fiber.Ctx) error {
	var input GuruInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Input tidak valid"})
	}

	// Validasi dasar
	if input.NIP == "" || input.NamaGuru == "" || input.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "NIP, Nama, dan Password wajib diisi"})
	}

	// Gunakan Transaction agar jika salah satu gagal, semuanya dibatalkan
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		// 1. Hash Password
		hash, errHash := utils.HashPassword(input.Password)
		if errHash != nil {
			return errHash
		}

		// 2. Buat User (Login Akun)
		newUser := models.User{
			Username: input.NIP, // NIP dijadikan Username
			Password: hash,
			Role:     "guru",
			IsActive: true,
		}

		if err := tx.Create(&newUser).Error; err != nil {
			return err
		}

		// 3. Buat Master Guru
		newGuru := models.MasterGuru{
			UserID:   newUser.ID,
			NIP:      input.NIP,
			NamaGuru: input.NamaGuru,
		}

		if err := tx.Create(&newGuru).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal menyimpan data guru (Mungkin NIP sudah terdaftar)"})
	}

	return c.JSON(fiber.Map{"message": "Guru berhasil ditambahkan"})
}

// UpdateGuru mengubah data guru (nama atau reset password)
func UpdateGuru(c *fiber.Ctx) error {
	id := c.Params("id")
	var input GuruInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Input tidak valid"})
	}

	var guru models.MasterGuru
	if err := database.DB.Preload("User").First(&guru, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Guru tidak ditemukan"})
	}

	// LOCKING: Cek apakah guru sedang mengawas ujian aktif
	var activeDutyCount int64
	database.DB.Model(&models.CBTJadwalUjian{}).
		Where("pengawas_id = ? AND status = ?", id, "Berlangsung").
		Count(&activeDutyCount)
	
	if activeDutyCount > 0 {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Data guru tidak dapat diubah karena sedang bertugas sebagai pengawas pada ujian yang sedang berlangsung."})
	}

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		// Update tabel guru
		if input.NamaGuru != "" {
			guru.NamaGuru = input.NamaGuru
			guru.NIP = input.NIP
			if err := tx.Save(&guru).Error; err != nil {
				return err
			}
		}

		// Update tabel user (jika ganti password atau NIP berubah)
		guru.User.Username = guru.NIP // Sesuaikan username dengan NIP baru
		if input.Password != "" {
			hash, errHash := utils.HashPassword(input.Password)
			if errHash != nil {
				return errHash
			}
			guru.User.Password = hash
		}
		
		if err := tx.Save(&guru.User).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal memperbarui data guru"})
	}

	return c.JSON(fiber.Map{"message": "Data guru berhasil diperbarui"})
}

// DeleteGuru menghapus guru dan akun loginnya secara permanen
func DeleteGuru(c *fiber.Ctx) error {
	id := c.Params("id")
	
	var guru models.MasterGuru
	if err := database.DB.First(&guru, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Guru tidak ditemukan"})
	}

	// LOCKING: Cek apakah guru sedang mengawas ujian aktif
	var activeDutyCount int64
	database.DB.Model(&models.CBTJadwalUjian{}).
		Where("pengawas_id = ? AND status = ?", id, "Berlangsung").
		Count(&activeDutyCount)
	
	if activeDutyCount > 0 {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Guru tidak dapat dihapus karena sedang bertugas sebagai pengawas pada ujian yang sedang berlangsung."})
	}

	// Gunakan transaksi untuk menghapus di kedua tabel
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		// 1. Hapus profil MasterGuru
		if err := tx.Delete(&guru).Error; err != nil {
			return err
		}
		// 2. Hapus akun User-nya
		if err := tx.Delete(&models.User{}, guru.UserID).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		errorMessage := "Gagal menghapus data guru secara permanen"
		if strings.Contains(err.Error(), "FOREIGN KEY") {
			errorMessage = "Guru tidak dapat dihapus karena masih memiliki data Bank Soal, penugasan Mapel, atau sedang ditugaskan sebagai Pengawas Ujian. Silakan pindahkan atau hapus keterikatan data tersebut terlebih dahulu."
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": errorMessage})
	}

	return c.JSON(fiber.Map{"message": "Guru berhasil dihapus secara permanen"})
}

// ToggleGuruStatus mengaktifkan atau menonaktifkan akun guru
func ToggleGuruStatus(c *fiber.Ctx) error {
	id := c.Params("id")
	var guru models.MasterGuru
	if err := database.DB.Preload("User").First(&guru, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Guru tidak ditemukan"})
	}

	// Balikkan status is_active
	guru.User.IsActive = !guru.User.IsActive
	database.DB.Save(&guru.User)

	status := "diaktifkan"
	if !guru.User.IsActive {
		status = "dinonaktifkan"
	}

	return c.JSON(fiber.Map{"message": "Akun guru berhasil " + status})
}

// ImportGuru memproses mass-upload guru via CSV (Turbo Parallel)
func ImportGuru(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return SendError(c, 400, "Gagal membaca file upload")
	}

	f, err := file.Open()
	if err != nil {
		return SendError(c, 500, "Gagal membuka file")
	}
	defer f.Close()

	reader := encoding_csv.NewReader(f)
	records, err := reader.ReadAll()
	if err != nil {
		return SendError(c, 400, "Format CSV tidak valid")
	}

	if len(records) < 2 {
		return SendError(c, 400, "File CSV kosong")
	}

	// 1. Kumpulkan data mentah
	type TempGuru struct {
		nip      string
		nama     string
		password string
	}
	var rawData []TempGuru
	for i, row := range records {
		if i == 0 { continue }
		if len(row) < 3 { continue }
		rawData = append(rawData, TempGuru{
			nip:      strings.TrimSpace(row[0]),
			nama:      strings.TrimSpace(row[1]),
			password:  strings.TrimSpace(row[2]),
		})
	}

	// 2. Hashing Paralel
	users := make([]models.User, len(rawData))
	var wg sync.WaitGroup
	for i := range rawData {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			hash, _ := utils.HashPassword(rawData[idx].password)
			users[idx] = models.User{
				Username: rawData[idx].nip,
				Password: hash,
				Role:     "guru",
				IsActive: true,
			}
		}(i)
	}
	wg.Wait()

	// 3. Simpan ke Database
	errTx := database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&users).Error; err != nil {
			return err
		}

		var gurus []models.MasterGuru
		for i, u := range users {
			t := rawData[i]
			gurus = append(gurus, models.MasterGuru{
				UserID:   u.ID,
				NIP:      t.nip,
				NamaGuru: t.nama,
			})
		}
		return tx.Create(&gurus).Error
	})

	if errTx != nil {
		return SendError(c, 500, "Gagal simpan (NIP duplikat atau data salah)")
	}

	return SendSuccess(c, fmt.Sprintf("Berhasil impor %d guru", len(users)), nil)
}
