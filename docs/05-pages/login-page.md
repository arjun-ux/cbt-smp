# Login Page

## Informasi Halaman
- **Nama Halaman**: Login Utama
- **Route/URL**: `/login`
- **Tujuan**: Pintu masuk utama bagi Admin dan Guru ke dalam sistem.
- **Role Akses**: Public (Unauthenticated)

## Teknis & State
- **Component**: `Login.vue`
- **Store**: `useAuthStore`
- **API Call**: `POST /api/auth/login`
- **State Lokal**:
    - `username`: String (Input user/NIP)
    - `password`: String (Input password)
    - `isLoading`: Boolean (Loading state saat request)
    - `errorMessage`: String (Menampung pesan error dari server)
- **Lifecycle**: `onMounted` tidak melakukan fetch data khusus.

## Elemen UI & Aksi

### Form Login
- **Fungsi**: Mengumpulkan kredensial user.
- **Validasi**: Input `required` pada HTML level.
- **Event**: `handleLogin` dipicu saat form di-submit.

### Tombol "Masuk"
- **Fungsi**: Memicu proses autentikasi.
- **Flow**:
    1. Klik tombol.
    2. `isLoading` menjadi `true`.
    3. Kirim request ke `/api/auth/login`.
    4. Jika sukses: Simpan token ke Pinia Store dan LocalStorage via `authStore.setAuth`.
    5. Redirect ke dashboard sesuai role (`AdminDashboard` atau `GuruDashboard`).
    6. Jika gagal: Tampilkan pesan error di `errorMessage`.
- **Database Effect**: Mencatat aktivitas ke tabel `log_sistems` (melalui backend).
- **Error Handling**: Menangkap pesan error dari response JSON backend (misal: "Username/Password salah").
