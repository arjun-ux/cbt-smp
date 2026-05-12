# Scoring System Overview

Sistem CBT menggunakan pendekatan **Hybrid Scoring** yang menggabungkan penilaian otomatis (untuk Pilihan Ganda) dan penilaian manual (untuk Essay).

## 1. Arsitektur Penilaian

Sistem penilaian terbagi menjadi tiga layer utama:
1.  **Transactional Layer (`cbt_jawaban_siswas`)**: Menyimpan skor per butir soal.
2.  **Aggregation Layer (`cbt_peserta_ujians`)**: Menyimpan akumulasi nilai PG, Essay, dan Total untuk setiap siswa pada jadwal tertentu.
3.  **Archive Layer (`cbt_rekap_nilais`)**: Brankas nilai permanen yang menyimpan snapshot hasil ujian setelah proses koreksi selesai.

## 2. Komponen Nilai

| Jenis Soal | Metode Penilaian | Pemicu (Trigger) |
| :--- | :--- | :--- |
| **Pilihan Ganda (PG)** | Otomatis | Saat siswa klik "Selesai" atau "Force Submit" oleh pengawas. |
| **Essay** | Manual | Dinilai oleh Guru/Admin melalui dashboard koreksi. |
| **Status Koreksi** | Flagging | Ditandai dengan kolom `is_koreksi` di tabel peserta. |

## 3. Fitur Pelacakan Koreksi (New)
Untuk mempermudah administrasi, sistem kini mendukung pelacakan status koreksi essay secara visual:
- **Pending**: Ditandai dengan ikon pensil biru di Rekap Nilai (berarti ada jawaban essay yang belum diberi nilai).
- **Completed**: Ditandai dengan ikon centang hijau setelah semua butir essay pada siswa tersebut telah memiliki skor.

## 4. Alur Lifecycle Nilai
1.  **Ujian Berlangsung**: Skor di `cbt_jawaban_siswas` masih 0.
2.  **Ujian Selesai**: Fungsi `HitungNilaiPG` dijalankan. `nilai_pg` diisi, `nilai_essay` masih 0.
3.  **Masa Koreksi**: Guru mengisi skor essay. Kolom `nilai_essay` diupdate, dan `total_nilai` dihitung ulang secara otomatis (`nilai_pg + nilai_essay`). Kolom `is_koreksi` diset menjadi `true` jika penilaian selesai.
4.  **Pengarsipan**: Admin mengunci nilai. Data disalin ke `CBTRekapNilai`. Status jadwal berubah menjadi `Diarsipkan`.

## 4. Keamanan & Integritas Nilai
- **Immutable Archive**: Setelah data dipindahkan ke Brankas Nilai (Archive), data di tabel transaksi (`peserta_ujian`) tidak lagi digunakan untuk pelaporan utama, mencegah manipulasi pada data aktif.
- **Role Restriction**: Hanya Admin yang memiliki otoritas untuk melakukan "Archive Nilai". Guru hanya dibatasi pada pemberian skor essay untuk bank soal milik mereka.
