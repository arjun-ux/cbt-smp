# Technology Stack

## Backend
Teknologi utama yang digunakan di sisi server:

- **Bahasa**: [Golang 1.23.0](https://go.dev/)
- **Framework HTTP**: [Fiber v2](https://gofiber.io/) - Framework web yang sangat cepat dan terinspirasi oleh Express.js.
- **ORM**: [GORM](https://gorm.io/) - Developer friendly ORM untuk Golang.
- **Database**: [SQLite](https://www.sqlite.org/) - Database file-based yang ringan dan portable, digunakan via driver `glebarez/sqlite` (Pure Go).
- **Authentication**: [JWT (JSON Web Token)](https://jwt.io/) - Menggunakan library `golang-jwt/jwt/v5`.
- **Security**: `bcrypt` untuk hashing password.
- **Environment**: `godotenv` untuk manajemen variabel lingkungan.

## Frontend
Teknologi utama yang digunakan di sisi client:

- **Framework**: [Vue 3](https://vuejs.org/) (Composition API)
- **Build Tool**: [Vite](https://vitejs.dev/) - Build tool generasi berikutnya yang sangat cepat.
- **Routing**: [Vue Router 4](https://router.vuejs.org/) - Router resmi untuk Vue.js.
- **State Management**: [Pinia](https://pinia.vuejs.org/) - Store management yang ringan dan intuitif.
- **Styling**: [Tailwind CSS 3](https://tailwindcss.com/) - Utility-first CSS framework.
- **Icons**: Heroicons (umum digunakan dengan Tailwind).

## Deployment & Build
- **Mode**: Standalone Binary.
- **Frontend Serving**: File frontend di-build menjadi aset statis dan di-embed ke dalam binary Go menggunakan `embed` package, sehingga aplikasi hanya berupa satu file executable tunggal.
- **Environment Management**: Menggunakan file `.env` untuk konfigurasi port dan database.
