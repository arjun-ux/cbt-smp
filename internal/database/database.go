package database

import (
	"log"
	"os"

	"cbt-smp/internal/models"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

// ConnectDB inisialisasi koneksi ke SQLite dan menjalankan AutoMigrate
func ConnectDB() {
	var err error

	// Ambil jalur DB dari .env (default cbt.db)
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "cbt.db"
	}

	// Gunakan WAL Mode agar Read dan Write tidak saling mengunci (Mencegah "Pending" di Frontend)
	dsn := dbPath + "?_pragma=foreign_keys(1)&_pragma=busy_timeout(5000)&_pragma=journal_mode(WAL)"
	DB, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Gagal koneksi ke database SQLite: \n", err)
	}

	log.Println("Koneksi ke database " + dbPath + " berhasil.")

	// Menjalankan AutoMigrate
	log.Println("Menjalankan AutoMigrate...")
	err = DB.AutoMigrate(
		&models.User{},
		&models.MasterKelas{},
		&models.MasterRuang{},
		&models.MasterSesi{},
		&models.MasterSiswa{},
		&models.MasterGuru{},
		&models.MasterMapel{},
		&models.CBTBankSoal{},
		&models.CBTSoal{},
		&models.CBTJadwalUjian{},
		&models.CBTPesertaUjian{},
		&models.CBTJawabanSiswa{},
		&models.LogSistem{},
		&models.LogUjian{},
		&models.CBTRekapNilai{},
	)

	if err != nil {
		log.Fatal("Gagal menjalankan AutoMigrate: \n", err)
	}
	log.Println("AutoMigrate selesai.")

	// Jalankan Seeder Minimal (Produksi)
	Seed()
}
