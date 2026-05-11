# System Flow

Dokumen ini menjelaskan alur kerja sistem CBT secara step-by-step untuk proses-proses krusial.

## 1. Login Flow (Umum & Siswa)
Proses masuk ke sistem untuk semua role.

```mermaid
sequenceDiagram
    participant User
    participant Frontend
    participant Backend
    participant DB

    User->>Frontend: Masukkan Username/NISN & Password
    Frontend->>Backend: POST /api/auth/login
    Backend->>DB: Query User berdasarkan Username/NISN
    DB-->>Backend: Data User (Hashed Pass)
    Backend->>Backend: bcrypt.CompareHash()
    Backend->>Backend: Generate JWT (Claims: id, role)
    Backend-->>Frontend: 200 OK (Token + User Data)
    Frontend->>Frontend: Simpan Token di LocalStorage
    Frontend->>User: Redirect ke Dashboard sesuai Role
```

## 2. Ujian Flow (Persiapan & Pengerjaan)
Alur siswa mulai masuk ke ruang ujian.

```mermaid
sequenceDiagram
    participant Siswa
    participant Frontend
    participant Backend
    participant DB

    Siswa->>Frontend: Klik Ujian & Masukkan Token Ujian
    Frontend->>Backend: POST /api/siswa/validate
    Backend->>DB: Cek Jadwal, Sesi, & Status Peserta
    DB-->>Backend: Valid/Invalid
    Backend-->>Frontend: 200 OK (peserta_id, sisa_waktu)
    
    Frontend->>Backend: GET /api/siswa/soal/:jadwalId
    Backend->>DB: Fetch Soal & Jawaban Eksis
    DB-->>Backend: List Soal (Acak jika diatur)
    Backend-->>Frontend: 200 OK (Data Soal + Jawaban)
    Frontend->>Siswa: Tampilkan Layar Pengerjaan
```

## 3. Submit Flow (Selesai Ujian)
Proses mengakhiri ujian dan perhitungan nilai.

```mermaid
sequenceDiagram
    participant Siswa
    participant Frontend
    participant Backend
    participant DB

    Siswa->>Frontend: Klik tombol "Selesai"
    Frontend->>Backend: POST /api/siswa/submit/:pesertaId
    Backend->>Backend: HitungNilaiPG (Compare dengan Kunci)
    Backend->>DB: Update Status=Selesai, NilaiPG, WaktuSelesai
    DB-->>Backend: Success
    Backend-->>Frontend: 200 OK (Hasil Nilai Sementara)
    Frontend->>Siswa: Tampilkan Notifikasi & Redirect Dashboard
```

## 4. Scoring Flow (Otomatis & Manual)
- **PG (Pilihan Ganda)**: Dihitung otomatis oleh server saat submit dengan membandingkan `jawaban_siswa` dengan `kunci_jawaban` pada tabel `cbt_soals`.
- **Essay**: Disimpan sebagai teks. Guru melakukan penilaian manual melalui halaman koreksi di Admin Panel, yang kemudian akan mengupdate kolom `nilai_essay` di tabel `cbt_peserta_ujians`.

## 5. Monitoring & Anti-Cheat Flow
Fitur pengawasan real-time oleh Admin/Guru.

```mermaid
sequenceDiagram
    participant Siswa
    participant Browser
    participant Backend
    participant Admin

    Siswa->>Browser: Pindah Tab / Minimize
    Browser->>Backend: POST /api/siswa/log (Keterangan: Pindah Tab)
    Backend->>Backend: Auto Update is_terblokir = true
    Admin->>Backend: GET /api/admin/monitor/:jadwalId (Polling)
    Backend-->>Admin: Data Peserta (Status: Terblokir)
    Admin->>Backend: POST /api/monitor/unblock/:pesertaId
    Backend->>Backend: Reset is_terblokir = false
    Admin-->>Siswa: Akses Terbuka Kembali
```

## 6. Data Integrity & Session Fencing (Anti-Pollution)
Alur untuk menjamin kebersihan data saat terjadi Reset Total oleh Admin.

```mermaid
sequenceDiagram
    participant Admin
    participant DB
    participant Siswa
    participant LocalStorage

    Admin->>DB: Klik "Reset Total" (AttemptID + 1)
    DB-->>Admin: Success
    
    Note over Siswa: Background Sync Berjalan
    Siswa->>DB: POST /api/siswa/sync (Old AttemptID)
    DB->>DB: Validasi AttemptID
    DB-->>Siswa: 200 OK (Status: OUTDATED_SESSION)
    
    Siswa->>LocalStorage: removeItem() & clearStore()
    Siswa->>Siswa: Force Page Reload
    Siswa->>DB: GET /api/siswa/soal (New AttemptID)
    DB-->>Siswa: Data Fresh (Kosong)
```
