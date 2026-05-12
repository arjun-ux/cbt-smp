package handlers

import (
	"cbt-smp/internal/database"
	"cbt-smp/internal/models"

	"github.com/gofiber/fiber/v2"
)

// GetSettings mengambil semua pengaturan dalam bentuk map
func GetSettings(c *fiber.Ctx) error {
	var settings []models.CBTSetting
	if err := database.DB.Find(&settings).Error; err != nil {
		return SendError(c, fiber.StatusInternalServerError, "Gagal mengambil pengaturan")
	}

	// Ubah slice ke map untuk memudahkan frontend
	settingsMap := make(map[string]string)
	for _, s := range settings {
		settingsMap[s.Key] = s.Value
	}

	return SendSuccess(c, "Daftar pengaturan", settingsMap)
}

// UpdateSettings memperbarui atau membuat pengaturan secara batch
func UpdateSettings(c *fiber.Ctx) error {
	var payload map[string]string
	if err := c.BodyParser(&payload); err != nil {
		return SendError(c, fiber.StatusBadRequest, "Format data tidak valid")
	}

	// Gunakan transaction untuk memastikan semua terupdate atau tidak sama sekali
	tx := database.DB.Begin()

	for key, value := range payload {
		var setting models.CBTSetting
		// Cari berdasarkan key, jika ada update, jika tidak buat baru
		err := tx.Where("key = ?", key).First(&setting).Error
		if err == nil {
			setting.Value = value
			if err := tx.Save(&setting).Error; err != nil {
				tx.Rollback()
				return SendError(c, fiber.StatusInternalServerError, "Gagal memperbarui key: "+key)
			}
		} else {
			newSetting := models.CBTSetting{
				Key:   key,
				Value: value,
			}
			if err := tx.Create(&newSetting).Error; err != nil {
				tx.Rollback()
				return SendError(c, fiber.StatusInternalServerError, "Gagal membuat key: "+key)
			}
		}
	}

	if err := tx.Commit().Error; err != nil {
		return SendError(c, fiber.StatusInternalServerError, "Gagal menyimpan perubahan")
	}

	return SendSuccess(c, "Pengaturan berhasil diperbarui", nil)
}
