package database

import (
	"cbt-smp/internal/models"
	"cbt-smp/pkg/utils"
	"log"
)

// Seed mengisi database dengan data awal
func Seed() {
	// 1. Seeder Admin
	var countAdmin int64
	DB.Model(&models.User{}).Where("role = ?", "admin").Count(&countAdmin)
	if countAdmin == 0 {
		hash, _ := utils.HashPassword("admin123")
		adminUser := models.User{
			Username: "admin",
			Password: hash,
			Role:     "admin",
			IsActive: true,
		}
		DB.Create(&adminUser)
		log.Println("Seeder: Berhasil membuat akun Admin default (username: admin, password: admin123)")
	}

	// 2. Seeder Master Kelas
	var countKelas int64
	DB.Model(&models.MasterKelas{}).Count(&countKelas)
	if countKelas == 0 {
		kelases := []models.MasterKelas{
			{Tingkat: "7", NamaKelas: "7A"},
			{Tingkat: "7", NamaKelas: "7B"},
			{Tingkat: "8", NamaKelas: "8A"},
			{Tingkat: "8", NamaKelas: "8B"},
			{Tingkat: "9", NamaKelas: "9A"},
			{Tingkat: "9", NamaKelas: "9B"},
		}
		DB.Create(&kelases)
		log.Println("Seeder: Berhasil membuat data Kelas default")
	}

	// 3. Seeder Master Ruang
	var countRuang int64
	DB.Model(&models.MasterRuang{}).Count(&countRuang)
	if countRuang == 0 {
		ruangs := []models.MasterRuang{
			{KodeRuang: "R01", NamaRuang: "LAB KOMPUTER 01"},
			{KodeRuang: "R02", NamaRuang: "LAB KOMPUTER 02"},
			{KodeRuang: "R03", NamaRuang: "RUANG KELAS 01"},
		}
		DB.Create(&ruangs)
		log.Println("Seeder: Berhasil membuat data Ruangan default")
	}

	// 4. Seeder Master Sesi
	var countSesi int64
	DB.Model(&models.MasterSesi{}).Count(&countSesi)
	if countSesi == 0 {
		sesis := []models.MasterSesi{
			{NamaSesi: "SESI 1 (07:30 - 09:30)", WaktuMulai: "07:30", WaktuSelesai: "09:30"},
			{NamaSesi: "SESI 2 (10:00 - 12:00)", WaktuMulai: "10:00", WaktuSelesai: "12:00"},
			{NamaSesi: "SESI 3 (13:00 - 15:00)", WaktuMulai: "13:00", WaktuSelesai: "15:00"},
		}
		DB.Create(&sesis)
		log.Println("Seeder: Berhasil membuat data Sesi default")
	}

	// 5. Seeder Master Mapel
	var countMapel int64
	DB.Model(&models.MasterMapel{}).Count(&countMapel)
	if countMapel == 0 {
		mapels := []models.MasterMapel{
			{KodeMapel: "IND", NamaMapel: "BAHASA INDONESIA"},
			{KodeMapel: "MAT", NamaMapel: "MATEMATIKA"},
			{KodeMapel: "ING", NamaMapel: "BAHASA INGGRIS"},
			{KodeMapel: "IPA", NamaMapel: "IPA"},
			{KodeMapel: "IPS", NamaMapel: "IPS"},
			{KodeMapel: "PAI", NamaMapel: "PENDIDIKAN AGAMA ISLAM & BUDI PEKERTI"},
			{KodeMapel: "PKN", NamaMapel: "PENDIDIKAN KEWARGAAN NEGARA"},
			{KodeMapel: "PJOK", NamaMapel: "PENDIDIKAN JASMANI, OLAHRAGA, DAN KESEHATAN"},
			{KodeMapel: "SBK", NamaMapel: "SENI BUDAYA DAN PRAKARYA"},
			{KodeMapel: "TIK", NamaMapel: "TEKNOLOGI INFORMASI DAN KOMUNIKASI"},
		}
		DB.Create(&mapels)
		log.Println("Seeder: Berhasil membuat data Mapel default")
	}
}
