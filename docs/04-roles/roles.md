# Roles & Responsibilities

Sistem CBT memiliki tiga role utama yang dikonfigurasi secara eksplisit dalam kolom `role` pada tabel `users`.

## 1. Admin
Role dengan otoritas tertinggi untuk mengelola seluruh aspek teknis dan data sistem.
- **Tanggung Jawab**:
    - Manajemen akun user (Admin lain, Guru, Siswa).
    - Konfigurasi data master (Kelas, Mapel, Ruang, Sesi).
    - Kontrol penuh terhadap seluruh Bank Soal dan Jadwal Ujian.
    - Monitoring aktivitas ujian secara real-time.
    - Pengarsipan nilai dan reset sistem.

## 2. Guru
Role fungsional yang berfokus pada konten pendidikan dan pengawasan.
- **Tanggung Jawab**:
    - Membuat dan mengelola Bank Soal miliknya sendiri.
    - Melihat jadwal ujian yang berkaitan dengan mata pelajarannya.
    - Menjadi pengawas (Proctor) dalam ujian.
    - Melakukan koreksi manual untuk soal tipe Essay.
    - Melihat rekap nilai dan analisis soal.

## 3. Siswa
Role pengguna akhir (end-user) sistem ujian.
- **Tanggung Jawab**:
    - Melihat daftar jadwal ujian yang aktif untuk kelasnya.
    - Mengikuti ujian dengan memvalidasi token.
    - Mengirimkan jawaban (sync) secara real-time.
    - Melihat hasil nilai (jika diizinkan oleh admin).

## 4. Pengawas (Fungsional)
Meskipun bukan role terpisah di tabel `users`, **Pengawas** adalah status fungsional yang diberikan kepada **Guru** melalui tabel `cbt_jadwal_ujians`. Seorang guru yang ditunjuk sebagai pengawas memiliki akses untuk memantau status pengerjaan siswa pada jadwal tersebut.

## 5. Role Lain (Planned/None)
- **Operator**: Saat ini fungsi operator digabung ke dalam role Admin.
- **Guest**: Tidak ada akses guest; seluruh halaman (selain login) bersifat terproteksi.
