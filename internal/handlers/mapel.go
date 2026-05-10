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
	return SendSuccess(c, "Berhasil mengambil daftar mapel", data)
}

func CreateMapel(c *fiber.Ctx) error {
	var data models.MasterMapel
	if err := c.BodyParser(&data); err != nil {
		return SendError(c, 400, "Input tidak valid")
	}
	database.DB.Create(&data)
	return SendSuccess(c, "Mapel berhasil dibuat", data)
}

func UpdateMapel(c *fiber.Ctx) error {
	id := c.Params("id")
	var data models.MasterMapel
	if err := database.DB.First(&data, id).Error; err != nil {
		return SendError(c, 404, "Data tidak ditemukan")
	}
	c.BodyParser(&data)
	database.DB.Save(&data)
	return SendSuccess(c, "Mapel berhasil diupdate", nil)
}

func DeleteMapel(c *fiber.Ctx) error {
	id := c.Params("id")
	
	if err := database.DB.Delete(&models.MasterMapel{}, id).Error; err != nil {
		errorMessage := "Gagal menghapus mata pelajaran"
		if strings.Contains(err.Error(), "FOREIGN KEY") {
			errorMessage = "Mata Pelajaran tidak dapat dihapus karena masih terikat dengan Bank Soal atau Jadwal Ujian."
		}
		return SendError(c, fiber.StatusInternalServerError, errorMessage)
	}
	
	return SendSuccess(c, "Mapel berhasil dihapus", nil)
}
