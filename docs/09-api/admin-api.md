# Administrative API Reference

Endpoint untuk pengelolaan data master, konten ujian, dan laporan.

## 1. Master Data Management (Admin Only)

### 1.1 Siswa CRUD
- **Create**: `POST /api/admin/siswa`
- **Update**: `PUT /api/admin/siswa/:id`
- **Delete**: `DELETE /api/admin/siswa/:id`
- **Import**: `POST /api/admin/siswa/import` (Upload CSV)
- **Status Toggle**: `PATCH /api/admin/siswa/:id/status`
- **Bulk Status Update**: `PATCH /api/admin/siswa/bulk-status`
- **Bulk Plotting**: `POST /api/admin/siswa/bulk-plot`
- **Side Effects**: Mengubah tabel `users` (untuk kredensial) dan `master_siswas` (untuk biodata).

### 1.2 Guru, Kelas, Mapel, Ruang, Sesi
Endpoint standar CRUD untuk data master sekolah.
- **Base Route**: `/api/admin/[resource]`
- **Methods**: `GET`, `POST`, `PUT`, `DELETE`.

---

## 2. Bank Soal & Soal

### 2.1 Bank Soal CRUD
- **Route**: `/api/[role]/bank-soal`
- **Access**: `admin`, `guru`.
- **Logic**: Guru hanya bisa memodifikasi Bank Soal yang `guru_id`-nya sesuai dengan ID mereka. Admin memiliki akses global.

### 2.2 Soal Detail
- **Route**: `/api/[role]/bank-soal/:bankSoalId/soal`
- **Import Soal**: `POST /api/[role]/bank-soal/:bankSoalId/import`
    - *Note*: Endpoint ini tersedia di backend, namun antarmuka saat ini menggunakan metode *Direct Looping POST* dari frontend setelah parsing dokumen di sisi klien.
    - Mendukung ekstraksi gambar otomatis via utilitas `ExtractBase64Images`.
    - Side Effect: Batch insert atau individual insert ke `cbt_soals`.

---

## 3. Jadwal Ujian

### 3.1 Create Jadwal
- **Route**: `POST /api/admin/jadwal`
- **Side Effect**: Menghasilkan `token_ujian` acak secara otomatis.

### 3.2 Update & Manage Status
- **Route**: `PUT /api/admin/jadwal/:id`
- **Tujuan**: Mengubah data jadwal atau mengakhiri ujian.
- **Side Effect (Critical)**: Jika status diubah menjadi **"Selesai"**, sistem secara otomatis melakukan **Batch Force Submit** untuk seluruh siswa yang belum menekan tombol selesai (Menghitung nilai & menutup sesi).

### 3.3 Refresh Token
- **Route**: `PATCH /api/admin/jadwal/:id/token`
- **Tujuan**: Mengganti token ujian jika terjadi kebocoran keamanan.

### 3.3 Archive Results
- **Route**: `POST /api/admin/jadwal/:id/archive`
- **Logic**: Menyalin data dari `cbt_peserta_ujians` ke `cbt_rekap_nilais`.
- **Status Update**: Mengubah status jadwal menjadi `Diarsipkan`.

---

## 4. Monitoring & Proctoring

### 4.1 Real-time Monitoring
- **Route**: `GET /api/[role]/monitor/:jadwalId`

### 4.2 Proctors Actions
- **Force Submit**: `POST /api/[role]/monitor/force-submit/:pesertaId`
- **Unblock**: `POST /api/[role]/monitor/unblock/:pesertaId`
- **Block Manual**: `POST /api/[role]/monitor/block/:pesertaId`
    - *Side Effect*: Mengubah `is_terblokir=true` dan `waktu_login=null` (Force Kick).
- **Reset Sesi**: `POST /api/[role]/monitor/reset-sesi/:pesertaId`
- **Reset Ujian (Full)**: `POST /api/[role]/monitor/reset-ujian/:pesertaId`

---

## 5. Scoring & Reporting

### 4.1 Update Koreksi Essay
- **Route**: `POST /api/[role]/monitor/koreksi/:pesertaId`
- **Request Body**: `[{ "soal_id": 1, "skor": 10.5 }]`
- **Side Effect**: Update `nilai_essay` di tabel peserta.

### 4.2 Analisis Butir Soal
- **Route**: `GET /api/[role]/analisis/:jadwalId`
- **Response**: Statistik tingkat kesukaran dan sebaran pilihan jawaban untuk setiap butir soal.
