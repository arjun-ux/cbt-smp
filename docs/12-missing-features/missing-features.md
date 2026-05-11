# Missing Features Analysis

Berdasarkan perbandingan dengan standar sistem CBT modern (seperti Moodle, CandyCBT, atau Exam.net), berikut adalah daftar fitur yang belum ada namun sangat krusial untuk meningkatkan kualitas dan kredibilitas sistem.

## 1. Security & Anti-Cheat (Advanced)
- **Safe Exam Browser (SEB) Integration**: Memaksa siswa menggunakan browser khusus yang mengunci seluruh OS agar tidak bisa membuka aplikasi lain.
- **Face Verification (AI Proctoring)**: Menggunakan kamera untuk memastikan orang yang mengerjakan adalah siswa yang bersangkutan dan mendeteksi jika ada orang lain di sekitar.
- **Clipboard Blocking**: Mencegah aksi Copy-Paste dari luar ke dalam sistem ujian.
- **Device Fingerprinting**: Mencatat ID perangkat unik untuk mencegah satu akun login di dua browser berbeda secara bergantian.

## 2. Real-time Monitoring (Optimization)
- **Remote Action Notification**: Mengirim pesan instruksi langsung ke layar siswa dari dashboard pengawas.

## 3. Advanced Question Bank
- **Question Versioning**: Menyimpan riwayat perubahan soal.
- **Taxonomy Mapping**: Pengelompokan soal berdasarkan tingkat kognitif (C1-C6) atau standar kompetensi.
- **Centralized Media Manager**: Galeri untuk mengelola gambar/video agar tidak terjadi duplikasi file yang sama di banyak soal.

## 4. Analytics & Insight
- **Discriminatory Index**: Analisis otomatis untuk mendeteksi soal yang "terlalu sulit bagi siswa pintar" (indikasi soal bermasalah).
- **Student Performance Trend**: Grafik progres nilai siswa antar ujian untuk melihat perkembangan belajar.
- **Class Heatmap**: Visualisasi materi mana yang paling tidak dikuasai oleh satu kelas.

## 5. Reliability & Scalability
- **Redis Session Storage**: Memindahkan sesi aktif ke Redis agar server bisa di-restart tanpa memutus koneksi siswa yang sedang ujian.
- **PostgreSQL/MySQL Support**: Dukungan database enterprise untuk menangani ribuan siswa secara simultan (SQLite memiliki limitasi konkurensi tulis).
- **S3 Storage Integration**: Menyimpan file gambar soal di cloud storage agar aplikasi bisa berjalan secara serverless/multi-instance.

## 6. Communication & Administration
- **Notification System**: Mengirimkan token ujian atau pengumuman melalui WhatsApp/Email secara otomatis.
- **Administrative Documents (Cetak PDF)**: Fitur ekspor dokumen standar ujian seperti Kartu Ujian Siswa, Daftar Hadir, Berita Acara, dan Jadwal Ujian.
- **Application Settings**: Pengaturan identitas sekolah (Nama, Logo, Alamat) dan konfigurasi sistem secara terpusat melalui UI Admin.
