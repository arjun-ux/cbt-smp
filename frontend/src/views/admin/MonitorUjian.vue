<script setup>
import { ref, onMounted, onUnmounted, computed } from 'vue'
import { useRoute } from 'vue-router'
import { useAuthStore } from '../../store/auth'
import { useAlertStore } from '../../store/alert'
import ConfirmModal from '../../components/ConfirmModal.vue'

const route = useRoute()
const authStore = useAuthStore()
const alertStore = useAlertStore()

const participants = ref([])
const jadwal = ref(null)
const isLoading = ref(true)
const pollingInterval = ref(null)
const showForceSubmitModal = ref(false)
const selectedPeserta = ref(null)
const isActionLoading = ref(false)
const searchQuery = ref('')

const showResetConfirm = ref(false)
const resetTargetId = ref(null)
const isResetting = ref(false)


const apiPrefix = computed(() => authStore.user?.role === 'admin' ? '/api/admin' : '/api/guru')

const isAuthorized = computed(() => {
  if (authStore.user?.role === 'admin') return true
  if (authStore.user?.role === 'guru' && jadwal.value?.pengawas_id === authStore.user?.guru_id) return true
  return false
})

const filteredParticipants = computed(() => {
  const list = participants.value || []
  if (!searchQuery.value) return list
  const q = searchQuery.value.toLowerCase()
  return list.filter(p => 
    p.nama_siswa?.toLowerCase().includes(q) || 
    p.kelas?.toString().toLowerCase().includes(q)
  )
})

const fetchData = async () => {
  try {
    const res = await fetch(`${apiPrefix.value}/monitor/${route.params.jadwalId}`, {
      headers: { 'Authorization': `Bearer ${authStore.token}` }
    })
    const data = await res.json()
    if (res.ok) {
      participants.value = data.data || []
    }
  } catch (error) {
    console.error("Gagal polling data monitoring")
  } finally {
    isLoading.value = false
  }
}

const fetchJadwalDetail = async () => {
  try {
    const res = await fetch(`${apiPrefix.value}/jadwal/${route.params.jadwalId}`, {
      headers: { 'Authorization': `Bearer ${authStore.token}` }
    })
    const data = await res.json()
    if (res.ok) {
      jadwal.value = data.data
    }
  } catch (error) {
    console.error("Gagal mengambil detail jadwal")
  }
}

const confirmForceSubmit = (peserta) => {
  selectedPeserta.value = peserta
  showForceSubmitModal.value = true
}

const handleForceSubmit = async () => {
  if (!selectedPeserta.value) return
  isActionLoading.value = true
  try {
    const res = await fetch(`${apiPrefix.value}/monitor/force-submit/${selectedPeserta.value.id}`, {
      method: 'POST',
      headers: { 'Authorization': `Bearer ${authStore.token}` }
    })
    if (res.ok) {
      alertStore.showAlert("Berhasil menghentikan ujian peserta", "success")
      fetchData()
    } else {
      alertStore.showAlert("Gagal menghentikan ujian", "error")
    }
  } catch (error) {
    alertStore.showAlert("Kesalahan koneksi", "error")
  } finally {
    isActionLoading.value = false
    showForceSubmitModal.value = false
  }
}

const handleUnblock = async (pesertaId) => {
  try {
    const res = await fetch(`${apiPrefix.value}/monitor/unblock/${pesertaId}`, {
      method: 'POST',
      headers: { 'Authorization': `Bearer ${authStore.token}` }
    })
    if (res.ok) {
      alertStore.showAlert("Blokir berhasil dibuka", "success")
      fetchData()
    } else {
      alertStore.showAlert("Gagal membuka blokir", "error")
    }
  } catch (error) {
    alertStore.showAlert("Kesalahan koneksi", "error")
  }
}



const confirmResetSesi = (pesertaId) => {
  resetTargetId.value = pesertaId
  showResetConfirm.value = true
}

const executeResetSesi = async () => {
  if (!resetTargetId.value) return
  
  isResetting.value = true
  try {
    const res = await fetch(`${apiPrefix.value}/monitor/reset-sesi/${resetTargetId.value}`, {
      method: 'POST',
      headers: { 'Authorization': `Bearer ${authStore.token}` }
    })
    if (res.ok) {
      alertStore.showAlert("Sesi berhasil di-reset", "success")
      showResetConfirm.value = false
      fetchData()
    } else {
      alertStore.showAlert("Gagal reset sesi", "error")
    }
  } catch (error) {
    alertStore.showAlert("Kesalahan koneksi", "error")
  } finally {
    isResetting.value = false
  }
}

