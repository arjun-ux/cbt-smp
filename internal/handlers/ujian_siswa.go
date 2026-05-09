package handlers

import (
	"cbt-smp/internal/database"
	"cbt-smp/internal/models"
	"cbt-smp/pkg/utils"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// GetSiswaJadwal mengambil daftar jadwal ujian untuk dashboard siswa
func GetSiswaJadwal(c *fiber.Ctx) error {
	var userID uint
	if val, ok := c.Locals("user_id").(float64); ok {
		userID = uint(val)
	} else if val, ok := c.Locals("user_id").(uint); ok {
		userID = val
	}

	var siswa models.MasterSiswa
	if err := database.DB.Preload("Kelas").Where("user_id = ?", userID).First(&siswa).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Data siswa tidak ditemukan"})
	}

	var jadwals []models.CBTJadwalUjian
	// Ambil jadwal yang tingkat kelasnya sama dengan siswa
	// Dan Sesi/Ruang harus cocok jika ditentukan di jadwal
	err := database.DB.Preload("BankSoal").Preload("BankSoal.Mapel").Preload("Ruang").Preload("Sesi").
		Joins("JOIN cbt_bank_soals ON cbt_bank_soals.id = cbt_jadwal_ujians.bank_soal_id").
		Where("cbt_bank_soals.tingkat_kelas = ? AND cbt_bank_soals.status = ?", siswa.Kelas.Tingkat, "Aktif").
		Where("(cbt_jadwal_ujians.sesi_id IS NULL OR cbt_jadwal_ujians.sesi_id = ?)", siswa.SesiID).
		Where("(cbt_jadwal_ujians.ruang_id IS NULL OR cbt_jadwal_ujians.ruang_id = ?)", siswa.RuangID).
		Order("tanggal_ujian desc, waktu_mulai asc").
		Find(&jadwals).Error

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal mengambil jadwal"})
	}

	// Cek status pengerjaan untuk setiap jadwal
	type JadwalSiswaView struct {
		models.CBTJadwalUjian
		StatusPengerjaan string `json:"status_pengerjaan"` // Belum, Selesai, Sedang Mengerjakan
	}

	var results []JadwalSiswaView
	for _, j := range jadwals {
		var peserta models.CBTPesertaUjian
		database.DB.Where("jadwal_id = ? AND siswa_id = ?", j.ID, siswa.ID).Limit(1).Find(&peserta)

		statusP := "Belum"
		if peserta.ID != 0 {
			statusP = peserta.StatusUjian
		}

		results = append(results, JadwalSiswaView{
			CBTJadwalUjian:   j,
			StatusPengerjaan: statusP,
		})
	}

	return c.JSON(fiber.Map{
		"data":  results,
		"siswa": siswa,
	})
}

type StudentLoginRequest struct {
	NISN  string `json:"nisn"`
	Token string `json:"token"`
}

