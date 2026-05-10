# Entity Relationship Diagram (ERD)

Diagram di bawah menunjukkan hubungan antar tabel dalam sistem CBT.

```mermaid
erDiagram
    USERS ||--o| MASTER_SISWA : "has profile"
    USERS ||--o| MASTER_GURU : "has profile"
    
    MASTER_KELAS ||--o{ MASTER_SISWA : "contains"
    MASTER_RUANG ||--o{ MASTER_SISWA : "assigned to"
    MASTER_SESI ||--o{ MASTER_SISWA : "assigned to"
    
    MASTER_MAPEL ||--o{ CBT_BANK_SOAL : "categorizes"
    MASTER_GURU ||--o{ CBT_BANK_SOAL : "creates"
    
    CBT_BANK_SOAL ||--o{ CBT_SOAL : "contains"
    CBT_BANK_SOAL ||--o{ CBT_JADWAL_UJIAN : "scheduled as"
    
    MASTER_RUANG ||--o{ CBT_JADWAL_UJIAN : "venue for"
    MASTER_SESI ||--o{ CBT_JADWAL_UJIAN : "time slot for"
    MASTER_GURU ||--o{ CBT_JADWAL_UJIAN : "supervises"
    
    CBT_JADWAL_UJIAN ||--o{ CBT_PESERTA_UJIAN : "has"
    MASTER_SISWA ||--o{ CBT_PESERTA_UJIAN : "participates in"
    
    CBT_PESERTA_UJIAN ||--o{ CBT_JAWABAN_SISWA : "provides"
    CBT_SOAL ||--o{ CBT_JAWABAN_SISWA : "is answered in"
    
    CBT_PESERTA_UJIAN ||--o{ LOG_UJIAN : "generates"
```

## Relational Flow
1.  **User Management**: Data User adalah root dari autentikasi. User dihubungkan ke `master_siswas` atau `master_gurus` melalui `user_id`.
2.  **Master Data Binding**: Siswa diikat ke Kelas, Ruang, dan Sesi untuk mempermudah distribusi jadwal.
3.  **Content Creation**: Guru membuat `BankSoal`, yang kemudian diisi dengan banyak `Soal`.
4.  **Scheduling**: `JadwalUjian` mengikat `BankSoal` ke waktu dan lokasi tertentu.
5.  **Examination**: Saat siswa mulai ujian, record `PesertaUjian` dibuat sebagai "session" ujian. Setiap jawaban yang dikirim disimpan di `JawabanSiswa` yang merujuk pada `PesertaUjian` dan `Soal`.
6.  **Monitoring**: Aktivitas mencurigakan dicatat di `LogUjian` yang merujuk pada sesi `PesertaUjian`.
