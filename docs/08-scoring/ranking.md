# Ranking & Analysis

Meskipun sistem tidak menyimpan "Rank" secara permanen di database, urutan peringkat dihasilkan secara dinamis untuk keperluan pelaporan.

## 1. Mekanisme Perankingan
Perankingan dilakukan di sisi **Frontend (Vue.js)** dan **Backend (Query Sorting)**:

- **Dashboard Rekap**: Data diambil menggunakan `ORDER BY total_nilai DESC`.
- **Logic Tie-breaker**: Jika ada dua siswa dengan `total_nilai` yang sama, sistem menggunakan `waktu_selesai_ujian` sebagai pemutus (Siswa yang selesai lebih awal berada di posisi lebih tinggi).

## 2. Analisis Butir Soal
Sistem menyediakan fitur Analisis Butir Soal (`GetAnalisisSoal`) untuk melihat kualitas soal berdasarkan statistik pengerjaan siswa:

| Indikator | Rumus / Logika |
| :--- | :--- |
| **Tingkat Kesukaran** | `(Jumlah Benar / Total Peserta) * 100` |
| **Kategori Sukar** | Persentase < 30% |
| **Kategori Sedang** | Persentase 30% - 70% |
| **Kategori Mudah** | Persentase > 70% |
| **Sebaran Jawaban** | Menghitung frekuensi setiap opsi (A, B, C, D) yang dipilih siswa untuk mendeteksi *distractor* (pengecoh) yang tidak efektif.

## 3. Statistik Ringkasan
Dashboard Admin/Guru menampilkan ringkasan statistik per jadwal ujian:
- **Rata-rata (Average)**: Mean dari `total_nilai` seluruh peserta yang sudah selesai.
- **Nilai Tertinggi (Max)**: Score tertinggi dalam satu jadwal.
- **Nilai Terendah (Min)**: Score terendah dalam satu jadwal.
- **Ketuntasan**: (Opsional/Planned) Persentase siswa yang mencapai skor di atas KKM.

## 4. Ranking per Kelas
Sistem memungkinkan filter per kelas pada halaman rekap. Saat filter diterapkan, perankingan otomatis menyesuaikan hanya untuk populasi kelas tersebut. Hal ini berguna bagi Wali Kelas untuk mengambil nilai rapor.