// ValidateExam memvalidasi NISN dan Token Ujian
func ValidateExam(c *fiber.Ctx) error {
	var input StudentLoginRequest
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Input tidak valid"})
	}

	// 1. Cek Jadwal berdasarkan Token (Case Insensitive & Hapus Spasi)
	cleanToken := strings.ReplaceAll(input.Token, " ", "")
	var jadwal models.CBTJadwalUjian
	query := database.DB.Preload("BankSoal").Preload("BankSoal.Mapel").Where("UPPER(token_ujian) = UPPER(?)", cleanToken)

	// Jika ada JadwalID (dari dashboard baru), gunakan itu agar lebih akurat
	if jId := c.Query("jadwalId"); jId != "" {
		query = query.Where("id = ?", jId)
	}

	if err := query.First(&jadwal).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Token ujian tidak valid untuk jadwal ini"})
	}

	// 2. Cek Status Jadwal
	if jadwal.Status == "Belum Mulai" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Ujian belum dimulai"})
	}
	if jadwal.Status == "Selesai" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Ujian telah berakhir"})
	}

	// 3. Cek Siswa (Utamakan Sesi Login JWT)
	var userID uint
	if val, ok := c.Locals("user_id").(float64); ok {
		userID = uint(val)
	} else if val, ok := c.Locals("user_id").(uint); ok {
		userID = val
	}

	var siswa models.MasterSiswa
	if userID != 0 {
		database.DB.Preload("User").Preload("Kelas").Preload("Ruang").Preload("Sesi").Where("user_id = ?", userID).First(&siswa)
	} else {
		// Fallback cara lama (jika NISN diisi di halaman depan)
		if err := database.DB.Preload("User").Preload("Kelas").Preload("Ruang").Preload("Sesi").Where("nisn = ?", input.NISN).First(&siswa).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "NISN tidak terdaftar"})
		}
	}

	if siswa.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Data siswa tidak ditemukan"})
	}

	// 4. Validasi Tingkat Kelas
	if siswa.Kelas.Tingkat != jadwal.BankSoal.TingkatKelas {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Tingkat kelas tidak sesuai dengan bank soal"})
	}

	// 5. Validasi Sesi (Hanya jika jadwal memiliki Sesi yang ditentukan)
	if jadwal.SesiID != nil && (siswa.SesiID == nil || *siswa.SesiID != *jadwal.SesiID) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Sesi ujian Anda tidak sesuai dengan jadwal ini"})
	}

	// 6. Validasi Ruang (Hanya jika jadwal memiliki Ruang yang ditentukan)
	if jadwal.RuangID != nil && (siswa.RuangID == nil || *siswa.RuangID != *jadwal.RuangID) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Lokasi/Ruang ujian Anda tidak sesuai dengan jadwal ini"})
	}

	// 7. Kelola Record Peserta Ujian
	var peserta models.CBTPesertaUjian
	err := database.DB.Where("jadwal_id = ? AND siswa_id = ?", jadwal.ID, siswa.ID).First(&peserta).Error

	now := time.Now()
	if err != nil {
		// Buat baru jika belum ada (Ujian Pertama Kali)
		peserta = models.CBTPesertaUjian{
			JadwalID:       jadwal.ID,
			SiswaID:        siswa.ID,
			StatusUjian:    "Sedang Mengerjakan",
			WaktuLogin:     &now,
			SisaWaktuDetik: jadwal.DurasiMenit * 60,
		}
		database.DB.Create(&peserta)
	} else {
		// 7a. Jika terblokir, izinkan lewat agar bisa melihat layar blokir
		if !peserta.IsTerblokir {
		// CEK SESI AKTIF (Cara Lama)
		if peserta.WaktuLogin != nil {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Anda sudah login di perangkat lain. Silakan hubungi pengawas untuk Reset Sesi.",
			})
		}
		}

		// Update status agar sedang mengerjakan
		database.DB.Model(&peserta).Updates(map[string]interface{}{
			"waktu_login":  &now,
			"status_ujian": "Sedang Mengerjakan",
		})

		// Cek jika sudah selesai secara permanen
		if peserta.StatusUjian == "Selesai" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Anda sudah menyelesaikan ujian ini"})
		}
	}

	// 8. Generate Token
	token, err := utils.GenerateJWT(siswa.User.ID, siswa.NamaLengkap, "siswa")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal generate token"})
	}

	return c.JSON(fiber.Map{
		"token": token,
		"data": fiber.Map{
			"peserta_id": peserta.ID,
			"siswa":      siswa,
			"jadwal":     jadwal,
			"sisa_waktu": peserta.SisaWaktuDetik,
		},
	})
}

// GetSoalUjian mengambil daftar soal (tanpa kunci)
func GetSoalUjian(c *fiber.Ctx) error {
	jadwalID := c.Params("jadwalId")

	var jadwal models.CBTJadwalUjian
	if err := database.DB.First(&jadwal, jadwalID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Jadwal tidak ditemukan"})
	}

	var soals []models.CBTSoal
	query := database.DB.Where("bank_soal_id = ?", jadwal.BankSoalID)
	if jadwal.AcakSoal {
		query = query.Order("RANDOM()")
	}
	query.Find(&soals)

	type SoalView struct {
		ID         uint    `json:"id"`
		JenisSoal  string  `json:"jenis_soal"`
		Pertanyaan string  `json:"pertanyaan"`
		OpsiA      string  `json:"opsi_a"`
		OpsiB      string  `json:"opsi_b"`
		OpsiC      string  `json:"opsi_c"`
		OpsiD      string  `json:"opsi_d"`
		BobotNilai float64 `json:"bobot_nilai"`
	}

	var results []SoalView
	for _, s := range soals {
		results = append(results, SoalView{
			ID:         s.ID,
			JenisSoal:  s.JenisSoal,
			Pertanyaan: s.Pertanyaan,
			OpsiA:      s.OpsiA,
			OpsiB:      s.OpsiB,
			OpsiC:      s.OpsiC,
			OpsiD:      s.OpsiD,
			BobotNilai: s.BobotNilai,
		})
	}

	// Ambil Peserta ID dari JWT (Claims) - Handle float64 conversion safety
	var userID uint
	if val, ok := c.Locals("user_id").(float64); ok {
		userID = uint(val)
	} else if val, ok := c.Locals("user_id").(uint); ok {
		userID = val
	}
	var peserta models.CBTPesertaUjian
	if err := database.DB.Where("jadwal_id = ? AND siswa_id = (SELECT id FROM master_siswas WHERE user_id = ?)", utils.StringToUint(jadwalID), userID).First(&peserta).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Data peserta tidak ditemukan. Pastikan Anda sudah login melalui halaman depan."})
	}

	// Sisa waktu dari database
	sisaWaktu := peserta.SisaWaktuDetik

	// Ambil jawaban yang sudah ada
	var jawabans []models.CBTJawabanSiswa
	database.DB.Where("peserta_ujian_id = ?", peserta.ID).Find(&jawabans)

	existingAnswers := make(map[uint]interface{})
	for _, j := range jawabans {
		existingAnswers[j.SoalID] = fiber.Map{
			"jawaban":   j.JawabanSiswa,
			"ragu_ragu": j.RaguRagu,
		}
	}

	return c.JSON(fiber.Map{
		"data":             results,
		"sisa_waktu":       sisaWaktu,
		"peserta_id":       peserta.ID,
		"is_terblokir":     peserta.IsTerblokir,
		"acak_jawaban":     jadwal.AcakJawaban,
		"existing_answers": existingAnswers,
	})
}

