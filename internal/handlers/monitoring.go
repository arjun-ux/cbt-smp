package handlers

import (
	"cbt-smp/internal/database"
	"cbt-smp/internal/models"
	"cbt-smp/pkg/utils"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

type MonitorSiswa struct {
	ID                uint   `json:"id"`
	SiswaID           uint   `json:"siswa_id"`
	NamaSiswa         string `json:"nama_siswa"`
	Kelas             string `json:"kelas"`
	StatusUjian       string `json:"status_ujian"`
	Progres           int64  `json:"progres"`
	TotalSoal         int64  `json:"total_soal"`
	SisaWaktu         int    `json:"sisa_waktu"`
	JumlahPelanggaran int64  `json:"jumlah_pelanggaran"`
	IsTerblokir       bool   `json:"is_terblokir"`
}

// GetMonitorUjian mengambil data progres real-time untuk dashboard pengawas
func GetMonitorUjian(c *fiber.Ctx) error {
	jadwalID := c.Params("jadwalId")
	role := GetUserRole(c)
	userID := GetUserID(c)

	// Proteksi Akses Guru: Guru hanya bisa monitor jika dia adalah Pengawas jadwal tsb
	if role == "guru" {
		var jadwal models.CBTJadwalUjian
		if err := database.DB.First(&jadwal, jadwalID).Error; err != nil {
			return SendError(c, fiber.StatusNotFound, "Jadwal tidak ditemukan")
		}

		guru, _ := GetGuruByUserID(userID)
		if jadwal.PengawasID == nil || *jadwal.PengawasID != guru.ID {
			return SendError(c, fiber.StatusForbidden, "Anda tidak ditugaskan sebagai pengawas untuk ujian ini")
		}
	}

	// Optimasi: Ambil semua data progres dalam 1 Query (Subqueries)
	// Ini jauh lebih cepat daripada menghitung di dalam loop (N+1 Query)
	type MonitorRow struct {
		ID               uint   `gorm:"column:id"`
		SiswaID          uint   `gorm:"column:siswa_id"`
		NamaSiswa        string `gorm:"column:nama_siswa"`
		Kelas            string `gorm:"column:kelas"`
		StatusUjian      string `gorm:"column:status_ujian"`
		SisaWaktuDetik   int    `gorm:"column:sisa_waktu_detik"`
		IsTerblokir      bool   `gorm:"column:is_terblokir"`
		CountJawaban     int64  `gorm:"column:count_jawaban"`
		CountPelanggaran int64  `gorm:"column:count_pelanggaran"`
		BankSoalID       uint   `gorm:"column:bank_soal_id"`
	}

	var rows []MonitorRow
	err := database.DB.Raw(`
		SELECT 
			p.id, p.siswa_id, p.status_ujian, p.sisa_waktu_detik, p.is_terblokir,
			s.nama_lengkap as nama_siswa,
			k.nama_kelas as kelas,
			j.bank_soal_id,
			(SELECT COUNT(*) FROM cbt_jawaban_siswas WHERE peserta_ujian_id = p.id) as count_jawaban,
			(SELECT COUNT(*) FROM log_ujians WHERE peserta_ujian_id = p.id AND keterangan_log LIKE '%Pindah Tab%') as count_pelanggaran
		FROM cbt_peserta_ujians p
		JOIN master_siswas s ON s.id = p.siswa_id
		LEFT JOIN master_kelas k ON k.id = s.kelas_id
		JOIN cbt_jadwal_ujians j ON j.id = p.jadwal_id
		WHERE p.jadwal_id = ?
	`, jadwalID).Scan(&rows).Error

	if err != nil {
		return SendError(c, fiber.StatusInternalServerError, "Gagal mengambil data monitoring: "+err.Error())
	}

	// Ambil info jadwal untuk mendapatkan BankSoalID
	var jadwal models.CBTJadwalUjian
	database.DB.First(&jadwal, jadwalID)

	// Ambil total soal secara akurat dari Bank Soal
	var totalSoal int64 = 0
	database.DB.Model(&models.CBTSoal{}).Where("bank_soal_id = ?", jadwal.BankSoalID).Count(&totalSoal)

	results := make([]MonitorSiswa, 0)
	for _, r := range rows {
		results = append(results, MonitorSiswa{
			ID:                r.ID,
			SiswaID:           r.SiswaID,
			NamaSiswa:         r.NamaSiswa,
			Kelas:             r.Kelas,
			StatusUjian:       r.StatusUjian,
			Progres:           r.CountJawaban,
			TotalSoal:         totalSoal,
			SisaWaktu:         r.SisaWaktuDetik,
			JumlahPelanggaran: r.CountPelanggaran,
			IsTerblokir:       r.IsTerblokir,
		})
	}

	return SendSuccess(c, "Data monitoring berhasil diambil", results)
}

// ForceSubmit menghentikan paksa ujian siswa
func ForceSubmit(c *fiber.Ctx) error {
	id := c.Params("pesertaId")
	pesertaID := utils.StringToUint(id)

	var peserta models.CBTPesertaUjian
	if err := database.DB.Preload("Jadwal").First(&peserta, pesertaID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Peserta tidak ditemukan"})
	}

	// Security Check
	role := GetUserRole(c)
	if role == "guru" {
		userID := GetUserID(c)
		guru, _ := GetGuruByUserID(userID)
		if peserta.Jadwal.PengawasID == nil || *peserta.Jadwal.PengawasID != guru.ID {
			return SendError(c, fiber.StatusForbidden, "Hanya pengawas jadwal ini yang boleh melakukan aksi ini")
		}
	}

	// Hitung Nilai PG Otomatis
	totalNilaiPG := HitungNilaiPG(peserta.ID, peserta.Jadwal.BankSoalID)

	// Update status menjadi Selesai dan masukkan Nilai
	now := time.Now()
	database.DB.Model(&peserta).Updates(map[string]interface{}{
		"status_ujian":        "Selesai",
		"waktu_selesai_ujian": &now,
		"nilai_pg":            totalNilaiPG,
		"total_nilai":         totalNilaiPG + peserta.NilaiEssay,
	})

	// Log aktivitas
	database.DB.Create(&models.LogUjian{
		PesertaUjianID: peserta.ID,
		KeteranganLog:  "Dihentikan Paksa oleh Pengawas",
	})

	return c.JSON(fiber.Map{"message": "Ujian berhasil dihentikan paksa"})
}

// UnblockSiswa membuka kembali akses ujian siswa
func UnblockSiswa(c *fiber.Ctx) error {
	id := c.Params("pesertaId")
	pesertaID := utils.StringToUint(id)

	var peserta models.CBTPesertaUjian
	if err := database.DB.Preload("Jadwal").First(&peserta, pesertaID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Peserta tidak ditemukan"})
	}

	// Security Check
	role := GetUserRole(c)
	if role == "guru" {
		userID := GetUserID(c)
		guru, _ := GetGuruByUserID(userID)
		if peserta.Jadwal.PengawasID == nil || *peserta.Jadwal.PengawasID != guru.ID {
			return SendError(c, fiber.StatusForbidden, "Hanya pengawas jadwal ini yang boleh melakukan aksi ini")
		}
	}

	// Gunakan Select() agar GORM tidak mengabaikan nilai NULL/nil pada waktu_login
	if err := database.DB.Model(&models.CBTPesertaUjian{}).Where("id = ?", pesertaID).
		Select("is_terblokir", "waktu_login").
		Updates(map[string]interface{}{
			"is_terblokir": false,
			"waktu_login":  nil,
		}).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal membuka blokir"})
	}

	// Log aktivitas
	database.DB.Create(&models.LogUjian{
		PesertaUjianID: uint(pesertaID),
		KeteranganLog:  "Blokir dibuka oleh Pengawas",
	})

	return c.JSON(fiber.Map{"message": "Blokir berhasil dibuka"})
}

