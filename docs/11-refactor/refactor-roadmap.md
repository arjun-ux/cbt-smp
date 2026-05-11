# Refactor Roadmap

Rencana perbaikan kode ini dibagi menjadi tiga fase berdasarkan urgensi dan dampak terhadap stabilitas sistem.

## Fase 1: Short Term (Security & Stability)
*Target: 1-2 Minggu*

1.  **Fix Otorisasi (IDOR)**: ✅ COMPLETED (May 2026) - Validasi kepemilikan ditambahkan pada endpoint Sync, Submit, dan Log.
2.  **Server-side Timer Validation**: ✅ COMPLETED (May 2026) - Implementasi Timer Hardening.
3.  **Unified Response Format**: ✅ COMPLETED (May 2026) - Seluruh handler backend telah menggunakan helper SendSuccess dan SendError untuk konsistensi API.
4.  **Database Concurrency Optimization**: ✅ COMPLETED (May 2026) - Implementasi **Bulk Upsert** (OnConflict) dan **Bulk SQL Update** untuk skor. Menghilangkan error `SQLITE_BUSY` secara total.

## Fase 2: Medium Term (Maintainability)
*Target: 1 Bulan*

1.  **Component Decomposition**: ✅ COMPLETED (May 2026) - `UjianPengerjaan.vue` dipecah menjadi komponen modular.
2.  **Pinia Migration**: ✅ COMPLETED (May 2026) - State ujian dipindahkan ke `store/exam.js`.
3.  **Extraction to Composables**: ✅ COMPLETED (May 2026) - Logika timer dan anti-cheat dipindahkan ke folder `composables`.
4.  **Optimasi Query**: ✅ COMPLETED (May 2026) - Implementasi Bulk Operations untuk mengurangi beban I/O database.

## Fase 3: Long Term (Architecture & Scalability)
*Target: 3+ Bulan*

1.  **Repository Pattern**: Melakukan abstraksi database untuk mendukung unit testing dan kemudahan migrasi database di masa depan.
2.  **Service Layer Implementation**: Memindahkan logika bisnis dari handler ke service layer untuk mendukung pemisahan tanggung jawab (Separation of Concerns).
3.  **Frontend Testing**: Implementasi unit testing untuk komponen-komponen kritis di frontend menggunakan Vitest.
