# Riwayat Nilai Page (Brankas Nilai)

## Informasi Halaman
- **Nama Halaman**: Riwayat Nilai (Brankas Nilai)
- **Route/URL**: `/admin/riwayat-nilai`
- **Tujuan**: Sebagai pusat arsip permanen untuk melihat hasil ujian yang sudah selesai dan diarsipkan.
- **Role Akses**: `admin`, `guru` (hanya ujian miliknya).

## Teknis & State
- **Component**: `RiwayatNilai.vue`
- **API Calls**:
    - `GET /api/[role]/rekap-arsip`: Mengambil seluruh data dari tabel `cbt_rekap_nilais`.
- **State Lokal**:
    - `results`: List seluruh data arsip yang ditarik dari server.
    - `filterMapel`: Menyaring data berdasarkan Mata Pelajaran.
    - `filterKelas`: Menyaring data berdasarkan Nama Kelas.
    - `searchQuery`: Pencarian teks bebas (Nama, NISN, Judul Ujian).

## Elemen UI & Aksi

### Header "Brankas Nilai"
- **Visual**: Menggunakan tema warna **Violet** untuk menandakan area arsip permanen.
- **Tombol Ekspor (CSV)**: Menghasilkan file `.csv` dari data yang saat ini tampil (terpengaruh oleh filter dan pencarian).

### Sistem Filter & Pencarian
- **Dropdown Mapel & Kelas**: Opsi pada dropdown ini di-generate secara dinamis menggunakan *Computed Property* dari data `results` yang tersedia.
- **Search Input**: Melakukan filter reaktif terhadap Nama Siswa, NISN, dan Judul Ujian.

### Tabel Riwayat
- **Kolom Siswa**: Menampilkan Nama Lengkap, NISN, dan Kelas.
- **Kolom Ujian**: Menampilkan Judul Ujian dan label Mata Pelajaran (Violet Badge).
- **Detail Nilai (PG)**: 
    - Menampilkan Skor PG.
    - Menampilkan Badge **Benar (B)** dan **Salah (S)**.
    - **Ikon Detail**: Memungkinkan Admin/Guru melihat perbandingan seluruh jawaban siswa meskipun data sudah diarsipkan.
- **Detail Nilai (Total)**: Menampilkan total skor akhir dalam format font **Monospace**.

## Alur Data & Keamanan
1.  **Archiving Flow**: Data masuk ke halaman ini HANYA setelah Admin menekan tombol **"Arsipkan Nilai"** pada halaman Rekap Nilai. Proses ini memindahkan data dari tabel operasional ke tabel arsip permanen.
2.  **Security (RBAC)**:
    *   **Admin**: Memiliki akses penuh ke seluruh riwayat nilai sekolah.
    *   **Guru**: Backend secara otomatis memfilter data sehingga Guru hanya bisa melihat riwayat nilai dari Bank Soal yang mereka buat sendiri.

## Pagination
- Menggunakan pagination di sisi klien (Client-side) dengan limit **20 item per halaman** untuk menjaga kecepatan navigasi tanpa harus melakukan request ulang ke server setiap perpindahan halaman.