// GetSiswaStatus hanya mengambil status terblokir dan sisa waktu (Payload sangat ringan)
func GetSiswaStatus(c *fiber.Ctx) error {
	jadwalID := c.Params("jadwalId")
	var userID uint
	if val, ok := c.Locals("user_id").(float64); ok {
		userID = uint(val)
	} else if val, ok := c.Locals("user_id").(uint); ok {
		userID = val
	}

	var peserta models.CBTPesertaUjian
	if err := database.DB.Where("jadwal_id = ? AND siswa_id = (SELECT id FROM master_siswas WHERE user_id = ?)", utils.StringToUint(jadwalID), userID).First(&peserta).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Data peserta tidak ditemukan"})
	}

	sisaWaktu := peserta.SisaWaktuDetik

	return c.JSON(fiber.Map{
		"is_terblokir": peserta.IsTerblokir,
		"sisa_waktu":   sisaWaktu,
		"status_ujian": peserta.StatusUjian,
	})
}

// SyncJawaban menyimpan jawaban siswa (mendukung batch/grosir)
func SyncJawaban(c *fiber.Ctx) error {
	type JawabanItem struct {
		SoalID      uint   `json:"soal_id"`
		JawabanTeks string `json:"jawaban_teks"`
		RaguRagu    bool   `json:"ragu_ragu"`
	}

	type JawabanBatchInput struct {
		PesertaUjianID uint          `json:"peserta_ujian_id"`
		Items          []JawabanItem `json:"items"`
		SisaWaktu      int           `json:"sisa_waktu"`
	}

	var input JawabanBatchInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Input tidak valid"})
	}

	// Cek Status Ujian
	var peserta models.CBTPesertaUjian
	if err := database.DB.First(&peserta, input.PesertaUjianID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Data peserta tidak ditemukan"})
	}

	if peserta.StatusUjian == "Selesai" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Ujian telah berakhir atau dihentikan paksa",
			"code":  "EXAM_FINISHED",
		})
	}

	if peserta.IsTerblokir {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Akses ujian Anda diblokir karena pelanggaran. Hubungi pengawas.",
			"code":  "ACCOUNT_BLOCKED",
		})
	}

	// Gunakan Transaksi untuk Batch Save
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		// Update Sisa Waktu
		if input.SisaWaktu > 0 {
			if err := tx.Model(&peserta).Update("sisa_waktu_detik", input.SisaWaktu).Error; err != nil {
				return err
			}
		}

		for _, item := range input.Items {
			var jawaban models.CBTJawabanSiswa
			res := tx.Where("peserta_ujian_id = ? AND soal_id = ?", input.PesertaUjianID, item.SoalID).Limit(1).Find(&jawaban)

			if res.RowsAffected == 0 {
				// Simpan Baru
				newJawaban := models.CBTJawabanSiswa{
					PesertaUjianID: input.PesertaUjianID,
					SoalID:         item.SoalID,
					JawabanSiswa:   item.JawabanTeks,
					RaguRagu:       item.RaguRagu,
				}
				if err := tx.Create(&newJawaban).Error; err != nil {
					return err
				}
			} else {
				// Update yang sudah ada
				if err := tx.Model(&jawaban).Updates(map[string]interface{}{
					"jawaban_siswa": item.JawabanTeks,
					"ragu_ragu":     item.RaguRagu,
				}).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal sinkronisasi jawaban"})
	}

	return c.JSON(fiber.Map{
		"message":    "Sinkronisasi berhasil",
		"sisa_waktu": input.SisaWaktu,
		"synced":     len(input.Items),
	})
}

