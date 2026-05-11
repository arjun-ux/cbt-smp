# Improvement Roadmap (Features)

Dokumen ini merangkum rencana pengembangan fitur baru untuk membawa sistem CBT ini ke level *Enterprise Ready*.

## Tahap 1: Core Hardening (Q2 2026)
*Fokus: Keamanan dan Reliabilitas Dasar*

- [x] **Clipboard & Right-click Lock**: ✅ COMPLETED (May 2026) - Proteksi dasar konten soal.
- [ ] **Administrative Documents**: Fitur cetak Kartu Ujian, Daftar Hadir, dan Berita Acara (PDF).
- [ ] **Global Application Settings**: UI untuk mengatur Nama Sekolah, Logo, dan Tahun Ajaran.
- [ ] **Auto-Backup Database**: Skrip backup otomatis `cbt.db` setiap jam ke storage eksternal.
- [ ] **Enhanced Logs**: Logging setiap klik navigasi siswa untuk audit jejak digital.

## Tahap 2: Advanced Proctoring (Q3 2026)
*Fokus: Integritas Ujian*

- [ ] **Face Recognition Login**: Validasi wajah saat awal ujian.
- [ ] **SEB Support**: Integrasi dengan Safe Exam Browser untuk ujian di Lab Sekolah.
- [ ] **IP Whitelisting**: Membatasi akses ujian hanya dari jaringan WiFi sekolah.
- [ ] **Activity Heatmap**: Analisis waktu pengerjaan per butir soal oleh siswa.

## Tahap 3: Intelligence & Scale (Q4 2026)
*Fokus: Analytics dan Performa Skala Besar*

- [ ] **AI-Powered Item Analysis**: Rekomendasi perbaikan soal secara otomatis berdasarkan data statistik.
- [ ] **Multi-Database Support**: Migrasi dari SQLite ke PostgreSQL untuk beban kerja >2000 siswa.
- [ ] **Automated PDF Report**: Generate rapor hasil ujian dalam format PDF yang siap cetak.
- [ ] **Question Import (QTI Standard)**: Mendukung format standar industri agar bisa bertukar soal dengan LMS lain.

## Tahap 4: Cloud & Multi-Tenancy (2027+)
*Fokus: Komersialisasi dan Ekspansi*

- [ ] **Multi-Tenant Architecture**: Satu aplikasi untuk banyak sekolah dengan isolasi data.
- [ ] **S3 Media Storage**: Skalabilitas aset gambar dan video.
- [ ] **Mobile App App**: Aplikasi Android/iOS dengan fitur *Locked-down mode*.
