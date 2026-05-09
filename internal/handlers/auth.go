package handlers

import (
	"cbt-smp/internal/database"
	"cbt-smp/internal/models"
	"cbt-smp/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

// LoginInput struktur data request login dari frontend
type LoginInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Login memproses autentikasi pengguna
func Login(c *fiber.Ctx) error {
	var input LoginInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Format request tidak valid",
		})
	}

	// Cari user di database
	var user models.User
	result := database.DB.Where("username = ?", input.Username).First(&user)
	if result.Error != nil {
		// Coba cari berdasarkan NISN jika tidak ketemu di Username
		var siswa models.MasterSiswa
		if err := database.DB.Where("nisn = ?", input.Username).First(&siswa).Error; err == nil {
			database.DB.First(&user, siswa.UserID)
		} else {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "NISN/Username atau Password salah",
			})
		}
	}

	// Pastikan akun aktif
	if !user.IsActive {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Akun Anda tidak aktif, hubungi administrator",
		})
	}

	// Verifikasi Password Hash
	if !utils.CheckPasswordHash(input.Password, user.Password) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Username atau Password salah",
		})
	}

	// Generate JWT
	token, err := utils.GenerateJWT(user.ID, user.Username, user.Role)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal membuat sesi login",
		})
	}

	// (Opsional) Catat log sistem jika login berhasil
	// Pada aplikasi yang sibuk, sebaiknya diletakkan di Goroutine
	go func(uID uint, ip string) {
		logEntry := models.LogSistem{
			UserID:    uID,
			Aktivitas: "Login berhasil ke dalam sistem",
			IPAddress: ip,
		}
		database.DB.Create(&logEntry)
	}(user.ID, c.IP())

	// Siapkan data tambahan berdasarkan role
	userData := fiber.Map{
		"id":       user.ID,
		"username": user.Username,
		"role":     user.Role,
		"nama":     "",
	}

	if user.Role == "siswa" {
		var siswa models.MasterSiswa
		if err := database.DB.Preload("Kelas").Preload("Ruang").Preload("Sesi").Where("user_id = ?", user.ID).First(&siswa).Error; err == nil {
			userData["nama"] = siswa.NamaLengkap
			userData["kelas"] = siswa.Kelas.NamaKelas
			if siswa.RuangID != nil {
				userData["ruang"] = siswa.Ruang.NamaRuang
			} else {
				userData["ruang"] = "-"
			}
			if siswa.SesiID != nil {
				userData["sesi"] = siswa.Sesi.NamaSesi
			} else {
				userData["sesi"] = "-"
			}
		}
	} else if user.Role == "guru" {
		var guru models.MasterGuru
		if err := database.DB.Where("user_id = ?", user.ID).First(&guru).Error; err == nil {
			userData["nama"] = guru.NamaGuru
			userData["guru_id"] = guru.ID
		}
	} else {
		userData["nama"] = "Administrator"
	}

	// Kembalikan token dan data user (tanpa password)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Login berhasil",
		"token":   token,
		"user":     userData,
	})
}

// Logout memproses penghentian sesi di server
func Logout(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "Berhasil logout"})
}