// LogSiswa mencatat aktivitas siswa (pindah tab dll)
func LogSiswa(c *fiber.Ctx) error {
	type LogInput struct {
		PesertaUjianID uint   `json:"peserta_ujian_id"`
		Keterangan     string `json:"keterangan"`
	}

	var input LogInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Input tidak valid"})
	}

	newLog := models.LogUjian{
		PesertaUjianID: input.PesertaUjianID,
		KeteranganLog:  input.Keterangan,
	}
	database.DB.Create(&newLog)

	// Auto Block jika pindah tab (hanya jika ujian sedang berlangsung)
	if strings.Contains(strings.ToLower(input.Keterangan), "pindah tab") {
		var p models.CBTPesertaUjian
		database.DB.Select("status_ujian").First(&p, input.PesertaUjianID)

		// Hanya blokir jika statusnya masih 'Sedang Mengerjakan'
		if p.StatusUjian == "Sedang Mengerjakan" {
			database.DB.Model(&models.CBTPesertaUjian{}).Where("id = ?", input.PesertaUjianID).Update("is_terblokir", true)
		}
	}

	return c.JSON(fiber.Map{"status": "logged"})
}

// HitungNilaiPG adalah helper untuk menghitung skor PG siswa
func HitungNilaiPG(pesertaID uint, bankSoalID uint) float64 {
	var soals []models.CBTSoal
	database.DB.Where("bank_soal_id = ? AND jenis_soal = 'PG'", bankSoalID).Find(&soals)

	var jawabans []models.CBTJawabanSiswa
	database.DB.Where("peserta_ujian_id = ?", pesertaID).Find(&jawabans)

	mapJawaban := make(map[uint]string)
	for _, j := range jawabans {
		mapJawaban[j.SoalID] = j.JawabanSiswa
	}

	var totalNilai float64 = 0
	database.DB.Transaction(func(tx *gorm.DB) error {
		for _, soal := range soals {
			jawabanSiswa := mapJawaban[soal.ID]
			if jawabanSiswa == soal.KunciJawaban {
				totalNilai += soal.BobotNilai
				tx.Model(&models.CBTJawabanSiswa{}).
					Where("peserta_ujian_id = ? AND soal_id = ?", pesertaID, soal.ID).
					Update("skor", soal.BobotNilai)
			} else {
				tx.Model(&models.CBTJawabanSiswa{}).
					Where("peserta_ujian_id = ? AND soal_id = ?", pesertaID, soal.ID).
					Update("skor", 0)
			}
		}
		return nil
	})
	
	return totalNilai
}

// SubmitUjian mengakhiri sesi ujian dan menghitung nilai PG otomatis
func SubmitUjian(c *fiber.Ctx) error {
	idStr := c.Params("pesertaId")
	pesertaID := utils.StringToUint(idStr)
	now := time.Now()

	var peserta models.CBTPesertaUjian
	if err := database.DB.Preload("Jadwal").Where("id = ?", pesertaID).First(&peserta).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Data peserta tidak ditemukan"})
	}

	// Hitung Nilai via Helper
	totalNilaiPG := HitungNilaiPG(peserta.ID, peserta.Jadwal.BankSoalID)

	// 5. Update Status dan Nilai
	err := database.DB.Model(&peserta).Updates(map[string]interface{}{
		"status_ujian":        "Selesai",
		"waktu_selesai_ujian": &now,
		"nilai_pg":            totalNilaiPG,
		"total_nilai":         totalNilaiPG + peserta.NilaiEssay, // Essay masih 0
	}).Error

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal menyimpan hasil ujian"})
	}

	// 6. Catat Log Selesai
	logSelesai := models.LogUjian{
		PesertaUjianID: peserta.ID,
		KeteranganLog:  "Selesai Ujian - Nilai PG Dihitung Otomatis",
	}
	database.DB.Create(&logSelesai)

	return c.JSON(fiber.Map{
		"message":  "Ujian berhasil diselesaikan",
		"nilai_pg": totalNilaiPG,
	})
}