// ResetSesiSiswa menghapus status login agar siswa bisa masuk kembali (misal jika PC crash)
func ResetSesiSiswa(c *fiber.Ctx) error {
	id := c.Params("pesertaId")
	pesertaID := utils.StringToUint(id)

	var peserta models.CBTPesertaUjian
	if err := database.DB.Preload("Jadwal").First(&peserta, pesertaID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Peserta tidak ditemukan"})
	}

	// Security Check
	role := GetUserRole(c)
	if role == "guru" {
		userID := GetUserID(c)
		guru, _ := GetGuruByUserID(userID)
		if peserta.Jadwal.PengawasID == nil || *peserta.Jadwal.PengawasID != guru.ID {
			return SendError(c, fiber.StatusForbidden, "Hanya pengawas jadwal ini yang boleh melakukan aksi ini")
		}
	}

	// Reset status agar bisa login ulang tanpa hambatan
	err := database.DB.Model(&models.CBTPesertaUjian{}).Where("id = ?", pesertaID).Updates(map[string]interface{}{
		"is_terblokir":        false,
		"waktu_login":         nil,
		"waktu_selesai_ujian": nil,
		"sisa_waktu_detik":    0,
		"status_ujian":        "Belum Mengerjakan",
	}).Error

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal reset sesi"})
	}

	// Log aktivitas
	database.DB.Create(&models.LogUjian{
		PesertaUjianID: uint(pesertaID),
		KeteranganLog:  "Sesi di-reset oleh Pengawas",
	})

	return c.JSON(fiber.Map{"message": "Sesi berhasil di-reset"})
}

