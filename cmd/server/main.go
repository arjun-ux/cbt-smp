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
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	"os"
	"flag"
	"os/exec"
	"runtime"
	"strings"
	"fmt"
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

	// --- SECURITY: JWT Secret Hardening ---
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("KRITIKAL: JWT_SECRET belum diatur di file .env! Aplikasi tidak dapat dijalankan demi keamanan.")
	}
	if len(jwtSecret) < 32 {
		log.Fatal("KRITIKAL: JWT_SECRET terlalu pendek (minimal 32 karakter)! Harap gunakan kunci yang lebih kuat untuk melindungi data.")
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

	// --- SECURITY: Rate Limiting ---
	// 1. Strict Limiter untuk Login & Validasi Token (Anti-Brute Force)
	authLimiter := limiter.New(limiter.Config{
		Max:               5, // Maksimal 5 percobaan
		Expiration:        1 * time.Minute,
		KeyGenerator:      func(c *fiber.Ctx) string { return c.IP() },
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "Terlalu banyak percobaan login/validasi. Silakan tunggu 1 menit.",
			})
		},
	})
	// Terapkan pada endpoint sensitif
	app.Use("/api/auth/login", authLimiter)
	app.Use("/api/siswa/validate", authLimiter)

	// API Routes
	api := app.Group("/api")

	// 2. Global API Limiter (Pencegahan DOS/Spam umum)
	api.Use(limiter.New(limiter.Config{
		Max:        300, // Menaikkan batas ke 300 untuk rute umum
		Expiration: 1 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string { return c.IP() },
		Next: func(c *fiber.Ctx) bool {
			path := c.Path()
			// JANGAN batasi rute polling agar tidak mengganggu ujian massal di Lab Sekolah
			return strings.Contains(path, "/sync") || 
			       strings.Contains(path, "/log") || 
			       strings.Contains(path, "/monitor")
		},
	}))
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

	// --- TERMINAL BANNER ---
	fmt.Println("\n  ________________________________________________")
	fmt.Println(" |                                                |")
	fmt.Println(" |          APLIKASI CBT SMP SUDAH AKTIF          |")
	fmt.Println(" |________________________________________________|")
	fmt.Println("")
	fmt.Printf("  Akses Server: http://localhost:%s\n", port)
	fmt.Println("  ________________________________________________")
	fmt.Println("\n  [TEKAN CTRL+C UNTUK BERHENTI]")
	fmt.Println("")

	log.Fatal(app.Listen(":" + port))
}
