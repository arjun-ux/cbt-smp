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
2. Guru mengisi soal, opsi jawaban (A, B, C, D), dan kunci jawaban ke dalam tabel template.
3. Guru mengunggah file tersebut melalui tombol **"Import Word"**.
4. **Proses Parsing**: Frontend menggunakan pustaka `mammoth.js` untuk membaca isi file Word secara lokal di browser.
5. **Feedback Visual**: Sistem menampilkan *Loading Overlay* dengan efek blur selama proses pembacaan file berlangsung.
6. **Backend Processing**: Gambar yang disisipkan di Word akan otomatis diekstrak dan disimpan sebagai file fisik oleh backend menggunakan utilitas `ExtractBase64Images`.

### Ketentuan Template Word:
- **Jenis Soal**: Diidentifikasi berdasarkan angka di kolom tipe (2/5 untuk Essay, lainnya PG).
- **Kunci Jawaban**: Untuk PG, ditandai dengan huruf 'v' atau 'V' pada kolom pilihan yang benar.
- **Gambar**: Disisipkan langsung di dalam sel tabel pertanyaan atau opsi.

## Fitur Tambahan:
- **Tambah Manual**: Input soal satu per satu melalui form editor (mendukung KaTeX untuk rumus matematika).
- **Edit Soal**: Mengubah teks pertanyaan, opsi, atau bobot nilai.
- **Hapus Soal**: Menghapus butir soal tertentu.
- **Pratinjau Impor**: Menampilkan daftar soal yang berhasil dibaca dari Word sebelum benar-benar disimpan ke database.

## Teknis & API Calls:
- `GET /api/[role]/bank-soal/:id/soal`: Mengambil daftar soal.
- `POST /api/[role]/bank-soal/:id/soal`: Digunakan secara looping oleh frontend saat proses impor masal.
- `DELETE /api/[role]/soal/:id`: Hapus butir soal.
- `DELETE /api/[role]/bank-soal/:id/soal/clear`: Membersihkan semua soal lama sebelum impor baru dimulai.

## Validasi:
- Sistem akan menolak file jika bukan berformat `.docx`.
- Jika format tabel di dalam Word tidak sesuai template, sistem akan memberikan pesan error pada baris soal yang bermasalah.
