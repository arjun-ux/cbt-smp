package handlers

import (
	"cbt-smp/internal/database"
	"cbt-smp/internal/models"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func GetRuang(c *fiber.Ctx) error {
	var data []models.MasterRuang
	database.DB.Find(&data)
	return SendSuccess(c, "Berhasil mengambil daftar ruang", data)
}

func CreateRuang(c *fiber.Ctx) error {
	var data models.MasterRuang
	if err := c.BodyParser(&data); err != nil {
		return SendError(c, 400, "Input tidak valid")
	}
	database.DB.Create(&data)
	return SendSuccess(c, "Ruang berhasil dibuat", data)
}

func UpdateRuang(c *fiber.Ctx) error {
	id := c.Params("id")
	var data models.MasterRuang
	if err := database.DB.First(&data, id).Error; err != nil {
		return SendError(c, 404, "Data tidak ditemukan")
	}

	// LOCKING: Cek apakah ruangan sedang digunakan ujian aktif
	var activeExamCount int64
	database.DB.Model(&models.CBTJadwalUjian{}).
		Where("ruang_id = ? AND status = ?", id, "Berlangsung").
		Count(&activeExamCount)

	if activeExamCount > 0 {
		return SendError(c, fiber.StatusForbidden, "Data ruangan tidak dapat diubah karena sedang digunakan untuk ujian yang sedang berlangsung.")
	}

	if err := c.BodyParser(&data); err != nil {
		return SendError(c, 400, "Input tidak valid")
	}
	database.DB.Save(&data)
	return SendSuccess(c, "Ruang berhasil diupdate", nil)
}

func DeleteRuang(c *fiber.Ctx) error {
	id := c.Params("id")
	
	if err := database.DB.Delete(&models.MasterRuang{}, id).Error; err != nil {
		errorMessage := "Gagal menghapus data ruangan"
		if strings.Contains(err.Error(), "FOREIGN KEY") {
			errorMessage = "Ruangan tidak dapat dihapus karena masih digunakan oleh data Siswa atau Jadwal Ujian."
		}
		return SendError(c, fiber.StatusInternalServerError, errorMessage)
	}
	
	return SendSuccess(c, "Ruang berhasil dihapus", nil)
}
