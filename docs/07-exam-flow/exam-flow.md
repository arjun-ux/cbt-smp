# Deep Exam Flow Analysis

Sistem CBT ini dirancang sebagai aplikasi *high-concurrency* yang memprioritaskan integritas data dan keamanan ujian. Berikut adalah bedah tuntas alur ujian dari sisi Siswa, Pengawas, dan Admin.

## 1. Alur Siswa (Student Journey)

### Tahap 1: Inisiasi & Validasi
1.  **Dashboard**: Siswa melihat daftar jadwal yang aktif.
2.  **Validasi Token**: Siswa memasukkan token ujian.
    - **Backend**: Validasi kesesuaian `JadwalID`, `SesiID`, `RuangID`, dan `TingkatKelas`.
    - **Session Check**: Memeriksa kolom `waktu_login` di tabel `cbt_peserta_ujians`. Jika tidak NULL, akses ditolak (Multi-login prevention).
3.  **Generate Session**: Jika valid, server membuat/mengupdate record `CBTPesertaUjian` dan mengembalikan JWT khusus ujian.

### Tahap 2: Pengerjaan (Live Session)
1.  **Fetch Questions**: Frontend mengambil soal. Jika `AcakSoal` aktif, DB mengirim urutan acak (`RANDOM()`).
2.  **Security Lock**: Browser masuk ke mode Fullscreen. Event `blur` dipasang untuk mendeteksi pindah tab.
3.  **Real-time Sync**: Jawaban disimpan secara lokal (LocalStorage) dan dikirim ke server dalam batch setiap 15 detik atau saat perpindahan soal.

### Tahap 3: Finalisasi
1.  **Submit Manual**: Siswa menekan tombol "Selesai" (Hanya aktif jika semua soal terisi).
2.  **Automatic Scoring**: Backend menghitung skor PG secara instan dan mengunci record peserta (`status_ujian = 'Selesai'`).

---

## 2. Alur Pengawas (Proctor Journey)

1.  **Monitoring**: Pengawas melihat dashboard real-time yang melakukan polling data setiap 10 detik.
2.  **Intervensi Pelanggaran**: Jika siswa terblokir (Auto-block), pengawas melakukan "Unblock" setelah memberikan peringatan.
3.  **Technical Support**: Jika PC siswa mati, pengawas melakukan "Reset Sesi" agar siswa bisa login kembali tanpa kehilangan jawaban (karena sudah tersinkronisasi di DB).
4.  **Force Submit**: Pengawas dapat menghentikan ujian siswa yang terdeteksi melakukan kecurangan berat atau waktu sudah habis namun tidak menekan submit.

---

## 3. Alur Admin (Master Flow)

1.  **Configuration**: Admin mengatur `Bank Soal` dan `Jadwal Ujian`.
2.  **Token Management**: Admin dapat mengubah atau me-refresh token ujian jika bocor.
3.  **Final Audit**: Setelah ujian selesai, Admin melakukan "Archive Nilai" untuk mengunci data secara permanen ke dalam Brankas Nilai (`cbt_rekap_nilais`).

---

## 4. Analisis Teknis Deep-Dive

### Database Updates during Exam
- `cbt_peserta_ujians`: Update `sisa_waktu_detik` dan `is_terblokir`.
- `cbt_jawaban_siswas`: Insert/Update jawaban (Batch Sync).
- `log_ujians`: Insert setiap aktivitas penting (Login, Sync, Pindah Tab).

### Anti-Cheat Mechanism
- **Client-Side**: Detection of `visibilitychange` and `blur` events.
- **Server-Side**: Stateless JWT validation dan strict Role-Based access.
- **Data Integrity**: Penutupan akses otomatis jika `status_ujian` sudah `Selesai`.

### Kemungkinan Race Condition
- **Double Sync**: Jika siswa menekan tombol navigasi dengan sangat cepat sambil interval auto-sync berjalan.
    - **Mitigasi**: Implementasi `syncStatus = 'syncing'` di frontend untuk mencegah request paralel.
- **Timer Skew**: Perbedaan waktu lokal browser dengan server.
    - **Mitigasi**: Backend selalu mengembalikan `sisa_waktu` yang otoritatif setiap kali sinkronisasi jawaban berhasil.
