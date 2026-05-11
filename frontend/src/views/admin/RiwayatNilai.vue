<script setup>
import { ref, onMounted, computed, watch } from 'vue'
import { useAuthStore } from '../../store/auth'
import { useAlertStore } from '../../store/alert'

const authStore = useAuthStore()
const alertStore = useAlertStore()

const results = ref([])
const isLoading = ref(false)
const searchQuery = ref('')
const filterMapel = ref('')
const filterKelas = ref('')
const isDetailPGModalOpen = ref(false)
const selectedPeserta = ref(null)
const jawabanLengkap = ref([])

const currentPage = ref(1)
const itemsPerPage = ref(20)

const apiPrefix = computed(() => authStore.user?.role === 'admin' ? '/api/admin' : '/api/guru')

const mapels = computed(() => {
  const unique = [...new Set(results.value.map(r => r.nama_mapel))]
  return unique.sort()
})

const classes = computed(() => {
  const unique = [...new Set(results.value.map(r => r.nama_kelas))]
  return unique.sort()
})

const filteredResults = computed(() => {
  let list = [...results.value]
  
  if (filterMapel.value) {
    list = list.filter(r => r.nama_mapel === filterMapel.value)
  }
  
  if (filterKelas.value) {
    list = list.filter(r => r.nama_kelas === filterKelas.value)
  }
  
  if (searchQuery.value) {
    const q = searchQuery.value.toLowerCase()
    list = list.filter(r => 
      r.nama_siswa?.toLowerCase().includes(q) || 
      r.nisn?.toLowerCase().includes(q) ||
      r.judul_ujian?.toLowerCase().includes(q) ||
      r.nama_kelas?.toLowerCase().includes(q)
    )
  }
  return list
})

const totalPages = computed(() => {
  return Math.ceil(filteredResults.value.length / itemsPerPage.value)
})

const paginatedResults = computed(() => {
  const start = (currentPage.value - 1) * itemsPerPage.value
  const end = start + itemsPerPage.value
  return filteredResults.value.slice(start, end)
})

// Reset ke halaman 1 jika filter berubah
watch([searchQuery, filterMapel, filterKelas], () => {
  currentPage.value = 1
})

