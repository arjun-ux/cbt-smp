# Student Dashboard

## Informasi Halaman
- **Nama Halaman**: Dashboard Siswa
- **Route/URL**: `/siswa`
- **Tujuan**: Menampilkan informasi diri siswa dan daftar jadwal ujian yang tersedia untuk diikuti.
- **Role Akses**: `siswa`

## Teknis & State
- **Component**: `StudentDashboard.vue`
- **Store**: `useAuthStore`, `useAlertStore`
- **API Calls**:
    - `GET /api/siswa/jadwal`: Mengambil daftar jadwal ujian yang relevan dengan kelas siswa tersebut.
    - `POST /api/siswa/validate`: Memvalidasi token ujian sebelum masuk ke layar pengerjaan.
- **State Lokal**:
    - `jadwals`: Array data jadwal.
    - `showTokenModal`: Boolean (Menampilkan pop-up input token).
    - `tokenInput`: String (Input token dari siswa).
    - `selectedJadwal`: Object (Jadwal yang sedang dipilih untuk dimasuki).

## Elemen UI & Aksi

### Kartu Informasi Diri (Header)
- **Fungsi**: Menampilkan Nama, NISN, Kelas, Ruang, dan Sesi.
- **Aesthetic**: Menggunakan gradient premium dan backdrop blur.

### Kartu Jadwal Ujian (Card List)
- **Status Badge**:
    - `Belum Mulai` (Amber): Ujian belum mencapai jam mulai.
    - `Berlangsung` (Blue): Ujian siap dikerjakan.
    - `Selesai` (Slate): Ujian sudah melewati batas waktu.
- **Indikator Pengerjaan**: Menampilkan ikon centang hijau jika siswa sudah mengirimkan jawaban (`Selesai`).

### Tombol "Mulai Kerjakan"
- **Fungsi**: Membuka modal validasi token.
- **Validasi**: Tombol dinonaktifkan (`disabled`) jika status bukan `Berlangsung` atau jika siswa sudah pernah mengerjakan.

### Modal Validasi Token
- **Flow**:
    1. Klik "Mulai Kerjakan".
    2. Masukkan token (case-insensitive & auto-trim).
    3. Klik "Verifikasi & Mulai".
    4. Jika valid: Server mengembalikan data sesi pengerjaan.
    5. Sesi disimpan di `sessionStorage` dengan key `exam_session`.
    6. Redirect ke `/siswa/pengerjaan/:id`.
- **Error Handling**: Menampilkan pesan "Token salah" atau "Sesi Anda berakhir".

## Keamanan & Validasi
- **Session Check**: Jika request ke API mengembalikan `401 Unauthorized`, sistem akan otomatis logout dan melempar user ke halaman login siswa.
- **Offline Defense**: Data sesi disimpan di `sessionStorage` untuk memastikan persistensi minimal saat reload halaman pengerjaan.