type RekapNilai struct {
	ID           uint    `json:"id"`
	NamaSiswa    string  `json:"nama_siswa"`
	NISN         string  `json:"nisn"`
	Kelas        string  `json:"kelas"`
	NilaiPG      float64 `json:"nilai_pg"`
	NilaiEssay   float64 `json:"nilai_essay"`
	TotalNilai   float64 `json:"total_nilai"`
	StatusUjian  string  `json:"status_ujian"`
	WaktuSelesai string  `json:"waktu_selesai"`
}

// GetPengawasJadwals mengambil daftar jadwal di mana guru ybs ditugaskan sebagai pengawas
func GetPengawasJadwals(c *fiber.Ctx) error {
	role := GetUserRole(c)
	userID := GetUserID(c)

	if role != "guru" {
		return SendError(c, fiber.StatusForbidden, "Hanya untuk role guru")
	}

	guru, err := GetGuruByUserID(userID)
	if err != nil {
		return SendError(c, fiber.StatusForbidden, "Data Guru tidak ditemukan")
	}

	var jadwals []models.CBTJadwalUjian
	err = database.DB.Preload("BankSoal").
		Preload("BankSoal.Mapel").
		Preload("Ruang").
		Preload("Sesi").
		Where("pengawas_id = ?", guru.ID).
		Find(&jadwals).Error

	if err != nil {
		return SendError(c, fiber.StatusInternalServerError, "Gagal mengambil daftar pengawasan")
	}

	return SendSuccess(c, "Berhasil mengambil daftar pengawasan", jadwals)
}

// GetLaporanJadwals mengambil daftar jadwal untuk pilihan di halaman laporan
func GetLaporanJadwals(c *fiber.Ctx) error {
	role := GetUserRole(c)
	userID := GetUserID(c)

	var jadwals []models.CBTJadwalUjian
	query := database.DB.Preload("BankSoal").Preload("BankSoal.Mapel")

	if role == "guru" {
		guru, _ := GetGuruByUserID(userID)
		// Filter jadwal yang bank soalnya milik guru ini
		query = query.Joins("JOIN cbt_bank_soals ON cbt_bank_soals.id = cbt_jadwal_ujians.bank_soal_id").
			Where("cbt_bank_soals.guru_id = ?", guru.ID)
	}

	if err := query.Order("tanggal_ujian desc").Find(&jadwals).Error; err != nil {
		return SendError(c, fiber.StatusInternalServerError, "Gagal mengambil daftar jadwal")
	}

	return SendSuccess(c, "Berhasil mengambil daftar jadwal laporan", jadwals)
}

