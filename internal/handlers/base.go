package handlers

import (
	"cbt-smp/internal/database"
	"cbt-smp/internal/models"
	"github.com/gofiber/fiber/v2"
)

// GetUserID mengambil User ID dari context (JWT Claims)
func GetUserID(c *fiber.Ctx) uint {
	if val, ok := c.Locals("user_id").(float64); ok {
		return uint(val)
	} else if val, ok := c.Locals("user_id").(uint); ok {
		return val
	}
	return 0
}

// GetUserRole mengambil Role User dari context
func GetUserRole(c *fiber.Ctx) string {
	if val, ok := c.Locals("role").(string); ok {
		return val
	}
	return ""
}

// GetSiswaByUserID mengambil data master siswa berdasarkan User ID
func GetSiswaByUserID(userID uint) (models.MasterSiswa, error) {
	var siswa models.MasterSiswa
	err := database.DB.Preload("Kelas").Preload("Ruang").Preload("Sesi").
		Where("user_id = ?", userID).First(&siswa).Error
	return siswa, err
}

// GetGuruByUserID mengambil data master guru berdasarkan User ID
func GetGuruByUserID(userID uint) (models.MasterGuru, error) {
	var guru models.MasterGuru
	err := database.DB.Where("user_id = ?", userID).First(&guru).Error
	return guru, err
}

// SendError mengirim response error standar
func SendError(c *fiber.Ctx, status int, message string) error {
	return c.Status(status).JSON(fiber.Map{
		"error": message,
	})
}

// SendSuccess mengirim response sukses standar
func SendSuccess(c *fiber.Ctx, message string, data interface{}) error {
	return c.JSON(fiber.Map{
		"message": message,
		"data":    data,
	})
}
