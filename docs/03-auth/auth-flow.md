# Authentication Flow

Sistem CBT menggunakan mekanisme **Stateless Authentication** berbasis **JSON Web Token (JWT)**. Tidak ada sesi yang disimpan di sisi server (sessionless), sehingga skalabilitas sistem terjaga.

## 1. Login Process
Alur autentikasi saat pengguna masuk ke sistem:

1.  **Request**: Pengguna mengirimkan `username` (atau `NISN` bagi siswa) dan `password` ke endpoint `/api/auth/login`.
2.  **Verification**:
    - Backend mencari user di tabel `users`.
    - Jika tidak ditemukan sebagai username, sistem mencoba mencari di tabel `master_siswas` berdasarkan NISN.
    - Password divalidasi menggunakan `bcrypt.CompareHashAndPassword`.
3.  **Token Generation**: Jika valid, backend menghasilkan JWT yang berisi:
    - `user_id`: ID unik user.
    - `username`: Identitas user.
    - `role`: Peran user (`admin`, `guru`, `siswa`).
4.  **Response**: Token dikembalikan ke client bersama metadata user.
5.  **Storage**: Frontend menyimpan token ini di `LocalStorage` dan `Pinia Store`.

## 2. Token Validation (Middleware)
Setiap request ke rute terproteksi (`/api/admin/*`, `/api/guru/*`, `/api/siswa/*`) harus menyertakan header:
`Authorization: Bearer <token>`

Middleware `Protected()` melakukan:
- **Extraction**: Mengambil token dari header.
- **Verification**: Memvalidasi tanda tangan digital JWT menggunakan `Secret Key`.
- **DB Check**: Memastikan `user_id` yang ada di token masih terdaftar di database (mencegah token lama dari DB yang berbeda tetap bisa digunakan).
- **Context injection**: Menyimpan data user ke dalam `c.Locals` untuk digunakan oleh handler selanjutnya.

## 3. Session Validation & Expiry
- **Stateless**: Server tidak melacak siapa yang sedang login secara aktif.
- **Expiry**: Token memiliki masa berlaku (TTL). Jika expired, backend akan mengembalikan status `401 Unauthorized`.
- **Frontend handling**: Jika mendapat error `401`, frontend secara otomatis menghapus token di LocalStorage dan mengarahkan pengguna kembali ke halaman login.

## 4. Logout Process
Karena menggunakan JWT, logout sebenarnya cukup dilakukan di sisi client dengan menghapus token. Namun, sistem menyediakan endpoint `/api/auth/logout` untuk standarisasi, yang saat ini hanya mengembalikan pesan sukses.

## 5. Unauthorized & Forbidden Flow
- **Unauthorized (401)**: Terjadi jika token tidak ada, tidak valid, atau expired. Pengguna diminta login ulang.
- **Forbidden (403)**: Terjadi jika token valid, tetapi role pengguna tidak memiliki izin untuk mengakses rute tersebut (misal: Siswa mencoba akses rute Admin).
