# Question Management Page (Detail Soal)

## Informasi Halaman
- **Nama Halaman**: Detail Butir Soal
- **Route/URL**: `/admin/bank-soal/:id/soal` atau `/guru/bank-soal/:id/soal`
- **Tujuan**: Mengelola butir-butir soal (PG & Essay) di dalam satu Bank Soal.
- **Role Akses**: `admin`, `guru`

## Fitur Utama: Impor Soal dari Word (.docx)
Ini adalah fitur paling krusial untuk memudahkan guru memasukkan soal dalam jumlah banyak sekaligus.

### Cara Kerja:
1. Guru mengunduh **Template Soal (.docx)** yang sudah disediakan.
2. Guru mengisi soal, opsi jawaban (A, B, C, D), dan kunci jawaban ke dalam tabel atau format yang ditentukan di Word.
3. Guru mengunggah file tersebut melalui tombol **"Import Word"**.
4. Backend akan membedah (parse) file Word, mengambil teks serta gambar, dan menyimpannya ke tabel `cbt_soals`.

### Ketentuan Template Word:
- **Jenis Soal**: Harus ditentukan (PG atau ESSAY).
- **Format Opsi**: A, B, C, D harus jelas.
- **Kunci Jawaban**: Ditulis di kolom/baris yang sesuai dalam template.
- **Gambar**: Bisa langsung disisipkan (Copy-Paste) ke dalam file Word.

## Fitur Tambahan:
- **Tambah Manual**: Input soal satu per satu melalui form editor (mendukung KaTeX untuk rumus matematika).
- **Edit Soal**: Mengubah teks pertanyaan, opsi, atau bobot nilai soal yang sudah ada.
- **Hapus Soal**: Menghapus butir soal tertentu.
- **Pratinjau Gambar**: Menampilkan gambar yang sudah diunggah untuk memastikan visual soal benar.

## Teknis & API Calls:
- `GET /api/[role]/bank-soal/:id/soal`: Mengambil daftar soal.
- `POST /api/[role]/bank-soal/:id/import`: Endpoint untuk unggah file .docx.
- `POST /api/[role]/bank-soal/:id/soal`: Simpan soal manual.
- `DELETE /api/[role]/soal/:id`: Hapus butir soal.

## Validasi:
- Sistem akan menolak file jika bukan berformat `.docx`.
- Jika format tabel di dalam Word tidak sesuai template, sistem akan memberikan pesan error pada baris soal yang bermasalah.
