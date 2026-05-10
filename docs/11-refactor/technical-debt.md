# Technical Debt Inventory

Dokumen ini mencatat "hutang teknis" yang ada di dalam proyek CBT. Hutang ini adalah keputusan desain atau implementasi yang diambil demi kecepatan pengembangan namun perlu diperbaiki agar tidak menjadi beban di masa depan.

## 1. Arsitektur & Backend

| Item | Deskripsi | Risiko |
| :--- | :--- | :--- |
| **Direct DB Access** | Handlers memanggil GORM Global DB secara langsung. | Sulit dilakukan Testing (Mocking). |
| **Logic in Handlers** | Logika bisnis (seperti hitung nilai) berada di layer handler. | Kode sulit dibaca dan tidak reusable. |
| **Insecure Secrets** | Penggunaan secret key default di kode program. | Potensi eksploitasi jika env tidak disetel. |
| **SQLite Bottleneck** | Penggunaan SQLite untuk pengerjaan massal. | Potensi database locked jika RPS sangat tinggi (>200 concurrent write). |

## 2. Frontend & UI

| Item | Deskripsi | Risiko |
| :--- | :--- | :--- |
| **Large Vue Files** | Komponen views yang sangat besar (fat components). | Perubahan UI berisiko tinggi merusak logika. |
| **Prop Drilling** | Pengiriman data melalui banyak layer komponen tanpa Store. | Sulit melacak aliran data. |
| **Hardcoded Style** | Penggunaan utility tailwind yang tidak konsisten. | Desain UI menjadi tidak seragam (Inconsistent). |
| **KaTeX Performance** | Rendering rumus matematika dilakukan berulang kali saat navigasi. | Browser menjadi lag pada perangkat spek rendah. |

## 3. Infrastruktur & DevOps

| Item | Deskripsi | Risiko |
| :--- | :--- | :--- |
| **No Auto-Testing** | Tidak ada pipeline CI/CD untuk testing otomatis. | Bug regresi mudah muncul kembali. |
| **Manual Deployment** | Proses build dan deploy masih manual. | Risiko human error saat deployment. |
| **Static Uploads** | Gambar disimpan di folder lokal aplikasi. | Sulit jika aplikasi ingin dideploy ke multiple server (Scalability). |

## Kesimpulan
Hutang teknis saat ini berada pada level **Manageable** namun perlu segera ditangani pada bagian Keamanan dan Efisiensi Query untuk menghindari kegagalan sistem saat ujian sebenarnya berlangsung.
