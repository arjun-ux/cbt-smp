package main

import (
	"log"
	"net/http"
	"time"

	"cbt-smp/frontend"
	"cbt-smp/internal/database"
	"cbt-smp/internal/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	"os"
	"flag"
	"os/exec"
	"runtime"
)

func openBrowser(url string) {
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	}
	if err != nil {
		log.Println("Gagal membuka browser:", err)
	}
}

func main() {
	// Definisikan Flag
	seedStress := flag.Bool("seed-stress", false, "Jalankan seeder untuk stress test 100 siswa")
	flag.Parse()

	// Muat file .env
	if err := godotenv.Load(); err != nil {
		log.Println("Info: File .env tidak ditemukan, menggunakan setelan default.")
	}

	// Inisialisasi Database
	database.ConnectDB()

	// Cek jika disuruh jalankan seeder stress test
	if *seedStress {
		database.SeedStressTestData()
		log.Println("Proses Seeder selesai. Silakan jalankan server tanpa flag untuk memulai.")
		return
	}

	app := fiber.New(fiber.Config{
		AppName: "CBT SMP App v1.0",
	})

	// Middleware
	app.Use(logger.New())

	// API Routes
	api := app.Group("/api")
	// Route Download Template (jika masih diperlukan atau folder statis lainnya)
	// app.Static("/template-soal.csv", "./template_soal.csv")

	// Setup registrasi route (Login, Admin, Guru, Siswa)
	routes.SetupRoutes(api)

	// Serve Static Files (Images, etc) with 24h Caching
	app.Static("/uploads", "./uploads", fiber.Static{
		MaxAge: 3600 * 24, // 24 jam (dalam detik)
	})

	// Serve Frontend (Vue SPA) via Embed
	app.Use("/", filesystem.New(filesystem.Config{
		Root:         http.FS(frontend.DistFS),
		PathPrefix:   "dist",
		Browse:       false,
		Index:        "index.html",
		NotFoundFile: "dist/index.html",
	}))


	// Ambil Port dari .env (default 3000)
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	// Buka browser otomatis
	go func() {
		time.Sleep(1 * time.Second)
		openBrowser("http://localhost:" + port)
	}()

	log.Fatal(app.Listen(":" + port))
}
