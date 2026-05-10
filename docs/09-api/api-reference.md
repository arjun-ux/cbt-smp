# API Reference Overview

Dokumentasi ini mencakup seluruh endpoint API yang tersedia pada sistem CBT. Semua request API harus memiliki prefix `/api`.

## 1. Struktur Request/Response Umum

### Request Header
Untuk endpoint yang terproteksi, wajib menyertakan token JWT:
```http
Authorization: Bearer <your_token>
Content-Type: application/json
```

### Format Response Sukses
```json
{
  "success": true,
  "message": "Pesan keberhasilan",
  "data": { ... }
}
```
*Catatan: Field `data` bersifat opsional tergantung pada endpoint.*

### Format Response Error
```json
{
  "success": false,
  "error": "Pesan kesalahan detail",
  "code": "ERROR_CODE_OPTIONAL"
}
```
*Catatan: Field `code` digunakan untuk identifikasi error spesifik di frontend (seperti `EXAM_FINISHED` atau `ACCOUNT_BLOCKED`).*

## 2. Pengelompokan API
API dikelompokkan berdasarkan fungsionalitas dan hak akses:

1.  **[Auth API](auth-api.md)**: Manajemen login, logout, dan sesi.
2.  **[Exam API](exam-api.md)**: Alur pengerjaan ujian bagi siswa dan monitoring bagi pengawas.
3.  **[Admin API](admin-api.md)**: Pengelolaan data master (Guru, Siswa, Kelas, Mapel), Bank Soal, dan Jadwal Ujian.

## 3. Status Code
| Code | Deskripsi |
| :--- | :--- |
| `200 OK` | Request berhasil. |
| `201 Created` | Resource baru berhasil dibuat. |
| `400 Bad Request` | Input tidak valid atau parameter kurang. |
| `401 Unauthorized` | Token tidak valid atau expired. |
| `403 Forbidden` | Role tidak memiliki akses ke endpoint tersebut. |
| `404 Not Found` | Resource tidak ditemukan. |
| `500 Internal Server Error` | Terjadi kesalahan pada server. |
