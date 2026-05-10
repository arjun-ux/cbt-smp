package routes

import (
	"cbt-smp/internal/handlers"
	"cbt-smp/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(router fiber.Router) {
	// router di sini adalah group "/api" dari main.go
	
	// 1. Auth (Public) - Tanpa prefix /api lagi
	router.Post("/auth/login", handlers.Login)

	// 2. Protected Group (Admin, Guru, Siswa)
	// Kita buat group baru agar middleware.Protected() tidak mengenai rute login di atas
	protected := router.Group("/")
	protected.Use(middleware.Protected())
	protected.Post("/auth/logout", handlers.Logout)

	// Admin Group -> /api/admin
	admin := protected.Group("/admin")
	admin.Use(middleware.RoleRequired("admin"))
	
	admin.Get("/dashboard", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Admin Dashboard"})
	})

	// Manajemen Admin
	admin.Get("/manage", handlers.GetAdmins)
	admin.Post("/manage", handlers.CreateAdmin)
	admin.Put("/manage/:id", handlers.UpdateAdmin)
	admin.Delete("/manage/:id", handlers.DeleteAdmin)

	// Master Data
	admin.Get("/guru", handlers.GetGurus)
	admin.Post("/guru", handlers.CreateGuru)
	admin.Post("/guru/import", handlers.ImportGuru)
	admin.Put("/guru/:id", handlers.UpdateGuru)
	admin.Delete("/guru/:id", handlers.DeleteGuru)
	admin.Patch("/guru/:id/status", handlers.ToggleGuruStatus)

	admin.Get("/siswa", handlers.GetSiswas)
	admin.Post("/siswa", handlers.CreateSiswa)
	admin.Post("/siswa/import", handlers.ImportSiswa)
	admin.Put("/siswa/:id", handlers.UpdateSiswa)
	admin.Delete("/siswa/:id", handlers.DeleteSiswa)
	admin.Patch("/siswa/:id/status", handlers.ToggleSiswaStatus)
	admin.Patch("/siswa/bulk-status", handlers.BulkUpdateSiswaStatus)
	admin.Post("/siswa/bulk-plot", handlers.BulkPlotSiswa)

	admin.Get("/kelas", handlers.GetKelas)
	admin.Post("/kelas", handlers.CreateKelas)
	admin.Put("/kelas/:id", handlers.UpdateKelas)
	admin.Delete("/kelas/:id", handlers.DeleteKelas)

	admin.Get("/mapel", handlers.GetMapel)
	admin.Post("/mapel", handlers.CreateMapel)
	admin.Put("/mapel/:id", handlers.UpdateMapel)
	admin.Delete("/mapel/:id", handlers.DeleteMapel)

	admin.Get("/ruang", handlers.GetRuang)
	admin.Post("/ruang", handlers.CreateRuang)
	admin.Put("/ruang/:id", handlers.UpdateRuang)
	admin.Delete("/ruang/:id", handlers.DeleteRuang)

	admin.Get("/sesi", handlers.GetSesi)
	admin.Post("/sesi", handlers.CreateSesi)
	admin.Put("/sesi/:id", handlers.UpdateSesi)
	admin.Delete("/sesi/:id", handlers.DeleteSesi)

	// Bank Soal
	admin.Get("/bank-soal", handlers.GetBankSoals)
	admin.Get("/bank-soal/:id", handlers.GetBankSoal)
	admin.Post("/bank-soal", handlers.CreateBankSoal)
	admin.Put("/bank-soal/:id", handlers.UpdateBankSoal)
	admin.Delete("/bank-soal/:id", handlers.DeleteBankSoal)

	// Jadwal Ujian
	admin.Get("/jadwal", handlers.GetJadwals)
	admin.Get("/jadwal/:id", handlers.GetJadwal)
	admin.Post("/jadwal", handlers.CreateJadwal)
	admin.Put("/jadwal/:id", handlers.UpdateJadwal)
	admin.Delete("/jadwal/:id", handlers.DeleteJadwal)
	admin.Patch("/jadwal/:id/token", handlers.RefreshToken)
	admin.Post("/jadwal/:id/pengawas", handlers.UpdatePengawas)
	admin.Post("/jadwal/:id/archive", handlers.ArchiveJadwalResults)

	// Monitoring
	admin.Get("/monitor/:jadwalId", handlers.GetMonitorUjian)
	admin.Post("/monitor/force-submit/:pesertaId", handlers.ForceSubmit)
	admin.Post("/monitor/unblock/:pesertaId", handlers.UnblockSiswa)
	admin.Post("/monitor/block/:pesertaId", handlers.BlockSiswa)
	admin.Post("/monitor/reset-sesi/:pesertaId", handlers.ResetSesiSiswa)
	admin.Post("/monitor/reset-ujian/:pesertaId", handlers.ResetUjianSiswa)
	admin.Get("/rekap-jadwal", handlers.GetLaporanJadwals)
	admin.Get("/rekap-nilai/:jadwalId", handlers.GetRekapNilai)
	admin.Get("/rekap-arsip", handlers.GetRiwayatNilai)
	admin.Get("/monitor/jawaban/:pesertaId", handlers.GetJawabanPeserta)
	admin.Post("/monitor/koreksi/:pesertaId", handlers.UpdateKoreksiEssay)
	admin.Get("/analisis/:jadwalId", handlers.GetAnalisisSoal)
	admin.Post("/migrate/images", handlers.MigrateBase64Images)

	admin.Get("/bank-soal/:bankSoalId/soal", handlers.GetSoals)
	admin.Post("/bank-soal/:bankSoalId/soal", handlers.CreateSoal)
	admin.Post("/bank-soal/:bankSoalId/import", handlers.ImportSoal)
	admin.Delete("/bank-soal/:bankSoalId/soal/clear", handlers.ClearSoals)
	admin.Put("/soal/:id", handlers.UpdateSoal)
	admin.Delete("/soal/:id", handlers.DeleteSoal)

	// Guru Group -> /api/guru
	guru := protected.Group("/guru")
	guru.Use(middleware.RoleRequired("admin", "guru"))
	
	// Guru juga butuh akses ke Bank Soal & Soal miliknya
	guru.Get("/bank-soal", handlers.GetBankSoals)
	guru.Get("/bank-soal/:id", handlers.GetBankSoal)
	guru.Post("/bank-soal", handlers.CreateBankSoal)
	guru.Put("/bank-soal/:id", handlers.UpdateBankSoal)
	guru.Delete("/bank-soal/:id", handlers.DeleteBankSoal)
	guru.Get("/mapel", handlers.GetMapel)
	
	guru.Get("/jadwal", handlers.GetJadwals)
	guru.Get("/jadwal/:id", handlers.GetJadwal)
	guru.Get("/monitor/:jadwalId", handlers.GetMonitorUjian)
	guru.Post("/monitor/force-submit/:pesertaId", handlers.ForceSubmit)
	guru.Post("/monitor/unblock/:pesertaId", handlers.UnblockSiswa)
	guru.Post("/monitor/block/:pesertaId", handlers.BlockSiswa)
	guru.Post("/monitor/reset-sesi/:pesertaId", handlers.ResetSesiSiswa)
	guru.Post("/monitor/reset-ujian/:pesertaId", handlers.ResetUjianSiswa)
	guru.Get("/pengawas-jadwal", handlers.GetPengawasJadwals)
	guru.Get("/rekap-jadwal", handlers.GetLaporanJadwals)
	guru.Get("/rekap-nilai/:jadwalId", handlers.GetRekapNilai)
	guru.Get("/rekap-arsip", handlers.GetRiwayatNilai)
	guru.Get("/monitor/jawaban/:pesertaId", handlers.GetJawabanPeserta)
	guru.Post("/monitor/koreksi/:pesertaId", handlers.UpdateKoreksiEssay)
	guru.Get("/analisis/:jadwalId", handlers.GetAnalisisSoal)

	guru.Get("/bank-soal/:bankSoalId/soal", handlers.GetSoals)
	guru.Post("/bank-soal/:bankSoalId/soal", handlers.CreateSoal)
	guru.Post("/bank-soal/:bankSoalId/import", handlers.ImportSoal)
	guru.Delete("/bank-soal/:bankSoalId/soal/clear", handlers.ClearSoals)
	guru.Get("/ruang", handlers.GetRuang)
	guru.Get("/sesi", handlers.GetSesi)
	guru.Put("/soal/:id", handlers.UpdateSoal)
	guru.Delete("/soal/:id", handlers.DeleteSoal)

	// Siswa Group -> /api/siswa (Protected)
	siswa := protected.Group("/siswa")
	siswa.Use(middleware.RoleRequired("siswa"))
	siswa.Post("/validate", handlers.ValidateExam)
	siswa.Get("/jadwal", handlers.GetSiswaJadwal)
	siswa.Get("/soal/:jadwalId", handlers.GetSoalUjian)
	siswa.Get("/status/:jadwalId", handlers.GetSiswaStatus)
	siswa.Post("/sync", handlers.SyncJawaban)
	siswa.Post("/log", handlers.LogSiswa)
	siswa.Post("/submit/:pesertaId", handlers.SubmitUjian)
}
