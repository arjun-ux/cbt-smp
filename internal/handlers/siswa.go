package handlers

import (
	"cbt-smp/internal/database"
	"cbt-smp/internal/models"
	"cbt-smp/pkg/utils"
	encoding_csv "encoding/csv"
	"fmt"
	"strings"
	"sync"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type SiswaInput struct {
	NISN         string `json:"nisn"`
	NomorPeserta string `json:"nomor_peserta"`
	NamaLengkap  string `json:"nama_lengkap"`
	Password     string `json:"password"`
	KelasID      uint   `json:"kelas_id"`
	RuangID      uint   `json:"ruang_id"`
	SesiID       uint   `json:"sesi_id"`
}

type BulkPlotInput struct {
	SiswaIDs []uint `json:"siswa_ids"`
	RuangID  uint   `json:"ruang_id"`
	SesiID   uint   `json:"sesi_id"`
}

// GetSiswas mengambil semua data siswa beserta relasinya
func GetSiswas(c *fiber.Ctx) error {
	var siswas []models.MasterSiswa
	
	database.DB.Preload("User").
		Preload("Kelas").
		Preload("Ruang").
		Preload("Sesi").
		Find(&siswas)

	return SendSuccess(c, "Berhasil mengambil daftar siswa", siswas)
}

// CreateSiswa menambah siswa baru (sekaligus membuatkan akun login)
func CreateSiswa(c *fiber.Ctx) error {
	var input SiswaInput
	if err := c.BodyParser(&input); err != nil {
		return SendError(c, fiber.StatusBadRequest, "Input tidak valid")
	}

	if input.NISN == "" || input.NamaLengkap == "" || input.Password == "" {
		return SendError(c, fiber.StatusBadRequest, "NISN, Nama, dan Password wajib diisi")
	}

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		// 1. Hash Password
		hash, errHash := utils.HashPassword(input.Password)
		if errHash != nil {
			return errHash
		}

		// 2. Buat User (Login Akun)
		newUser := models.User{
			Username:      input.NISN, // NISN dijadikan Username
			Password:      hash,
			PasswordPlain: input.Password, // Simpan plaintext khusus siswa untuk cetak kartu
			Role:          "siswa",
			IsActive:      true,
		}

		if err := tx.Create(&newUser).Error; err != nil {
			return err
		}

		// 3. Buat Master Siswa
		var kelasID, ruangID, sesiID *uint
		if input.KelasID > 0 { kelasID = &input.KelasID }
		if input.RuangID > 0 { ruangID = &input.RuangID }
		if input.SesiID > 0 { sesiID = &input.SesiID }

		newSiswa := models.MasterSiswa{
			UserID:       newUser.ID,
			NISN:         input.NISN,
			NomorPeserta: input.NomorPeserta,
			NamaLengkap:  input.NamaLengkap,
			KelasID:      kelasID,
			RuangID:      ruangID,
			SesiID:       sesiID,
		}

		if err := tx.Create(&newSiswa).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return SendError(c, fiber.StatusInternalServerError, "Gagal menyimpan data siswa (Periksa apakah NISN sudah terdaftar)")
	}

	return SendSuccess(c, "Siswa berhasil ditambahkan", nil)
}

