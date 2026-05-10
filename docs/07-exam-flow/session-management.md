# Session & Connection Management

Sistem CBT menggunakan strategi **Stateful Participant** namun **Stateless Auth** untuk mengelola sesi ujian.

## 1. Lifecycle Sesi Peserta

### Login & Lock (Sesi Dimulai)
Saat siswa memvalidasi token, backend melakukan:
- Set `waktu_login = time.Now()`.
- Record ini berfungsi sebagai "Soft Lock". Jika ada request login lain dengan NISN yang sama, sistem akan melihat `waktu_login` tidak NULL dan menolak akses (mencegah pengerjaan satu akun oleh dua orang).

### Sinkronisasi (Sesi Berjalan)
Selama ujian, setiap request `SyncJawaban` berfungsi sebagai "Heartbeat". 
- Backend mengetahui siswa masih aktif mengerjakan.
- Jika dalam monitoring siswa terlihat tidak ada progres (jawaban tidak bertambah), pengawas bisa mencurigai adanya kendala teknis.

### Logout & Reset (Sesi Berakhir)
- **Normal**: Saat Submit, `waktu_login` tetap ada (sebagai histori) namun `status_ujian` menjadi `Selesai`, menutup pintu masuk kembali.
- **Abnormal (Reset Sesi)**: Jika PC siswa rusak, pengawas melakukan **Reset Sesi**.
    - Backend mengubah `waktu_login` menjadi `NULL`.
    - Ini membuka kunci agar siswa bisa login kembali dari PC lain.

## 2. Penanganan Koneksi (Disconnect/Reconnect)

### Deteksi Offline di Client
Frontend memonitor status request API. Jika gagal (timeout/no network):
1. Mengubah indikator UI menjadi **"Koneksi Terputus"**.
2. **Offline-First Mode**: Semua input jawaban tetap diterima dan disimpan ke `localStorage.setItem('answers_backup', ...)`.
3. **Queueing**: ID soal yang diubah dimasukkan ke dalam antrean (Set) `dirtyQuestions`.

### Pemulihan (Recovery)
Saat internet kembali:
1. Siklus auto-sync berikutnya (15 detik) akan mencoba mengirimkan data.
2. Jika sukses, `localStorage` dibersihkan dan status kembali **Synced**.
3. Jika siswa melakukan **Refresh Halaman** saat internet mati:
    - Data soal hilang dari memori.
    - Saat internet menyala dan siswa login kembali, frontend akan membandingkan data dari Server vs data di `LocalStorage`.
    - Jika ada data di `LocalStorage` yang lebih baru/tidak ada di server, frontend akan memicu sync paksa untuk memulihkan jawaban yang tertunda (Recovery Logic).

## 3. Keamanan Sesi
- **JWT Expiry**: Token ujian memiliki masa aktif terbatas (biasanya durasi ujian + 30 menit).
- **Tab Switching**: Pelanggaran dicatat secara permanen di tabel `log_ujians`. Membuka blokir memerlukan intervensi manual pengawas untuk menjamin integritas.

## 4. Rekomendasi Infrastruktur
Karena sistem melakukan polling/sync setiap 15 detik per siswa:
- Jika ada 1000 siswa, akan ada ~66 request per detik (RPS).
- **Optimasi**: Database SQLite menggunakan mode WAL (*Write-Ahead Logging*) untuk menangani konkurensi tulis yang tinggi saat sinkronisasi massal.
