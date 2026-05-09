package handlers

import (
	"cbt-smp/internal/database"
	"cbt-smp/internal/models"
	"cbt-smp/pkg/utils"
	encoding_csv "encoding/csv"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type SoalInput struct {
	JenisSoal    string  `json:"jenis_soal"` // PG, ESSAY
	Pertanyaan   string  `json:"pertanyaan"`
	OpsiA        string  `json:"opsi_a"`
	OpsiB        string  `json:"opsi_b"`
	OpsiC        string  `json:"opsi_c"`
	OpsiD        string  `json:"opsi_d"`
	KunciJawaban string  `json:"kunci_jawaban"`
	BobotNilai   float64 `json:"bobot_nilai"`
}

// GetSoals mengambil daftar soal berdasarkan BankSoalID
func GetSoals(c *fiber.Ctx) error {
	bankSoalId := c.Params("bankSoalId")
	role := GetUserRole(c)
	userID := GetUserID(c)

	// Proteksi Guru: Pastikan bank soal ini miliknya
	if role == "guru" {
		guru, _ := GetGuruByUserID(userID)
		var bs models.CBTBankSoal
		if err := database.DB.Where("id = ? AND guru_id = ?", bankSoalId, guru.ID).First(&bs).Error; err != nil {
			return SendError(c, fiber.StatusForbidden, "Akses ditolak: Ini bukan Bank Soal Anda")
		}
	}

	var soals []models.CBTSoal
	database.DB.Where("bank_soal_id = ?", bankSoalId).Order("id asc").Find(&soals)

	return SendSuccess(c, "Berhasil mengambil daftar soal", soals)
}

// CreateSoal menambah soal ke bank soal tertentu
func CreateSoal(c *fiber.Ctx) error {
	bankSoalId := c.Params("bankSoalId")

	// LOCKING: Cek apakah bank soal sedang aktif ujian
	var activeJadwal int64
	database.DB.Model(&models.CBTJadwalUjian{}).Where("bank_soal_id = ? AND status = ?", bankSoalId, "Berlangsung").Count(&activeJadwal)
	if activeJadwal > 0 {
		return SendError(c, fiber.StatusForbidden, "Wadah soal tidak dapat dimodifikasi karena sedang digunakan dalam ujian yang sedang berlangsung.")
	}

	role := GetUserRole(c)
	userID := GetUserID(c)

	// Proteksi Guru
	if role == "guru" {
		guru, _ := GetGuruByUserID(userID)
		var bs models.CBTBankSoal
		if err := database.DB.Where("id = ? AND guru_id = ?", bankSoalId, guru.ID).First(&bs).Error; err != nil {
			return SendError(c, fiber.StatusForbidden, "Akses ditolak: Anda tidak bisa menambah soal ke wadah ini")
		}
	}

	var input SoalInput
	if err := c.BodyParser(&input); err != nil {
		return SendError(c, fiber.StatusBadRequest, "Input tidak valid")
	}

	if input.Pertanyaan == "" {
		return SendError(c, fiber.StatusBadRequest, "Pertanyaan wajib diisi")
	}

	newSoal := models.CBTSoal{
		BankSoalID:   utils.StringToUint(bankSoalId),
		JenisSoal:    input.JenisSoal,
		Pertanyaan:   utils.ExtractBase64Images(input.Pertanyaan),
		OpsiA:        utils.ExtractBase64Images(input.OpsiA),
		OpsiB:        utils.ExtractBase64Images(input.OpsiB),
		OpsiC:        utils.ExtractBase64Images(input.OpsiC),
		OpsiD:        utils.ExtractBase64Images(input.OpsiD),
		KunciJawaban: utils.ExtractBase64Images(input.KunciJawaban),
		BobotNilai:   input.BobotNilai,
	}

	if err := database.DB.Create(&newSoal).Error; err != nil {
		return SendError(c, fiber.StatusInternalServerError, "Gagal menyimpan soal")
	}

	return SendSuccess(c, "Soal berhasil ditambahkan", newSoal)
}

// UpdateSoal mengubah data soal
func UpdateSoal(c *fiber.Ctx) error {
	id := c.Params("id")
	var soal models.CBTSoal
	if err := database.DB.First(&soal, id).Error; err != nil {
		return SendError(c, fiber.StatusNotFound, "Soal tidak ditemukan")
	}

	// LOCKING: Cek apakah bank soal sedang aktif ujian
	var activeJadwal int64
	database.DB.Model(&models.CBTJadwalUjian{}).Where("bank_soal_id = ? AND status = ?", soal.BankSoalID, "Berlangsung").Count(&activeJadwal)
	if activeJadwal > 0 {
		return SendError(c, fiber.StatusForbidden, "Soal tidak dapat diubah karena sedang digunakan dalam ujian yang sedang berlangsung.")
	}

	role := GetUserRole(c)
	userID := GetUserID(c)

	// Proteksi Guru
	if role == "guru" {
		guru, _ := GetGuruByUserID(userID)
		var bs models.CBTBankSoal
		if err := database.DB.Where("id = ? AND guru_id = ?", soal.BankSoalID, guru.ID).First(&bs).Error; err != nil {
			return SendError(c, fiber.StatusForbidden, "Akses ditolak: Anda tidak berhak mengedit soal ini")
		}
	}

	var input SoalInput
	if err := c.BodyParser(&input); err != nil {
		return SendError(c, fiber.StatusBadRequest, "Input tidak valid")
	}

	soal.JenisSoal = input.JenisSoal
	soal.Pertanyaan = utils.ExtractBase64Images(input.Pertanyaan)
	soal.OpsiA = utils.ExtractBase64Images(input.OpsiA)
	soal.OpsiB = utils.ExtractBase64Images(input.OpsiB)
	soal.OpsiC = utils.ExtractBase64Images(input.OpsiC)
	soal.OpsiD = utils.ExtractBase64Images(input.OpsiD)
	soal.KunciJawaban = utils.ExtractBase64Images(input.KunciJawaban)
	soal.BobotNilai = input.BobotNilai

	if err := database.DB.Save(&soal).Error; err != nil {
		return SendError(c, fiber.StatusInternalServerError, "Gagal memperbarui soal")
	}

	return SendSuccess(c, "Soal berhasil diperbarui", soal)
}

// DeleteSoal menghapus soal
func DeleteSoal(c *fiber.Ctx) error {
	id := c.Params("id")
	var soal models.CBTSoal
	if err := database.DB.First(&soal, id).Error; err != nil {
		return SendError(c, fiber.StatusNotFound, "Soal tidak ditemukan")
	}

	// LOCKING: Cek apakah bank soal sedang aktif ujian
	var activeJadwal int64
	database.DB.Model(&models.CBTJadwalUjian{}).Where("bank_soal_id = ? AND status = ?", soal.BankSoalID, "Berlangsung").Count(&activeJadwal)
	if activeJadwal > 0 {
		return SendError(c, fiber.StatusForbidden, "Soal tidak dapat dihapus karena sedang digunakan dalam ujian yang sedang berlangsung.")
	}

	role := GetUserRole(c)
	userID := GetUserID(c)

	// Proteksi Guru
	if role == "guru" {
		guru, _ := GetGuruByUserID(userID)
		var bs models.CBTBankSoal
		if err := database.DB.Where("id = ? AND guru_id = ?", soal.BankSoalID, guru.ID).First(&bs).Error; err != nil {
			return SendError(c, fiber.StatusForbidden, "Akses ditolak: Anda tidak berhak menghapus soal ini")
		}
	}

	// HAPUS GAMBAR FISIK: Bersihkan gambar dari disk
	utils.DeleteImagesFromHTML(soal.Pertanyaan)
	utils.DeleteImagesFromHTML(soal.OpsiA)
	utils.DeleteImagesFromHTML(soal.OpsiB)
	utils.DeleteImagesFromHTML(soal.OpsiC)
	utils.DeleteImagesFromHTML(soal.OpsiD)
	utils.DeleteImagesFromHTML(soal.KunciJawaban)

	database.DB.Delete(&soal)
	return SendSuccess(c, "Soal berhasil dihapus", nil)
}

// ClearSoals menghapus seluruh soal dalam satu bank soal (Sapu Bersih)
func ClearSoals(c *fiber.Ctx) error {
	bankSoalId := c.Params("bankSoalId")

	// LOCKING: Cek apakah bank soal sedang aktif ujian
	var activeJadwal int64
	database.DB.Model(&models.CBTJadwalUjian{}).Where("bank_soal_id = ? AND status = ?", bankSoalId, "Berlangsung").Count(&activeJadwal)
	if activeJadwal > 0 {
		return SendError(c, fiber.StatusForbidden, "Wadah soal tidak dapat dikosongkan karena sedang digunakan dalam ujian yang sedang berlangsung.")
	}

	role := GetUserRole(c)
	userID := GetUserID(c)

	// Proteksi Guru
	if role == "guru" {
		guru, _ := GetGuruByUserID(userID)
		var bs models.CBTBankSoal
		if err := database.DB.Where("id = ? AND guru_id = ?", bankSoalId, guru.ID).First(&bs).Error; err != nil {
			return SendError(c, fiber.StatusForbidden, "Akses ditolak: Anda tidak berhak mengosongkan wadah ini")
		}
	}

	// Mulai Transaksi untuk memastikan semua terhapus atau tidak sama sekali
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		// 1. Ambil semua soal dulu untuk didata gambarnya
		var soals []models.CBTSoal
		tx.Where("bank_soal_id = ?", bankSoalId).Find(&soals)
		
		for _, s := range soals {
			// 2. HAPUS GAMBAR FISIK dari disk
			utils.DeleteImagesFromHTML(s.Pertanyaan)
			utils.DeleteImagesFromHTML(s.OpsiA)
			utils.DeleteImagesFromHTML(s.OpsiB)
			utils.DeleteImagesFromHTML(s.OpsiC)
			utils.DeleteImagesFromHTML(s.OpsiD)
			utils.DeleteImagesFromHTML(s.KunciJawaban)

			// 3. HAPUS JAWABAN SISWA terkait soal ini (Manual bypass agar tidak error constraint)
			if err := tx.Where("soal_id = ?", s.ID).Delete(&models.CBTJawabanSiswa{}).Error; err != nil {
				return err
			}
		}

		// 4. Hapus seluruh soal dari database
		if err := tx.Where("bank_soal_id = ?", bankSoalId).Delete(&models.CBTSoal{}).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return SendError(c, fiber.StatusInternalServerError, "Gagal mengosongkan Bank Soal: " + err.Error())
	}
	
	return SendSuccess(c, "Seluruh soal dan jawaban lama berhasil dibersihkan", nil)
}

