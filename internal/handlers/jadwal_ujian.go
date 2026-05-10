package handlers

import (
	"cbt-smp/internal/database"
	"cbt-smp/internal/models"
	"math/rand"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type JadwalInput struct {
	BankSoalID   uint   `json:"bank_soal_id"`
	TanggalUjian string `json:"tanggal_ujian"`
	WaktuMulai   string `json:"waktu_mulai"`
	DurasiMenit  int    `json:"durasi_menit"`
	RuangID      uint   `json:"ruang_id"`
	SesiID       uint   `json:"sesi_id"`
	AcakSoal     bool   `json:"acak_soal"`
	AcakJawaban  bool   `json:"acak_jawaban"`
	Status       string `json:"status"`
}

// GetJadwals mengambil daftar jadwal ujian (Admin lihat semua, Guru lihat jadwal yang menggunakan Bank Soal miliknya)
func GetJadwals(c *fiber.Ctx) error {
	role := GetUserRole(c)
	userID := GetUserID(c)

	var jadwals []models.CBTJadwalUjian
	query := database.DB.Preload("BankSoal").
		Preload("BankSoal.Mapel").
		Preload("BankSoal.Guru").
		Preload("Pengawas").
		Preload("Ruang").
		Preload("Sesi")

	if role == "guru" {
		guru, _ := GetGuruByUserID(userID)
		// Guru hanya melihat jadwal di mana dia ditugaskan sebagai PENGAWAS
		query = query.Where("pengawas_id = ?", guru.ID)
	}

	if err := query.Order("tanggal_ujian desc, waktu_mulai desc").Find(&jadwals).Error; err != nil {
		return SendError(c, fiber.StatusInternalServerError, err.Error())
	}

	return SendSuccess(c, "Berhasil mengambil daftar jadwal", jadwals)
}

// GetJadwal mengambil satu detail jadwal
func GetJadwal(c *fiber.Ctx) error {
	id := c.Params("id")
	role := GetUserRole(c)
	userID := GetUserID(c)

	var jadwal models.CBTJadwalUjian
	if err := database.DB.Preload("BankSoal").
		Preload("BankSoal.Mapel").
		Preload("Pengawas").
		Preload("Ruang").
		Preload("Sesi").
		First(&jadwal, id).Error; err != nil {
		return SendError(c, fiber.StatusNotFound, "Jadwal tidak ditemukan")
	}

	// Proteksi Guru
	if role == "guru" {
		guru, _ := GetGuruByUserID(userID)
		var bs models.CBTBankSoal
		database.DB.First(&bs, jadwal.BankSoalID)

		isOwner := bs.GuruID != nil && *bs.GuruID == guru.ID
		isPengawas := jadwal.PengawasID != nil && *jadwal.PengawasID == guru.ID

		if !isOwner && !isPengawas {
			return SendError(c, fiber.StatusForbidden, "Akses ditolak: Anda bukan pemilik soal atau pengawas jadwal ini")
		}
	}

	return SendSuccess(c, "Berhasil mengambil detail jadwal", jadwal)
}

// UpdatePengawas khusus untuk mengubah pengawas pada jadwal
func UpdatePengawas(c *fiber.Ctx) error {
	id := c.Params("id")
	type Input struct {
		PengawasID *uint `json:"pengawas_id"`
	}
	var input Input
	if err := c.BodyParser(&input); err != nil {
		return SendError(c, fiber.StatusBadRequest, "Input tidak valid")
	}

	var jadwal models.CBTJadwalUjian
	if err := database.DB.First(&jadwal, id).Error; err != nil {
		return SendError(c, fiber.StatusNotFound, "Jadwal tidak ditemukan")
	}

	if jadwal.Status == "Berlangsung" {
		return SendError(c, fiber.StatusForbidden, "Tidak dapat mengubah pengawas saat ujian sedang berlangsung.")
	}

	if err := database.DB.Model(&jadwal).Update("pengawas_id", input.PengawasID).Error; err != nil {
		return SendError(c, fiber.StatusInternalServerError, "Gagal update pengawas")
	}

	return SendSuccess(c, "Pengawas berhasil diperbarui", nil)
}

// CreateJadwal membuat jadwal ujian baru
func CreateJadwal(c *fiber.Ctx) error {
	role := GetUserRole(c)
	userID := GetUserID(c)

	var input JadwalInput
	if err := c.BodyParser(&input); err != nil {
		return SendError(c, fiber.StatusBadRequest, "Input tidak valid")
	}

	if input.BankSoalID == 0 || input.TanggalUjian == "" || input.WaktuMulai == "" || input.DurasiMenit == 0 {
		return SendError(c, fiber.StatusBadRequest, "Semua field wajib diisi")
	}

	// Proteksi Guru: Pastikan BankSoal milik guru ini
	if role == "guru" {
		guru, _ := GetGuruByUserID(userID)
		var bs models.CBTBankSoal
		if err := database.DB.Where("id = ? AND guru_id = ?", input.BankSoalID, guru.ID).First(&bs).Error; err != nil {
			return SendError(c, fiber.StatusForbidden, "Akses ditolak: Anda tidak bisa membuat jadwal untuk Bank Soal milik orang lain")
		}
	}

	// Generate Token acak (5 karakter huruf besar)
	token := generateToken(5)

	var ruangID, sesiID *uint
	if input.RuangID > 0 {
		ruangID = &input.RuangID
	}
	if input.SesiID > 0 {
		sesiID = &input.SesiID
	}

	newJadwal := models.CBTJadwalUjian{
		BankSoalID:   input.BankSoalID,
		TanggalUjian: input.TanggalUjian,
		WaktuMulai:   input.WaktuMulai,
		DurasiMenit:  input.DurasiMenit,
		RuangID:      ruangID,
		SesiID:       sesiID,
		AcakSoal:     input.AcakSoal,
		AcakJawaban:  input.AcakJawaban,
		TokenUjian:   token,
		Status:       "Belum Mulai",
	}

	if err := database.DB.Create(&newJadwal).Error; err != nil {
		return SendError(c, fiber.StatusInternalServerError, "Gagal membuat jadwal")
	}

	return SendSuccess(c, "Jadwal berhasil dibuat", newJadwal)
}

// UpdateJadwal mengubah data jadwal
func UpdateJadwal(c *fiber.Ctx) error {
	id := c.Params("id")
	role := GetUserRole(c)
	userID := GetUserID(c)

	var input JadwalInput
	if err := c.BodyParser(&input); err != nil {
		return SendError(c, fiber.StatusBadRequest, "Input tidak valid")
	}

	var jadwal models.CBTJadwalUjian
	if err := database.DB.First(&jadwal, id).Error; err != nil {
		return SendError(c, fiber.StatusNotFound, "Jadwal tidak ditemukan")
	}

	// Proteksi Guru
	if role == "guru" {
		guru, _ := GetGuruByUserID(userID)
		var bs models.CBTBankSoal
		database.DB.First(&bs, jadwal.BankSoalID)
		if bs.GuruID == nil || *bs.GuruID != guru.ID {
			return SendError(c, fiber.StatusForbidden, "Akses ditolak: Ini bukan jadwal Anda")
		}
	}

	// PROTEKSI: Jangan izinkan edit jika sedang berlangsung
	if jadwal.Status == "Berlangsung" && input.Status != "Selesai" {
		return SendError(c, fiber.StatusBadRequest, "Jadwal sedang berlangsung. Anda tidak dapat mengubah data jadwal kecuali untuk mengakhiri ujian.")
	}

	var ruangID, sesiID *uint
	if input.RuangID > 0 {
		ruangID = &input.RuangID
	}
	if input.SesiID > 0 {
		sesiID = &input.SesiID
	}

	jadwal.BankSoalID = input.BankSoalID
	jadwal.TanggalUjian = input.TanggalUjian
	jadwal.WaktuMulai = input.WaktuMulai
	jadwal.DurasiMenit = input.DurasiMenit
	jadwal.RuangID = ruangID
	jadwal.SesiID = sesiID
	jadwal.AcakSoal = input.AcakSoal
	jadwal.AcakJawaban = input.AcakJawaban
	if input.Status != "" {
		jadwal.Status = input.Status
	}

	// Gunakan transaksi untuk menjamin integritas data saat force-submit masal
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		// 1. Simpan perubahan Jadwal
		if err := tx.Save(&jadwal).Error; err != nil {
			return err
		}

		// 2. Jika status diubah menjadi 'Selesai', lakukan Force Submit otomatis untuk semua siswa
		if input.Status == "Selesai" {
			var pesertaBelumSelesai []models.CBTPesertaUjian
			// Ambil hanya peserta yang tidak terblokir untuk menghindari tabrakan logika/database
			if err := tx.Where("jadwal_id = ? AND status_ujian = ? AND is_terblokir = 0", jadwal.ID, "Sedang Mengerjakan").Find(&pesertaBelumSelesai).Error; err != nil {
				return err
			}

			now := time.Now()
			for _, p := range pesertaBelumSelesai {
				// Hitung Nilai (Helper) - Sekarang menggunakan transaksi 'tx' yang sama
				nilaiPG := HitungNilaiPG(tx, p.ID, jadwal.BankSoalID)

				// Update data peserta
				if err := tx.Model(&p).Updates(map[string]interface{}{
					"status_ujian":        "Selesai",
					"waktu_selesai_ujian": &now,
					"nilai_pg":            nilaiPG,
					"total_nilai":         nilaiPG + p.NilaiEssay,
				}).Error; err != nil {
					return err
				}

				// Catat Log System
				tx.Create(&models.LogUjian{
					PesertaUjianID: p.ID,
					KeteranganLog:  "Force Submit Otomatis (Jadwal Diakhiri oleh Admin/Guru)",
				})
			}
		}

		return nil
	})

	if err != nil {
		return SendError(c, fiber.StatusInternalServerError, "Gagal mengupdate jadwal dan memproses finalisasi siswa: "+err.Error())
	}

	return SendSuccess(c, "Jadwal berhasil diupdate", jadwal)
}

