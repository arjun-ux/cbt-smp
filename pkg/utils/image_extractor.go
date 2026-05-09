package utils

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/google/uuid"
)

// ExtractBase64Images mencari gambar base64 dalam HTML, menyimpannya ke disk,
// dan mengganti src dengan path file statis.
func ExtractBase64Images(htmlContent string) string {
	if !strings.Contains(htmlContent, "data:image") {
		return htmlContent
	}

	// Regex yang lebih longgar: mencari data:image/...;base64, lalu mengambil isinya sampai bertemu tanda petik
	// Ini lebih aman untuk menangani data base64 dari Word yang sangat panjang dan berantakan
	re := regexp.MustCompile(`data:image/(?P<ext>[a-zA-Z0-9]+);base64,(?P<data>[^"']+)`)
	
	// Pastikan folder upload ada
	uploadDir := filepath.Join("uploads", "soal")
	_ = os.MkdirAll(uploadDir, 0755)

	newContent := htmlContent
	
	// Kita gunakan FindAllStringSubmatch untuk mendapatkan semua gambar
	matches := re.FindAllStringSubmatch(htmlContent, -1)
	
	for _, match := range matches {
		if len(match) < 3 {
			continue
		}

		fullBase64 := match[0]
		ext := match[1]
		b64Data := match[2]

		// Membersihkan karakter yang tidak valid (spasi, baris baru, dll) 
		// yang sering muncul saat pembedahan file Word
		b64Data = strings.ReplaceAll(b64Data, " ", "")
		b64Data = strings.ReplaceAll(b64Data, "\n", "")
		b64Data = strings.ReplaceAll(b64Data, "\r", "")
		b64Data = strings.TrimSpace(b64Data)

		// Decode base64
		decoded, err := base64.StdEncoding.DecodeString(b64Data)
		if err != nil {
			// Jika gagal decode standar, coba dengan URLEncoding
			decoded, err = base64.URLEncoding.DecodeString(b64Data)
			if err != nil {
				continue
			}
		}

		// Buat nama file unik
		fileName := fmt.Sprintf("%s.%s", uuid.New().String(), ext)
		filePath := filepath.Join(uploadDir, fileName)

		// Simpan ke disk
		err = os.WriteFile(filePath, decoded, 0644)
		if err != nil {
			continue
		}

		// Ganti Base64 dengan URL path yang bisa diakses via browser
		urlPath := fmt.Sprintf("/uploads/soal/%s", fileName)
		newContent = strings.Replace(newContent, fullBase64, urlPath, -1)
	}

	return newContent
}

// DeleteImagesFromHTML mencari semua path gambar /uploads/soal/ dalam HTML 
// dan menghapus file fisiknya dari disk.
func DeleteImagesFromHTML(htmlContent string) {
	if htmlContent == "" || !strings.Contains(htmlContent, "/uploads/soal/") {
		return
	}

	// Regex untuk menangkap nama file setelah /uploads/soal/
	re := regexp.MustCompile(`/uploads/soal/([a-zA-Z0-9\-\.]+)`)
	matches := re.FindAllStringSubmatch(htmlContent, -1)

	for _, match := range matches {
		if len(match) < 2 {
			continue
		}
		fileName := match[1]
		filePath := filepath.Join("uploads", "soal", fileName)

		// Hapus file dari disk
		_ = os.Remove(filePath)
	}
}
