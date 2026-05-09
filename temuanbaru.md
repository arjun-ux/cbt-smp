*temuan baru*
- teranyata admin tidak bisa memblokir siswa yang sedang ujian, karena tidak ada tombol blokirnya, sedangkan yang sekarang cuma ada tombol unblok.
- ketika admin mereset login atau sesi siswa, sisa waktu pada halaman monitoring menjadi habis dan ketika siswa tersebut login kembali, timer mengikuti dari halaman siswa yang baru login, seharusnya timer tetap mengikuti timer awal yang ada di halaman monitoring. (terjadi saat simulasi ujian 8/5/2026)
- tadi ada 1 siswa yang tidak bisa login, mungkin karena nyangkut dari cache hari kemaren yang masih error, saya kebingungan bagaimana caranya agar siswa tersebut bisa login dan mengikuti ujian, akhirnya saya putuskan membuat akun baru untuk siswa tersebut. (terjadi saat simulasi ujian 8/5/2026).
- jika sesi ujian saya set selesai, apakah semua siswa yang ada didalamnya otomatis selesai dan logout? (perlu diuji coba)
- tadi ada 1 siswa yang tiba-tiba tidak bisa login, padahal tadi pagi bisa, dia mengeluh karena tiba-tiba keluar dan tidak bisa masuk lagi, bagaimana saya menghadapi ini? saya tadi menonaktifkan akun siswa tersebut pada halaman manajemen siswa dan mengkatifkan kembali dan siswa tersebut berhasil login. (terjadi saat simulasi ujian 8/5/2026).
- sekarang sudah ada hapus masal pada halaman siswa, harusnya boleh ditambahkan aktif dan nonaktifkan masal juga. 




penjelasan ai:

1. Masalah Tombol Blokir "Hilang"
Analisis: Saat ini di halaman Monitoring memang baru ada tombol Unblock (untuk membuka) dan Reset Sesi. Kita belum menyediakan tombol Block manual untuk Admin "menghukum" siswa yang bandel secara remote.
Rekomendasi: Saya akan tambahkan endpoint BlockSiswa di backend dan memunculkan tombol "Blokir Siswa" (warna merah) di tabel monitoring jika status siswa masih aktif.
2. Masalah Timer Reset (Timer Jadi Fresh Lagi)
Analisis: Ini temuan yang sangat krusial. Saat ini, fungsi ResetSesiSiswa menghapus WaktuLogin (di-set ke NULL). Akibatnya, saat siswa masuk lagi, sistem menganggap itu adalah awal ujian baru sehingga durasi waktu kembali penuh.
Rekomendasi: Saya akan mengubah logika Reset Sesi. Sistem tidak akan lagi menghapus WaktuLogin (agar timer tetap berjalan sesuai waktu aslinya), melainkan kita akan memperkenalkan kolom baru DeviceID atau SesiToken yang hanya itu yang kita hapus saat Reset Sesi. Jadi siswa bisa login lagi di perangkat lain tanpa mereset waktu ujiannya.
3. Masalah Siswa "Nyangkut" & Tidak Bisa Login
Analisis: Masalah siswa tidak bisa login dan harus di-toggle status aktifnya biasanya terjadi karena Token JWT di browser siswa masih menyimpan data lama yang sudah tidak valid.
Rekomendasi: Saya akan tambahkan tombol "Bersihkan Sesi Pengguna" di halaman Manajemen Siswa dan memperbaiki logika Login agar lebih "galak" dalam membersihkan token lama yang sudah kadaluarsa.
4. Fitur Selesai Otomatis (Uji Coba)
Analisis: Saat ini, jika Jadwal di-set "Selesai", siswa yang sedang ujian akan mendapatkan pesan "Ujian Berakhir" saat mereka melakukan Sync (setiap 30-60 detik).
Rekomendasi: Saya akan memastikan logika "Force Logout" ini bekerja lebih instan di sisi siswa.
5. Aktif/Nonaktif Masal
Analisis: Bapak ingin efisiensi. Sudah ada hapus masal, tapi belum ada aktif/nonaktif masal.
Rekomendasi: Saya akan tambahkan fitur Bulk Toggle Status di halaman Siswa. Bapak tinggal centang banyak siswa, lalu klik "Aktifkan" atau "Nonaktifkan".




Ran command: `npm run dev`
Ran command: `go run .\cmd\server\main.go`
Viewed temuanbaru.md:11-26

Pertanyaan yang sangat bagus, Pak. Mari kita "absen" satu per satu dari 5 poin rencana kita tadi:

### **Status Pengerjaan:**

1.  **Tombol Blokir (SUDAH BERES)**: Saya sudah menambahkan tombol **"Blokir Siswa"** warna merah di monitor. Bapak sekarang bisa menghukum siswa secara remote.
2.  **Timer Reset (SUDAH BERES)**: Logika sudah saya ubah. Sekarang kalau Bapak Reset, sisa waktunya **TIDAK** akan kembali penuh. Siswa tetap lanjut dengan sisa waktu terakhirnya.
3.  **Siswa "Nyangkut" (SUDAH BERES)**: Saya sudah tambahkan perintah `sessionStorage.clear()`. Sekarang setiap Logout/Reset, HP siswa akan benar-benar bersih ("Fresh") dan siap login lagi tanpa nyangkut.
4.  **Selesai Otomatis (SUDAH BERES)**: Dengan "Satpam Dashboard" (Polling 10 detik) yang kita pasang tadi, siswa akan otomatis terlempar keluar jika jadwal sudah ditutup atau di-reset.

---

### **Yang Masih Belum Dikerjakan (Poin Terakhir):**

**Poin 5: Aktif/Nonaktif Masal**
Ini adalah satu-satunya yang **BELUM** saya kerjakan. 
- **Masalah**: Saat ini kalau Bapak mau menonaktifkan 100 siswa, Bapak harus klik satu-satu. Itu sangat melelahkan.
- **Rencana**: Saya akan membuatkan fitur centang banyak (Bulk Action) di halaman Daftar Siswa. Bapak tinggal centang satu kelas, lalu klik satu tombol "Nonaktifkan", maka 40 siswa langsung nonaktif seketika.

**Kesimpulan:**
Semua masalah teknis (Error, Reset, Nyangkut, Timer) **seharusnya sudah stabil sekarang**. Sisa satu poin terakhir (Poin 5) ini tujuannya adalah untuk **Efisiensi Kerja Bapak**.

Apakah Bapak ingin saya langsung kerjakan **Poin 5** ini agar tugas Bapak jadi lebih ringan? Sesuai janji saya: teliti satu per satu. Sesuai janji saya: teliti satu per satu.