// DeleteJadwal menghapus jadwal
func DeleteJadwal(c *fiber.Ctx) error {
	id := c.Params("id")
	role := GetUserRole(c)
	userID := GetUserID(c)

	var jadwal models.CBTJadwalUjian
	if err := database.DB.First(&jadwal, id).Error; err != nil {
		return SendError(c, fiber.StatusNotFound, "Jadwal tidak ditemukan")
	}

	// Proteksi Guru
	if role == "guru" {
		guru, _ := GetGuruByUserID(userID)
		var bs models.CBTBankSoal
		database.DB.First(&bs, jadwal.BankSoalID)
		if bs.GuruID == nil || *bs.GuruID != guru.ID {
			return SendError(c, fiber.StatusForbidden, "Akses ditolak: Anda tidak berhak menghapus jadwal ini")
		}
	}

	if jadwal.Status == "Berlangsung" {
		return SendError(c, fiber.StatusForbidden, "Tidak dapat menghapus jadwal yang sedang berlangsung")
	}

	// SAFE-DELETE: Jika jadwal sudah selesai, pastikan sudah di-rekap
	if jadwal.Status == "Selesai" {
		var count int64
		database.DB.Model(&models.CBTRekapNilai{}).Where("jadwal_id = ?", id).Count(&count)
		if count == 0 {
			return SendError(c, fiber.StatusForbidden, "Jadwal ini belum diarsipkan. Silakan lakukan Rekap/Arsip Nilai terlebih dahulu agar data nilai tidak hilang secara permanen.")
		}
	}

	if err := database.DB.Delete(&jadwal).Error; err != nil {
		return SendError(c, fiber.StatusInternalServerError, "Gagal menghapus jadwal")
	}
	return SendSuccess(c, "Jadwal berhasil dihapus", nil)
}

