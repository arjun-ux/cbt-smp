# Security Audit Report

Laporan ini merangkum hasil audit keamanan menyeluruh terhadap sistem CBT (Computer Based Test). Audit dilakukan pada lapisan Autentikasi, Otorisasi, Integritas Data, dan Infrastruktur.

## 1. Executive Summary

Secara umum, sistem telah mengimplementasikan dasar-dasar keamanan seperti JWT, Role-Based Access Control (RBAC), dan parameter binding untuk mencegah SQL Injection. Namun, ditemukan beberapa celah kritikal pada level logika bisnis (Otorisasi) yang memungkinkan manipulasi data antar pengguna (IDOR).

## 2. Metodologi Audit
Audit dilakukan dengan teknik:
- **Static Analysis (SAST)**: Review kode sumber Golang dan Vue.
- **Business Logic Review**: Menelusuri alur pengerjaan ujian dan sinkronisasi jawaban.
- **Endpoint Testing**: Memeriksa validasi parameter pada API.

## 3. Scorecard Keamanan

| Kategori | Status | Keterangan |
| :--- | :--- | :--- |
| **Authentication** | 🟡 Medium | Menggunakan JWT standar, namun rawan Brute Force. |
| **Authorization** | 🟢 Good | ✅ FIXED (May 2026): IDOR pada endpoint pengerjaan siswa sudah dipatch. |
| **Data Integrity** | 🟢 Good | ✅ FIXED (May 2026): Implementasi Server-side Timer Hardening. |
| **Input Validation** | 🟢 Good | Menggunakan GORM (Safe from SQLi) dan sanitasi path gambar. |
| **Session Mgmt** | 🟡 Medium | Durasi token 24 jam dianggap terlalu lama untuk ujian. |

## 4. Cakupan Audit
- `internal/middleware/auth.go`: Logika pengecekan token dan role.
- `internal/handlers/ujian_siswa.go`: Logika sinkronisasi dan submit jawaban.
- `pkg/utils/image_extractor.go`: Penanganan file gambar dan path.
- `internal/handlers/auth.go`: Mekanisme login dan hashing password.
