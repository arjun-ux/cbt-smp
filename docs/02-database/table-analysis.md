# Table Analysis

Analisis mendalam terhadap fungsi dan siklus hidup data pada setiap tabel utama.

## 1. Tabel Inti (Core)

### `users`
- **Fungsi**: Single source of truth untuk kredensial login.
- **Lifecycle**: Data dibuat saat pendaftaran (Admin/Import). Status `is_active` menentukan apakah user bisa masuk ke sistem.
- **Relasi**: One-to-One dengan `master_siswas` atau `master_gurus`.

### `master_siswas`
- **Fungsi**: Profil lengkap siswa termasuk penempatan akademik.
- **Lifecycle**: Data biasanya di-import secara bulk di awal semester. Terikat secara permanen dengan `users`.
- **Flow**: Digunakan oleh sistem saat login siswa untuk menentukan jadwal mana yang tersedia berdasarkan `kelas_id`, `ruang_id`, dan `sesi_id`.

## 2. Tabel Konten Ujian

### `cbt_bank_soals`
- **Fungsi**: Wadah koleksi soal per mata pelajaran dan tingkat.
- **Lifecycle**: `Draft` (sedang dibuat) -> `Aktif` (siap dijadwalkan). Tidak bisa dihapus jika sudah ada `cbt_jadwal_ujians` yang merujuk padanya (tergantung constraint database).
- **Pengguna**: Guru (Pemilik) dan Admin.

### `cbt_soals`
- **Fungsi**: Butir soal ujian.
- **Lifecycle**: Data tetap. Berubah hanya jika ada revisi soal sebelum ujian dimulai.
- **Penting**: Kolom `kunci_jawaban` adalah data sensitif yang hanya boleh diakses oleh backend saat proses scoring.

## 3. Tabel Pelaksanaan Ujian (Transaction)

### `cbt_jadwal_ujians`
- **Fungsi**: Konfigurasi waktu dan syarat akses ujian.
- **Lifecycle**: `Belum Mulai` -> `Berlangsung` (token aktif) -> `Selesai` (akses ditutup).
- **Relasi**: Menghubungkan `BankSoal` dengan `Sesi` dan `Ruang`.

### `cbt_peserta_ujians`
- **Fungsi**: Mencatat sesi aktif pengerjaan ujian tiap siswa.
- **Lifecycle**: Dibuat saat siswa validasi token. Berakhir saat siswa klik submit atau waktu habis.
- **Penting**: Kolom `sisa_waktu_detik` diupdate secara berkala (autosave) untuk mencegah data loss.
- **Data Integrity**: Menyimpan snapshot nilai (`nilai_pg`, `nilai_essay`) setelah ujian selesai.

### `cbt_jawaban_siswas`
- **Fungsi**: Penyimpanan detail jawaban per butir soal.
- **Lifecycle**: Insert/Update secara dinamis selama ujian berlangsung (`SyncJawaban`).
- **Relasi**: Foreign Key ke `peserta_ujian_id` dan `soal_id`.

## 4. Tabel Monitoring & Log

### `log_ujians`
- **Fungsi**: Audit trail aktivitas siswa.
- **Flow**: Frontend mengirim sinyal (pindah tab, koneksi putus) -> Backend simpan ke log -> Admin Dashboard baca log secara real-time.
- **Lifecycle**: Data hanya insert (append-only), biasanya diarsipkan setelah periode ujian selesai.

## Analisis Query & Indexing
- **Indexing**: GORM secara otomatis membuat index pada Primary Key (`id`) dan kolom dengan tag `uniqueIndex` (seperti `username`, `nisn`, `kode_ruang`).
- **Foreign Keys**: Database dikonfigurasi dengan `PRAGMA foreign_keys = ON`. Relasi menggunakan `CASCADE` pada operasi delete untuk memastikan tidak ada data yatim (orphan data), kecuali pada relasi master data tertentu yang menggunakan `SET NULL`.
- **Transactions**: Digunakan secara ketat pada modul `SyncJawaban` dan `SubmitUjian` untuk memastikan sinkronisasi antara sisa waktu dan jawaban yang tersimpan tetap konsisten.
