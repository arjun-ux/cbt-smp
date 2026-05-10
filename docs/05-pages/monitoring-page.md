# Monitoring Page

## Informasi Halaman
- **Nama Halaman**: Monitor Ujian Real-time
- **Route/URL**: `/admin/monitor/:jadwalId`
- **Tujuan**: Memantau aktivitas dan progres seluruh siswa dalam satu jadwal ujian secara live.
- **Role Akses**: `admin`, `guru` (hanya yang ditugaskan sebagai pengawas).

## Teknis & State
- **Component**: `MonitorUjian.vue`
- **API Calls**:
    - `GET /api/[role]/monitor/:jadwalId`: Polling data status peserta.
    - `POST /api/[role]/monitor/force-submit/:id`: Hentikan ujian siswa secara paksa.
    - `POST /api/[role]/monitor/unblock/:id`: Buka blokir siswa.
    - `POST /api/[role]/monitor/reset-sesi/:id`: Reset login siswa (Crash Recovery).
    - `POST /api/[role]/monitor/reset-ujian/:id`: Reset total progres ujian siswa.
- **Lifecycle**: `onMounted` memulai polling interval setiap 10 detik. `onUnmounted` membersihkan interval.

## Elemen UI & Aksi

### Kartu Statistik (Header)
- **Fungsi**: Ringkasan jumlah peserta (Total, Aktif, Selesai).
- **Live State**: Data diperbarui setiap siklus polling.

### Tabel Monitor (Body)
- **Progres Bar**: Visualisasi persentase soal yang sudah dijawab (`progres / total_soal`).
- **Badge Status**: 
    - `TERBLOKIR`: Muncul dengan animasi pulse merah jika siswa melanggar (pindah tab).
    - `Sedang Mengerjakan` / `Selesai`: Status dasar pengerjaan.

### Manajemen Aksi (Pencegahan Human Error)
Untuk meminimalisir kesalahan klik, aksi dibagi menjadi dua zona:

#### 1. Zona Aman (Ikon Langsung)
- **Buka Blokir (Ikon Gembok Orange)**: Mengizinkan siswa melanjutkan ujian jika terblokir.
- **Reset Login (Ikon Pintu Biru)**: **PENTING!** Digunakan jika PC siswa crash/hang. Aksi ini hanya mengeluarkan sesi login agar siswa bisa masuk kembali **tanpa menghapus jawaban** dan **tanpa mereset waktu**.

#### 2. Zona Berbahaya (Menu Dropdown Titik Tiga)
Aksi berisiko tinggi disembunyikan di dalam menu tambahan:
- **Hentikan Paksa**: Memaksa siswa selesai ujian. Digunakan jika waktu habis atau siswa curang secara fisik.
- **Reset Total**: **Sangat Berbahaya!** Menghapus seluruh jawaban di database, membersihkan LocalStorage di browser siswa, dan mengembalikan durasi waktu ke awal. Digunakan hanya jika siswa diizinkan mengulang dari nol.

## Validasi & Filtering
- **Search**: Memfilter baris tabel berdasarkan nama atau kelas siswa secara lokal (computed property).
- **Autorisasi**: Pengawas hanya bisa melakukan aksi jika `pengawas_id` pada jadwal sesuai dengan `guru_id` mereka.
