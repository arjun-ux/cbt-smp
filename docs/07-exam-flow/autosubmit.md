# Auto-Submit & Force Submit

Sistem menjamin bahwa setiap sesi ujian akan berakhir (Finalized) baik melalui aksi siswa maupun kontrol sistem.

## 1. Auto-Submit (Client-Side)
Terjadi ketika timer di sisi siswa mencapai angka nol.

- **Trigger**: `timeLeft.value <= 0`.
- **Flow**:
    1. Interval timer mendeteksi waktu habis.
    2. Frontend secara otomatis memanggil fungsi `finishExam(true)`.
    3. Parameter `true` menandakan auto-submit, melewati konfirmasi modal.
    4. Seluruh jawaban yang ada di state dikirim (final sync) dan sesi ditutup.

## 2. Force Submit (Server-Side/Proctor)
Sistem menyediakan dua cara untuk mengakhiri sesi siswa secara paksa dari sisi server:

### A. Manual Force Submit (Per Siswa)
Dilakukan melalui Dashboard Monitoring untuk siswa tertentu (misal: karena pelanggaran).
- **Flow**: Pengawas klik ikon "Force Submit" -> Backend hitung nilai -> Status DB jadi `Selesai`.

### B. Batch Force Submit (Auto-Finalize)
Dilakukan secara otomatis ketika status **Jadwal Ujian** diubah menjadi **"Selesai"** oleh Admin/Guru.
- **Tujuan**: Menjamin seluruh siswa yang lupa klik selesai tetap mendapatkan nilai berdasarkan jawaban terakhir mereka.
- **Flow**:
    1. Status Jadwal diubah ke "Selesai".
    2. Backend mencari seluruh peserta yang masih berstatus "Sedang Mengerjakan".
    3. Untuk setiap peserta, sistem memicu `HitungNilaiPG` dan mengupdate status menjadi "Selesai".
    4. **Catatan**: Peserta yang berstatus **Terblokir** tidak akan diproses secara otomatis untuk menghindari konflik data; Pengawas harus melakukan penanganan manual jika diperlukan.
    5. Seluruh proses dibungkus dalam transaksi database untuk menjamin keamanan data.

## 3. Penanganan Scoring Otomatis
Setiap proses Submit (Auto/Manual/Force) akan memicu fungsi `HitungNilaiPG` di backend:

```go
func HitungNilaiPG(db *gorm.DB, pesertaID uint, bankSoalID uint) float64 {
    // 1. Ambil Kunci Jawaban dari cbt_soals
    // 2. Ambil Jawaban Siswa dari cbt_jawaban_siswas
    // 3. Bandingkan Soal per Soal
    // 4. Akumulasikan BobotNilai jika Benar
    // 5. Update kolom nilai_pg di cbt_peserta_ujians
}
```

## 4. Analisis Kegagalan (Failure Handling)

| Skenario | Dampak | Penanganan |
| :--- | :--- | :--- |
| Internet Mati saat Timer 0 | Auto-submit gagal kirim | Jawaban tersimpan di LocalStorage. Siswa harus mencari sinyal/hotspot dan membuka kembali dashboard untuk pengerjaan ulang (data akan auto-sync saat login kembali). |
| Server Down saat Submit | Response 500 | Frontend menampilkan alert error. Sesi tetap dianggap "Sedang Mengerjakan" di DB. Siswa bisa mencoba submit ulang setelah server normal. |
| Browser Crash saat Auto-submit | Proses terhenti | Sisa waktu di DB sudah mendekati 0. Saat siswa login lagi, sistem akan mendeteksi waktu habis dan langsung memicu submit ulang. |
| Listrik Mati Total | Data belum di-submit | Jawaban aman di DB hingga titik sinkronisasi terakhir (setiap 15 detik). Pengawas bisa melakukan Force Submit dari PC server untuk memproses nilai yang sudah masuk. |
