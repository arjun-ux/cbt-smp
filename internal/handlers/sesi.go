package handlers

import (
	"cbt-smp/internal/database"
	"cbt-smp/internal/models"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func GetSesi(c *fiber.Ctx) error {
	var data []models.MasterSesi
	database.DB.Find(&data)
	return c.JSON(fiber.Map{"data": data})
}

func CreateSesi(c *fiber.Ctx) error {
	var data models.MasterSesi
	if err := c.BodyParser(&data); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Input tidak valid"})
	}

	// Validasi Overlap Waktu
	var count int64
	database.DB.Model(&models.MasterSesi{}).
		Where("NOT (waktu_mulai >= ? OR waktu_selesai <= ?)", data.WaktuSelesai, data.WaktuMulai).
		Count(&count)
	
	if count > 0 {
		return c.Status(400).JSON(fiber.Map{"error": "Waktu sesi bentrok dengan sesi lain yang sudah ada"})
	}

	database.DB.Create(&data)
	return c.JSON(fiber.Map{"message": "Sesi berhasil dibuat", "data": data})
}

func UpdateSesi(c *fiber.Ctx) error {
	id := c.Params("id")
	var data models.MasterSesi
	if err := database.DB.First(&data, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Data tidak ditemukan"})
	}
	
	var input models.MasterSesi
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Input tidak valid"})
	}

	// Validasi Overlap Waktu (Kecuali sesi itu sendiri)
	var count int64
	database.DB.Model(&models.MasterSesi{}).
		Where("id <> ?", id).
		Where("NOT (waktu_mulai >= ? OR waktu_selesai <= ?)", input.WaktuSelesai, input.WaktuMulai).
		Count(&count)
	
	if count > 0 {
		return c.Status(400).JSON(fiber.Map{"error": "Update gagal! Waktu sesi bentrok dengan sesi lain"})
	}

	// LOCKING: Cek apakah sesi sedang digunakan ujian aktif
	var activeExamCount int64
	database.DB.Model(&models.CBTJadwalUjian{}).
		Where("sesi_id = ? AND status = ?", id, "Berlangsung").
		Count(&activeExamCount)

	if activeExamCount > 0 {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Data sesi tidak dapat diubah karena sedang digunakan untuk ujian yang sedang berlangsung."})
	}

	data.NamaSesi = input.NamaSesi
	data.WaktuMulai = input.WaktuMulai
	data.WaktuSelesai = input.WaktuSelesai
	
	database.DB.Save(&data)
	return c.JSON(fiber.Map{"message": "Sesi berhasil diupdate"})
}

func DeleteSesi(c *fiber.Ctx) error {
	id := c.Params("id")
	
	if err := database.DB.Delete(&models.MasterSesi{}, id).Error; err != nil {
		errorMessage := "Gagal menghapus data sesi"
		if strings.Contains(err.Error(), "FOREIGN KEY") {
			errorMessage = "Sesi tidak dapat dihapus karena masih digunakan oleh data Siswa atau Jadwal Ujian."
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": errorMessage})
	}
	
	return c.JSON(fiber.Map{"message": "Sesi berhasil dihapus"})
}
