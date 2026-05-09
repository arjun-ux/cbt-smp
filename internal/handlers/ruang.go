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
	return c.JSON(fiber.Map{"data": data})
}

func CreateRuang(c *fiber.Ctx) error {
	var data models.MasterRuang
	if err := c.BodyParser(&data); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Input tidak valid"})
	}
	database.DB.Create(&data)
	return c.JSON(fiber.Map{"message": "Ruang berhasil dibuat", "data": data})
}

func UpdateRuang(c *fiber.Ctx) error {
	id := c.Params("id")
	var data models.MasterRuang
	if err := database.DB.First(&data, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Data tidak ditemukan"})
	}

	// LOCKING: Cek apakah ruangan sedang digunakan ujian aktif
	var activeExamCount int64
	database.DB.Model(&models.CBTJadwalUjian{}).
		Where("ruang_id = ? AND status = ?", id, "Berlangsung").
		Count(&activeExamCount)

	if activeExamCount > 0 {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Data ruangan tidak dapat diubah karena sedang digunakan untuk ujian yang sedang berlangsung."})
	}

	if err := c.BodyParser(&data); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Input tidak valid"})
	}
	database.DB.Save(&data)
	return c.JSON(fiber.Map{"message": "Ruang berhasil diupdate"})
}

func DeleteRuang(c *fiber.Ctx) error {
	id := c.Params("id")
	
	if err := database.DB.Delete(&models.MasterRuang{}, id).Error; err != nil {
		errorMessage := "Gagal menghapus data ruangan"
		if strings.Contains(err.Error(), "FOREIGN KEY") {
			errorMessage = "Ruangan tidak dapat dihapus karena masih digunakan oleh data Siswa atau Jadwal Ujian."
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": errorMessage})
	}
	
	return c.JSON(fiber.Map{"message": "Ruang berhasil dihapus"})
}
