# Layar Pengerjaan Ujian (Exam Page)

## Informasi Halaman
- **Nama Halaman**: Layar Pengerjaan Ujian (Modern Professional UI)
- **Route/URL**: `/siswa/pengerjaan/:jadwalId`
- **Tujuan**: Fokus utama siswa dalam mengerjakan butir-butir soal dengan antarmuka yang tenang dan responsif.
- **Role Akses**: `siswa`

## Desain & Antarmuka (Aesthetic Upgrade)
Halaman ini menggunakan bahasa desain **"Indigo-Glass"** yang berfokus pada profesionalisme dan konsentrasi:
- **Glassmorphism Header/Footer**: Menggunakan efek transparansi kabur (`backdrop-blur`) untuk kesan premium.
- **Indigo Palette**: Skema warna biru indigo untuk memberikan kesan aplikasi enterprise yang stabil.
- **Compact Sidebar**: Informasi peserta dan status sinkronisasi dikemas ringkas.
- **Interactive Question Box**: Kartu soal dengan tipografi yang dioptimalkan.

## Teknis & State
- **Component**: `UjianPengerjaan.vue`
- **Store**: `useAuthStore`, `useAlertStore`, `useExamStore`
- **API Calls**:
    - `GET /api/siswa/soal/:jadwalId`: Mengambil soal dan sisa waktu.
    - `POST /api/siswa/sync`: Menyimpan jawaban batch & **Polling Status Blokir**.
    - `POST /api/siswa/log`: Mencatat pelanggaran keamanan (pindah tab).
    - `POST /api/siswa/submit/:pesertaId`: Finalisasi ujian.

## Elemen UI & Aksi

### Navigasi Soal (Modern Navigator)
- **Status Indikator**: Indigo (Sudah dijawab), Amber (Ragu-ragu), Slate (Belum dikerjakan).
- **Active State**: Memiliki bingkai indigo tebal dan indikator titik aktif.

### Tombol "Selesai" (Smart Finish)
- Memiliki efek *shimmer* visual dan hanya aktif jika seluruh soal telah dijawab (`isAllAnswered`).

## Logika Keamanan & Validasi

### 1. Smart Auto-Unlock
Aplikasi memiliki kemampuan untuk membuka blokir secara otomatis:
- Meskipun layar terkunci (`isTerblokir: true`), sistem tetap melakukan polling ke server via `/api/siswa/sync` secara periodik.
- Jika Admin melepas blokir di Dashboard, respon server akan berubah menjadi `200 OK`.
- Frontend menangkap respon sukses tersebut dan otomatis menghilangkan layar kunci tanpa perlu refresh.

### 2. Anti-Cheat (Client-Side)
- **Pindah Tab**: Mendengarkan event `blur`. Memasukkan log pelanggaran dan mengunci layar seketika.
- **Fullscreen**: Dipaksa aktif saat masuk halaman pengerjaan.

### 3. Sync Status
- **Cloud Synced**: Semua jawaban aman di server.
- **Queueing**: Ada jawaban di antrean yang sedang menunggu pengiriman otomatis.
- **Offline Mode**: Koneksi terputus, sistem menggunakan LocalStorage sebagai fallback.
