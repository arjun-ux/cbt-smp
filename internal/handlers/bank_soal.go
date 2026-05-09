package handlers

import (
	"cbt-smp/internal/database"
	"cbt-smp/internal/models"
	"cbt-smp/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type BankSoalInput struct {
	GuruID            *uint   `json:"guru_id"`
	MapelID           uint    `json:"mapel_id"`
	TingkatKelas      string  `json:"tingkat_kelas"`
	JudulBankSoal     string  `json:"judul_bank_soal"`
	DefaultBobotPG    float64 `json:"default_bobot_pg"`
	DefaultBobotEssay float64 `json:"default_bobot_essay"`
	Status            string  `json:"status"`
}

// GetBankSoals mengambil daftar bank soal (Admin lihat semua, Guru lihat miliknya sendiri)
func GetBankSoals(c *fiber.Ctx) error {
	role := GetUserRole(c)
	userID := GetUserID(c)

	var bankSoals []models.CBTBankSoal
	query := database.DB.Preload("Guru").Preload("Mapel")

	if role == "guru" {
		guru, err := GetGuruByUserID(userID)
		if err != nil {
			return SendError(c, fiber.StatusForbidden, "Data Guru tidak ditemukan")
		}
		query = query.Where("guru_id = ?", guru.ID)
	}

	query.Find(&bankSoals)
	return SendSuccess(c, "Berhasil mengambil daftar bank soal", bankSoals)
}

// GetBankSoal mengambil detail satu bank soal
func GetBankSoal(c *fiber.Ctx) error {
	id := c.Params("id")
	role := GetUserRole(c)
	userID := GetUserID(c)

	var bs models.CBTBankSoal
	if err := database.DB.Preload("Guru").Preload("Mapel").First(&bs, id).Error; err != nil {
		return SendError(c, fiber.StatusNotFound, "Bank Soal tidak ditemukan")
	}

	// Proteksi: Guru hanya boleh melihat miliknya sendiri
	if role == "guru" {
		guru, _ := GetGuruByUserID(userID)
		if bs.GuruID == nil || *bs.GuruID != guru.ID {
			return SendError(c, fiber.StatusForbidden, "Akses ditolak: Ini bukan Bank Soal Anda")
		}
	}

	return SendSuccess(c, "Berhasil mengambil detail bank soal", bs)
}

// CreateBankSoal membuat wadah bank soal baru
func CreateBankSoal(c *fiber.Ctx) error {
	role := GetUserRole(c)
	userID := GetUserID(c)

	var input BankSoalInput
	if err := c.BodyParser(&input); err != nil {
		return SendError(c, fiber.StatusBadRequest, "Input tidak valid")
	}

	if input.MapelID == 0 || input.JudulBankSoal == "" || input.TingkatKelas == "" {
		return SendError(c, fiber.StatusBadRequest, "Mapel, Judul, dan Tingkat wajib diisi")
	}

	// Jika role guru, paksa guru_id ke ID guru si pembuat
	guruID := input.GuruID
	if role == "guru" {
		guru, err := GetGuruByUserID(userID)
		if err != nil {
			return SendError(c, fiber.StatusForbidden, "Data Guru tidak ditemukan")
		}
		guruID = &guru.ID
	}

	newBS := models.CBTBankSoal{
		GuruID:            guruID,
		MapelID:           input.MapelID,
		TingkatKelas:      input.TingkatKelas,
		JudulBankSoal:     input.JudulBankSoal,
		DefaultBobotPG:    input.DefaultBobotPG,
		DefaultBobotEssay: input.DefaultBobotEssay,
		Status:            "Draft",
	}

	if err := database.DB.Create(&newBS).Error; err != nil {
		return SendError(c, fiber.StatusInternalServerError, "Gagal membuat Bank Soal")
	}

	return SendSuccess(c, "Bank Soal berhasil dibuat", newBS)
}

