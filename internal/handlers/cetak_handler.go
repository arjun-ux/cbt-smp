package handlers

import (
	"cbt-smp/internal/database"
	"cbt-smp/internal/models"
	"github.com/gofiber/fiber/v2"
)

type KartuResponse struct {
	Sekolah map[string]string `json:"sekolah"`
	Siswa   []SiswaKartu      `json:"siswa"`
}

type SiswaKartu struct {
	ID           uint   `json:"id"`
	NamaLengkap  string `json:"nama_lengkap"`
	NomorPeserta string `json:"nomor_peserta"`
	NISN         string `json:"nisn"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	Kelas        string `json:"kelas"`
	Ruang        string `json:"ruang"`
	Sesi         string `json:"sesi"`
}

// GetCetakKartu mengambil data lengkap untuk pencetakan kartu ujian
func GetCetakKartu(c *fiber.Ctx) error {
	kelasID := c.Query("kelas_id")
	ruangID := c.Query("ruang_id")
	sesiID := c.Query("sesi_id")

	// 1. Ambil Pengaturan Sekolah
	var settings []models.CBTSetting
	database.DB.Find(&settings)
	settingsMap := make(map[string]string)
	for _, s := range settings {
		settingsMap[s.Key] = s.Value
	}

	// 2. Ambil Detail Kepala Sekolah dari tabel Guru jika ada kepala_sekolah_id
	if kepsekID, ok := settingsMap["kepala_sekolah_id"]; ok && kepsekID != "" {
		var guru models.MasterGuru
		if err := database.DB.Where("id = ?", kepsekID).First(&guru).Error; err == nil {
			settingsMap["kepala_sekolah_nama"] = guru.NamaGuru
			settingsMap["kepala_sekolah_nip"] = guru.NIP
		}
	}

	// 3. Ambil Data Siswa berdasarkan filter
	var siswas []models.MasterSiswa
	query := database.DB.Preload("User").Preload("Kelas").Preload("Ruang").Preload("Sesi")

	if kelasID != "" {
		query = query.Where("kelas_id = ?", kelasID)
	}
	if ruangID != "" {
		query = query.Where("ruang_id = ?", ruangID)
	}
	if sesiID != "" {
		query = query.Where("sesi_id = ?", sesiID)
	}

	if err := query.Find(&siswas).Error; err != nil {
		return SendError(c, fiber.StatusInternalServerError, "Gagal mengambil data siswa")
	}

	// 4. Mapping ke struct response
	var listSiswa []SiswaKartu
	for _, s := range siswas {
		listSiswa = append(listSiswa, SiswaKartu{
			ID:           s.ID,
			NamaLengkap:  s.NamaLengkap,
			NomorPeserta: s.NomorPeserta,
			NISN:         s.NISN,
			Username:     s.User.Username,
			Password:     s.User.PasswordPlain, // Mengambil password asli yang disimpan khusus siswa
			Kelas:        s.Kelas.NamaKelas,
			Ruang:        s.Ruang.NamaRuang,
			Sesi:         s.Sesi.NamaSesi,
		})
	}

	return SendSuccess(c, "Data kartu ujian", KartuResponse{
		Sekolah: settingsMap,
		Siswa:   listSiswa,
	})
}
