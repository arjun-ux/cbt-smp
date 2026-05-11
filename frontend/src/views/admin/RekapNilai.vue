<script setup>
import { ref, onMounted, computed, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '../../store/auth'
import { useAlertStore } from '../../store/alert'
import BaseModal from '../../components/BaseModal.vue'
import ConfirmModal from '../../components/ConfirmModal.vue'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()
const alertStore = useAlertStore()

const results = ref([])
const jadwals = ref([])
const classes = ref([])
const isLoading = ref(false)
const selectedJadwalId = ref(route.params.jadwalId || '')
const filterKelas = ref('')
const searchQuery = ref('')

const isKoreksiModalOpen = ref(false)
const isDetailPGModalOpen = ref(false)
const selectedPeserta = ref(null)
const jawabanPeserta = ref([])
const jawabanLengkap = ref([])
const isSaving = ref(false)
const isArchiving = ref(false)
const showArchiveConfirm = ref(false)

const apiPrefix = computed(() => authStore.user?.role === 'admin' ? '/api/admin' : '/api/guru')

const currentJadwal = computed(() => {
  return jadwals.value.find(j => j.id == selectedJadwalId.value)
})

const filteredResults = computed(() => {
  const listRaw = results.value || []
  let list = Array.isArray(listRaw) ? [...listRaw] : []
  
  if (filterKelas.value) {
    list = list.filter(r => r.kelas === filterKelas.value)
  }
  if (searchQuery.value) {
    const q = searchQuery.value.toLowerCase()
    list = list.filter(r => 
      r.nama_siswa?.toLowerCase().includes(q) || 
      r.nisn?.toLowerCase().includes(q)
    )
  }
  return list
})

const stats = computed(() => {
  const activeList = filteredResults.value.filter(r => r.status_ujian === 'Selesai')
  if (activeList.length === 0) return { avg: 0, max: 0, min: 0 }
  
  const scores = activeList.map(r => r.total_nilai)
  const sum = scores.reduce((a, b) => a + b, 0)
  return {
    avg: (sum / scores.length).toFixed(1),
    max: Math.max(...scores),
    min: Math.min(...scores)
  }
})

const fetchJadwals = async () => {
  try {
    const res = await fetch(`${apiPrefix.value}/rekap-jadwal`, {
      headers: { 'Authorization': `Bearer ${authStore.token}` }
    })
    const d = await res.json()
    if (res.ok) {
      jadwals.value = d.data || []
    }
  } catch (e) {
    console.error("Gagal ambil daftar jadwal")
  }
}

const fetchData = async () => {
  if (!selectedJadwalId.value) {
    results.value = []
    return
  }
  
  isLoading.value = true
  try {
    const res = await fetch(`${apiPrefix.value}/rekap-nilai/${selectedJadwalId.value}`, {
      headers: { 'Authorization': `Bearer ${authStore.token}` }
    })
    if (res.ok) {
      const d = await res.json()
      results.value = d.data || []
      // Ambil daftar kelas unik dari hasil
      const uniqueClasses = [...new Set(results.value.map(r => r.kelas))]
      classes.value = uniqueClasses.sort()
    } else {
      results.value = []
    }
  } catch (error) {
    alertStore.showAlert("Gagal mengambil data rekap", "error")
  } finally {
    isLoading.value = false
  }
}

const openKoreksi = async (peserta) => {
  selectedPeserta.value = peserta
  isKoreksiModalOpen.value = true
  jawabanPeserta.value = []
  
  try {
    const res = await fetch(`${apiPrefix.value}/monitor/jawaban/${peserta.id}`, {
      headers: { 'Authorization': `Bearer ${authStore.token}` }
    })
    const d = await res.json()
    if (res.ok) {
      // Hanya tampilkan Essay untuk dikoreksi (case-insensitive)
      jawabanPeserta.value = d.data.filter(j => j.jenis_soal && j.jenis_soal.toUpperCase() === 'ESSAY')
    }
  } catch (e) {
    alertStore.showAlert("Gagal mengambil jawaban", "error")
  }
}

const openDetailPG = async (peserta) => {
  selectedPeserta.value = peserta
  isDetailPGModalOpen.value = true
  jawabanLengkap.value = []
  
  try {
    const res = await fetch(`${apiPrefix.value}/monitor/jawaban/${peserta.id}`, {
      headers: { 'Authorization': `Bearer ${authStore.token}` }
    })
    const d = await res.json()
    if (res.ok) {
      // Ambil seluruh jawaban untuk dianalisis (fokus PG di UI)
      jawabanLengkap.value = d.data || []
    }
  } catch (e) {
    alertStore.showAlert("Gagal mengambil detail jawaban", "error")
  }
}

const getOpsiText = (jawaban, item) => {
  if (!jawaban) return '(Tidak menjawab)'
  const key = jawaban.toUpperCase()
  if (key === 'A') return item.opsi_a
  if (key === 'B') return item.opsi_b
  if (key === 'C') return item.opsi_c
  if (key === 'D') return item.opsi_d
  return jawaban
}

const saveKoreksi = async () => {
  isSaving.value = true
  try {
    const payload = jawabanPeserta.value.map(j => ({
      soal_id: j.soal_id,
      skor: parseFloat(j.skor) || 0
    }))
    
    const res = await fetch(`${apiPrefix.value}/monitor/koreksi/${selectedPeserta.value.id}`, {
      method: 'POST',
      headers: { 
        'Authorization': `Bearer ${authStore.token}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(payload)
    })
    
    if (res.ok) {
      alertStore.showAlert("Koreksi berhasil disimpan", "success")
      isKoreksiModalOpen.value = false
      fetchData() // Refresh tabel
    } else {
      const d = await res.json()
      alertStore.showAlert(d.error || "Gagal simpan koreksi", "error")
    }
  } catch (e) {
    alertStore.showAlert("Terjadi kesalahan sistem", "error")
  } finally {
    isSaving.value = false
  }
}

const exportExcel = () => {
  if (!selectedJadwalId.value) return
  
  // Simple CSV Export
  let csv = 'No,Nama Siswa,NISN,Kelas,Nilai PG,Benar,Salah,Nilai Essay,Total Nilai,Status,Waktu Selesai\n'
  filteredResults.value.forEach((r, i) => {
    csv += `${i+1},"${r.nama_siswa}","${r.nisn}","${r.kelas}",${r.nilai_pg},${r.jumlah_benar},${r.jumlah_salah},${r.nilai_essay},${r.total_nilai},"${r.status_ujian}","${r.waktu_selesai}"\n`
  })
  
  const blob = new Blob([csv], { type: 'text/csv' })
  const url = window.URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.setAttribute('hidden', '')
  a.setAttribute('href', url)
  a.setAttribute('download', `Rekap_Nilai_${currentJadwal.value?.bank_soal?.judul_bank_soal}.csv`)
  document.body.appendChild(a)
  a.click()
  document.body.removeChild(a)
  
  alertStore.showAlert("Laporan berhasil diunduh (CSV)", "success")
}

const executeArchive = async () => {
  if (!selectedJadwalId.value) return
  
  isArchiving.value = true
  try {
    const res = await fetch(`/api/admin/jadwal/${selectedJadwalId.value}/archive`, {
      method: 'POST',
      headers: { 'Authorization': `Bearer ${authStore.token}` }
    })
    const data = await res.json()
    if(res.ok) {
      alertStore.showAlert("Nilai berhasil diarsipkan secara permanen", "success")
      showArchiveConfirm.value = false
      await fetchJadwals() // Refresh list jadwal agar statusnya update
      fetchData()
    } else {
      alertStore.showAlert(data.error || "Gagal mengarsipkan", "error")
    }
  } catch (error) {
    alertStore.showAlert("Kesalahan sistem saat mengarsipkan", "error")
  } finally {
    isArchiving.value = false
  }
}

const formatDate = (dateStr) => {
  if (!dateStr) return '-'
  const d = new Date(dateStr)
  return d.toLocaleDateString('id-ID', { day: '2-digit', month: '2-digit', year: 'numeric' })
}

watch(selectedJadwalId, (newVal) => {
  if (newVal) {
    fetchData()
    const name = authStore.user?.role === 'admin' ? 'AdminRekap' : 'GuruRekap'
    router.replace({ name, params: { jadwalId: newVal } })
  } else {
    results.value = []
    classes.value = []
    router.replace({ name: authStore.user?.role === 'admin' ? 'AdminRekap' : 'GuruRekap' })
  }
})

onMounted(async () => {
  await fetchJadwals()
  if (selectedJadwalId.value) fetchData()
})
</script>

<template>
  <div class="space-y-6 pb-12">
    <!-- Header & Selectors -->
    <div class="bg-white p-6 rounded-2xl shadow-sm border border-slate-100">
      <div class="flex flex-col lg:flex-row lg:items-center justify-between gap-6">
        <div class="flex items-center gap-4">
          <div>
            <h3 class="text-xl font-bold text-slate-800">Laporan & Nilai</h3>
            <p class="text-sm text-slate-500 mt-0.5">Pantau hasil ujian dan lakukan koreksi essay</p>
          </div>
        </div>

        <div class="flex flex-col sm:flex-row items-center gap-3">
          <div class="w-full sm:w-72">
            <label class="block text-[10px] font-bold text-slate-400 uppercase tracking-widest mb-1 ml-1">Pilih Jadwal Ujian</label>
            <select 
              v-model="selectedJadwalId"
              class="w-full px-4 py-2.5 bg-slate-50 border border-slate-200 rounded-xl text-sm font-bold focus:ring-2 focus:ring-blue-500 transition-all outline-none"
            >
              <option value="">-- Pilih Ujian --</option>
              <option v-for="j in jadwals" :key="j.id" :value="j.id">
                [{{ formatDate(j.tanggal_ujian) }}] {{ j.bank_soal?.judul_bank_soal }}
              </option>
            </select>
          </div>
          
          <div class="w-full sm:w-40">
            <label class="block text-[10px] font-bold text-slate-400 uppercase tracking-widest mb-1 ml-1">Filter Kelas</label>
            <select 
              v-model="filterKelas"
              :disabled="!selectedJadwalId"
              class="w-full px-4 py-2.5 bg-slate-50 border border-slate-200 rounded-xl text-sm font-bold focus:ring-2 focus:ring-blue-500 transition-all outline-none disabled:opacity-50"
            >
              <option value="">Semua Kelas</option>
              <option v-for="c in classes" :key="c" :value="c">{{ c }}</option>
            </select>
          </div>
        </div>
      </div>
    </div>

    <template v-if="selectedJadwalId">
      <!-- Stats Cards -->
      <div class="grid grid-cols-1 md:grid-cols-4 gap-4">
        <div class="bg-white p-5 rounded-2xl border border-slate-100 shadow-sm hover:shadow-md transition-all">
          <p class="text-[10px] font-black text-slate-400 uppercase tracking-widest mb-1">Rata-rata</p>
          <p class="text-3xl font-black text-blue-600 font-mono">{{ stats.avg }}</p>
        </div>
        <div class="bg-white p-5 rounded-2xl border border-slate-100 shadow-sm hover:shadow-md transition-all">
          <p class="text-[10px] font-black text-slate-400 uppercase tracking-widest mb-1">Tertinggi</p>
          <p class="text-3xl font-black text-emerald-600 font-mono">{{ stats.max }}</p>
        </div>
        <div class="bg-white p-5 rounded-2xl border border-slate-100 shadow-sm hover:shadow-md transition-all">
          <p class="text-[10px] font-black text-slate-400 uppercase tracking-widest mb-1">Terendah</p>
          <p class="text-3xl font-black text-rose-600 font-mono">{{ stats.min }}</p>
        </div>
        <div class="bg-slate-800 p-5 rounded-2xl border border-slate-700 shadow-lg text-white">
          <p class="text-[10px] font-black text-slate-400 uppercase tracking-widest mb-1">Total Selesai</p>
          <p class="text-3xl font-black">{{ results.filter(r => r.status_ujian === 'Selesai').length }} <span class="text-xs text-slate-400">/ {{ results.length }}</span></p>
        </div>
      </div>

      <!-- Table Section -->
      <div class="bg-white rounded-2xl shadow-sm border border-slate-100 overflow-hidden">
        <div class="p-4 border-b border-slate-50 bg-slate-50/30 flex flex-col md:flex-row justify-between items-center gap-4">
          <div class="relative w-full md:w-64">
            <svg class="absolute left-4 top-1/2 -translate-y-1/2 w-4 h-4 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"></path></svg>
            <input v-model="searchQuery" type="text" placeholder="Cari nama siswa..." class="w-full pl-11 pr-4 py-2.5 bg-white border border-slate-200 rounded-xl text-sm focus:ring-2 focus:ring-blue-500 transition-all outline-none">
          </div>
          <div class="flex flex-wrap items-center gap-2">
            <button 
              v-if="selectedJadwalId"
              @click="router.push({ name: authStore.user?.role === 'admin' ? 'AdminAnalisis' : 'GuruAnalisis', params: { jadwalId: selectedJadwalId } })"
              class="flex items-center gap-2 bg-blue-50 text-blue-600 hover:bg-blue-100 px-4 py-2 rounded-xl text-xs font-bold border border-blue-100 transition-all"
            >
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z"></path></svg>
              Analisis Soal
            </button>
            <div v-if="currentJadwal?.status === 'Diarsipkan'" class="px-4 py-2 bg-emerald-50 text-emerald-600 rounded-xl text-xs font-bold border border-emerald-100 flex items-center gap-2">
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"></path></svg>
              Sudah Diarsipkan
            </div>
            <button 
              v-else-if="currentJadwal?.status === 'Selesai' && authStore.user?.role === 'admin'"
              @click="showArchiveConfirm = true" 
              class="bg-violet-600 hover:bg-violet-700 text-white px-4 py-2 rounded-xl text-xs font-bold transition-all shadow-lg shadow-violet-100"
            >
              Arsipkan Nilai
            </button>
            <button @click="exportExcel" class="flex items-center gap-2 bg-emerald-600 hover:bg-emerald-700 text-white px-4 py-2 rounded-xl text-xs font-bold transition-all shadow-lg shadow-emerald-100">
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 10v6m0 0l-3-3m3 3l3-3m2 8H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"></path></svg>
              Ekspor CSV
            </button>
          </div>
        </div>

        <div class="overflow-x-auto">
          <table class="min-w-full divide-y divide-slate-100">
            <thead class="bg-slate-50">
              <tr>
                <th class="px-6 py-4 text-left text-xs font-bold text-slate-500 uppercase tracking-wider">Nama Siswa</th>
                <th class="px-6 py-4 text-center text-xs font-bold text-slate-500 uppercase tracking-wider">Status</th>
                <th class="px-6 py-4 text-center text-xs font-bold text-slate-500 uppercase tracking-wider">PG</th>
                <th class="px-6 py-4 text-center text-xs font-bold text-slate-500 uppercase tracking-wider">Essay</th>
                <th class="px-6 py-4 text-center text-xs font-bold text-slate-500 uppercase tracking-wider">Total</th>
                <th class="px-6 py-4 text-right text-xs font-bold text-slate-500 uppercase tracking-wider w-20">Aksi</th>
              </tr>
            </thead>
            <tbody class="bg-white divide-y divide-slate-50">
              <tr v-if="isLoading">
                <td colspan="6" class="px-6 py-12 text-center text-slate-400">
                  <div class="flex flex-col items-center gap-2 animate-pulse">
                    <div class="w-8 h-8 border-4 border-blue-500 border-t-transparent rounded-full animate-spin"></div>
                    <span class="text-sm font-medium">Memuat data...</span>
                  </div>
                </td>
              </tr>
              <tr v-for="r in filteredResults" :key="r.id" class="hover:bg-slate-50/50 transition-colors group">
                <td class="px-6 py-4">
                  <div class="flex flex-col">
                    <span class="text-sm font-bold text-slate-800">{{ r.nama_siswa }}</span>
                    <span class="text-[10px] font-black text-slate-400 uppercase tracking-tighter">NISN: {{ r.nisn }} • KELAS {{ r.kelas }}</span>
                  </div>
                </td>
                <td class="px-6 py-4 text-center">
                  <span :class="[
                    r.status_ujian === 'Selesai' ? 'bg-emerald-100 text-emerald-700 border-emerald-200' : 'bg-blue-100 text-blue-700 border-blue-200'
                  ]" class="px-2.5 py-1 rounded-md text-[10px] font-black uppercase border">
                    {{ r.status_ujian }}
                  </span>
                </td>
                <td class="px-6 py-4 text-center font-bold text-slate-600 font-mono">
                  <div class="flex flex-col items-center">
                    <span class="text-sm">{{ r.nilai_pg }}</span>
                    <div class="flex items-center gap-2 mt-1">
                      <span class="text-[10px] text-emerald-600 bg-emerald-50 px-1.5 py-0.5 rounded border border-emerald-100" title="Jawaban Benar">B: {{ r.jumlah_benar }}</span>
                      <span class="text-[10px] text-rose-600 bg-rose-50 px-1.5 py-0.5 rounded border border-rose-100" title="Jawaban Salah">S: {{ r.jumlah_salah }}</span>
                      <button @click="openDetailPG(r)" class="p-0.5 text-blue-400 hover:text-blue-600 transition-colors" title="Lihat Detail Jawaban">
                        <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path></svg>
                      </button>
                    </div>
                  </div>
                </td>
                <td class="px-6 py-4 text-center font-bold text-slate-600 font-mono">{{ r.nilai_essay }}</td>
                <td class="px-6 py-4 text-center">
                  <span class="text-lg font-black text-blue-600 font-mono">{{ r.total_nilai }}</span>
                </td>
                <td class="px-6 py-4 text-right">
                   <button 
                    @click="openKoreksi(r)" 
                    :class="[
                      r.is_koreksi 
                        ? 'bg-emerald-50 text-emerald-600 border-emerald-100' 
                        : 'bg-blue-50 text-blue-600 border-blue-100'
                    ]"
                    class="p-2.5 rounded-xl transition-all border shadow-sm group-hover:scale-110" 
                    :title="r.is_koreksi ? 'Update Nilai' : 'Koreksi Jawaban'"
                   >
                     <svg v-if="!r.is_koreksi" class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"></path></svg>
                     <svg v-else class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"></path></svg>
                   </button>
                </td>
              </tr>
              <tr v-if="filteredResults.length === 0 && !isLoading">
                <td colspan="6" class="px-6 py-12 text-center text-slate-400 italic">Data tidak ditemukan.</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </template>

    <template v-else>
      <div class="bg-white rounded-2xl py-24 text-center border-2 border-dashed border-slate-100">
        <div class="w-16 h-16 bg-slate-50 rounded-full flex items-center justify-center mx-auto mb-4">
          <svg class="w-8 h-8 text-slate-200" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 17v-2m3 2v-4m3 4v-6m2 10H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"></path></svg>
        </div>
        <h3 class="text-lg font-bold text-slate-800">Pilih Jadwal Ujian</h3>
        <p class="text-slate-400 text-sm mt-1">Silakan pilih jadwal untuk melihat rekap nilai.</p>
      </div>
    </template>

    <!-- Modal Koreksi -->
    <BaseModal 
      :show="isKoreksiModalOpen"
      :title="'Koreksi: ' + selectedPeserta?.nama_siswa"
      size="max-w-5xl"
      confirmText="Simpan Semua Nilai"
      :isLoading="isSaving"
      @close="isKoreksiModalOpen = false"
      @confirm="saveKoreksi"
    >
      <div class="space-y-4">
        <div v-if="jawabanPeserta.length === 0" class="text-center py-12 text-slate-400 italic">
          Tidak ada soal essay untuk dikoreksi pada ujian ini.
        </div>
        
        <div v-for="(j, idx) in jawabanPeserta" :key="j.id" class="p-4 bg-slate-50 rounded-2xl border border-slate-100 group">
          <div class="grid grid-cols-12 gap-6">
            <!-- Soal & Jawaban Area -->
            <div class="col-span-12 lg:col-span-9 space-y-3">
              <div class="flex items-start gap-3">
                <span class="flex-shrink-0 w-6 h-6 rounded-lg bg-blue-100 text-blue-600 text-[10px] font-black flex items-center justify-center">{{ idx + 1 }}</span>
                <div class="text-sm font-bold text-slate-800 leading-relaxed prose prose-slate max-w-none" v-html="j.pertanyaan"></div>
              </div>
              
              <div class="ml-9 p-3 bg-white rounded-xl border border-slate-100 text-xs text-slate-600 whitespace-pre-wrap italic">
                <span class="text-[10px] font-black text-emerald-500 uppercase tracking-tighter block mb-1">Jawaban Siswa:</span>
                <div v-html="j.jawaban_siswa || '(Tidak menjawab)'" class="prose prose-sm max-w-none"></div>
              </div>
            </div>

            <!-- Input Skor Area -->
            <div class="col-span-12 lg:col-span-3 flex flex-col justify-center gap-2 border-t lg:border-t-0 lg:border-l border-slate-200 lg:pl-6 pt-4 lg:pt-0">
               <div class="flex items-center justify-between lg:justify-start lg:gap-3">
                 <span class="text-[10px] font-black text-slate-400 uppercase">Bobot Max</span>
                 <span class="text-sm font-black text-slate-700 font-mono">{{ j.bobot_maks }}</span>
               </div>
               <div class="flex flex-col gap-1">
                 <label class="text-[10px] font-black text-slate-500 uppercase">Berikan Skor</label>
                 <input 
                    v-model="j.skor"
                    type="number" 
                    :max="j.bobot_maks"
                    min="0"
                    step="0.5"
                    class="w-full px-4 py-2 bg-white border-2 border-blue-100 rounded-xl font-black text-blue-600 focus:border-blue-500 outline-none transition-all"
                 >
               </div>
            </div>
          </div>
        </div>
      </div>
    </BaseModal>

    <!-- Modal Detail Jawaban PG -->
    <BaseModal 
      :show="isDetailPGModalOpen"
      :title="'Detail Jawaban PG: ' + selectedPeserta?.nama_siswa"
      size="max-w-5xl"
      @close="isDetailPGModalOpen = false"
    >
      <template #footer>
         <button @click="isDetailPGModalOpen = false" class="px-6 py-2.5 bg-slate-100 hover:bg-slate-200 text-slate-600 rounded-xl text-sm font-bold transition-all">Tutup</button>
      </template>

      <div class="space-y-4">
        <div v-if="jawabanLengkap.filter(j => j.jenis_soal === 'PG').length === 0" class="text-center py-12 text-slate-400 italic">
          Tidak ada data jawaban pilihan ganda.
        </div>
        
        <div v-for="(j, idx) in jawabanLengkap.filter(j => j.jenis_soal === 'PG')" :key="j.soal_id" 
          :class="[
            j.skor > 0 ? 'bg-emerald-50/50 border-emerald-100' : 'bg-rose-50/50 border-rose-100'
          ]"
          class="p-4 rounded-2xl border transition-all"
        >
          <div class="flex items-start gap-4">
            <span :class="j.skor > 0 ? 'bg-emerald-100 text-emerald-600' : 'bg-rose-100 text-rose-600'" 
              class="flex-shrink-0 w-8 h-8 rounded-xl text-xs font-black flex items-center justify-center border border-white shadow-sm">
              {{ idx + 1 }}
            </span>
            
            <div class="flex-grow space-y-3">
              <div class="text-sm font-bold text-slate-800 leading-relaxed prose prose-slate max-w-none" v-html="j.pertanyaan"></div>
              
              <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
                <!-- Jawaban Siswa -->
                <div class="p-3 bg-white/60 rounded-xl border border-white/80">
                  <span class="text-[9px] font-black text-slate-400 uppercase tracking-widest block mb-1">Jawaban Siswa</span>
                  <div class="flex items-center gap-2">
                    <span :class="j.skor > 0 ? 'text-emerald-600' : 'text-rose-600'" class="text-sm font-black font-mono">{{ j.jawaban_siswa || '-' }}</span>
                    <span class="text-xs text-slate-600" v-html="getOpsiText(j.jawaban_siswa, j)"></span>
                  </div>
                </div>

                <!-- Kunci Jawaban -->
                <div class="p-3 bg-white/60 rounded-xl border border-white/80">
                  <span class="text-[9px] font-black text-slate-400 uppercase tracking-widest block mb-1">Kunci Jawaban</span>
                  <div class="flex items-center gap-2">
                    <span class="text-sm font-black text-blue-600 font-mono">{{ j.kunci_jawaban }}</span>
                    <span class="text-xs text-slate-600" v-html="getOpsiText(j.kunci_jawaban, j)"></span>
                  </div>
                </div>
              </div>
            </div>

            <!-- Status Icon -->
            <div class="flex-shrink-0 pt-1">
              <div v-if="j.skor > 0" class="w-6 h-6 bg-emerald-500 text-white rounded-full flex items-center justify-center shadow-sm">
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="3" d="M5 13l4 4L19 7"></path></svg>
              </div>
              <div v-else class="w-6 h-6 bg-rose-500 text-white rounded-full flex items-center justify-center shadow-sm">
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="3" d="M6 18L18 6M6 6l12 12"></path></svg>
              </div>
            </div>
          </div>
        </div>
      </div>
    </BaseModal>

    <!-- Archive Confirmation -->
    <ConfirmModal 
      :show="showArchiveConfirm"
      title="Arsipkan Nilai Permanen"
      message="Apakah Anda yakin ingin mengarsipkan nilai ini? Data akan disimpan secara permanen ke Brankas Nilai. Pastikan semua essay sudah dikoreksi karena setelah diarsipkan data tidak dapat diubah lagi."
      confirmText="Ya, Arsipkan Sekarang"
      :isLoading="isArchiving"
      variant="success"
      @close="showArchiveConfirm = false"
      @confirm="executeArchive"
    />
  </div>
</template>

<style scoped>
.font-mono {
  font-family: 'JetBrains Mono', 'Fira Code', monospace;
}
</style>