// UpdateSiswa mengubah data siswa
func UpdateSiswa(c *fiber.Ctx) error {
	id := c.Params("id")
	var input SiswaInput
	if err := c.BodyParser(&input); err != nil {
		return SendError(c, fiber.StatusBadRequest, "Input tidak valid")
	}

	var siswa models.MasterSiswa
	if err := database.DB.Preload("User").First(&siswa, id).Error; err != nil {
		return SendError(c, fiber.StatusNotFound, "Siswa tidak ditemukan")
	}

	// LOCKING: Cek apakah siswa sedang ikut ujian aktif
	var activeExamCount int64
	database.DB.Model(&models.CBTPesertaUjian{}).
		Joins("JOIN cbt_jadwal_ujians ON cbt_jadwal_ujians.id = cbt_peserta_ujians.jadwal_id").
		Where("cbt_peserta_ujians.siswa_id = ? AND cbt_jadwal_ujians.status = ?", id, "Berlangsung").
		Count(&activeExamCount)
	
	if activeExamCount > 0 {
		return SendError(c, fiber.StatusForbidden, "Data siswa tidak dapat diubah karena siswa sedang mengikuti ujian yang sedang berlangsung.")
	}

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		// Update tabel siswa
		siswa.NISN = input.NISN
		siswa.NomorPeserta = input.NomorPeserta
		siswa.NamaLengkap = input.NamaLengkap
		
		// Update FK jika dikirimkan (>0)
		if input.KelasID > 0 { 
			siswa.KelasID = &input.KelasID 
		} else {
			siswa.KelasID = nil
		}
		
		if input.RuangID > 0 { 
			siswa.RuangID = &input.RuangID 
		} else {
			siswa.RuangID = nil
		}
		
		if input.SesiID > 0 { 
			siswa.SesiID = &input.SesiID 
		} else {
			siswa.SesiID = nil
		}
		
		if err := tx.Save(&siswa).Error; err != nil {
			return err
		}

		// Update tabel user (Username & Password)
		siswa.User.Username = siswa.NISN
		if input.Password != "" {
			hash, errHash := utils.HashPassword(input.Password)
			if errHash != nil {
				return errHash
			}
			siswa.User.Password = hash
			siswa.User.PasswordPlain = input.Password
		}
		
		if err := tx.Save(&siswa.User).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return SendError(c, fiber.StatusInternalServerError, "Gagal memperbarui data siswa")
	}

	return SendSuccess(c, "Data siswa berhasil diperbarui", nil)
}
// DeleteSiswa menghapus siswa dan akun loginnya secara permanen
func DeleteSiswa(c *fiber.Ctx) error {
	id := c.Params("id")
	
	var siswa models.MasterSiswa
	if err := database.DB.First(&siswa, id).Error; err != nil {
		return SendError(c, fiber.StatusNotFound, "Siswa tidak ditemukan")
	}

	// Gunakan transaksi untuk menghapus di kedua tabel
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		// 1. Hapus profil MasterSiswa
		if err := tx.Delete(&siswa).Error; err != nil {
			return err
		}
		// 2. Hapus akun User-nya
		if err := tx.Delete(&models.User{}, siswa.UserID).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		errorMessage := "Gagal menghapus data siswa secara permanen"
		if strings.Contains(err.Error(), "FOREIGN KEY") {
			errorMessage = "Siswa tidak dapat dihapus karena sudah memiliki data ujian/nilai. Silakan hapus data ujian terkait terlebih dahulu."
		}
		return SendError(c, fiber.StatusInternalServerError, errorMessage)
	}

	return SendSuccess(c, "Siswa berhasil dihapus secara permanen", nil)
}

// ToggleSiswaStatus mengaktifkan atau menonaktifkan akun siswa
func ToggleSiswaStatus(c *fiber.Ctx) error {
	id := c.Params("id")
	var siswa models.MasterSiswa
	if err := database.DB.Preload("User").First(&siswa, id).Error; err != nil {
		return SendError(c, fiber.StatusNotFound, "Siswa tidak ditemukan")
	}

	// Balikkan status is_active
	siswa.User.IsActive = !siswa.User.IsActive
	database.DB.Save(&siswa.User)

	status := "diaktifkan"
	if !siswa.User.IsActive {
		status = "dinonaktifkan"
	}

	return SendSuccess(c, "Akun siswa berhasil " + status, nil)
}

