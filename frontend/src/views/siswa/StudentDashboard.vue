<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../../store/auth'
import { useAlertStore } from '../../store/alert'
import BaseModal from '../../components/BaseModal.vue'

const router = useRouter()
const authStore = useAuthStore()
const alertStore = useAlertStore()

const jadwals = ref([])
const isLoading = ref(true)
const showTokenModal = ref(false)
const selectedJadwal = ref(null)
const tokenInput = ref('')
const isValidating = ref(false)

const fetchData = async () => {
  isLoading.value = true
  try {
    const res = await fetch('/api/siswa/jadwal', {
      headers: { 'Authorization': `Bearer ${authStore.token}` }
    })
    if (res.status === 401) {
      authStore.logout()
      router.push({ name: 'StudentLogin' })
      alertStore.showAlert("Sesi Anda berakhir atau akun tidak ditemukan. Silakan login kembali.", "error")
      return
    }

    const d = await res.json()
    if (res.ok) {
      jadwals.value = d.data.items || []
    } else {
      alertStore.showAlert(d.error || "Gagal memuat jadwal", "error")
    }
  } catch (e) {
    alertStore.showAlert("Koneksi gagal", "error")
  } finally {
    isLoading.value = false
  }
}

const handleStartClick = (j) => {
  if (j.status_pengerjaan === 'Selesai') {
    alertStore.showAlert("Anda sudah menyelesaikan ujian ini", "warning")
    return
  }
  if (j.status !== 'Berlangsung') {
    alertStore.showAlert("Ujian belum dimulai", "info")
    return
  }
  selectedJadwal.value = j
  tokenInput.value = ''
  showTokenModal.value = true
}

const validateToken = async () => {
  if (!tokenInput.value) return
  
  isValidating.value = true
  try {
    const res = await fetch(`/api/siswa/validate?jadwalId=${selectedJadwal.value.id}`, {
      method: 'POST',
      headers: { 
        'Authorization': `Bearer ${authStore.token}`,
        'Content-Type': 'application/json' 
      },
      body: JSON.stringify({ token: tokenInput.value.trim().replace(/\s+/g, '').toUpperCase() })
    })
    
    if (res.status === 401) {
      await authStore.logout()
      router.push({ name: 'StudentLogin' })
      alertStore.showAlert("Sesi Anda berakhir atau akun tidak ditemukan. Silakan login kembali.", "error")
      return
    }

    const d = await res.json()
    if (res.ok) {
      // Jika ini adalah sesi baru (Reset Total), bersihkan sampah LocalStorage segera!
      if (d.data.is_new_session) {
        const key = `answers_${d.data.peserta_id}`
        console.log("!!! RESET TOTAL DETECTED !!!", "Cleaning key:", key)
        localStorage.removeItem(key)
        // Pastikan memori juga bersih jika ada state tersisa
        const examStore = (await import('../../store/exam')).useExamStore()
        examStore.clearStore()
      }

      // Simpan data sesi ujian
      sessionStorage.setItem('exam_session', JSON.stringify(d.data))
      showTokenModal.value = false
      alertStore.showAlert("Token Valid! Memulai ujian...", "success")
      router.push({ name: 'StudentPengerjaan', params: { jadwalId: selectedJadwal.value.id } })
    } else {
      alertStore.showAlert(d.error || "Token salah", "error")
    }
  } catch (e) {
    alertStore.showAlert("Gagal validasi token", "error")
  } finally {
    isValidating.value = false
  }
}

const formatDate = (dateStr) => {
  if (!dateStr) return '-'
  
  if (dateStr.includes('-') && !dateStr.includes('T')) {
    const [y, m, d] = dateStr.split('-')
    const date = new Date(y, m - 1, d)
    return date.toLocaleDateString('id-ID', {
      weekday: 'long',
      year: 'numeric',
      month: 'long',
      day: 'numeric'
    })
  }

  const date = new Date(dateStr)
  return isNaN(date.getTime()) ? dateStr : date.toLocaleDateString('id-ID', {
    weekday: 'long',
    year: 'numeric',
    month: 'long',
    day: 'numeric'
  })
}

const formatTime = (timeStr) => {
  return timeStr ? timeStr.substring(0, 5) : '--:--'
}

onMounted(fetchData)
</script>

