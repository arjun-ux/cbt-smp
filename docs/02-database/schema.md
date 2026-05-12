# Database Schema

Dokumen ini mendefinisikan skema tabel database SQLite yang digunakan dalam sistem CBT. Database dikelola menggunakan GORM AutoMigrate.

## 1. User & Authentication
### Tabel: `users`
| Kolom | Tipe | Deskripsi |
| :--- | :--- | :--- |
| `id` | uint (PK) | Primary Key |
| `username` | string (Unique) | Username untuk login atau NISN untuk siswa |
| `password` | string | Password yang telah di-hash (bcrypt) |
| `role` | string | Peran user: `admin`, `guru`, `siswa` |
| `is_active` | boolean | Status aktif akun |
| `created_at` | datetime | Waktu pembuatan |
| `updated_at` | datetime | Waktu pembaruan terakhir |

## 2. Master Data
### Tabel: `master_kelas`
| Kolom | Tipe | Deskripsi |
| :--- | :--- | :--- |
| `id` | uint (PK) | Primary Key |
| `tingkat` | string (2) | Tingkat kelas (7, 8, 9) |
| `nama_kelas` | string (20) | Nama kelas (misal: 7A, 8B) |

### Tabel: `master_ruangs`
| Kolom | Tipe | Deskripsi |
| :--- | :--- | :--- |
| `id` | uint (PK) | Primary Key |
| `kode_ruang` | string (Unique) | Kode identifikasi ruang |
| `nama_ruang` | string | Nama lengkap ruang |

### Tabel: `master_sesis`
| Kolom | Tipe | Deskripsi |
| :--- | :--- | :--- |
| `id` | uint (PK) | Primary Key |
| `nama_sesi` | string | Nama sesi (misal: Sesi 1) |
| `waktu_mulai` | string | Jam mulai (format HH:mm) |
| `waktu_selesai` | string | Jam selesai (format HH:mm) |

### Tabel: `master_siswas`
| Kolom | Tipe | Deskripsi |
| :--- | :--- | :--- |
| `id` | uint (PK) | Primary Key |
| `user_id` | uint (FK) | Relasi ke `users.id` |
| `nisn` | string (Unique) | Nomor Induk Siswa Nasional |
| `nama_lengkap` | string | Nama lengkap siswa |
| `kelas_id` | uint (FK) | Relasi ke `master_kelas.id` |
| `ruang_id` | uint (FK) | Relasi ke `master_ruangs.id` |
| `sesi_id` | uint (FK) | Relasi ke `master_sesis.id` |

### Tabel: `master_gurus`
| Kolom | Tipe | Deskripsi |
| :--- | :--- | :--- |
| `id` | uint (PK) | Primary Key |
| `user_id` | uint (FK) | Relasi ke `users.id` |
| `nip` | string (Unique) | Nomor Induk Pegawai |
| `nama_guru` | string | Nama lengkap guru |

### Tabel: `master_mapels`
| Kolom | Tipe | Deskripsi |
| :--- | :--- | :--- |
| `id` | uint (PK) | Primary Key |
| `kode_mapel` | string (Unique) | Kode mata pelajaran |
| `nama_mapel` | string | Nama mata pelajaran |

## 3. CBT Content
### Tabel: `cbt_bank_soals`
| Kolom | Tipe | Deskripsi |
| :--- | :--- | :--- |
| `id` | uint (PK) | Primary Key |
| `guru_id` | uint (FK) | Pemilik/Pembuat bank soal |
| `mapel_id` | uint (FK) | Relasi ke `master_mapels.id` |
| `tingkat_kelas` | string | Filter tingkat kelas (7, 8, 9) |
| `judul_bank_soal` | string | Judul koleksi soal |
| `default_bobot_pg` | float64 | Bobot nilai default untuk PG |
| `default_bobot_essay` | float64 | Bobot nilai default untuk Essay |
| `status` | string | Status: `Draft`, `Aktif` |

### Tabel: `cbt_soals`
| Kolom | Tipe | Deskripsi |
| :--- | :--- | :--- |
| `id` | uint (PK) | Primary Key |
| `bank_soal_id` | uint (FK) | Relasi ke `cbt_bank_soals.id` |
| `jenis_soal` | string | `PG` atau `ESSAY` |
| `pertanyaan` | text | Konten pertanyaan (HTML supported) |
| `opsi_a` s/d `opsi_d` | text | Pilihan jawaban untuk PG |
| `kunci_jawaban` | text | Jawaban benar (A/B/C/D atau teks essay) |
| `bobot_nilai` | float64 | Nilai jika jawaban benar |

## 4. Exam Execution
### Tabel: `cbt_jadwal_ujians`
| Kolom | Tipe | Deskripsi |
| :--- | :--- | :--- |
| `id` | uint (PK) | Primary Key |
| `bank_soal_id` | uint (FK) | Bank soal yang digunakan |
| `tanggal_ujian` | date | Tanggal pelaksanaan |
| `waktu_mulai` | time | Jam mulai akses |
| `durasi_menit` | int | Lama pengerjaan dalam menit |
| `ruang_id` | uint (FK) | Filter lokasi (opsional) |
| `sesi_id` | uint (FK) | Filter sesi (opsional) |
| `pengawas_id` | uint (FK) | Guru yang bertugas mengawas |
| `acak_soal` | boolean | Flag untuk mengacak urutan soal |
| `acak_jawaban` | boolean | Flag untuk mengacak opsi jawaban |
| `token_ujian` | string | Token unik untuk masuk ujian |
| `status` | string | `Belum Mulai`, `Berlangsung`, `Selesai` |

### Tabel: `cbt_peserta_ujians`
| Kolom | Tipe | Deskripsi |
| :--- | :--- | :--- |
| `id` | uint (PK) | Primary Key |
| `jadwal_id` | uint (FK) | Relasi ke `cbt_jadwal_ujians.id` |
| `siswa_id` | uint (FK) | Relasi ke `master_siswas.id` |
| `waktu_login` | datetime | Waktu mulai pengerjaan |
| `waktu_selesai_ujian` | datetime | Waktu submit pengerjaan |
| `sisa_waktu_detik` | int | Durasi tersisa (untuk autosave) |
| `status_ujian` | string | `Sedang Mengerjakan`, `Selesai` |
| `nilai_pg` | float64 | Total nilai pilihan ganda |
| `nilai_essay` | float64 | Total nilai essay (manual/otomatis) |
| `total_nilai` | float64 | Akumulasi nilai |
| `is_koreksi` | boolean | Flag apakah essay sudah dikoreksi semua |
| `is_terblokir` | boolean | Status blokir jika curang |
| `attempt_id` | int | ID sesi pengerjaan (untuk reset data) |

### Tabel: `cbt_jawaban_siswas`
| Kolom | Tipe | Deskripsi |
| :--- | :--- | :--- |
| `id` | uint (PK) | Primary Key |
| `peserta_ujian_id` | uint (FK) | Relasi ke `cbt_peserta_ujians.id` |
| `soal_id` | uint (FK) | Relasi ke `cbt_soals.id` |
| `jawaban_siswa` | text | Jawaban yang dipilih/diinput |
| `ragu_ragu` | boolean | Status ragu-ragu siswa |
| `skor` | float64 | Nilai yang didapat untuk soal ini |

## 5. Monitoring & Logs
### Tabel: `log_sistems`
Mencatat aktivitas login dan perubahan data krusial.

### Tabel: `log_ujians`
Mencatat aktivitas siswa selama ujian (login, pindah tab, submit).

### Tabel: `cbt_rekap_nilais`
Snapshot hasil ujian untuk arsip permanen setelah jadwal ditutup.
