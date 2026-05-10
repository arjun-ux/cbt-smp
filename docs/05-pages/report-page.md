# Report Page

## Informasi Halaman
- **Nama Halaman**: Laporan & Nilai
- **Route/URL**: `/admin/rekap-nilai/:jadwalId?`
- **Tujuan**: Melihat hasil akhir ujian, melakukan koreksi manual, dan ekspor data.
- **Role Akses**: `admin`, `guru`

## Teknis & State
- **Component**: `RekapNilai.vue`
- **API Calls**:
    - `GET /api/[role]/rekap-nilai/:jadwalId`: Ambil data nilai seluruh peserta.
    - `GET /api/[role]/monitor/jawaban/:id`: Ambil detail jawaban per butir untuk koreksi.
    - `POST /api/[role]/monitor/koreksi/:id`: Simpan nilai essay.
    - `POST /api/admin/jadwal/:id/archive`: Kunci nilai secara permanen.
- **State Lokal**:
    - `results`: List nilai siswa.
    - `isKoreksiModalOpen`: Toggle modal koreksi essay.
    - `jawabanPeserta`: List butir jawaban (filter khusus Essay) untuk dinilai.

## Elemen UI & Aksi

### Selektor Jadwal (Dropdown)
- **Fungsi**: Memilih jadwal ujian yang ingin dilihat laporannya.
- **Event**: Perubahan nilai memicu `fetchData` dan update URL via `router.replace`.

### Tombol "Koreksi" (Action Table)
- **Fungsi**: Membuka modal yang berisi pertanyaan dan jawaban essay siswa.
- **Flow**:
    1. Klik ikon koreksi.
    2. Fetch detail jawaban dari server.
    3. Guru memasukkan angka skor pada input number.
    4. Klik "Simpan Semua Nilai".
- **Database Effect**: Mengupdate kolom `skor` di tabel `cbt_jawaban_siswas` dan otomatis mengupdate `nilai_essay` & `total_nilai` di tabel `cbt_peserta_ujians` via API backend.

### Tombol "Ekspor CSV"
- **Fungsi**: Mengunduh data nilai dalam format `.csv`.
- **Flow**: Memproses data reaktif `filteredResults` menjadi string CSV dan memicu download browser secara instan.

### Tombol "Arsipkan Nilai" (Hanya Admin)
- **Fungsi**: Memindahkan data nilai ke tabel `cbt_rekap_nilais` (Brankas Nilai) dan mengunci jadwal.
- **Efek**: Jadwal tidak bisa lagi diubah atau dikoreksi.

## Ringkasan Statistik (Stats Cards)
- **Average/Max/Min**: Kalkulasi otomatis (computed) dari data nilai yang ditampilkan, mempermudah guru melihat sebaran kemampuan siswa.
