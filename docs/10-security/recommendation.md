# Security Recommendations

Berdasarkan temuan audit, berikut adalah langkah-langkah perbaikan yang direkomendasikan untuk meningkatkan keamanan sistem CBT.

## 1. Perbaikan Mendesak (Immediate Action)

### Fix IDOR Otorisasi (✅ IMPLEMENTED)
Wajib menambahkan pengecekan kepemilikan pada setiap endpoint siswa:
```go
// Perbaikan telah diterapkan di internal/handlers/ujian_siswa.go
var peserta models.CBTPesertaUjian
database.DB.Where("id = ? AND siswa_id = (SELECT id FROM master_siswas WHERE user_id = ?)", 
    input.PesertaUjianID, currentUserID).First(&peserta)
```

### Validasi Timer di Server (✅ IMPLEMENTED)
Jangan mempercayai input `sisa_waktu` dari client secara mentah.
- **Status**: Terimplementasi via `Timer Hardening` (Pencegahan nilai sisa waktu bertambah dari nilai DB).

### Unified Error Handling (✅ IMPLEMENTED)
Menyeragamkan seluruh respons API untuk mencegah kebocoran informasi teknis melalui pesan error yang tidak terduga.
- **Status**: Seluruh handler telah dimigrasikan menggunakan helper `SendError` dengan format JSON yang aman.

## 2. Penguatan Autentikasi

### Implementasi Rate Limiting
Gunakan middleware rate-limiter (misal: `fiber/middleware/limiter`) pada endpoint sensitif:
- `/api/auth/login` (Maks 5 attempt per menit per IP).
- `/api/siswa/validate` (Maks 10 attempt per menit per IP).

### JWT Security
- Hapus default secret key dari kode. Paksa aplikasi berhenti jika `JWT_SECRET` kosong.
- Kurangi durasi token siswa menjadi seumur ujian saja (misal: 2-3 jam).

## 3. Sanitasi Input & Output

### XSS Prevention
- Gunakan library sanitasi HTML (seperti `bluemonday` di Go) sebelum menyimpan pertanyaan soal ke database.
- Pastikan hanya tag HTML yang aman (seperti `<b>`, `<i>`, `<img>`, `<table>`) yang diizinkan.

### Secure Uploads
- Tambahkan validasi MIME type dan file signature pada fitur import Word (.docx).
- Batasi ukuran maksimal file upload di level middleware.

## 4. Monitoring & Logging
- Tambahkan alert sistem jika terdeteksi satu akun melakukan request dari IP yang berbeda dalam waktu singkat.
- Simpan log akses sensitif (perubahan nilai, hapus soal) ke file log eksternal yang tidak dapat dihapus oleh aplikasi.
