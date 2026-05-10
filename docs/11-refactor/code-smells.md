# Code Smells Analysis

Berikut adalah hasil audit kualitas kode pada sistem CBT, dikategorikan berdasarkan jenis permasalahan dan tingkat keparahannya.

## 1. God Component (`UjianPengerjaan.vue`)
- **Severity**: 🟢 FIXED
- **Status**: ✅ FIXED
- **Alasan**: Komponen ini menangani terlalu banyak hal: Rendering soal, Timer, Logika Anti-Cheat, Sinkronisasi data ke server, dan Rendering KaTeX.
- **Dampak**: File sangat besar (>600 baris), sulit dibaca, dan perubahan kecil pada UI berisiko merusak logika sinkronisasi.
- **Solusi**: Pecah menjadi komponen-komponen kecil: `ExamTimer.vue`, `ExamNavigator.vue`, `QuestionItem.vue`. Pindahkan logika state ke Pinia Store.

## 2. Duplicate Logic & Route Overlap
- **Severity**: 🟡 MEDIUM
- **Alasan**: Terdapat duplikasi rute antara `/api/admin` dan `/api/guru` yang memanggil handler yang sama namun dengan pengecekan role manual di dalam handler.
- **Dampak**: Duplikasi kode di `routes.go` dan logika otorisasi yang tersebar di banyak tempat (Hard to maintain).
- **Solusi**: Gunakan rute terpadu (misal: `/api/manage/...`) dan biarkan middleware atau Service Layer menangani filter data berdasarkan role.

## 3. Mixed Responsibility in Handlers
- **Severity**: 🟡 MEDIUM
- **Alasan**: File seperti `monitoring.go` mencampur logika monitoring real-time, koreksi essay, dan analisis butir soal.
- **Dampak**: Melanggar *Single Responsibility Principle*. Handler menjadi terlalu gemuk dan sulit di-unit test.
- **Solusi**: Pisahkan logika bisnis ke dalam Service Layer (misal: `MonitorService`, `ScoringService`).

## 4. Inefficient Database Queries (N+1 Problem)
- **Severity**: 🟢 OPTIMIZED
- **Status**: ✅ COMPLETED (BULK SQL OPTIMIZED)
- **Alasan**: Fungsi penghitungan nilai kini menggunakan **Bulk SQL Update (`UPDATE ... CASE`)** dan fungsi simpan jawaban menggunakan **Atomic Bulk Upsert (`ON CONFLICT`)**.
- **Dampak**: Menghilangkan ribuan operasi tulis (write) berulang dalam satu transaksi.
- **Hasil**: Menghilangkan error *SQLITE_BUSY* secara total baik saat sinkronisasi jawaban masal maupun saat pengakhiran jadwal serentak.

## 5. Tight Coupling to GORM
- **Severity**: 🔵 LOW
- **Alasan**: Hampir seluruh handler memanggil `database.DB` secara langsung.
- **Dampak**: Sulit untuk melakukan mocking database saat pengujian otomatis (Unit Testing).
- **Solusi**: Implementasikan **Repository Pattern** untuk mengabstraksi akses data.

## 6. Mixed Logic in Frontend Views
- **Severity**: 🟡 MEDIUM
- **Alasan**: Logika sinkronisasi batch dan penanganan error 401 ditulis langsung di dalam method Vue.
- **Dampak**: Logika tidak dapat digunakan kembali (non-reusable) dan mengotori kode presentasi UI.
- **Solusi**: Ekstrak ke dalam Vue Composables (misal: `useExamSync.js`).
