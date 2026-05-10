package handlers

import (
	"cbt-smp/internal/database"
	"cbt-smp/internal/models"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func GetKelas(c *fiber.Ctx) error {
	var data []models.MasterKelas
	database.DB.Find(&data)
	return SendSuccess(c, "Berhasil mengambil daftar kelas", data)
}

func CreateKelas(c *fiber.Ctx) error {
	var data models.MasterKelas
	if err := c.BodyParser(&data); err != nil {
		return SendError(c, 400, "Input tidak valid")
	}
	database.DB.Create(&data)
	return SendSuccess(c, "Kelas berhasil dibuat", data)
}

func UpdateKelas(c *fiber.Ctx) error {
	id := c.Params("id")
	var data models.MasterKelas
	if err := database.DB.First(&data, id).Error; err != nil {
		return SendError(c, 404, "Data tidak ditemukan")
	}

	// LOCKING: Cek apakah ada siswa di kelas ini yang sedang ujian aktif
	var activeExamCount int64
	database.DB.Model(&models.CBTPesertaUjian{}).
		Joins("JOIN master_siswas ON master_siswas.id = cbt_peserta_ujians.siswa_id").
		Joins("JOIN cbt_jadwal_ujians ON cbt_jadwal_ujians.id = cbt_peserta_ujians.jadwal_id").
		Where("master_siswas.kelas_id = ? AND cbt_jadwal_ujians.status = ?", id, "Berlangsung").
		Count(&activeExamCount)

	if activeExamCount > 0 {
		return SendError(c, fiber.StatusForbidden, "Data kelas tidak dapat diubah karena ada siswa di kelas ini yang sedang mengikuti ujian aktif.")
	}

	if err := c.BodyParser(&data); err != nil {
		return SendError(c, 400, "Input tidak valid")
	}
	
	database.DB.Save(&data)
	return SendSuccess(c, "Kelas berhasil diupdate", nil)
}

func DeleteKelas(c *fiber.Ctx) error {
	id := c.Params("id")
	
	if err := database.DB.Delete(&models.MasterKelas{}, id).Error; err != nil {
		errorMessage := "Gagal menghapus data kelas"
		if strings.Contains(err.Error(), "FOREIGN KEY") {
			errorMessage = "Kelas tidak dapat dihapus karena masih memiliki data siswa di dalamnya. Kosongkan kelas terlebih dahulu sebelum dihapus."
		}
		return SendError(c, fiber.StatusInternalServerError, errorMessage)
	}
	
	return SendSuccess(c, "Kelas berhasil dihapus", nil)
}
