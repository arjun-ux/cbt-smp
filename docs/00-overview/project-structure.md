# Project Structure

## Overview
Proyek ini adalah sistem Computer Based Test (CBT) yang dibangun menggunakan Golang untuk Backend dan Vue 3 untuk Frontend. Frontend disajikan sebagai file statis yang di-embed ke dalam binary Golang.

## Backend (Golang)
Struktur direktori backend mengikuti pola yang cukup modular dengan pemisahan antara entry point, logic bisnis (handlers), dan data layer (models/database).

- **`cmd/server/main.go`**: Entry point utama aplikasi. Bertanggung jawab untuk inisialisasi database, middleware, route, dan menjalankan server Fiber.
- **`internal/database/`**:
    - `database.go`: Konfigurasi dan koneksi ke SQLite menggunakan GORM.
    - `seeder.go`: Data awal untuk sistem (admin default, dll).
    - `seeder_stress.go`: Data untuk pengujian beban (stress test).
- **`internal/handlers/`**: Berisi logika penanganan HTTP request.
    - `auth.go`: Autentikasi dan JWT.
    - `admin.go`, `guru.go`, `siswa.go`: Manajemen user berdasarkan role.
    - `ujian_siswa.go`: Logika inti pengerjaan ujian oleh siswa.
    - `monitoring.go`: Fitur pengawasan ujian real-time.
- **`internal/middleware/`**:
    - `auth.go`: Proteksi rute dengan JWT.
    - `role.go` (dalam handlers/middleware): Validasi role user.
- **`internal/models/`**:
    - `models.go`: Definisi skema database menggunakan GORM struct tags.
- **`internal/routes/`**:
    - `routes.go`: Definisi seluruh endpoint API.
- **`pkg/utils/`**: Helper fungsi umum (hash password, format data, dll).

## Frontend (Vue 3)
Struktur direktori frontend menggunakan standar Vite + Vue 3.

- **`frontend/src/main.js`**: Entry point frontend.
- **`frontend/src/App.vue`**: Komponen root.
- **`frontend/src/router/`**: Konfigurasi Vue Router untuk SPA.
- **`frontend/src/store/`**: State management menggunakan Pinia.
    - `auth.js`: Manajemen sesi login.
    - `exam.js`: Logika inti pengerjaan ujian, sinkronisasi, dan timer.
- **`frontend/src/views/`**:
    - `Login.vue`: Halaman login utama.
    - `admin/`: Halaman untuk dashboard dan manajemen data oleh Admin/Guru.
    - `siswa/`: Halaman pengerjaan ujian modular (terdiri dari komponen `Timer`, `Navigator`, dll).
- **`frontend/src/components/`**: Komponen UI yang dapat digunakan kembali.

## Root & Assets
- **`uploads/`**: Tempat penyimpanan file yang diunggah (gambar soal, dll).
- **`cbt.db`**: Database SQLite yang digunakan secara default.
- **`.env`**: Konfigurasi environment (PORT, SECRET, dll).
