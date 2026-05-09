package handlers

import (
	"cbt-smp/internal/database"
	"cbt-smp/internal/models"
	"cbt-smp/pkg/utils"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type AdminInput struct {
	Username string `json:"username"`
	Password string `json:"password"` // Opsional saat update
	IsActive bool   `json:"is_active"`
}

// GetAdmins mengambil semua data user dengan role admin
func GetAdmins(c *fiber.Ctx) error {
	var admins []models.User
	database.DB.Where("role = ?", "admin").Find(&admins)

	return c.JSON(fiber.Map{
		"data": admins,
	})
}

// CreateAdmin menambah admin baru
func CreateAdmin(c *fiber.Ctx) error {
	var input AdminInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Input tidak valid"})
	}

	if input.Username == "" || input.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Username dan Password wajib diisi"})
	}

	hash, err := utils.HashPassword(input.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal enkripsi password"})
	}

	newAdmin := models.User{
		Username: input.Username,
		Password: hash,
		Role:     "admin",
		IsActive: true,
	}

	if err := database.DB.Create(&newAdmin).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal simpan admin (Username mungkin sudah ada)"})
	}

	return c.JSON(fiber.Map{"message": "Admin berhasil ditambahkan"})
}

// UpdateAdmin mengubah password atau status admin
func UpdateAdmin(c *fiber.Ctx) error {
	id := c.Params("id")
	var input AdminInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Input tidak valid"})
	}

	var admin models.User
	if err := database.DB.Where("id = ? AND role = ?", id, "admin").First(&admin).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Admin tidak ditemukan"})
	}

	if input.Username != "" {
		admin.Username = input.Username
	}
	
	if input.Password != "" {
		hash, err := utils.HashPassword(input.Password)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal enkripsi password"})
		}
		admin.Password = hash
	}

	admin.IsActive = input.IsActive

	if err := database.DB.Save(&admin).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal memperbarui admin"})
	}

	return c.JSON(fiber.Map{"message": "Admin berhasil diperbarui"})
}

// DeleteAdmin menghapus admin
func DeleteAdmin(c *fiber.Ctx) error {
	id := c.Params("id")
	
	var admin models.User
	if err := database.DB.Where("id = ? AND role = ?", id, "admin").First(&admin).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Admin tidak ditemukan"})
	}

	if err := database.DB.Delete(&admin).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal menghapus admin"})
	}

	return c.JSON(fiber.Map{"message": "Admin berhasil dihapus"})
}
// MigrateBase64Images menyisir seluruh soal dan memindahkan base64 ke file fisik
func MigrateBase64Images(c *fiber.Ctx) error {
	var soals []models.CBTSoal
	if err := database.DB.Find(&soals).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal mengambil data soal"})
	}

	count := 0
	for i := range soals {
		// Cek apakah ada base64 di pertanyaan atau opsi manapun
		hasBase64 := strings.Contains(soals[i].Pertanyaan, "data:image") ||
			strings.Contains(soals[i].OpsiA, "data:image") ||
			strings.Contains(soals[i].OpsiB, "data:image") ||
			strings.Contains(soals[i].OpsiC, "data:image") ||
			strings.Contains(soals[i].OpsiD, "data:image")

		if hasBase64 {
			soals[i].Pertanyaan = utils.ExtractBase64Images(soals[i].Pertanyaan)
			soals[i].OpsiA = utils.ExtractBase64Images(soals[i].OpsiA)
			soals[i].OpsiB = utils.ExtractBase64Images(soals[i].OpsiB)
			soals[i].OpsiC = utils.ExtractBase64Images(soals[i].OpsiC)
			soals[i].OpsiD = utils.ExtractBase64Images(soals[i].OpsiD)

			database.DB.Save(&soals[i])
			count++
		}
	}

	return c.JSON(fiber.Map{
		"message": "Migrasi selesai",
		"data": fiber.Map{
			"total_soal_dibersihkan": count,
		},
	})
}
