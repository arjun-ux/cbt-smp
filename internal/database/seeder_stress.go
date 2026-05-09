package database

import (
	"cbt-smp/internal/models"
	"cbt-smp/pkg/utils"
	"fmt"
	"log"
	"time"
)

// SeedStressTestData membuat 100 siswa dan jadwal ujian untuk simulasi beban
func SeedStressTestData() {
	log.Println("Memulai Seeder Stress Test...")

	// 1. Buat Kelas
	kelas := models.MasterKelas{Tingkat: "9", NamaKelas: "STRESS-TEST"}
	DB.FirstOrCreate(&kelas, models.MasterKelas{NamaKelas: "STRESS-TEST"})

	// 2. Buat Ruang & Sesi
	ruang := models.MasterRuang{KodeRuang: "LAB-STRESS", NamaRuang: "LAB KOMPUTER STRESS TEST"}
	DB.FirstOrCreate(&ruang, models.MasterRuang{KodeRuang: "LAB-STRESS"})

	sesi := models.MasterSesi{NamaSesi: "SESI-STRESS", WaktuMulai: "07:30", WaktuSelesai: "09:30"}
	DB.FirstOrCreate(&sesi, models.MasterSesi{NamaSesi: "SESI-STRESS"})

	// 3. Buat Mapel & Bank Soal
	mapel := models.MasterMapel{KodeMapel: "ST-01", NamaMapel: "Simulasi Beban"}
	DB.FirstOrCreate(&mapel, models.MasterMapel{KodeMapel: "ST-01"})

	bankSoal := models.CBTBankSoal{
		MapelID:       mapel.ID,
		TingkatKelas:  "9",
		JudulBankSoal: "Ujian Simulasi 100 Siswa",
		Status:        "Aktif",
	}
	DB.FirstOrCreate(&bankSoal, models.CBTBankSoal{JudulBankSoal: "Ujian Simulasi 100 Siswa"})

	// 3. Buat 5 Soal PG
	for i := 1; i <= 5; i++ {
		soal := models.CBTSoal{
			BankSoalID:   bankSoal.ID,
			JenisSoal:    "PG",
			Pertanyaan:   fmt.Sprintf("Pertanyaan Simulasi Nomor %d", i),
			OpsiA:        "Jawaban A",
			OpsiB:        "Jawaban B",
			OpsiC:        "Jawaban C",
			OpsiD:        "Jawaban D",
			KunciJawaban: "A",
		}
		DB.FirstOrCreate(&soal, models.CBTSoal{Pertanyaan: soal.Pertanyaan})
	}

	// 5. Buat Jadwal Ujian
	now := time.Now()
	jadwal := models.CBTJadwalUjian{
		BankSoalID:   bankSoal.ID,
		TanggalUjian: now.Format("2006-01-02"),
		WaktuMulai:   now.Format("15:04"),
		DurasiMenit:  120,
		RuangID:      &ruang.ID,
		SesiID:       &sesi.ID,
		TokenUjian:   "TEST100",
		Status:       "Berlangsung",
	}
	DB.FirstOrCreate(&jadwal, models.CBTJadwalUjian{TokenUjian: "TEST100"})

	// 5. Buat 100 Siswa
	passwordHash, _ := utils.HashPassword("siswa123")
	for i := 1; i <= 100; i++ {
		username := fmt.Sprintf("siswa%d", i)

		// Create User Account
		user := models.User{
			Username: username,
			Password: passwordHash,
			Role:     "siswa",
			IsActive: true,
		}

		if err := DB.Where("username = ?", username).First(&user).Error; err != nil {
			DB.Create(&user)

			// Create Student Data
			siswa := models.MasterSiswa{
				UserID:      user.ID,
				NISN:        fmt.Sprintf("1000%03d", i),
				NamaLengkap: fmt.Sprintf("Siswa Simulasi %d", i),
				KelasID:     &kelas.ID,
				RuangID:     &ruang.ID,
				SesiID:      &sesi.ID,
			}
			DB.Create(&siswa)
		}
	}

	log.Println("Seeder Stress Test Selesai! 100 Siswa & Jadwal 'TEST100' siap digunakan.")
}
