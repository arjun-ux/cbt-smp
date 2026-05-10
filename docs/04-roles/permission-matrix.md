# Permission Matrix

Tabel berikut merangkum hak akses untuk setiap aksi dalam sistem.

| Fitur / Modul | Aksi | Admin | Guru | Siswa |
| :--- | :--- | :---: | :---: | :---: |
| **User Management** | Create/Update/Delete Admin, Guru, Siswa | ✅ | ❌ | ❌ |
| **Master Data** | Manage Kelas, Ruang, Sesi, Mapel | ✅ | ❌ | ❌ |
| **Bank Soal** | Create/Update/Delete Bank Soal & Soal | ✅ | ✅ | ❌ |
| | Import Soal dari Word (.docx) | ✅ | ✅ | ❌ |
| **Jadwal Ujian** | Create/Update/Delete Jadwal | ✅ | ❌ | ❌ |
| | Lihat Daftar Jadwal Aktif | ✅ | ✅ | ✅ |
| | Refresh Token Ujian | ✅ | ❌ | ❌ |
| **Ujian (Siswa)** | Masuk Ujian (Token Validation) | ❌ | ❌ | ✅ |
| | Sync Jawaban & Auto-save | ❌ | ❌ | ✅ |
| | Submit Ujian | ❌ | ❌ | ✅ |
| **Monitoring** | Lihat Status Siswa Real-time | ✅ | ✅ | ❌ |
| | Unblock & Reset Sesi Siswa | ✅ | ✅ | ❌ |
| | Force Submit Siswa | ✅ | ✅ | ❌ |
| **Laporan** | Lihat Rekap Nilai & Analisis Soal | ✅ | ✅ | ❌ |
| | Download Laporan (Excel/PDF) | ✅ | ✅ | ❌ |
| | Arsip Hasil Ujian | ✅ | ❌ | ❌ |

## Keterangan
- **Admin**: Memiliki akses penuh ke seluruh endpoint API di bawah `/api/admin` dan `/api/guru`.
- **Guru**: Dibatasi pada pengelolaan konten (Bank Soal) dan pengawasan (Monitoring). Tidak bisa mengubah struktur user atau jadwal.
- **Siswa**: Hanya memiliki akses ke endpoint di bawah `/api/siswa`.
