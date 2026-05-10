# Authentication API

Manajemen akses pengguna ke dalam sistem.

## 1. Login
Digunakan oleh Admin, Guru, dan Siswa untuk mendapatkan token JWT.

- **Route**: `POST /api/auth/login`
- **Method**: `POST`
- **Auth**: Public
- **Request Body**:
    ```json
    {
      "username": "admin",
      "password": "password123"
    }
    ```
- **Response (200 OK)**:
    ```json
    {
      "token": "eyJhbGciOi...",
      "user": {
        "id": 1,
        "username": "admin",
        "role": "admin"
      }
    }
    ```
- **Validation**: 
    - Username dan Password wajib diisi.
    - Cek status `is_active` di tabel `users`.
- **Side Effect**: Mencatat login ke `log_sistems`.

## 2. Logout
Membersihkan sesi pada sisi server (jika diperlukan) dan sisi client.

- **Route**: `POST /api/auth/logout`
- **Method**: `POST`
- **Auth**: Required (Bearer Token)
- **Response (200 OK)**:
    ```json
    {
      "message": "Logout berhasil"
    }
    ```
- **Side Effect**: (Optional) Invalidate token jika menggunakan redis/blacklist. Saat ini sistem bersifat stateless.