// GetRekapNilai mengambil seluruh hasil nilai siswa untuk satu jadwal
func GetRekapNilai(c *fiber.Ctx) error {
	jadwalID := c.Params("jadwalId")
	role := GetUserRole(c)
	userID := GetUserID(c)

	// Proteksi Akses Guru
	if role == "guru" {
		guru, _ := GetGuruByUserID(userID)
		var count int64
		database.DB.Model(&models.CBTJadwalUjian{}).
			Joins("JOIN cbt_bank_soals ON cbt_bank_soals.id = cbt_jadwal_ujians.bank_soal_id").
			Where("cbt_jadwal_ujians.id = ? AND cbt_bank_soals.guru_id = ?", jadwalID, guru.ID).
			Count(&count)

		if count == 0 {
			return SendError(c, fiber.StatusForbidden, "Anda tidak memiliki akses ke laporan ini")
		}
	}

	var peserta []models.CBTPesertaUjian
	err := database.DB.Preload("Siswa").Preload("Siswa.Kelas").
		Where("jadwal_id = ?", jadwalID).
		Find(&peserta).Error

	if err != nil {
		return SendError(c, fiber.StatusInternalServerError, "Gagal mengambil rekap nilai")
	}

	results := make([]RekapNilai, 0)
	for _, p := range peserta {
		waktuSelesai := "-"
		if p.WaktuSelesaiUjian != nil {
			waktuSelesai = p.WaktuSelesaiUjian.Format("02/01/2006 15:04")
		}

		results = append(results, RekapNilai{
			ID:           p.ID,
			NamaSiswa:    p.Siswa.NamaLengkap,
			NISN:         p.Siswa.NISN,
			Kelas:        p.Siswa.Kelas.NamaKelas,
			NilaiPG:      p.NilaiPG,
			NilaiEssay:   p.NilaiEssay,
			TotalNilai:   p.TotalNilai,
			StatusUjian:  p.StatusUjian,
			WaktuSelesai: waktuSelesai,
		})
	}

	return SendSuccess(c, "Berhasil mengambil rekap nilai", results)
}

// GetJawabanPeserta mengambil detail jawaban untuk dikoreksi
func GetJawabanPeserta(c *fiber.Ctx) error {
	pesertaID := c.Params("pesertaId")

	// 1. Ambil info peserta & bank_soal_id
	var peserta models.CBTPesertaUjian
	if err := database.DB.Preload("Jadwal").First(&peserta, pesertaID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Peserta tidak ditemukan"})
	}

	// 2. Ambil semua soal dan join dengan jawaban (jika ada)
	// Kita ambil semua soal dari bank soal yang bersangkutan
	var results []struct {
		SoalID       uint    `json:"soal_id"`
		Pertanyaan   string  `json:"pertanyaan"`
		JenisSoal    string  `json:"jenis_soal"`
		BobotMaks    float64 `json:"bobot_maks"`
		JawabanSiswa string  `json:"jawaban_siswa"`
		Skor         float64 `json:"skor"`
	}

	err := database.DB.Table("cbt_soals").
		Select("cbt_soals.id as soal_id, cbt_soals.pertanyaan, cbt_soals.jenis_soal, cbt_soals.bobot_nilai as bobot_maks, cbt_jawaban_siswas.jawaban_siswa, cbt_jawaban_siswas.skor").
		Joins("LEFT JOIN cbt_jawaban_siswas ON cbt_jawaban_siswas.soal_id = cbt_soals.id AND cbt_jawaban_siswas.peserta_ujian_id = ?", pesertaID).
		Where("cbt_soals.bank_soal_id = ?", peserta.Jadwal.BankSoalID).
		Scan(&results).Error

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal mengambil detail jawaban"})
	}

	return c.JSON(fiber.Map{"data": results})
}

// UpdateKoreksiEssay menyimpan hasil koreksi manual dari guru
func UpdateKoreksiEssay(c *fiber.Ctx) error {
	pesertaID := c.Params("pesertaId")

	type ItemKoreksi struct {
		SoalID uint    `json:"soal_id"`
		Skor   float64 `json:"skor"`
	}
	var input []ItemKoreksi
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Input tidak valid"})
	}

	tx := database.DB.Begin()
	pID, _ := strconv.ParseUint(pesertaID, 10, 32)

	var totalEssay float64 = 0
	for _, item := range input {
		var jawaban models.CBTJawabanSiswa
		// Cari atau buat record jika belum ada (jika siswa tidak menjawab)
		if err := tx.Where("peserta_ujian_id = ? AND soal_id = ?", pID, item.SoalID).
			FirstOrCreate(&jawaban, models.CBTJawabanSiswa{
				PesertaUjianID: uint(pID),
				SoalID:         item.SoalID,
			}).Error; err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal akses data jawaban"})
		}

		jawaban.Skor = item.Skor
		if err := tx.Save(&jawaban).Error; err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal update skor"})
		}
		totalEssay += item.Skor
	}

	// Update Nilai Total di PesertaUjian
	var peserta models.CBTPesertaUjian
	if err := tx.First(&peserta, pesertaID).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Peserta tidak ditemukan"})
	}

	peserta.NilaiEssay = totalEssay
	peserta.TotalNilai = peserta.NilaiPG + totalEssay
	tx.Save(&peserta)

	tx.Commit()

	return c.JSON(fiber.Map{
		"message":     "Koreksi berhasil disimpan",
		"nilai_essay": totalEssay,
		"total_nilai": peserta.TotalNilai,
	})
}

