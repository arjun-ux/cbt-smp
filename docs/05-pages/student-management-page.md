# Student Management Page

## Informasi Halaman
- **Nama Halaman**: Manajemen Siswa
- **Route/URL**: `/admin/siswa`
- **Tujuan**: Pengelolaan data master peserta ujian (CRUD + Import).
- **Role Akses**: `admin`

## Teknis & State
- **Component**: `SiswaMaster.vue`
- **API Calls**:
    - `GET /api/admin/siswa`: Ambil list seluruh siswa.
    - `POST/PUT /api/admin/siswa`: Simpan/Update profil siswa.
    - `DELETE /api/admin/siswa/:id`: Hapus siswa.
    - `PATCH /api/admin/siswa/:id/status`: Toggle aktif/nonaktif akun (Tunggal).
    - `PATCH /api/admin/siswa/bulk-status`: Update status aktif/nonaktif banyak siswa sekaligus.
    - `POST /api/admin/siswa/import`: Upload file CSV massal.
- **State Lokal**:
    - `siswas`: Array data siswa.
    - `selectedIds`: Array (Untuk seleksi bulk action).
    - `previewData`: Array (Data pratinjau sebelum konfirmasi import CSV).

## Elemen UI & Aksi

### Tombol "Tambah Siswa"
- **Fungsi**: Membuka modal form pendaftaran siswa tunggal.
- **Data Penempatan**: Dropdown Kelas, Ruang, dan Sesi yang diambil dari master data.

### Fitur "Import CSV"
- **Fungsi**: Pendaftaran siswa secara masal menggunakan file spreadsheet.
- **Flow**:
    1. Klik "Import".
    2. Pilih file `.csv`.
    3. Sistem mem-parsing file dan menampilkan modal `Preview`.
    4. Admin memvalidasi data di layar preview.
    5. Klik "Konfirmasi & Impor".
- **Database Effect**: Membuat record baru di tabel `users` dan `master_siswas`.

### Tombol "Status" (Badge Aktif/Nonaktif)
- **Fungsi**: Mematikan akses login siswa tanpa menghapus datanya secara tunggal.
- **Toggle**: Mengubah kolom `is_active` di tabel `users`.

### Aksi Masal (Bulk Action)
Muncul saat admin mencentang satu atau lebih siswa pada tabel:
- **Hapus**: Menghapus banyak siswa sekaligus.
- **Aktifkan**: Mengaktifkan kembali akses login banyak siswa terpilih.
- **Nonaktifkan**: Mematikan akses login banyak siswa terpilih (Berguna untuk kelulusan/angkatan).
- **Atur Ruang/Sesi (Bulk Plot)**: Menentukan lokasi ujian (Ruang) dan waktu (Sesi) untuk banyak siswa sekaligus.

### Database Effect (Bulk Status)
Melakukan *batch update* pada tabel `users` melalui relasi `user_id` di tabel `master_siswas` menggunakan subquery SQL yang efisien.

## Validasi & Pagination
- **Search**: Pencarian lintas kolom (NISN, Nama, Kelas) menggunakan computed property.
- **Pagination**: Menampilkan data per 10 item untuk menjaga performa rendering browser.