// ImportSoal memproses file CSV untuk menambah banyak soal sekaligus
func ImportSoal(c *fiber.Ctx) error {
	bankSoalId := c.Params("bankSoalId")

	// LOCKING: Cek apakah bank soal sedang aktif ujian
	var activeJadwal int64
	database.DB.Model(&models.CBTJadwalUjian{}).Where("bank_soal_id = ? AND status = ?", bankSoalId, "Berlangsung").Count(&activeJadwal)
	if activeJadwal > 0 {
		return SendError(c, fiber.StatusForbidden, "Tidak dapat mengimpor soal karena wadah soal sedang digunakan dalam ujian yang sedang berlangsung.")
	}

	role := GetUserRole(c)
	userID := GetUserID(c)

	// Proteksi Guru
	if role == "guru" {
		guru, _ := GetGuruByUserID(userID)
		var bs models.CBTBankSoal
		if err := database.DB.Where("id = ? AND guru_id = ?", bankSoalId, guru.ID).First(&bs).Error; err != nil {
			return SendError(c, fiber.StatusForbidden, "Akses ditolak: Anda tidak berhak mengimpor ke wadah ini")
		}
	}

	file, err := c.FormFile("file")
	if err != nil {
		return SendError(c, fiber.StatusBadRequest, "File tidak ditemukan")
	}

	f, err := file.Open()
	if err != nil {
		return SendError(c, fiber.StatusInternalServerError, "Gagal membuka file")
	}
	defer f.Close()

	reader := encoding_csv.NewReader(f)
	records, err := reader.ReadAll()
	if err != nil {
		return SendError(c, fiber.StatusInternalServerError, "Format CSV tidak valid")
	}

	bsID := utils.StringToUint(bankSoalId)
	var newSoals []models.CBTSoal

	// Lewati header (i=0)
	for i, record := range records {
		if i == 0 || len(record) < 8 {
			continue
		}

		bobot, _ := strconv.ParseFloat(record[7], 64)
		
		newSoals = append(newSoals, models.CBTSoal{
			BankSoalID:   bsID,
			JenisSoal:    strings.ToUpper(record[0]), // PG / ESSAY
			Pertanyaan:   utils.ExtractBase64Images(record[1]),
			OpsiA:        utils.ExtractBase64Images(record[2]),
			OpsiB:        utils.ExtractBase64Images(record[3]),
			OpsiC:        utils.ExtractBase64Images(record[4]),
			OpsiD:        utils.ExtractBase64Images(record[5]),
			KunciJawaban: utils.ExtractBase64Images(record[6]),
			BobotNilai:   bobot,
		})
	}

	if err := database.DB.Create(&newSoals).Error; err != nil {
		return SendError(c, fiber.StatusInternalServerError, "Gagal menyimpan data import")
	}

	return SendSuccess(c, "Berhasil mengimpor "+strconv.Itoa(len(newSoals))+" soal", nil)
}
