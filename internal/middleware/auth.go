package middleware

import (
	"cbt-smp/internal/database"
	"cbt-smp/internal/models"
	"cbt-smp/pkg/utils"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// Protected adalah middleware untuk memvalidasi token JWT
func Protected() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Ambil header Authorization
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Akses ditolak: Token tidak ditemukan",
			})
		}

		// Format token biasanya: "Bearer <token_string>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Akses ditolak: Format token salah",
			})
		}

		tokenString := parts[1]

		// Verifikasi JWT
		claims, err := utils.VerifyJWT(tokenString)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Akses ditolak: Token tidak valid atau sudah kadaluarsa",
			})
		}

		// VALIDASI DATABASE: Pastikan user benar-benar ada di DB yang sekarang
		userID := uint(claims["user_id"].(float64))
		var user models.User
		if err := database.DB.First(&user, userID).Error; err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Akses ditolak: Akun Anda tidak ditemukan di database ini. Silakan login kembali.",
				"code":  "USER_NOT_FOUND",
			})
		}

		// Simpan klaim token ke Locals agar bisa diakses di handler selanjutnya
		c.Locals("user_id", userID)
		c.Locals("username", user.Username)
		c.Locals("role", user.Role)

		return c.Next()
	}
}

// RoleRequired memastikan endpoint hanya diakses oleh peran tertentu (Role-Based Access Control)
func RoleRequired(allowedRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Ambil peran user dari Locals (diset oleh middleware Protected)
		userRole, ok := c.Locals("role").(string)
		if !ok {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Akses ditolak: Gagal memverifikasi hak akses",
			})
		}

		// Periksa apakah peran user ada di daftar yang diizinkan
		isAllowed := false
		for _, role := range allowedRoles {
			if userRole == role {
				isAllowed = true
				break
			}
		}

		if !isAllowed {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Akses ditolak: Hak akses tidak mencukupi untuk rute ini",
			})
		}

		return c.Next()
	}
}