const formatSisaWaktu = (detik) => {
  if (detik <= 0) return "Habis"
  const m = Math.floor(detik / 60)
  const s = detik % 60
  return `${m}m ${s}s`
}

onMounted(() => {
  fetchData()
  fetchJadwalDetail()
  // Polling setiap 10 detik
  pollingInterval.value = setInterval(fetchData, 10000)
})

onUnmounted(() => {
  if (pollingInterval.value) clearInterval(pollingInterval.value)
})
</script>

<template>
  <div class="space-y-6 pb-12">
    <!-- Header Page (Seragam dengan Master Data) -->
    <div class="bg-white p-6 rounded-2xl shadow-sm border border-slate-100 flex flex-col md:flex-row md:items-center justify-between gap-4">
      <div>
        <div class="flex items-center gap-3">
          <h3 class="text-xl font-bold text-slate-800">Monitor Ujian Real-time</h3>
          <div class="flex items-center gap-2 px-2 py-0.5 bg-blue-50 text-blue-600 rounded border border-blue-100">
            <span class="relative flex h-1.5 w-1.5">
              <span class="animate-ping absolute inline-flex h-full w-full rounded-full bg-blue-400 opacity-75"></span>
              <span class="relative inline-flex rounded-full h-1.5 w-1.5 bg-blue-500"></span>
            </span>
            <span class="text-[9px] font-black uppercase tracking-widest">Live</span>
          </div>
        </div>
        <p class="text-sm text-slate-500 mt-1">
          {{ jadwal?.bank_soal?.judul_bank_soal || 'Memuat...' }} <span class="mx-2 text-slate-300">|</span> {{ jadwal?.bank_soal?.mapel?.nama_mapel || '' }}
        </p>
      </div>
      <div class="flex flex-col sm:flex-row items-center gap-3">
        <div class="relative w-full sm:w-64">
          <svg class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"></path></svg>
          <input 
            v-model="searchQuery"
            type="text" 
            placeholder="Cari nama atau kelas..." 
            class="w-full pl-10 pr-4 py-2 bg-slate-50 border border-slate-200 rounded-xl text-sm focus:ring-2 focus:ring-blue-500 transition-all outline-none"
          >
        </div>
        <button 
          @click="fetchData" 
          :disabled="isLoading"
          class="flex items-center gap-2 bg-slate-100 hover:bg-slate-200 text-slate-700 px-4 py-2 rounded-xl text-xs font-bold transition-all active:scale-95 disabled:opacity-50"
        >
          <svg :class="{ 'animate-spin': isLoading }" class="w-3.5 h-3.5 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"></path></svg>
          Refresh
        </button>
      </div>
    </div>

    <!-- Stats Summary & Token (Premium but Uniform) -->
    <div class="grid grid-cols-1 md:grid-cols-4 gap-4">
      <div class="bg-slate-800 p-5 rounded-2xl shadow-lg text-white relative overflow-hidden group">
        <div class="absolute -right-2 -top-2 w-16 h-16 bg-white/5 rounded-full blur-xl group-hover:scale-150 transition-transform"></div>
        <p class="text-[9px] font-black text-slate-400 uppercase tracking-widest mb-2 flex items-center gap-2">
          <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 7a2 2 0 012 2m4 0a6 6 0 01-7.743 5.743L11 17H9v2H7v2H4a1 1 0 01-1-1v-2.586a1 1 0 01.293-.707l5.964-5.964A6 6 0 1121 9z"></path></svg>
          Token Ujian
        </p>
        <div class="flex items-center justify-between">
          <p class="text-2xl font-black tracking-[0.2em] font-mono leading-none">{{ jadwal?.token_ujian || '------' }}</p>
          <button @click="fetchJadwalDetail" class="p-1.5 bg-white/10 hover:bg-white/20 rounded-lg transition-colors">
            <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"></path></svg>
          </button>
        </div>
      </div>

      <div class="bg-white p-5 rounded-2xl border border-slate-100 shadow-sm flex flex-col justify-center">
        <p class="text-[9px] font-black text-slate-400 uppercase tracking-widest mb-1">Total Peserta</p>
        <p class="text-3xl font-black text-slate-800 leading-none">{{ (participants || []).length }}</p>
      </div>
      <div class="bg-white p-5 rounded-2xl border border-slate-100 shadow-sm flex flex-col justify-center">
        <p class="text-[9px] font-black text-slate-400 uppercase tracking-widest mb-1">Aktif</p>
        <div class="flex items-center gap-2">
          <p class="text-3xl font-black text-blue-600 leading-none">{{ (participants || []).filter(p => p.status_ujian === 'Sedang Mengerjakan').length }}</p>
          <span class="w-1.5 h-1.5 bg-blue-500 rounded-full animate-pulse"></span>
        </div>
      </div>
      <div class="bg-white p-5 rounded-2xl border border-slate-100 shadow-sm flex flex-col justify-center">
        <p class="text-[9px] font-black text-slate-400 uppercase tracking-widest mb-1">Selesai</p>
        <p class="text-3xl font-black text-emerald-600 leading-none">{{ (participants || []).filter(p => p.status_ujian === 'Selesai').length }}</p>
      </div>
    </div>

    <!-- Table Monitor -->
    <div class="bg-white rounded-[2.5rem] shadow-xl shadow-slate-200/50 border border-slate-100 overflow-hidden">
      <div class="overflow-x-auto">
        <table class="w-full text-left border-collapse">
          <thead>
            <tr class="bg-slate-50/50 border-b border-slate-100">
              <th class="px-8 py-6 text-[10px] font-black text-slate-400 uppercase tracking-widest">Peserta</th>
              <th class="px-6 py-6 text-[10px] font-black text-slate-400 uppercase tracking-widest text-center">Status</th>
              <th class="px-6 py-6 text-[10px] font-black text-slate-400 uppercase tracking-widest text-center">Progres</th>
              <th class="px-6 py-6 text-[10px] font-black text-slate-400 uppercase tracking-widest text-center">Sisa Waktu</th>
              <th class="px-6 py-6 text-[10px] font-black text-slate-400 uppercase tracking-widest text-center">Pelanggaran</th>
              <th class="px-8 py-6 text-[10px] font-black text-slate-400 uppercase tracking-widest text-right">Aksi</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-50">
            <tr v-for="p in filteredParticipants" :key="p.id" class="hover:bg-slate-50/50 transition-colors group">
              <td class="px-8 py-5">
                <div class="flex flex-col">
                  <span class="text-sm font-bold text-slate-700">{{ p.nama_siswa }}</span>
                  <span class="text-[11px] font-medium text-slate-400 mt-0.5">Kelas {{ p.kelas }}</span>
                </div>
              </td>
              <td class="px-6 py-5">
                <div class="flex items-center justify-center gap-2">
                  <span 
                    v-if="p.is_terblokir"
                    class="px-4 py-1.5 rounded-xl text-[10px] font-bold uppercase tracking-wider bg-red-600 text-white border border-red-700 animate-pulse shadow-lg shadow-red-100"
                  >
                    TERBLOKIR
                  </span>
                  <span 
                    v-else
                    :class="[
                      p.status_ujian === 'Selesai' ? 'bg-emerald-50 text-emerald-600 border-emerald-100' : 'bg-blue-50 text-blue-600 border-blue-100'
                    ]"
                    class="px-4 py-1.5 rounded-xl text-[10px] font-bold uppercase tracking-wider border"
                  >
                    {{ p.status_ujian }}
                  </span>
                </div>
              </td>
              <td class="px-6 py-5">
                <div class="flex flex-col items-center">
                  <div class="flex items-center gap-2 mb-1.5">
                    <span class="text-[11px] font-bold text-slate-700">{{ Math.round((p.progres / p.total_soal) * 100) }}%</span>
                    <span class="text-[10px] font-medium text-slate-300">({{ p.progres }}/{{ p.total_soal }})</span>
                  </div>
                  <div class="h-2 w-24 bg-slate-100 rounded-full overflow-hidden">
                    <div 
                      class="h-full bg-blue-500 rounded-full transition-all duration-1000"
                      :style="{ width: `${(p.progres / p.total_soal) * 100}%` }"
                    ></div>
                  </div>
                </div>
              </td>
              <td class="px-6 py-5 text-center">
                <span class="font-bold text-slate-600 text-xs tracking-tight bg-slate-50 px-3 py-1 rounded-lg border border-slate-100">
                  {{ formatSisaWaktu(p.sisa_waktu) }}
                </span>
              </td>
              <td class="px-6 py-5 text-center">
                <div class="flex items-center justify-center">
                  <span 
                    v-if="p.jumlah_pelanggaran > 0"
                    class="inline-flex items-center justify-center w-8 h-8 bg-rose-50 text-rose-600 rounded-2xl text-xs font-black border border-rose-100 animate-pulse shadow-sm"
                    :title="`${p.jumlah_pelanggaran} kali pindah tab`"
                  >
                    {{ p.jumlah_pelanggaran }}
                  </span>
                  <span v-else class="text-slate-200 font-bold">-</span>
                </div>
              </td>
              <td class="px-8 py-5 text-right">
                <div class="flex items-center justify-end gap-2">
                  <button 
                    v-if="isAuthorized"
                    @click="confirmResetSesi(p.id)"
                    class="w-10 h-10 flex items-center justify-center text-blue-500 hover:bg-blue-50 rounded-2xl transition-all active:scale-90 border border-transparent hover:border-blue-100"
                    title="Reset Login / Sesi"
                  >
                    <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"></path></svg>
                  </button>
                  <button 
                    v-if="isAuthorized && p.is_terblokir"
                    @click="handleUnblock(p.id)"
                    class="w-10 h-10 flex items-center justify-center text-orange-400 hover:bg-orange-50 rounded-2xl transition-all active:scale-90 border border-transparent hover:border-orange-100"
                    title="Buka Blokir"
                  >
                    <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 11V7a4 4 0 118 0m-4 8v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2z"></path></svg>
                  </button>
                  <button 
                    v-if="isAuthorized && p.status_ujian !== 'Selesai'"
                    @click="confirmForceSubmit(p)"
                    class="w-10 h-10 flex items-center justify-center text-slate-400 hover:text-rose-500 hover:bg-rose-50 rounded-2xl transition-all active:scale-90 border border-transparent hover:border-rose-100"
                    title="Hentikan Paksa Ujian"
                  >
                    <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z"></path></svg>
                  </button>
                  <div v-if="!isAuthorized" class="bg-slate-50 px-3 py-1.5 rounded-xl border border-slate-100">
                    <span class="text-[10px] font-black text-slate-400 uppercase tracking-widest">View Only</span>
                  </div>
                </div>
              </td>
            </tr>
            <tr v-if="filteredParticipants.length === 0 && !isLoading">
              <td colspan="6" class="px-8 py-24 text-center">
                <div class="flex flex-col items-center gap-4">
                  <div class="w-24 h-24 bg-slate-50 rounded-[2.5rem] flex items-center justify-center text-slate-200">
                    <svg class="w-12 h-12" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"></path></svg>
                  </div>
                  <p class="text-sm font-black text-slate-300 uppercase tracking-widest">
                    {{ searchQuery ? 'Siswa tidak ditemukan' : 'Belum ada peserta yang masuk' }}
                  </p>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Confirm Force Submit Modal -->
    <ConfirmModal 
      :show="showForceSubmitModal"
      title="Hentikan Paksa Ujian?"
      :message="`Apakah Anda yakin ingin menghentikan paksa ujian ${selectedPeserta?.nama_siswa}? Siswa tidak akan bisa melanjutkan lagi.`"
      confirmText="Ya, Hentikan"
      variant="danger"
      :isLoading="isActionLoading"
      @close="showForceSubmitModal = false"
      @confirm="handleForceSubmit"
    />
    <!-- Reset Session Confirmation -->
    <ConfirmModal 
      :show="showResetConfirm"
      title="Reset Sesi Siswa"
      message="Apakah Anda yakin ingin me-reset sesi siswa ini? Gunakan jika siswa mengalami kendala teknis (PC mati/crash)."
      :isLoading="isResetting"
      confirmText="Ya, Reset Sesi"
      variant="primary"
      @close="showResetConfirm = false"
      @confirm="executeResetSesi"
    />
  </div>
</template>