<template>
  <div class="min-h-screen bg-slate-50 font-sans text-slate-900">
    <div class="max-w-[1200px] mx-auto px-6 py-8">
      <!-- Info Data Diri (MODERN) -->
      <div class="mb-10 grid grid-cols-1 lg:grid-cols-3 gap-6">
        <div class="lg:col-span-2 bg-gradient-to-br from-indigo-600 to-blue-700 p-8 lg:p-10 rounded-3xl text-white shadow-xl shadow-indigo-100 flex flex-col sm:flex-row items-center gap-8 relative overflow-hidden group">
          <div class="absolute -right-10 -top-10 w-40 h-40 bg-white/10 rounded-full blur-3xl group-hover:scale-150 transition-transform duration-1000"></div>
          
          <div class="w-24 h-24 bg-white/20 rounded-2xl flex items-center justify-center text-4xl font-black backdrop-blur-md flex-shrink-0 shadow-inner">
            {{ authStore.user?.nama?.charAt(0) }}
          </div>
          
          <div class="text-center sm:text-left z-10">
            <p class="text-[9px] font-black text-indigo-200 uppercase tracking-[0.3em] mb-2">Selamat Datang,</p>
            <h2 class="text-3xl font-black tracking-tight mb-4">{{ authStore.user?.nama }}</h2>
            
            <div class="flex flex-wrap justify-center sm:justify-start gap-3">
              <div class="px-3 py-1.5 bg-white/10 rounded-xl flex items-center gap-2 border border-white/10">
                <span class="text-[8px] font-black uppercase opacity-60">NISN</span>
                <span class="text-xs font-bold">{{ authStore.user?.username }}</span>
              </div>
              <div class="px-3 py-1.5 bg-white/10 rounded-xl flex items-center gap-2 border border-white/10">
                <span class="text-[8px] font-black uppercase opacity-60">Kelas</span>
                <span class="text-xs font-bold">{{ authStore.user?.kelas_nama || authStore.user?.kelas }}</span>
              </div>
              <div class="px-3 py-1.5 bg-white/10 rounded-xl flex items-center gap-2 border border-white/10">
                <span class="text-[8px] font-black uppercase opacity-60">Sesi</span>
                <span class="text-xs font-bold">{{ authStore.user?.sesi }}</span>
              </div>
            </div>
          </div>
        </div>
        
        <div class="bg-white p-8 rounded-3xl border border-slate-100 shadow-sm flex flex-col justify-center relative overflow-hidden">
           <div class="absolute -right-4 -bottom-4 w-20 h-20 bg-indigo-50 rounded-full opacity-50 blur-xl"></div>
           <h3 class="text-xl font-black text-slate-800 relative z-10">Status Ujian</h3>
           <p class="text-slate-500 mt-2 text-sm font-medium leading-relaxed relative z-10">
            Pastikan koneksi internet Anda stabil sebelum memulai pengerjaan soal.
           </p>
           <div class="mt-6 flex items-center gap-2 relative z-10">
             <div class="w-2 h-2 bg-emerald-500 rounded-full animate-pulse"></div>
             <span class="text-[10px] font-black text-slate-400 uppercase tracking-widest">Sistem Siap</span>
           </div>
        </div>
      </div>

      <!-- Section Title -->
      <div class="mb-8 flex items-center gap-3">
        <div class="w-1.5 h-5 bg-indigo-600 rounded-full"></div>
        <h2 class="text-lg font-black text-slate-800 uppercase tracking-tight">Jadwal Ujian Hari Ini</h2>
      </div>

      <!-- Loading State -->
      <div v-if="isLoading" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        <div v-for="i in 3" :key="i" class="h-64 bg-white rounded-3xl border border-slate-100 animate-pulse"></div>
      </div>

      <!-- Empty State -->
      <div v-else-if="jadwals.length === 0" class="text-center py-24 bg-white rounded-3xl border border-slate-100 shadow-sm">
        <div class="w-20 h-20 bg-slate-50 rounded-2xl flex items-center justify-center mx-auto mb-6">
          <svg class="w-10 h-10 text-slate-200" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"></path></svg>
        </div>
        <h3 class="text-lg font-bold text-slate-800">Tidak ada ujian</h3>
        <p class="text-slate-400 mt-1 text-sm">Belum ada jadwal ujian yang tersedia saat ini.</p>
      </div>

      <!-- Jadwal List -->
      <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        <div v-for="j in jadwals" :key="j.id" 
          :class="[
            j.status_pengerjaan === 'Selesai' ? 'opacity-70 grayscale-[0.5]' : '',
            j.status === 'Berlangsung' ? 'border-indigo-200 ring-4 ring-indigo-50' : 'border-slate-100'
          ]"
          class="bg-white rounded-3xl p-7 border transition-all hover:shadow-2xl hover:-translate-y-1 group relative flex flex-col">
          
          <div class="flex justify-between items-start mb-6">
            <span :class="{
              'bg-amber-50 text-amber-600 border-amber-100': j.status === 'Belum Mulai',
              'bg-indigo-600 text-white border-indigo-600 shadow-lg shadow-indigo-100': j.status === 'Berlangsung',
              'bg-slate-100 text-slate-400 border-slate-200': j.status === 'Selesai'
            }" class="px-4 py-1.5 rounded-xl text-[9px] font-black uppercase tracking-widest border">
              {{ j.status }}
            </span>
            <div v-if="j.status_pengerjaan === 'Selesai'" class="w-8 h-8 bg-emerald-50 text-emerald-500 rounded-lg flex items-center justify-center">
              <svg class="w-5 h-5" fill="currentColor" viewBox="0 0 20 20"><path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd"></path></svg>
            </div>
          </div>

          <h3 class="text-lg font-black text-slate-800 leading-tight mb-2 group-hover:text-indigo-600 transition-colors">
            {{ j.bank_soal?.judul_bank_soal }}
          </h3>
          <p class="text-[9px] font-black text-slate-400 uppercase tracking-widest mb-6">{{ j.bank_soal?.mapel?.nama_mapel }}</p>

          <div class="space-y-3 mb-8 flex-1">
            <div class="flex items-center gap-3 text-xs text-slate-600 font-bold">
              <div class="w-7 h-7 bg-slate-50 rounded-lg flex items-center justify-center text-slate-300">
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"></path></svg>
              </div>
              {{ formatDate(j.tanggal_ujian) }}
            </div>
            <div class="flex items-center gap-3 text-xs text-slate-600 font-bold">
              <div class="w-7 h-7 bg-slate-50 rounded-lg flex items-center justify-center text-slate-300">
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"></path></svg>
              </div>
              {{ formatTime(j.waktu_mulai) }} • {{ j.durasi_menit }} Menit
            </div>
          </div>

          <button 
            @click="handleStartClick(j)"
            :disabled="j.status === 'Selesai' || j.status_pengerjaan === 'Selesai'"
            :class="[
              j.status === 'Berlangsung' && j.status_pengerjaan !== 'Selesai' 
              ? 'bg-indigo-600 hover:bg-indigo-700 text-white shadow-xl shadow-indigo-100' 
              : 'bg-slate-100 text-slate-300 cursor-not-allowed'
            ]"
            class="w-full py-4 rounded-2xl font-black text-[10px] uppercase tracking-widest transition-all active:scale-95"
          >
            {{ j.status_pengerjaan === 'Selesai' ? 'SUDAH SELESAI' : (j.status === 'Berlangsung' ? 'MULAI UJIAN' : 'BELUM TERSEDIA') }}
          </button>
        </div>
      </div>
    </div>

    <!-- Modal Token (STYLISH) -->
    <BaseModal 
      :show="showTokenModal" 
      title="Validasi Token Keamanan"
      confirmText="Masuk Ruang Ujian"
      :isLoading="isValidating"
      @close="showTokenModal = false"
      @confirm="validateToken"
    >
      <div class="p-2 text-center">
        <div class="w-16 h-16 bg-indigo-50 text-indigo-600 rounded-2xl flex items-center justify-center mx-auto mb-6">
          <svg class="w-8 h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 7a2 2 0 012 2m4 0a6 6 0 01-7.743 5.743L11 17H9v2H7v2H4a1 1 0 01-1-1v-2.586a1 1 0 01.293-.707l5.964-5.964A6 6 0 1121 9z"></path></svg>
        </div>
        <p class="text-sm font-medium text-slate-500 mb-8 leading-relaxed px-4">
          Masukkan kode token dari pengawas untuk mengonfirmasi identitas Anda di mata pelajaran 
          <span class="font-black text-slate-800">{{ selectedJadwal?.bank_soal?.judul_bank_soal }}</span>.
        </p>
        <div class="relative group">
          <label class="absolute -top-2.5 left-6 px-2 bg-white text-[9px] font-black text-indigo-500 uppercase tracking-widest z-10">TOKEN UJIAN</label>
          <input 
            v-model="tokenInput"
            type="text" 
            placeholder="......"
            class="w-full px-6 py-5 bg-slate-50 border-2 border-slate-100 rounded-2xl text-center font-black text-3xl uppercase tracking-[0.6em] focus:border-indigo-600 focus:bg-white transition-all outline-none placeholder-slate-200"
            maxlength="10"
          >
        </div>
      </div>
    </BaseModal>
  </div>
</template>
