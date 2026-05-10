# Bank Soal Management Page

## Informasi Halaman
- **Nama Halaman**: Manajemen Bank Soal
- **Route/URL**: `/admin/bank-soal` atau `/guru/bank-soal`
- **Tujuan**: Mengelola wadah (folder) koleksi soal per mata pelajaran.
- **Role Akses**: `admin`, `guru`

## Teknis & State
- **Component**: `BankSoalAdmin.vue`
- **API Calls**:
    - `GET /api/[role]/bank-soal`: Ambil daftar wadah.
    - `POST/PUT /api/[role]/bank-soal`: Simpan/Update wadah.
    - `DELETE /api/[role]/bank-soal/:id`: Hapus wadah.
- **State Lokal**:
    - `bankSoals`: Array data bank soal.
    - `form`: Object data input (judul, mapel, tingkat, bobot).
    - `showForm`: Boolean (Toggle modal input).

## Elemen UI & Aksi

### Tombol "Buat Wadah Baru"
- **Fungsi**: Membuka modal form kosong.
- **Reset**: Mengosongkan objek `form` ke nilai default.

### Tombol "Kelola Soal" (Action Table)
- **Fungsi**: Navigasi ke detail butir soal.
- **Redirect**: Ke `/admin/bank-soal/:id/soal` (Sesuai role).

### Dropdown Mapel & Guru (Form)
- **Fungsi**: Mengaitkan bank soal dengan data master.
- **Data Source**: Diambil dari API master data saat `onMounted`.

### Tombol "Hapus" (Action Table)
- **Flow**:
    1. Klik hapus.
    2. Muncul `ConfirmModal` (Peringatan bahwa soal di dalamnya akan ikut terhapus).
    3. Jika konfirmasi: Kirim `DELETE` request.
    4. Refresh data tabel via `fetchData`.
- **Database Effect**: `cbt_bank_soals` record dihapus. `cbt_soals` yang terkait akan terhapus via database `ON DELETE CASCADE`.

## Validasi & Error Handling
- **Constraint**: Bank soal yang sudah terikat dengan `JadwalUjian` mungkin gagal dihapus tergantung pada constraint database (Foreign Key). Sistem menangkap pesan error "Data sedang digunakan" dari backend.
