package handlers

import (
	"cbt-smp/internal/database"
	"cbt-smp/internal/models"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func GetMapel(c *fiber.Ctx) error {
	var data []models.MasterMapel
	database.DB.Find(&data)
	return c.JSON(fiber.Map{"data": data})
}

func CreateMapel(c *fiber.Ctx) error {
	var data models.MasterMapel
	if err := c.BodyParser(&data); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Input tidak valid"})
	}
	database.DB.Create(&data)
	return c.JSON(fiber.Map{"message": "Mapel berhasil dibuat", "data": data})
}

func UpdateMapel(c *fiber.Ctx) error {
	id := c.Params("id")
	var data models.MasterMapel
	if err := database.DB.First(&data, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Data tidak ditemukan"})
	}
	c.BodyParser(&data)
	database.DB.Save(&data)
	return c.JSON(fiber.Map{"message": "Mapel berhasil diupdate"})
}

func DeleteMapel(c *fiber.Ctx) error {
	id := c.Params("id")
	
	if err := database.DB.Delete(&models.MasterMapel{}, id).Error; err != nil {
		errorMessage := "Gagal menghapus mata pelajaran"
		if strings.Contains(err.Error(), "FOREIGN KEY") {
			errorMessage = "Mata Pelajaran tidak dapat dihapus karena masih terikat dengan Bank Soal atau Jadwal Ujian."
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": errorMessage})
	}
	
	return c.JSON(fiber.Map{"message": "Mapel berhasil dihapus"})
}