// UpdateBankSoal mengubah data bank soal
func UpdateBankSoal(c *fiber.Ctx) error {
	id := c.Params("id")

	var bs models.CBTBankSoal
	if err := database.DB.First(&bs, id).Error; err != nil {
		return SendError(c, fiber.StatusNotFound, "Bank Soal tidak ditemukan")
	}

	// LOCKING: Cek apakah sedang digunakan ujian aktif
	var activeJadwal int64
	database.DB.Model(&models.CBTJadwalUjian{}).Where("bank_soal_id = ? AND status = ?", id, "Berlangsung").Count(&activeJadwal)
	if activeJadwal > 0 {
		return SendError(c, fiber.StatusForbidden, "Bank Soal tidak dapat diubah karena sedang digunakan dalam ujian yang sedang berlangsung.")
	}

	role := GetUserRole(c)
	userID := GetUserID(c)

	var input BankSoalInput
	if err := c.BodyParser(&input); err != nil {
		return SendError(c, fiber.StatusBadRequest, "Input tidak valid")
	}

	// Proteksi Guru
	if role == "guru" {
		guru, _ := GetGuruByUserID(userID)
		if bs.GuruID == nil || *bs.GuruID != guru.ID {
			return SendError(c, fiber.StatusForbidden, "Akses ditolak: Ini bukan Bank Soal Anda")
		}
		input.GuruID = bs.GuruID
	}

	bs.GuruID = input.GuruID
	bs.MapelID = input.MapelID
	bs.TingkatKelas = input.TingkatKelas
	bs.JudulBankSoal = input.JudulBankSoal
	bs.DefaultBobotPG = input.DefaultBobotPG
	bs.DefaultBobotEssay = input.DefaultBobotEssay
	if input.Status != "" {
		bs.Status = input.Status
	}

	if err := database.DB.Save(&bs).Error; err != nil {
		return SendError(c, fiber.StatusInternalServerError, "Gagal mengupdate Bank Soal")
	}

	return SendSuccess(c, "Bank Soal berhasil diupdate", bs)
}

// DeleteBankSoal menghapus bank soal
func DeleteBankSoal(c *fiber.Ctx) error {
	id := c.Params("id")

	var bs models.CBTBankSoal
	if err := database.DB.First(&bs, id).Error; err != nil {
		return SendError(c, fiber.StatusNotFound, "Bank Soal tidak ditemukan")
	}

	// LOCKING: Cek apakah sedang digunakan ujian aktif
	var activeJadwal int64
	database.DB.Model(&models.CBTJadwalUjian{}).Where("bank_soal_id = ? AND status = ?", id, "Berlangsung").Count(&activeJadwal)
	if activeJadwal > 0 {
		return SendError(c, fiber.StatusForbidden, "Bank Soal tidak dapat dihapus karena sedang digunakan dalam ujian yang sedang berlangsung.")
	}

	role := GetUserRole(c)
	userID := GetUserID(c)

	// Proteksi Guru
	if role == "guru" {
		guru, _ := GetGuruByUserID(userID)
		if bs.GuruID == nil || *bs.GuruID != guru.ID {
			return SendError(c, fiber.StatusForbidden, "Akses ditolak: Ini bukan Bank Soal Anda")
		}
	}

	// Mulai Transaksi Sapu Bersih agar soal di dalamnya ikut terhapus
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		// 1. HAPUS GAMBAR FISIK: Ambil semua soal dulu untuk didata gambarnya
		var soals []models.CBTSoal
		tx.Where("bank_soal_id = ?", id).Find(&soals)
		for _, s := range soals {
			utils.DeleteImagesFromHTML(s.Pertanyaan)
			utils.DeleteImagesFromHTML(s.OpsiA)
			utils.DeleteImagesFromHTML(s.OpsiB)
			utils.DeleteImagesFromHTML(s.OpsiC)
			utils.DeleteImagesFromHTML(s.OpsiD)
			utils.DeleteImagesFromHTML(s.KunciJawaban)
		}

		// 2. Hapus semua butir soal dari database
		if err := tx.Where("bank_soal_id = ?", id).Delete(&models.CBTSoal{}).Error; err != nil {
			return err
		}
		// 3. Hapus wadah bank soalnya
		if err := tx.Delete(&bs).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return SendError(c, fiber.StatusBadRequest, "Gagal menghapus: Bank Soal masih digunakan dalam jadwal ujian.")
	}

	return SendSuccess(c, "Bank Soal berhasil dihapus", nil)
}