// ImportSiswa memproses mass-upload siswa via CSV (Turbo Parallel)
func ImportSiswa(c *fiber.Ctx) error {
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
	type TempSiswa struct {
		nisn         string
		nomorPeserta string
		nama         string
		password     string
		namaKelas    string
	}
	var rawData []TempSiswa
	for i, row := range records {
		if i == 0 { continue }
		if len(row) < 4 { continue } // Minimal NISN, NoPeserta, Nama, Password
		rawData = append(rawData, TempSiswa{
			nisn:         strings.TrimSpace(row[0]),
			nomorPeserta: strings.TrimSpace(row[1]),
			nama:         strings.TrimSpace(row[2]),
			password:     strings.TrimSpace(row[3]),
			namaKelas:    "",
		})
		if len(row) >= 5 {
			rawData[len(rawData)-1].namaKelas = strings.TrimSpace(row[4])
		}
	}

	// 2. Hashing Paralel (Inilah yang bikin CEPAT)
	users := make([]models.User, len(rawData))
	var wg sync.WaitGroup
	for i := range rawData {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			hash, _ := utils.HashPassword(rawData[idx].password)
			users[idx] = models.User{
				Username:      rawData[idx].nisn,
				Password:      hash,
				PasswordPlain: rawData[idx].password, // Simpan password asli untuk cetak kartu
				Role:          "siswa",
				IsActive:      true,
			}
		}(i)
	}
	wg.Wait()

	// 3. Simpan ke Database (Batch)
	errTx := database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&users).Error; err != nil {
			return err
		}

		var allKelas []models.MasterKelas
		tx.Find(&allKelas)
		mapKelas := make(map[string]uint)
		for _, k := range allKelas {
			mapKelas[strings.ToLower(k.NamaKelas)] = k.ID
		}

		var siswas []models.MasterSiswa
		for i, u := range users {
			t := rawData[i]
			var kID *uint
			if id, ok := mapKelas[strings.ToLower(t.namaKelas)]; ok {
				kID = &id
			}
			siswas = append(siswas, models.MasterSiswa{
				UserID:       u.ID,
				NISN:         t.nisn,
				NomorPeserta: t.nomorPeserta,
				NamaLengkap:  t.nama,
				KelasID:      kID,
			})
		}
		return tx.Create(&siswas).Error
	})

	if errTx != nil {
		return SendError(c, 500, "Gagal simpan (NISN duplikat atau data salah)")
	}

	return SendSuccess(c, fmt.Sprintf("Berhasil impor %d siswa", len(users)), nil)
}

// BulkPlotSiswa memperbarui Ruang dan Sesi banyak siswa sekaligus
func BulkPlotSiswa(c *fiber.Ctx) error {
	var input BulkPlotInput
	if err := c.BodyParser(&input); err != nil {
		return SendError(c, fiber.StatusBadRequest, "Input tidak valid")
	}

	if len(input.SiswaIDs) == 0 {
		return SendError(c, fiber.StatusBadRequest, "Pilih minimal satu siswa")
	}

	var ruangID, sesiID *uint
	if input.RuangID > 0 { ruangID = &input.RuangID }
	if input.SesiID > 0 { sesiID = &input.SesiID }

	// Update secara massal menggunakan GORM
	err := database.DB.Model(&models.MasterSiswa{}).
		Where("id IN ?", input.SiswaIDs).
		Updates(map[string]interface{}{
			"ruang_id": ruangID,
			"sesi_id":  sesiID,
		}).Error

	if err != nil {
		return SendError(c, fiber.StatusInternalServerError, "Gagal memplot siswa secara massal")
	}

	return SendSuccess(c, "Berhasil memplot siswa secara massal", nil)
}

type BulkStatusInput struct {
	SiswaIDs []uint `json:"siswa_ids" validate:"required"`
	IsActive bool   `json:"is_active"`
}

// BulkUpdateSiswaStatus mengaktifkan atau menonaktifkan banyak akun siswa sekaligus
func BulkUpdateSiswaStatus(c *fiber.Ctx) error {
	var input BulkStatusInput
	if err := c.BodyParser(&input); err != nil {
		return SendError(c, fiber.StatusBadRequest, "Input tidak valid")
	}

	if len(input.SiswaIDs) == 0 {
		return SendError(c, fiber.StatusBadRequest, "Pilih minimal satu siswa")
	}

	// Update tabel users melalui subquery agar efisien
	err := database.DB.Model(&models.User{}).
		Where("id IN (SELECT user_id FROM master_siswas WHERE id IN ?)", input.SiswaIDs).
		Update("is_active", input.IsActive).Error

	if err != nil {
		return SendError(c, fiber.StatusInternalServerError, "Gagal memperbarui status secara masal")
	}

	status := "diaktifkan"
	if !input.IsActive {
		status = "dinonaktifkan"
	}

	return SendSuccess(c, fmt.Sprintf("Berhasil %s %d siswa", status, len(input.SiswaIDs)), nil)
}
