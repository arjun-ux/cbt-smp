<script setup>
import { computed, ref, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../store/auth'
import ConfirmModal from './ConfirmModal.vue'

const router = useRouter()
const authStore = useAuthStore()

const userRole = computed(() => authStore.user?.role)
const isSidebarOpen = ref(false)
const showLogoutConfirm = ref(false)
const currentTime = ref('')

const updateTime = () => {
  const now = new Date()
  currentTime.value = now.toLocaleTimeString('id-ID', { 
    hour: '2-digit', 
    minute: '2-digit', 
    second: '2-digit' 
  })
}

let timer
onMounted(() => {
  updateTime()
  timer = setInterval(updateTime, 1000)
})

onUnmounted(() => {
  clearInterval(timer)
})

const triggerLogout = () => {
  showLogoutConfirm.value = true
}

const handleLogout = () => {
  authStore.logout()
  router.push({ name: 'Login' })
}

const toggleSidebar = () => {
  isSidebarOpen.value = !isSidebarOpen.value
}
</script>

<template>
  <div class="h-screen flex overflow-hidden font-sans bg-slate-50">
    
    <!-- Mobile Sidebar Overlay -->
    <div 
      v-if="isSidebarOpen" 
      @click="toggleSidebar"
      class="fixed inset-0 bg-slate-900/60 backdrop-blur-sm z-40 lg:hidden transition-opacity"
    ></div>

    <!-- Sidebar -->
    <aside 
      :class="[
        isSidebarOpen ? 'translate-x-0' : '-translate-x-full',
        'fixed inset-y-0 left-0 z-50 w-72 bg-gradient-to-b from-slate-950 to-slate-900 text-white transition-transform duration-300 ease-in-out lg:translate-x-0 lg:static lg:flex-shrink-0 shadow-2xl overflow-hidden flex flex-col'
      ]"
    >
      <!-- App Logo -->
      <div class="flex items-center justify-between p-6 border-b border-white/10 flex-shrink-0">
        <div>
          <h2 class="text-2xl font-black bg-clip-text text-transparent bg-gradient-to-r from-blue-400 to-indigo-300 tracking-tight">CBT System</h2>
          <p class="text-[10px] text-slate-500 mt-1 uppercase tracking-[0.2em] font-black flex items-center gap-2">
            <span class="w-2 h-2 rounded-full bg-blue-500 shadow-[0_0_8px_rgba(59,130,246,0.5)]"></span>
            {{ userRole }} Workspace
          </p>
        </div>
        <button @click="toggleSidebar" class="lg:hidden text-slate-400 hover:text-white transition-colors">
          <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path></svg>
        </button>
      </div>
      
      <nav class="flex-1 mt-4 px-4 space-y-1.5 overflow-y-auto scrollbar-hide">
        
        <!-- SECTION: UTAMA -->
        <div class="pt-4 pb-2 px-4">
          <p class="text-[10px] font-black text-slate-600 uppercase tracking-widest">Utama</p>
        </div>
        <router-link :to="{ name: userRole === 'admin' ? 'AdminDashboard' : 'GuruDashboard' }" @click="isSidebarOpen = false" class="nav-link group" exact-active-class="active">
          <div class="nav-icon-box">
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6"></path></svg>
          </div>
          <span class="flex-1 text-xs uppercase tracking-wider font-bold">Dashboard</span>
        </router-link>

        <!-- SECTION: DATA MASTER (Admin Only) -->
        <template v-if="userRole === 'admin'">
          <div class="pt-6 pb-2 px-4">
            <p class="text-[10px] font-black text-slate-600 uppercase tracking-widest">Master Data</p>
          </div>
          <router-link :to="{ name: 'AdminMapel' }" @click="isSidebarOpen = false" class="nav-link group" active-class="active">
            <div class="nav-icon-box">
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253"></path></svg>
            </div>
            <span class="flex-1 text-xs uppercase tracking-wider font-bold">Mata Pelajaran</span>
          </router-link>
          <router-link :to="{ name: 'AdminKelas' }" @click="isSidebarOpen = false" class="nav-link group" active-class="active">
            <div class="nav-icon-box">
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4"></path></svg>
            </div>
            <span class="flex-1 text-xs uppercase tracking-wider font-bold">Data Kelas</span>
          </router-link>
          <router-link :to="{ name: 'AdminRuang' }" @click="isSidebarOpen = false" class="nav-link group" active-class="active">
            <div class="nav-icon-box">
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6"></path></svg>
            </div>
            <span class="flex-1 text-xs uppercase tracking-wider font-bold">Ruang Ujian</span>
          </router-link>
          <router-link :to="{ name: 'AdminSesi' }" @click="isSidebarOpen = false" class="nav-link group" active-class="active">
            <div class="nav-icon-box">
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"></path></svg>
            </div>
            <span class="flex-1 text-xs uppercase tracking-wider font-bold">Sesi Ujian</span>
          </router-link>

          <!-- SECTION: PENGGUNA (Admin Only) -->
          <div class="pt-6 pb-2 px-4">
            <p class="text-[10px] font-black text-slate-600 uppercase tracking-widest">Pengguna</p>
          </div>
          <router-link :to="{ name: 'AdminGuru' }" @click="isSidebarOpen = false" class="nav-link group" active-class="active">
            <div class="nav-icon-box">
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 14l9-5-9-5-9 5 9 5zm0 0l6.16-3.422a12.083 12.083 0 01.665 6.479A11.952 11.952 0 0012 20.055a11.952 11.952 0 00-6.824-2.998 12.078 12.078 0 01.665-6.479L12 14zm-4 6v-7.5l4-2.222"></path></svg>
            </div>
            <span class="flex-1 text-xs uppercase tracking-wider font-bold">Data Guru</span>
          </router-link>
          <router-link :to="{ name: 'AdminSiswa' }" @click="isSidebarOpen = false" class="nav-link group" active-class="active">
            <div class="nav-icon-box">
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z"></path></svg>
            </div>
            <span class="flex-1 text-xs uppercase tracking-wider font-bold">Data Siswa</span>
          </router-link>
          <router-link :to="{ name: 'AdminPlotting' }" @click="isSidebarOpen = false" class="nav-link group" active-class="active">
            <div class="nav-icon-box">
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-3 7h3m-3 4h3m-6-4h.01m-.01 4h.01"></path></svg>
            </div>
            <span class="flex-1 text-xs uppercase tracking-wider font-bold">Plotting Siswa</span>
          </router-link>
        </template>

        <!-- SECTION: MODUL UJIAN -->
        <div class="pt-6 pb-2 px-4">
          <p class="text-[10px] font-black text-slate-600 uppercase tracking-widest">Ujian</p>
        </div>
        <router-link :to="{ name: userRole === 'admin' ? 'AdminBankSoal' : 'GuruBankSoal' }" @click="isSidebarOpen = false" class="nav-link group" active-class="active">
          <div class="nav-icon-box">
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 14v3m4-3v3m4-3v3M3 21h18M3 10h18M3 7l9-4 9 4M4 10h16v11H4V10z"></path></svg>
          </div>
          <span class="flex-1 text-xs uppercase tracking-wider font-bold">Bank Soal</span>
        </router-link>
        
        <!-- Conditional Routes based on role -->
        <template v-if="userRole === 'admin'">
          <router-link :to="{ name: 'AdminJadwal' }" @click="isSidebarOpen = false" class="nav-link group" active-class="active">
            <div class="nav-icon-box">
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"></path></svg>
            </div>
            <span class="flex-1 text-xs uppercase tracking-wider font-bold">Master Jadwal</span>
          </router-link>
          <router-link :to="{ name: 'AdminPlottingPengawas' }" @click="isSidebarOpen = false" class="nav-link group" active-class="active">
            <div class="nav-icon-box">
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z"></path></svg>
            </div>
            <span class="flex-1 text-xs uppercase tracking-wider font-bold">Plotting Pengawas</span>
          </router-link>
        </template>
        <router-link v-else :to="{ name: 'GuruJadwal' }" @click="isSidebarOpen = false" class="nav-link group" active-class="active">
          <div class="nav-icon-box">
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z"></path></svg>
          </div>
          <span class="flex-1 text-xs uppercase tracking-wider font-bold">Jadwal & Monitor</span>
        </router-link>

        <router-link :to="{ name: userRole === 'admin' ? 'AdminRekap' : 'GuruRekap' }" @click="isSidebarOpen = false" class="nav-link group" active-class="active">
          <div class="nav-icon-box">
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 17v-2m3 2v-4m3 4v-6m2 10H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"></path></svg>
          </div>
          <span class="flex-1 text-xs uppercase tracking-wider font-bold">Laporan & Nilai</span>
        </router-link>

        <router-link :to="{ name: userRole === 'admin' ? 'AdminRiwayat' : 'GuruRiwayat' }" @click="isSidebarOpen = false" class="nav-link group" active-class="active">
          <div class="nav-icon-box">
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253"></path></svg>
          </div>
          <span class="flex-1 text-xs uppercase tracking-wider font-bold">Riwayat Nilai</span>
        </router-link>

        <!-- SECTION: ADMINISTRASI (Admin Only) -->
        <template v-if="userRole === 'admin'">
          <div class="pt-6 pb-2 px-4">
            <p class="text-[10px] font-black text-slate-600 uppercase tracking-widest">Administrasi</p>
          </div>
          <router-link :to="{ name: 'AdminCetakKartu' }" @click="isSidebarOpen = false" class="nav-link group" active-class="active">
            <div class="nav-icon-box">
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 10h18M7 15h1m4 0h1m-7 4h12a3 3 0 003-3V8a3 3 0 00-3-3H6a3 3 0 00-3 3v8a3 3 0 003 3z"></path></svg>
            </div>
            <span class="flex-1 text-xs uppercase tracking-wider font-bold">Cetak Kartu Peserta</span>
          </router-link>
        </template>

        <!-- SECTION: SISTEM (Admin Only) -->
        <template v-if="userRole === 'admin'">
          <div class="pt-6 pb-2 px-4">
            <p class="text-[10px] font-black text-slate-600 uppercase tracking-widest">Sistem</p>
          </div>
          <router-link :to="{ name: 'AdminMaster' }" @click="isSidebarOpen = false" class="nav-link group" active-class="active">
            <div class="nav-icon-box">
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z"></path><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"></path></svg>
            </div>
            <span class="flex-1 text-xs uppercase tracking-wider font-bold">Manajemen Admin</span>
          </router-link>
          <router-link :to="{ name: 'AdminSettings' }" @click="isSidebarOpen = false" class="nav-link group" active-class="active">
            <div class="nav-icon-box">
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4"></path></svg>
            </div>
            <span class="flex-1 text-xs uppercase tracking-wider font-bold">Pengaturan Sekolah</span>
          </router-link>
        </template>

      </nav>
      
      <!-- User Profile & Logout at Bottom -->
      <div class="p-4 border-t border-white/5 bg-black/20 flex-shrink-0">
        <div class="flex items-center justify-between gap-2 p-2 rounded-2xl bg-white/5 border border-white/5">
          <div class="flex items-center gap-3 min-w-0">
            <div class="w-10 h-10 rounded-xl bg-gradient-to-tr from-blue-600 to-indigo-600 flex-shrink-0 flex items-center justify-center text-white font-black shadow-lg">
              {{ authStore.user?.username?.charAt(0).toUpperCase() }}
            </div>
            <div class="min-w-0 flex-1">
              <p class="text-sm font-bold text-white truncate">{{ authStore.user?.nama || authStore.user?.username }}</p>
              <p class="text-[10px] text-slate-500 uppercase tracking-widest font-black leading-none mt-0.5">{{ userRole }}</p>
            </div>
          </div>
          
          <button 
            @click="triggerLogout" 
            class="w-10 h-10 flex-shrink-0 flex items-center justify-center rounded-xl bg-rose-500/10 hover:bg-rose-500 text-rose-500 hover:text-white transition-all duration-300 group shadow-lg shadow-rose-500/5"
            title="Keluar dari Aplikasi"
          >
            <svg class="w-5 h-5 group-hover:scale-110 transition-transform" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1"></path></svg>
          </button>
        </div>
      </div>
    </aside>

    <!-- Main Content Wrapper -->
    <div class="flex-1 flex flex-col h-full overflow-hidden">
      
      <!-- Top header -->
      <header class="bg-white/80 backdrop-blur-md border-b border-slate-200 sticky top-0 z-40 flex-shrink-0">
        <div class="flex items-center justify-between px-4 sm:px-6 py-4">
          
          <div class="flex items-center gap-4">
            <button @click="toggleSidebar" class="lg:hidden p-2 -ml-2 text-slate-600 hover:bg-slate-100 rounded-lg transition-colors">
              <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h7"></path></svg>
            </button>
            <h2 class="text-xl font-bold text-slate-800 hidden sm:block italic tracking-tight">
              CBT <span class="text-blue-600">Sytem</span>
            </h2>
          </div>

          <!-- Real-time Clock (Super Minimalist) -->
          <div class="flex items-center gap-2 px-3 py-1 bg-slate-50/50 rounded-xl border border-transparent hover:border-slate-100 transition-all">
            <div class="flex flex-col items-end">
              <span class="text-[8px] uppercase tracking-widest font-black text-slate-400 leading-none">Server Time</span>
              <span class="text-sm font-black tracking-tight text-slate-600 font-mono">{{ currentTime }}</span>
            </div>
            <svg class="w-4 h-4 text-blue-500/50" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"></path></svg>
          </div>

        </div>
      </header>

      <!-- Page Content Area -->
      <main class="flex-1 overflow-y-auto bg-slate-50 p-4 sm:p-8 custom-scrollbar">
        <div class="max-w-7xl mx-auto">
          <router-view v-slot="{ Component }">
            <transition name="page" mode="out-in">
              <component :is="Component" />
            </transition>
          </router-view>
        </div>
      </main>

    </div>

    <!-- Modal Konfirmasi Logout -->
    <ConfirmModal 
      :show="showLogoutConfirm"
      title="Konfirmasi Logout"
      message="Apakah Anda yakin ingin mengakhiri sesi ini?"
      confirmText="Ya, Keluar"
      @close="showLogoutConfirm = false"
      @confirm="handleLogout"
    />
  </div>
</template>

<style>
/* Sidebar Link Base Styles */
.nav-link {
  @apply flex items-center gap-3 px-4 py-3 rounded-2xl text-[13px] font-bold text-slate-400 transition-all duration-300 hover:text-white hover:bg-white/5;
}

.nav-link.active {
  @apply bg-blue-600/10 text-blue-400 shadow-[inset_4px_0_0_0_rgba(59,130,246,1)];
}

.nav-link.active .nav-icon-box {
  @apply text-blue-400;
}

.nav-icon-box {
  @apply flex items-center justify-center transition-colors;
}

/* Custom Scrollbar */
.custom-scrollbar::-webkit-scrollbar {
  width: 6px;
}
.custom-scrollbar::-webkit-scrollbar-track {
  background: transparent;
}
.custom-scrollbar::-webkit-scrollbar-thumb {
  @apply bg-slate-200 rounded-full hover:bg-slate-300;
}

/* Page Transitions */
.page-enter-active,
.page-leave-active {
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.page-enter-from {
  opacity: 0;
  transform: translateX(20px);
}

.page-leave-to {
  opacity: 0;
  transform: translateX(-20px);
}

.scrollbar-hide::-webkit-scrollbar {
  display: none;
}
</style>
