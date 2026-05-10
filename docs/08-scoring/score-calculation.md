# Score Calculation Logic

Halaman ini menjelaskan rumus dan logika teknis yang digunakan backend untuk menghitung nilai.

## 1. Rumus Pilihan Ganda (PG)

Penilaian PG bersifat biner per butir soal. Tidak ada nilai parsial (misal: benar 2 dari 3 pilihan).

**Logika Program:**
```go
if jawabanSiswa == soal.KunciJawaban {
    skor = soal.BobotNilai
} else {
    skor = 0
}
```

**Rumus Akumulasi:**
$$NilaiPG = \sum (SkorPG_{i})$$

> [!NOTE]
> Sistem tidak menggunakan rumus persentase (Skor/Total * 100) di level database, melainkan menggunakan **Bobot Nilai Langsung**. Jika admin ingin nilai maksimal 100 dan ada 50 soal, maka setiap soal harus diberi bobot 2.0.

## 2. Rumus Essay

Penilaian essay sepenuhnya subjektif dan dilakukan oleh Guru.

**Logika Program:**
1. Guru memasukkan `skor` manual untuk setiap butir essay.
2. Backend memvalidasi agar `skor <= bobot_maks` soal tersebut.
3. `NilaiEssay` adalah jumlahan dari seluruh skor essay yang diberikan.

## 3. Nilai Akhir (Total Score)

Nilai akhir adalah penjumlahan sederhana dari dua komponen nilai:
$$TotalNilai = NilaiPG + NilaiEssay$$

## 4. Analisis Edge Cases

### Jawaban Kosong
- **PG**: Jika siswa tidak menjawab, `HitungNilaiPG` akan memberikan skor 0 karena `"" != KunciJawaban`.
- **Essay**: Jika siswa tidak menjawab, record jawaban tetap dibuat dengan string kosong, dan Guru dapat memberikan skor 0.

### Bobot Desimal
Sistem menggunakan tipe data `float64`. Ini memungkinkan penggunaan bobot seperti `1.5` atau `2.25`. Namun, perlu diperhatikan potensi *floating point precision* pada laporan jika tidak dibatasi jumlah angka di belakang koma (saat ini frontend membatasi ke 1 angka desimal).

### Manipulasi Nilai
- **Risiko**: Karena nilai PG dihitung di server *setelah* submit, manipulasi di sisi browser (client-side) tidak akan mempengaruhi hasil akhir.
- **Risiko**: Nilai tersimpan di database dalam bentuk mentah. Akses langsung ke database (`cbt.db`) dapat mengubah nilai tanpa melalui sistem aplikasi.
- **Mitigasi**: Implementasi `LogUjian` mencatat waktu submit, sehingga jika ada perubahan nilai tanpa ada log aktivitas yang sesuai, dapat dideteksi sebagai anomali.
