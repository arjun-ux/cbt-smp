<script setup>
import { ref, onMounted, computed } from 'vue'
import { useAuthStore } from '../../store/auth'

const authStore = useAuthStore()
const stats = ref({
  guru: 0,
  siswa: 0,
  activeExams: 0
})
const activeExams = ref([])
const apiPrefix = computed(() => authStore.user?.role === 'admin' ? '/api/admin' : '/api/guru')

const fetchData = async () => {
  if (!authStore.token) return

  try {
    const headers = { 'Authorization': `Bearer ${authStore.token}` }
    const isGuru = authStore.user?.role === 'guru'
    
    // Fetch Jadwal (Disesuaikan per role)
    const urlJadwal = isGuru ? '/api/guru/pengawas-jadwal' : '/api/admin/jadwal'
    const resJadwal = await fetch(urlJadwal, { headers })
    if (resJadwal.ok) {
      const dataJadwal = await resJadwal.json()
      const listJadwal = dataJadwal.data || []
      activeExams.value = listJadwal.filter(j => j.status === 'Berlangsung')
      stats.value.activeExams = activeExams.value.length
    }

    if (!isGuru) {
      // Fetch Guru (Hanya Admin)
      const resGuru = await fetch('/api/admin/guru', { headers })
      if (resGuru.ok) {
        const dataGuru = await resGuru.json()
        stats.value.guru = (dataGuru.data || dataGuru).length
      }

      // Fetch Siswa (Hanya Admin)
      const resSiswa = await fetch('/api/admin/siswa', { headers })
      if (resSiswa.ok) {
        const dataSiswa = await resSiswa.json()
        stats.value.siswa = (dataSiswa.data || dataSiswa).length
      }
    }

  } catch (error) {
    console.error("Dashboard Error:", error)
  }
}

onMounted(fetchData)
</script>

