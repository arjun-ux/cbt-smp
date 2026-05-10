# Strategi Optimasi Integritas & Performa Ujian

Dokumen ini mencatat rencana strategis untuk menangani kendala teknis saat pelaksanaan ujian CBT skala 30-100 siswa (High Traffic).

## 1. Perlindungan Data Jawaban (Anti-Lost Data)
*   **Offline-First (LocalStorage):** Setiap kali siswa menjawab, simpan jawaban di memori browser (LocalStorage) terlebih dahulu sebelum dikirim ke server.
*   **Auto-Retry Sync:** Jika pengiriman jawaban gagal (karena WiFi sekolah lemot/RTO), sistem akan mencoba mengirim ulang secara otomatis tanpa mengganggu siswa yang sedang mengerjakan.
*   **Indikator Sinkronisasi:** Menampilkan icon kecil di layar siswa (Awan Hijau = Tersimpan di Server, Awan Kuning = Menunggu Kirim) agar siswa tenang.

## 2. Validasi Sesi & Status (Anti-Status "Selesai" Misterius)
*   **Gatekeeper Validasi:** Saat siswa klik "Mulai Ujian", Backend melakukan pengecekan berlapis:
    *   Apakah siswa benar-benar belum selesai?
    *   Apakah siswa sedang login di perangkat lain?
    *   Apakah waktu ujian masih tersedia?
*   **Pembersihan Sesi Gantung:** Fitur otomatis untuk membersihkan sesi yang "nyangkut" jika siswa menutup browser tanpa logout.

## 3. Optimasi Trafik (Anti-Server Lemot)
*   **Lazy Loading Gambar:** Gambar soal hanya didownload saat soal tersebut muncul di layar, bukan didownload semua di awal.
*   **Batch Saving:** (Opsional) Mengirimkan beberapa jawaban sekaligus dalam satu tarikan napas ke server untuk mengurangi jumlah antrean (request) ke Database SQLite.

## 4. Monitoring & Audit (Investigasi Bug)
*   **Log Aktivitas Detail:** Mencatat setiap langkah krusial siswa:
    *   `10:00:01` - Siswa Login (Device: Chrome Windows)
    *   `10:00:05` - Klik Mulai Ujian
    *   `10:45:10` - Koneksi Terputus (Gagal kirim jawaban soal no 5)
    *   `10:45:15` - Koneksi Kembali (Berhasil kirim ulang soal no 5)
    *   `11:00:00` - Klik Selesai / Waktu Habis

## 5. Arsip Nilai Abadi (Long-term Data Integrity)
*   **Tabel Rekap Statis:** Membuat tabel `cbt_rekap_nilais` yang menyimpan Nama Siswa, NISN, Kelas, dan Nilai dalam bentuk teks statis (bukan relasi ID).
*   **Snapshot Policy:** Data dipindahkan ke arsip setelah ujian selesai dan diverifikasi.
*   **Deletion Safety:** Mencegah penghapusan Jadwal Ujian yang datanya belum masuk ke Arsip Nilai.
*   **Data Longevity:** Menjamin nilai tetap ada meskipun data Master Siswa atau Master Kelas sudah dihapus (misal setelah siswa lulus).
