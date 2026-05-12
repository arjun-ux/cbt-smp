# Dokumentasi Modul: Cetak Kartu Peserta Ujian (Modernized)

Modul **Cetak Kartu Peserta Ujian** dirancang untuk memberikan pengalaman administrasi yang profesional, cepat, dan presisi. Modul ini menggabungkan konfigurasi desain *real-time* dengan penarikan data otomatis dari sistem.

## 1. Arsitektur Antarmuka (Professional Parallel Layout)

Antarmuka menggunakan sistem **Dual-Panel Workflow** untuk efisiensi navigasi:

*   **Panel Kiri (Configuration)**: Panel statis (non-sticky) untuk mengatur header dan footer kartu secara langsung.
    *   **Live Edit**: Nama Ujian, Nama Instansi, Alamat, Kota, dan Tanggal.
    *   **Data Integrity**: Informasi Kepala Sekolah (Nama & NIP) ditarik otomatis dari database untuk mencegah kesalahan input manual.
*   **Panel Kanan (Preview & Filter)**:
    *   **Live Preview**: Visualisasi kartu tunggal yang responsif terhadap perubahan di panel konfigurasi.
    *   **Smart Filtering**: Filter dinamis berdasarkan Kelas, Ruang, dan Sesi untuk menentukan daftar siswa yang akan dicetak.
*   **Bottom Section (Print Results)**: Kontainer lebar penuh (*full-width*) yang menampilkan daftar seluruh kartu siswa hasil generate, terpisah dari kontrol navigasi untuk memaksimalkan ruang pandang.

## 2. Integrasi Data Otomatis

Sistem mengeliminasi input manual untuk data administratif melalui integrasi `CBTSetting`:

1.  **Identitas Sekolah**: Diambil dari tabel `cbt_settings` (Key-Value Store).
2.  **Leadership Linkage**: Nama Kepala Sekolah dan NIP tidak diinput di halaman ini, melainkan ditarik melalui logika *join* antara `kepala_sekolah_id` di pengaturan dengan tabel `MasterGuru`.
3.  **Dynamic Mapping**: Saat memuat data cetak (`/api/admin/cetak/kartu`), backend secara otomatis menyertakan detail kepala sekolah yang aktif untuk disematkan di footer kartu.

## 3. Spesifikasi Teknis Pencetakan (A4 Layout)

Modul ini dioptimalkan untuk kertas **A4 Portrait** dengan akurasi dimensi tinggi:

| Properti | Spesifikasi | Keterangan |
| :--- | :--- | :--- |
| **Dimensi Kartu** | 100mm x 70mm | Dioptimalkan agar 8 kartu muat dalam 1 lembar A4 |
| **Total Width (Grid)** | 200mm | Memberikan margin aman 5mm di kiri & kanan halaman |
| **Grid Layout** | 2 Kolom (Portrait) | Menggunakan CSS Grid dengan `gap-y: 2mm` dan `gap-x: 0` |
| **CSS Print Control** | `@media print` | Menyembunyikan Navbar, Sidebar, dan UI controls secara total |
| **Color Adjust** | `exact` | Memaksa browser mencetak warna border dan background logo |

## 4. Keamanan & Performa

*   **Zero Manual Entry**: Mengurangi risiko "typo" pada dokumen resmi negara (Kartu Ujian).
*   **Base64 Logo Handling**: Mendukung logo dinamis (Kiri & Kanan) yang langsung di-render dari database.
*   **Print Visibility**: Menggunakan class `.no-print` untuk memastikan elemen web tidak mengganggu dokumen fisik.

---
**Status**: Production Ready ✅  
**Terakhir Diperbarui**: 2026-05-13  
**Konteks**: Modernisasi UI & Otomatisasi Data Kepala Sekolah.