<template>
  <div class="space-y-6 pb-12">
    <!-- Header Page (Seragam dengan Master Data) -->
    <div class="bg-white p-6 rounded-2xl shadow-sm border border-slate-100 flex flex-col md:flex-row md:items-center justify-between gap-4">
      <div>
        <h3 class="text-xl font-bold text-slate-800">Ringkasan Dashboard</h3>
        <p class="text-sm text-slate-500 mt-1">Selamat datang kembali, {{ authStore.user?.nama }}. Pantau aktivitas ujian hari ini.</p>
      </div>
      <div class="flex-shrink-0">
        <div class="w-12 h-12 bg-blue-50 text-blue-600 rounded-xl flex items-center justify-center border border-blue-100">
          <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6"></path></svg>
        </div>
      </div>
    </div>

    <!-- Stats Grid (Premium but Uniform) -->
    <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
      <template v-if="authStore.user?.role === 'admin'">
        <div class="bg-white p-6 rounded-2xl shadow-sm border border-slate-100 border-l-4 border-l-blue-500 group hover:shadow-md transition-all">
          <div class="flex items-center justify-between mb-2">
            <p class="text-[10px] font-bold text-slate-400 uppercase tracking-widest">Total Guru</p>
            <svg class="w-5 h-5 text-blue-500 opacity-20 group-hover:opacity-100 transition-opacity" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z"></path></svg>
          </div>
          <p class="text-3xl font-bold text-slate-800">{{ stats.guru }}</p>
        </div>
        
        <div class="bg-white p-6 rounded-2xl shadow-sm border border-slate-100 border-l-4 border-l-emerald-500 group hover:shadow-md transition-all">
          <div class="flex items-center justify-between mb-2">
            <p class="text-[10px] font-bold text-slate-400 uppercase tracking-widest">Total Siswa</p>
            <svg class="w-5 h-5 text-emerald-500 opacity-20 group-hover:opacity-100 transition-opacity" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z"></path></svg>
          </div>
          <p class="text-3xl font-bold text-slate-800">{{ stats.siswa }}</p>
        </div>
      </template>
 
      <div :class="authStore.user?.role === 'admin' ? '' : 'md:col-span-3'" class="bg-white p-6 rounded-2xl shadow-sm border border-slate-100 border-l-4 border-l-violet-500 group hover:shadow-md transition-all">
        <div class="flex items-center justify-between mb-2">
          <p class="text-[10px] font-bold text-slate-400 uppercase tracking-widest">Ujian Aktif</p>
          <div class="relative flex h-2 w-2">
            <span class="animate-ping absolute inline-flex h-full w-full rounded-full bg-violet-400 opacity-75"></span>
            <span class="relative inline-flex rounded-full h-2 w-2 bg-violet-500"></span>
          </div>
        </div>
        <p class="text-3xl font-bold text-slate-800">{{ stats.activeExams }}</p>
      </div>
    </div>

    <!-- Active Exams & Status -->
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
      <div class="bg-white p-6 rounded-2xl shadow-sm border border-slate-100 flex flex-col h-full">
        <div class="flex items-center justify-between mb-6">
          <h2 class="text-lg font-bold text-slate-800">Ujian Berlangsung</h2>
          <span class="px-2 py-1 bg-blue-50 text-blue-600 rounded text-[10px] font-bold uppercase tracking-widest border border-blue-100">Live</span>
        </div>
        
        <div v-if="activeExams.length > 0" class="space-y-3 flex-1">
          <div v-for="ex in activeExams" :key="ex.id" class="flex items-center justify-between p-4 bg-slate-50 rounded-xl border border-slate-100 group hover:border-blue-200 transition-all">
            <div class="flex flex-col">
              <span class="font-bold text-slate-700">{{ ex.bank_soal?.judul_bank_soal }}</span>
              <div class="flex items-center gap-2">
                <span class="text-[10px] font-bold text-slate-400 uppercase tracking-widest">{{ ex.bank_soal?.mapel?.nama_mapel }}</span>
                <span class="w-1 h-1 bg-slate-300 rounded-full"></span>
                <span class="text-[10px] font-bold text-slate-400 uppercase tracking-widest">{{ ex.ruang?.nama_ruang || 'N/A' }}</span>
              </div>
            </div>
            <button 
              @click="$router.push({ name: authStore.user?.role === 'admin' ? 'AdminMonitor' : 'GuruMonitor', params: { jadwalId: String(ex.id) } })"
              class="px-4 py-2 bg-white text-blue-600 border border-slate-200 rounded-lg text-[10px] font-black uppercase tracking-widest hover:bg-blue-600 hover:text-white transition-all shadow-sm active:scale-95"
            >
              Monitor
            </button>
          </div>
        </div>
        <div v-else class="flex-1 flex flex-col items-center justify-center py-10 opacity-50">
          <p class="text-xs font-bold text-slate-400 uppercase tracking-widest italic">Tidak ada ujian aktif</p>
        </div>
      </div>

      <div class="bg-white p-6 rounded-2xl shadow-sm border border-slate-100 flex flex-col h-full">
        <h2 class="text-lg font-bold text-slate-800 mb-6">Status Sistem</h2>
        <div class="space-y-3">
          <div class="flex items-center justify-between p-4 bg-slate-50 rounded-xl border border-slate-100">
            <div class="flex items-center gap-3">
              <div class="w-8 h-8 bg-white rounded-lg flex items-center justify-center shadow-sm text-emerald-600 border border-emerald-100">
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 7v10c0 2.21 3.582 4 8 4s8-1.79 8-4V7M4 7c0 2.21 3.582 4 8 4s8-1.79 8-4M4 7c0-2.21 3.582-4 8-4s8 1.79 8 4"></path></svg>
              </div>
              <span class="text-sm font-bold text-slate-600">Database Engine</span>
            </div>
            <span class="text-[10px] font-black text-emerald-600 uppercase tracking-widest bg-emerald-50 px-2 py-1 rounded border border-emerald-100">Connected</span>
          </div>
          <div class="flex items-center justify-between p-4 bg-slate-50 rounded-xl border border-slate-100">
            <div class="flex items-center gap-3">
              <div class="w-8 h-8 bg-white rounded-lg flex items-center justify-center shadow-sm text-blue-600 border border-blue-100">
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8.111 16.404a5.5 5.5 0 017.778 0M12 20h.01m-7.08-7.071a10 10 0 0114.142 0M2.435 8.929a16 16 0 0123.13 0"></path></svg>
              </div>
              <span class="text-sm font-bold text-slate-600">Server Network</span>
            </div>
            <span class="text-[10px] font-black text-blue-600 uppercase tracking-widest bg-blue-50 px-2 py-1 rounded border border-blue-100">Active</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