// ArchiveJadwalResults memindahkan nilai dari transaksi ke rekap permanen
func ArchiveJadwalResults(c *fiber.Ctx) error {
	id := c.Params("id")

	var jadwal models.CBTJadwalUjian
	if err := database.DB.Preload("BankSoal.Mapel").First(&jadwal, id).Error; err != nil {
		return SendError(c, 404, "Jadwal tidak ditemukan")
	}

	if jadwal.Status != "Selesai" {
		return SendError(c, 400, "Jadwal harus berstatus 'Selesai' sebelum bisa diarsipkan")
	}

	// Ambil semua peserta
	var peserta []models.CBTPesertaUjian
	database.DB.Preload("Siswa.Kelas").Where("jadwal_id = ?", id).Find(&peserta)

	if len(peserta) == 0 {
		return SendError(c, 400, "Tidak ada data peserta untuk diarsipkan")
	}

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		for _, p := range peserta {
			// Cek apakah sudah ada rekap
			var existing models.CBTRekapNilai
			if err := tx.Where("jadwal_id = ? AND nisn = ?", id, p.Siswa.NISN).First(&existing).Error; err == nil {
				continue
			}

			rekap := models.CBTRekapNilai{
				JadwalID:     jadwal.ID,
				NISN:         p.Siswa.NISN,
				NamaSiswa:    p.Siswa.NamaLengkap,
				NamaKelas:    p.Siswa.Kelas.NamaKelas,
				NamaMapel:    jadwal.BankSoal.Mapel.NamaMapel,
				JudulUjian:   jadwal.BankSoal.JudulBankSoal,
				TanggalUjian: jadwal.TanggalUjian,
				NilaiPG:      p.NilaiPG,
				NilaiEssay:   p.NilaiEssay,
				TotalNilai:   p.TotalNilai,
				CreatedAt:    time.Now(),
			}
			if err := tx.Create(&rekap).Error; err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return SendError(c, 500, "Gagal mengarsipkan nilai: "+err.Error())
	}

	// Update Status Jadwal menjadi 'Diarsipkan'
	database.DB.Model(&jadwal).Update("status", "Diarsipkan")

	return SendSuccess(c, "Nilai berhasil diarsipkan secara permanen", nil)
}