type AnalisisButir struct {
	SoalID     uint           `json:"soal_id"`
	Pertanyaan string         `json:"pertanyaan"`
	JenisSoal  string         `json:"jenis_soal"`
	Benar      int            `json:"benar"`
	Salah      int            `json:"salah"`
	Kosong     int            `json:"kosong"`
	Persentase float64        `json:"persentase"`
	Kesulitan  string         `json:"kesulitan"`
	Sebaran    map[string]int `json:"sebaran"`
}

func GetAnalisisSoal(c *fiber.Ctx) error {
	jadwalID := c.Params("jadwalId")

	// 1. Ambil Jadwal & Soal
	var jadwal models.CBTJadwalUjian
	if err := database.DB.Preload("BankSoal").First(&jadwal, jadwalID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Jadwal tidak ditemukan"})
	}

	var soals []models.CBTSoal
	database.DB.Where("bank_soal_id = ?", jadwal.BankSoalID).Find(&soals)

	// 2. Ambil Semua Jawaban untuk Jadwal ini
	var jawabans []models.CBTJawabanSiswa
	database.DB.Table("cbt_jawaban_siswas").
		Joins("JOIN cbt_peserta_ujians ON cbt_peserta_ujians.id = cbt_jawaban_siswas.peserta_ujian_id").
		Where("cbt_peserta_ujians.jadwal_id = ?", jadwalID).
		Find(&jawabans)

	// 3. Ambil Total Peserta yang SUDAH SELESAI/SEDANG UJIAN
	var totalPeserta int64
	database.DB.Model(&models.CBTPesertaUjian{}).Where("jadwal_id = ?", jadwalID).Count(&totalPeserta)

	if totalPeserta == 0 {
		return c.JSON([]AnalisisButir{})
	}

	// 4. Map Jawaban untuk mempercepat akses
	// Map: SoalID -> List Jawaban
	mapJawabans := make(map[uint][]models.CBTJawabanSiswa)
	for _, j := range jawabans {
		mapJawabans[j.SoalID] = append(mapJawabans[j.SoalID], j)
	}

	// 5. Hitung Agregasi
	analisis := make([]AnalisisButir, 0)
	for _, s := range soals {
		jList := mapJawabans[s.ID]

		benar := 0
		salah := 0
		kosong := int(totalPeserta) - len(jList)
		sebaran := make(map[string]int)

		for _, j := range jList {
			if j.Skor > 0 {
				benar++
			} else if j.JawabanSiswa != "" {
				salah++
			} else {
				kosong++
			}

			if s.JenisSoal == "PG" && j.JawabanSiswa != "" {
				sebaran[j.JawabanSiswa]++
			}
		}

		persentase := (float64(benar) / float64(totalPeserta)) * 100
		kesulitan := "Sedang"
		if persentase > 70 {
			kesulitan = "Mudah"
		} else if persentase < 30 {
			kesulitan = "Sukar"
		}

		analisis = append(analisis, AnalisisButir{
			SoalID:     s.ID,
			Pertanyaan: s.Pertanyaan,
			JenisSoal:  s.JenisSoal,
			Benar:      benar,
			Salah:      salah,
			Kosong:     kosong,
			Persentase: persentase,
			Kesulitan:  kesulitan,
			Sebaran:    sebaran,
		})
	}

	return c.JSON(analisis)
}