const openDetailPG = async (peserta) => {
  selectedPeserta.value = peserta
  isDetailPGModalOpen.value = true
  jawabanLengkap.value = []
  
  try {
    const apiPrefix = authStore.user?.role === 'admin' ? '/api/admin' : '/api/guru'
    const res = await fetch(`${apiPrefix}/monitor/jawaban/${peserta.id}`, {
      headers: { 'Authorization': `Bearer ${authStore.token}` }
    })
    const d = await res.json()
    if (res.ok) {
      jawabanLengkap.value = d.data || []
    }
  } catch (e) {
    console.error("Gagal ambil detail jawaban", e)
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

const fetchData = async () => {
  isLoading.value = true
  try {
    const res = await fetch(`${apiPrefix.value}/rekap-arsip`, {
      headers: { 'Authorization': `Bearer ${authStore.token}` }
    })
    if (res.ok) {
      const d = await res.json()
      results.value = d.data || []
    }
  } catch (error) {
    alertStore.showAlert("Gagal mengambil riwayat nilai", "error")
  } finally {
    isLoading.value = false
  }
}

const exportExcel = () => {
  if (filteredResults.value.length === 0) return
  
  let csv = 'No,Tanggal,Ujian,Mata Pelajaran,Nama Siswa,NISN,Kelas,Nilai PG,Nilai Essay,Total Nilai\n'
  filteredResults.value.forEach((r, i) => {
    csv += `${i+1},"${formatDate(r.tanggal_ujian)}","${r.judul_ujian}","${r.nama_mapel}","${r.nama_siswa}","${r.nisn}","${r.nama_kelas}",${r.nilai_pg},${r.nilai_essay},${r.total_nilai}\n`
  })
  
  const blob = new Blob([csv], { type: 'text/csv' })
  const url = window.URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.setAttribute('hidden', '')
  a.setAttribute('href', url)
  a.setAttribute('download', `Riwayat_Nilai_CBT_${new Date().getTime()}.csv`)
  document.body.appendChild(a)
  a.click()
  document.body.removeChild(a)
}

const formatDate = (dateStr) => {
  if (!dateStr) return '-'
  const d = new Date(dateStr)
  return d.toLocaleDateString('id-ID', { day: '2-digit', month: '2-digit', year: 'numeric' })
}

onMounted(fetchData)
</script>

<template>
  <div class="space-y-6 pb-12">
    <!-- Header -->
    <div class="flex flex-col sm:flex-row justify-between items-start sm:items-center bg-white p-6 rounded-2xl shadow-sm border border-slate-100 gap-4">
      <div class="flex items-center gap-4">
        <div>
          <h3 class="text-xl font-bold text-slate-800">Brankas Riwayat Nilai</h3>
          <p class="text-sm text-slate-500 mt-0.5">Pusat arsip nilai permanen yang telah diarsipkan</p>
        </div>
      </div>
      <button @click="exportExcel" class="flex items-center gap-2 bg-emerald-600 hover:bg-emerald-700 text-white px-6 py-2.5 rounded-xl text-sm font-bold transition-all shadow-lg shadow-emerald-100">
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 10v6m0 0l-3-3m3 3l3-3m2 8H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"></path></svg>
        Ekspor Semua (CSV)
      </button>
    </div>

    <!-- Filters & Search -->
    <div class="bg-white p-4 rounded-2xl border border-slate-100 shadow-sm flex flex-col md:flex-row gap-4 bg-slate-50/30">
      <div class="flex-1 relative">
        <svg class="absolute left-4 top-1/2 -translate-y-1/2 w-4 h-4 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"></path></svg>
        <input v-model="searchQuery" type="text" placeholder="Cari Nama, NISN, atau Ujian..." class="w-full pl-11 pr-4 py-2.5 bg-white border border-slate-200 rounded-xl text-sm focus:ring-2 focus:ring-violet-500 transition-all outline-none">
      </div>
      <div class="w-full md:w-48">
        <select v-model="filterKelas" class="w-full px-4 py-2.5 bg-white border border-slate-200 rounded-xl text-sm font-bold focus:ring-2 focus:ring-violet-500 transition-all outline-none">
          <option value="">Semua Kelas</option>
          <option v-for="c in classes" :key="c" :value="c">{{ c }}</option>
        </select>
      </div>
      <div class="w-full md:w-56">
        <select v-model="filterMapel" class="w-full px-4 py-2.5 bg-white border border-slate-200 rounded-xl text-sm font-bold focus:ring-2 focus:ring-violet-500 transition-all outline-none">
          <option value="">Semua Mapel</option>
          <option v-for="m in mapels" :key="m" :value="m">{{ m }}</option>
        </select>
      </div>
    </div>

    <!-- Table Section -->
    <div class="bg-white rounded-2xl shadow-sm border border-slate-100 overflow-hidden">
      <div class="overflow-x-auto">
        <table class="min-w-full divide-y divide-slate-100">
          <thead class="bg-slate-50">
            <tr>
              <th class="px-6 py-4 text-left text-xs font-bold text-slate-500 uppercase tracking-wider">Siswa</th>
              <th class="px-6 py-4 text-left text-xs font-bold text-slate-500 uppercase tracking-wider">Ujian & Mapel</th>
              <th class="px-6 py-4 text-center text-xs font-bold text-slate-500 uppercase tracking-wider">Tanggal</th>
              <th class="px-6 py-4 text-center text-xs font-bold text-slate-500 uppercase tracking-wider">PG</th>
              <th class="px-6 py-4 text-center text-xs font-bold text-slate-500 uppercase tracking-wider">Essay</th>
              <th class="px-6 py-4 text-center text-xs font-bold text-slate-500 uppercase tracking-wider w-24">Total</th>
            </tr>
          </thead>
          <tbody class="bg-white divide-y divide-slate-50">
            <tr v-if="isLoading">
              <td colspan="6" class="px-6 py-12 text-center text-slate-400">
                <div class="flex flex-col items-center gap-2 animate-pulse">
                  <div class="w-8 h-8 border-4 border-violet-500 border-t-transparent rounded-full animate-spin"></div>
                  <span class="text-sm font-medium">Membuka Brankas...</span>
                </div>
              </td>
            </tr>
            <tr v-for="r in paginatedResults" :key="r.id" class="hover:bg-slate-50/50 transition-colors group">
              <td class="px-6 py-4">
                <div class="flex flex-col">
                  <span class="text-sm font-bold text-slate-800">{{ r.nama_siswa }}</span>
                  <span class="text-[10px] font-black text-slate-400 uppercase tracking-tighter">NISN: {{ r.nisn }} • KELAS {{ r.nama_kelas }}</span>
                </div>
              </td>
              <td class="px-6 py-4">
                <div class="flex flex-col">
                  <span class="text-xs font-bold text-slate-700 leading-tight">{{ r.judul_ujian }}</span>
                  <span class="text-[10px] font-black text-violet-500 uppercase tracking-tighter">{{ r.nama_mapel }}</span>
                </div>
              </td>
              <td class="px-6 py-4 text-center text-xs font-bold text-slate-500">{{ formatDate(r.tanggal_ujian) }}</td>
              <td class="px-6 py-4 text-center font-bold text-slate-600 font-mono">
                <div class="flex flex-col items-center">
                  <span class="text-xs">{{ r.nilai_pg }}</span>
                  <div class="flex items-center gap-1.5 mt-1">
                    <span class="text-[9px] text-emerald-600 bg-emerald-50 px-1.5 py-0.5 rounded border border-emerald-100" title="Jawaban Benar">B: {{ r.jumlah_benar }}</span>
                    <span class="text-[9px] text-rose-600 bg-rose-50 px-1.5 py-0.5 rounded border border-rose-100" title="Jawaban Salah">S: {{ r.jumlah_salah }}</span>
                    <button @click="openDetailPG(r)" class="p-0.5 text-blue-400 hover:text-blue-600 transition-colors" title="Lihat Detail Jawaban">
                      <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path></svg>
                    </button>
                  </div>
                </div>
              </td>
              <td class="px-6 py-4 text-center font-bold text-slate-600 font-mono">{{ r.nilai_essay }}</td>
              <td class="px-6 py-4 text-center">
                <div class="inline-flex items-center justify-center w-10 h-10 bg-violet-50 text-violet-600 rounded-xl font-black font-mono shadow-sm border border-violet-100">
                  {{ r.total_nilai }}
                </div>
              </td>
            </tr>
            <tr v-if="filteredResults.length === 0 && !isLoading">
              <td colspan="6" class="px-6 py-20 text-center text-slate-300 font-black uppercase tracking-widest">Brankas Kosong</td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- Pagination -->
      <div v-if="totalPages > 1" class="p-4 border-t border-slate-50 bg-slate-50/30 flex flex-col sm:flex-row items-center justify-between gap-4">
        <div class="text-[10px] font-bold text-slate-400 uppercase tracking-widest">
          Hal {{ currentPage }} dari {{ totalPages }}
        </div>
        <div class="flex items-center gap-2">
          <button @click="currentPage--" :disabled="currentPage === 1" class="px-4 py-2 bg-white border border-slate-200 rounded-xl text-xs font-bold text-slate-600 disabled:opacity-30 transition-all hover:bg-slate-50">PREV</button>
          <div class="flex items-center gap-1">
            <button 
              v-for="p in totalPages" 
              :key="p"
              @click="currentPage = p"
              v-show="p === 1 || p === totalPages || (p >= currentPage - 1 && p <= currentPage + 1)"
              :class="[currentPage === p ? 'bg-violet-600 text-white border-violet-600 shadow-lg shadow-violet-100' : 'bg-white text-slate-400 border-slate-200']"
              class="w-8 h-8 border rounded-xl text-xs font-bold transition-all"
            >
              {{ p }}
            </button>
          </div>
          <button @click="currentPage++" :disabled="currentPage === totalPages" class="px-4 py-2 bg-white border border-slate-200 rounded-xl text-xs font-bold text-slate-600 disabled:opacity-30 transition-all hover:bg-slate-50">NEXT</button>
        </div>
      </div>
    </div>

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
                <div class="p-3 bg-white/60 rounded-xl border border-white/80">
                  <span class="text-[9px] font-black text-slate-400 uppercase tracking-widest block mb-1">Jawaban Siswa</span>
                  <div class="flex items-center gap-2">
                    <span :class="j.skor > 0 ? 'text-emerald-600' : 'text-rose-600'" class="text-sm font-black font-mono">{{ j.jawaban_siswa || '-' }}</span>
                    <span class="text-xs text-slate-600" v-html="getOpsiText(j.jawaban_siswa, j)"></span>
                  </div>
                </div>

                <div class="p-3 bg-white/60 rounded-xl border border-white/80">
                  <span class="text-[9px] font-black text-slate-400 uppercase tracking-widest block mb-1">Kunci Jawaban</span>
                  <div class="flex items-center gap-2">
                    <span class="text-sm font-black text-blue-600 font-mono">{{ j.kunci_jawaban }}</span>
                    <span class="text-xs text-slate-600" v-html="getOpsiText(j.kunci_jawaban, j)"></span>
                  </div>
                </div>
              </div>
            </div>

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
  </div>
</template>

<style scoped>
.font-mono {
  font-family: 'JetBrains Mono', 'Fira Code', monospace;
}
.ml-13 {
  margin-left: 3.25rem;
}
</style>
