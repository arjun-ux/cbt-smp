# Panduan & Rencana Refactoring Total: Sistem CBT SMP

Dokumen ini berisi pedoman standar, rencana strategis, dan daftar tugas untuk merapikan kode (Refactoring) aplikasi demi stabilitas data dan performa maksimal.

---

## 1. Pedoman Standar Refactoring

### Backend (Golang)
- **Slim Handlers**: Handler hanya bertugas menangani request dan response. Logika bisnis harus dipindah ke fungsi helper.
- **Database Transactions**: Gunakan transaksi (`db.Begin()`) untuk operasi yang melibatkan banyak tabel atau perhitungan penting (seperti Submit Ujian).
- **Standardized Response**: Gunakan format JSON yang seragam untuk sukses dan error.
- **Optimasi Query**: Hindari `Preload` massal dan loop database. Gunakan `Find` atau `Select` spesifik.

### Frontend (Vue.js)
- **Atomic Design**: Pecah file besar menjadi komponen-komponen kecil yang reusable.
- **Single Source of Truth**: Pastikan state data (seperti jawaban siswa) dikelola secara terpusat dan sinkron dengan server.
- **Error Handling UI**: Tampilkan pesan error yang informatif kepada pengguna saat API gagal.

---

## 2. Rencana Strategis (Refactor Plan)

### Fase 1: Backend Data Integrity (Prioritas Tinggi)
Fokus pada standarisasi alur data di server untuk menghilangkan masalah nilai tidak sinkron atau status gantung.
- Penyatuan logika sesi siswa.
- Implementasi Transaksi pada proses simpan jawaban dan perhitungan nilai.
- Pembersihan log debug.

### Fase 2: Frontend Restructuring
Fokus pada stabilitas UI dan efisiensi memori di HP siswa.
- Pemecahan `UjianPengerjaan.vue` menjadi sub-komponen.
- Perbaikan sistem auto-save agar lebih tangguh terhadap koneksi buruk.

### Fase 3: Monitoring & Finalisasi
Fokus pada akurasi laporan admin dan pembersihan akhir.
- Perbaikan dashboard monitoring.
- Finalisasi sistem ekspor nilai.

---

## 3. Daftar Tugas (Task List)

### [ ] Fase 1: Backend (Integritas Data)
- [ ] Buat helper `internal/handlers/utils.go` untuk standarisasi session & response.
- [ ] Refactor `ujian_siswa.go`: Terapkan helper dan DB Transactions pada `SubmitUjian`.
- [ ] Optimasi fungsi `HitungNilaiPG` (Gunakan bulk logic).
- [ ] Perbaiki logika `SyncJawaban` untuk menangani konkurensi (race condition).

### [ ] Fase 2: Frontend (Stabilitas UI)
- [ ] Buat folder komponen `siswa/components/`.
- [ ] Ekstrak Header ke `ExamHeader.vue`.
- [ ] Ekstrak Area Soal ke `QuestionContent.vue`.
- [ ] Ekstrak Navigasi ke `NavigationDock.vue`.
- [ ] Refactor state management di `UjianPengerjaan.vue` menggunakan Props/Events.

### [ ] Fase 3: Monitoring & Rekap
- [ ] Refactor `monitoring.go`: Perbaiki akurasi statistik real-time.
- [ ] Perbaiki alur ekspor nilai ke Excel/PDF agar sesuai dengan format sekolah.
- [ ] Hapus semua file sampah dan komentar kode yang tidak terpakai.

---
*Status: Dalam Pengerjaan (Fase 1)*
