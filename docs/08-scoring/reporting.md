# Reporting & Data Export

Sistem CBT menyediakan berbagai cara untuk mengekstrak data nilai guna keperluan administrasi sekolah (Rapor, Leger, dll).

## 1. Rekap Nilai Real-time
Dapat diakses saat ujian berlangsung atau segera setelah selesai.
- **Tujuan**: Melihat progres pengerjaan dan nilai sementara.
- **Format**: Tabel Interaktif di Browser.

## 2. Ekspor Data (CSV/Excel)
Fitur utama untuk integrasi dengan aplikasi lain (seperti ERapor).
- **Format**: Comma Separated Values (.csv).
- **Kolom Export**:
    - No
    - Nama Siswa
    - NISN
    - Kelas
    - Nilai PG
    - Nilai Essay
    - Total Nilai
    - Status (Selesai/Sedang Mengerjakan)
    - Waktu Selesai

## 3. Brankas Nilai (Archives)
Laporan dari data yang sudah diarsipkan.
- **Keunggulan**: Lebih cepat diakses karena tidak melakukan join ke tabel transaksi yang besar. Data bersifat statis (historical).
- **Filter**: Pencarian berdasarkan Tahun, Semester, atau Mata Pelajaran.

## 4. Analisis Butir Soal (PDF/Print)
Laporan mendalam per soal yang ditujukan untuk evaluasi bank soal oleh Guru.
- Menampilkan soal yang paling banyak dijawab salah (indikasi soal terlalu sulit atau kunci jawaban salah).
- Menampilkan grafik sebaran jawaban.

## 5. Alur Pembuatan Laporan (Workflow)
1. **Verifikasi**: Guru memastikan semua essay sudah dikoreksi.
2. **Snapshot**: Admin menekan tombol "Arsipkan".
3. **Download**: Guru/Admin mengunduh file CSV.
4. **Distribusi**: File CSV diolah di Excel atau diimpor ke sistem manajemen nilai sekolah.
