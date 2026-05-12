import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '../store/auth'

const routes = [
  {
    path: '/',
    name: 'StudentLogin',
    component: () => import('../views/siswa/UjianLogin.vue')
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('../views/Login.vue')
  },
  {
    path: '/siswa',
    component: () => import('../components/StudentLayout.vue'),
    meta: { requiresAuth: true, role: 'siswa' },
    children: [
      {
        path: 'dashboard',
        name: 'StudentDashboard',
        component: () => import('../views/siswa/StudentDashboard.vue')
      },
      {
        path: 'konfirmasi',
        name: 'StudentConfirm',
        component: () => import('../views/siswa/UjianKonfirmasi.vue')
      },
      {
        path: 'pengerjaan/:jadwalId',
        name: 'StudentPengerjaan',
        component: () => import('../views/siswa/UjianPengerjaan.vue'),
        props: true
      }
    ]
  },
  {
    path: '/admin',
    component: () => import('../components/DashboardLayout.vue'),
    meta: { requiresAuth: true, role: 'admin' },
    children: [
      {
        path: '',
        name: 'AdminDashboard',
        component: () => import('../views/admin/Dashboard.vue')
      },
      {
        path: 'manage-admin',
        name: 'AdminMaster',
        component: () => import('../views/admin/AdminMaster.vue')
      },
      {
        path: 'bank-soal',
        name: 'AdminBankSoal',
        component: () => import('../views/admin/BankSoalAdmin.vue')
      },
      {
        path: 'bank-soal/:bankSoalId/soal',
        name: 'AdminSoalDetail',
        component: () => import('../views/admin/SoalDetail.vue'),
        props: true
      },
      {
        path: 'guru',
        name: 'AdminGuru',
        component: () => import('../views/admin/GuruMaster.vue')
      },
      {
        path: 'siswa',
        name: 'AdminSiswa',
        component: () => import('../views/admin/SiswaMaster.vue')
      },
      {
        path: 'plotting-siswa',
        name: 'AdminPlotting',
        component: () => import('../views/admin/PlottingSiswa.vue')
      },
      {
        path: 'kelas',
        name: 'AdminKelas',
        component: () => import('../views/admin/KelasMaster.vue')
      },
      {
        path: 'mapel',
        name: 'AdminMapel',
        component: () => import('../views/admin/MapelMaster.vue')
      },
      {
        path: 'ruang',
        name: 'AdminRuang',
        component: () => import('../views/admin/RuangMaster.vue')
      },
      {
        path: 'sesi',
        name: 'AdminSesi',
        component: () => import('../views/admin/SesiMaster.vue')
      },
      {
        path: 'jadwal',
        name: 'AdminJadwal',
        component: () => import('../views/admin/JadwalUjian.vue')
      },
      {
        path: 'plotting-pengawas',
        name: 'AdminPlottingPengawas',
        component: () => import('../views/admin/PlottingPengawas.vue')
      },
      {
        path: 'monitor/:jadwalId',
        name: 'AdminMonitor',
        component: () => import('../views/admin/MonitorUjian.vue')
      },
      {
        path: 'rekap-nilai/:jadwalId?',
        name: 'AdminRekap',
        component: () => import('../views/admin/RekapNilai.vue')
      },
      {
        path: 'analisis-soal/:jadwalId',
        name: 'AdminAnalisis',
        component: () => import('../views/admin/AnalisisSoal.vue'),
        props: true
      },
      {
        path: 'riwayat-nilai',
        name: 'AdminRiwayat',
        component: () => import('../views/admin/RiwayatNilai.vue')
      },
      {
        path: 'pengaturan',
        name: 'AdminSettings',
        component: () => import('../views/admin/SchoolSettings.vue')
      },
      {
        path: 'cetak-kartu',
        name: 'AdminCetakKartu',
        component: () => import('../views/admin/CetakKartu.vue')
      }
    ]
  },
  {
    path: '/guru',
    component: () => import('../components/DashboardLayout.vue'),
    meta: { requiresAuth: true, role: 'guru' },
    children: [
      {
        path: '',
        name: 'GuruDashboard',
        component: () => import('../views/admin/Dashboard.vue')
      },
      {
        path: 'bank-soal',
        name: 'GuruBankSoal',
        component: () => import('../views/admin/BankSoalAdmin.vue')
      },
      {
        path: 'bank-soal/:bankSoalId/soal',
        name: 'GuruSoalDetail',
        component: () => import('../views/admin/SoalDetail.vue'),
        props: true
      },
      {
        path: 'jadwal',
        name: 'GuruJadwal',
        component: () => import('../views/admin/JadwalUjian.vue')
      },
      {
        path: 'monitor/:jadwalId',
        name: 'GuruMonitor',
        component: () => import('../views/admin/MonitorUjian.vue')
      },
      {
        path: 'rekap-nilai/:jadwalId?',
        name: 'GuruRekap',
        component: () => import('../views/admin/RekapNilai.vue')
      },
      {
        path: 'analisis-soal/:jadwalId',
        name: 'GuruAnalisis',
        component: () => import('../views/admin/AnalisisSoal.vue'),
        props: true
      },
      {
        path: 'riwayat-nilai',
        name: 'GuruRiwayat',
        component: () => import('../views/admin/RiwayatNilai.vue')
      }
    ]
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// Navigation Guard
router.beforeEach((to, from, next) => {
  const authStore = useAuthStore()
  const userRole = authStore.user?.role

  // 1. Jika butuh auth tapi belum login
  if (to.meta.requiresAuth && !authStore.token) {
    if (to.path.startsWith('/admin') || to.path.startsWith('/guru')) {
      return next({ name: 'Login' })
    }
    return next({ name: 'StudentLogin' })
  }

  // 2. Jika sudah login dan mencoba ke halaman login
  if ((to.name === 'Login' || to.name === 'StudentLogin') && authStore.token) {
    if (userRole === 'admin') return next({ name: 'AdminDashboard' })
    if (userRole === 'guru') return next({ name: 'GuruDashboard' })
    if (userRole === 'siswa') return next({ name: 'StudentDashboard' })
  }

  // 3. Proteksi Role
  if (to.meta.role && to.meta.role !== userRole) {
    if (userRole === 'admin') return next() // Admin bisa akses semuanya
    if (userRole === 'guru' && to.meta.role === 'guru') return next()
    
    // Selain itu, lempar ke dashboard masing-masing
    if (userRole === 'admin') return next({ name: 'AdminDashboard' })
    if (userRole === 'guru') return next({ name: 'GuruDashboard' })
    if (userRole === 'siswa') return next({ name: 'StudentDashboard' })
  }

  next()
})

export default router
