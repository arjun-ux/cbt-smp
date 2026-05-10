# Exam API Reference

Endpoint khusus untuk alur pengerjaan ujian dan pengawasan.

## 1. Alur Siswa (Student Side)

### 1.1 Validate Token
Memulai sesi ujian dengan memvalidasi token dari pengawas.

- **Route**: `POST /api/siswa/validate?jadwalId=1`
- **Method**: `POST`
- **Auth**: `siswa`
- **Request Body**: `{ "token": "ABCDE" }`
- **Database Effect**:
    - Mencari/Membuat record di `cbt_peserta_ujians`.
    - Update `waktu_login = NOW`.
- **Response**: Mengembalikan data sesi ujian dan `sisa_waktu`.

### 1.2 Get Soal Ujian
Mengambil daftar soal untuk jadwal tertentu.

- **Route**: `GET /api/siswa/soal/:jadwalId`
- **Method**: `GET`
- **Response**: List soal (pertanyaan + opsi) tanpa kunci jawaban. Jika `acak_soal` aktif, urutan dibalik secara random oleh database.

### 1.3 Sync Jawaban (Autosave)
Menyimpan jawaban siswa secara berkala.

- **Route**: `POST /api/siswa/sync`
- **Method**: `POST`
- **Request Body**:
    ```json
    {
      "peserta_ujian_id": 10,
      "items": [
        { "soal_id": 1, "jawaban_teks": "A", "ragu_ragu": false }
      ],
      "sisa_waktu": 3500
    }
    ```
- **Database Effect**: Update tabel `cbt_jawaban_siswas` dan update `sisa_waktu_detik` di `cbt_peserta_ujians`.
- **Security Validation**:
    - **Ownership Check**: Server memverifikasi bahwa `peserta_ujian_id` milik siswa yang sedang login. Mengembalikan `403 Forbidden` jika tidak cocok.
    - **Timer Hardening**: Server menolak penambahan sisa waktu yang lebih besar dari nilai sebelumnya di database untuk mencegah manipulasi timer.

### 1.4 Submit Ujian
Finalisasi pengerjaan ujian.

- **Route**: `POST /api/siswa/submit/:pesertaId`
- **Method**: `POST`
- **Database Effect**:
    - Update `status_ujian = 'Selesai'`.
    - Menghitung nilai PG otomatis via function `HitungNilaiPG`.
    - Simpan log selesai.
    - **Ownership Check**: Verifikasi kepemilikan sesi sebelum finalisasi nilai.

---

## 2. Alur Pengawas (Proctor Side)

### 2.1 Get Monitor Status
Melihat progres seluruh siswa secara real-time.

- **Route**: `GET /api/[role]/monitor/:jadwalId`
- **Access**: `admin`, `guru` (sebagai pengawas)
- **Response**: Statistik progres (persentase dijawab) dan status pelanggaran (pindah tab).

### 2.2 Unblock Siswa
Membuka blokir otomatis akibat pindah tab.

- **Route**: `POST /api/[role]/monitor/unblock/:pesertaId`
- **Database Effect**: Set `is_terblokir = false` dan `waktu_login = NULL`. Siswa dapat login kembali tanpa kehilangan progres.

### 2.3 Reset Sesi Login (Crash Recovery)
Menghapus jejak login agar siswa bisa masuk kembali di perangkat/browser lain.

- **Route**: `POST /api/[role]/monitor/reset-sesi/:pesertaId`
- **Database Effect**: Set `waktu_login = NULL` dan `is_terblokir = false`.
- **Note**: Progres jawaban dan sisa waktu **TIDAK** berubah. Gunakan ini jika PC siswa mati/hang.

### 2.4 Reset Total Ujian
Menghapus seluruh progres dan memulai ujian dari nol.

- **Route**: `POST /api/[role]/monitor/reset-ujian/:pesertaId`
- **Database Effect**:
    - Menghapus seluruh jawaban di `cbt_jawaban_siswas`.
    - Mengembalikan `sisa_waktu_detik` ke durasi default jadwal.
    - Set `status_ujian = 'Belum Mengerjakan'`.
    - Set `nilai_pg = 0`, `nilai_essay = 0`.

### 2.5 Force Submit
Menghentikan paksa ujian siswa.

- **Route**: `POST /api/[role]/monitor/force-submit/:pesertaId`
- **Side Effect**: Menghitung nilai terakhir dan mengunci status ujian menjadi `Selesai`.