// GetRiwayatNilai mengambil seluruh data dari Brankas Nilai (CBTRekapNilai)
func GetRiwayatNilai(c *fiber.Ctx) error {
	role := GetUserRole(c)
	userID := GetUserID(c)

	var results []models.CBTRekapNilai
	query := database.DB.Order("tanggal_ujian desc, nama_siswa asc")

	if role == "guru" {
		guru, _ := GetGuruByUserID(userID)
		// Guru hanya bisa melihat riwayat nilai (brankas) untuk Bank Soal miliknya
		query = query.Where("jadwal_id IN (SELECT id FROM cbt_jadwal_ujians WHERE bank_soal_id IN (SELECT id FROM cbt_bank_soals WHERE guru_id = ?))", guru.ID)
	}

	if err := query.Find(&results).Error; err != nil {
		return SendError(c, 500, "Gagal mengambil riwayat nilai")
	}

	return SendSuccess(c, "Berhasil mengambil riwayat nilai", results)
}

// RefreshToken menghasilkan token baru untuk jadwal tertentu
func RefreshToken(c *fiber.Ctx) error {
	id := c.Params("id")
	role := GetUserRole(c)
	userID := GetUserID(c)

	var jadwal models.CBTJadwalUjian
	if err := database.DB.First(&jadwal, id).Error; err != nil {
		return SendError(c, fiber.StatusNotFound, "Jadwal tidak ditemukan")
	}

	// Proteksi Guru
	if role == "guru" {
		guru, _ := GetGuruByUserID(userID)
		var bs models.CBTBankSoal
		database.DB.First(&bs, jadwal.BankSoalID)
		if bs.GuruID == nil || *bs.GuruID != guru.ID {
			return SendError(c, fiber.StatusForbidden, "Akses ditolak")
		}
	}

	jadwal.TokenUjian = generateToken(5)
	database.DB.Save(&jadwal)

	return SendSuccess(c, "Token berhasil diperbarui", fiber.Map{"token": jadwal.TokenUjian})
}

// Helper untuk generate token
func generateToken(n int) string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	seed := rand.NewSource(time.Now().UnixNano())
	random := rand.New(seed)

	b := make([]byte, n)
	for i := range b {
		b[i] = charset[random.Intn(len(charset))]
	}
	return string(b)
}
