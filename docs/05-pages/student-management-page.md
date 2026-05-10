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
    - `PATCH /api/admin/siswa/:id/status`: Toggle aktif/nonaktif akun.
    - `POST /api/admin/siswa/import`: Upload file CSV massal.
- **State Lokal**:
    - `siswas`: Array data siswa.
    - `selectedIds`: Array (Untuk hapus masal/bulk delete).
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
- **Fungsi**: Mematikan akses login siswa tanpa menghapus datanya.
- **Toggle**: Mengubah kolom `is_active` di tabel `users`.

### Hapus Masal (Bulk Action)
- **Fungsi**: Menghapus banyak siswa sekaligus yang telah dicentang pada checkbox.
- **Flow**: Iterasi penghapusan via API untuk setiap ID yang dipilih.

### Bulk Plot (Atur Ruang & Sesi Masal)
- **Fungsi**: Menentukan lokasi ujian (Ruang) dan waktu (Sesi) untuk banyak siswa sekaligus. Sangat berguna untuk pembagian gelombang ujian.
- **Flow**:
    1. Admin mencentang satu atau lebih siswa pada tabel.
    2. Klik tombol "Atur Ruang/Sesi" yang muncul di atas tabel.
    3. Pilih Ruang dan Sesi yang diinginkan pada modal.
    4. Klik "Simpan Perubahan".
- **Database Effect**: Melakukan *batch update* pada kolom `ruang_id` dan `sesi_id` di tabel `master_siswas` menggunakan query `WHERE id IN (...)`.

## Validasi & Pagination
- **Search**: Pencarian lintas kolom (NISN, Nama, Kelas) menggunakan computed property.
- **Pagination**: Menampilkan data per 10 item untuk menjaga performa rendering browser.
