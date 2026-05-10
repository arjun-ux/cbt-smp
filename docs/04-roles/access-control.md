# Access Control Implementation

Kontrol akses diterapkan secara berlapis di sisi Backend (API Protection) dan Frontend (UI Navigation).

## 1. Backend Route Protection (Golang/Fiber)
Proteksi dilakukan di `internal/routes/routes.go` menggunakan middleware.

### Layer 1: Authentication Guard
Semua rute di bawah group `protected` menggunakan middleware `Protected()`.
```go
protected := router.Group("/")
protected.Use(middleware.Protected())
```

### Layer 2: Role Authorization
Menggunakan middleware `RoleRequired()` untuk membatasi akses group rute ke role tertentu.
- **Admin Group**: `admin.Use(middleware.RoleRequired("admin"))`
- **Guru Group**: `guru.Use(middleware.RoleRequired("admin", "guru"))`
- **Siswa Group**: `siswa.Use(middleware.RoleRequired("siswa"))`

> [!IMPORTANT]
> Role **Admin** secara eksplisit diberikan izin untuk mengakses rute **Guru**, memberikan fleksibilitas bagi admin untuk membantu guru mengelola soal atau memantau ujian.

## 2. Frontend Navigation Guard (Vue/Vite)
Proteksi dilakukan di `frontend/src/router/index.js` menggunakan global navigation guard `router.beforeEach`.

- **`requiresAuth`**: Jika rute memiliki meta `requiresAuth: true`, router akan memeriksa keberadaan token di Pinia Store. Jika tidak ada, pengguna dilempar ke halaman Login.
- **Role Verification**: Router membandingkan role user di store dengan `meta.role` pada definisi rute.
    - Jika user mencoba akses halaman Admin tetapi role-nya adalah `siswa`, maka akan diarahkan ke Dashboard Siswa.
- **Redirection Logic**: Pengguna yang sudah login dilarang mengakses halaman login (`/login`) dan akan diarahkan ke dashboard masing-masing.

## 3. UI Component Authorization
Beberapa elemen UI disembunyikan atau dinonaktifkan berdasarkan role user menggunakan direktif `v-if` di Vue:
- **Menu Sidebar**: Sidebar hanya merender item menu yang sesuai dengan role user.
- **Action Buttons**: Tombol "Hapus" atau "Edit" data master hanya muncul di dashboard Admin.

## 4. Unauthorized Flow Detail
1.  **Client** mengirim request tanpa token / token expired.
2.  **Backend** membalas dengan `401 Unauthorized`.
3.  **Frontend (Axios Interceptor)** menangkap status `401`.
4.  **Frontend** memicu aksi `authStore.logout()`.
5.  **Frontend** melakukan `window.location.reload()` atau `router.push('/login')`.
6.  User melihat pesan "Sesi telah berakhir, silakan login kembali".
