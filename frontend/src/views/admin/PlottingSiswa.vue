<script setup>
import { ref, onMounted, computed, watch } from 'vue'
import { useAuthStore } from '../../store/auth'
import { useAlertStore } from '../../store/alert'
import ConfirmModal from '../../components/ConfirmModal.vue'

const authStore = useAuthStore()
const alertStore = useAlertStore()
const siswas = ref([])
const kelases = ref([])
const ruangs = ref([])
const sesis = ref([])

const isLoading = ref(false)
const isProcessing = ref(false)
const selectedKelas = ref('all')
const searchQuery = ref('')

// Selection State
const selectedSiswaIds = ref([])
const bulkRuangId = ref('')
const bulkSesiId = ref('')
const showConfirm = ref(false)

const fetchData = async () => {
  isLoading.value = true
  const headers = { 'Authorization': `Bearer ${authStore.token}` }
  try {
    const [resS, resK, resR, resSesi] = await Promise.all([
      fetch('/api/admin/siswa', { headers }),
      fetch('/api/admin/kelas', { headers }),
      fetch('/api/admin/ruang', { headers }),
      fetch('/api/admin/sesi', { headers })
    ])
    
    const [ds, dk, dr, dsesi] = await Promise.all([
      resS.json(), resK.json(), resR.json(), resSesi.json()
    ])
    
    siswas.value = ds.data || []
    kelases.value = dk.data || []
    ruangs.value = dr.data || []
    sesis.value = dsesi.data || []
  } catch (error) {
    alertStore.showAlert("Gagal memuat data", "error")
  } finally {
    isLoading.value = false
  }
}

const filteredSiswas = computed(() => {
  let data = siswas.value
  
  if (selectedKelas.value !== 'all') {
    data = data.filter(s => s.kelas_id === parseInt(selectedKelas.value))
  }
  
  if (searchQuery.value) {
    const q = searchQuery.value.toLowerCase()
    data = data.filter(s => 
      s.nama_lengkap.toLowerCase().includes(q) || 
      s.nisn.toLowerCase().includes(q)
    )
  }
  
  return data
})

const toggleSelectAll = (event) => {
  if (event.target.checked) {
    selectedSiswaIds.value = filteredSiswas.value.map(s => s.id)
  } else {
    selectedSiswaIds.value = []
  }
}

const isAllSelected = computed(() => {
  return filteredSiswas.value.length > 0 && selectedSiswaIds.value.length === filteredSiswas.value.length
})

const handleApplyClick = () => {
  if (selectedSiswaIds.value.length === 0) return
  if (!bulkRuangId.value || !bulkSesiId.value) {
    alertStore.showAlert("Silakan pilih Ruang dan Sesi terlebih dahulu", "warning")
    return
  }
  showConfirm.value = true
}

const applyBulkPlot = async () => {
  isProcessing.value = true
  try {
    const res = await fetch('/api/admin/siswa/bulk-plot', {
      method: 'POST',
      headers: { 
        'Authorization': `Bearer ${authStore.token}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        siswa_ids: selectedSiswaIds.value,
        ruang_id: parseInt(bulkRuangId.value),
        sesi_id: parseInt(bulkSesiId.value)
      })
    })
    
    if (res.ok) {
      alertStore.showAlert("Berhasil memperbarui data siswa secara massal", "success")
      selectedSiswaIds.value = []
      showConfirm.value = false
      fetchData()
    } else {
      const data = await res.json()
      alertStore.showAlert(data.error || "Gagal memproses", "error")
    }
  } catch (error) {
    alertStore.showAlert("Terjadi kesalahan sistem", "error")
  } finally {
    isProcessing.value = false
  }
}

onMounted(fetchData)

watch(selectedKelas, () => {
  selectedSiswaIds.value = []
})
</script>

<template>
  <div class="space-y-6 pb-32">
    <!-- Header -->
    <div class="bg-white p-6 rounded-2xl shadow-sm border border-slate-100 flex flex-col md:flex-row justify-between items-start md:items-center gap-4">
      <div>
        <h3 class="text-xl font-bold text-slate-800">Plotting Ruang & Sesi</h3>
        <p class="text-sm text-slate-500 mt-1">Atur pembagian ruang dan waktu ujian secara massal</p>
      </div>
      
      <div class="flex flex-wrap items-center gap-3 w-full md:w-auto">
        <div class="relative flex-1 md:w-64">
          <span class="absolute inset-y-0 left-0 pl-3 flex items-center text-slate-400">
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"></path></svg>
          </span>
          <input v-model="searchQuery" type="text" placeholder="Cari nama/NISN..." class="block w-full pl-10 pr-4 py-2 bg-slate-50 border border-slate-200 rounded-xl text-sm focus:ring-2 focus:ring-blue-500 transition-all">
        </div>
        
        <select v-model="selectedKelas" class="bg-slate-50 border border-slate-200 rounded-xl px-4 py-2 text-sm focus:ring-2 focus:ring-blue-500 outline-none">
          <option value="all">Semua Kelas</option>
          <option v-for="k in kelases" :key="k.id" :value="k.id">{{ k.nama_kelas }}</option>
        </select>
      </div>
    </div>

    <!-- Alert Information -->
    <div class="bg-blue-50 border border-blue-100 p-4 rounded-2xl flex items-start gap-3">
      <div class="p-2 bg-blue-100 rounded-lg text-blue-600">
        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path></svg>
      </div>
      <div>
        <h5 class="text-sm font-bold text-blue-900">Tips Plotting Massal</h5>
        <p class="text-xs text-blue-700 mt-0.5">Filter berdasarkan kelas terlebih dahulu, lalu klik checkbox di header tabel untuk memilih seluruh siswa dalam kelas tersebut.</p>
      </div>
    </div>

    <!-- Table -->
    <div class="bg-white rounded-2xl shadow-sm border border-slate-100 overflow-hidden">
      <div class="overflow-x-auto">
        <table class="min-w-full divide-y divide-slate-100">
          <thead class="bg-slate-50/50">
            <tr>
              <th class="px-6 py-4 text-left">
                <input type="checkbox" @change="toggleSelectAll" :checked="isAllSelected" class="w-4 h-4 rounded border-slate-300 text-blue-600 focus:ring-blue-500">
              </th>
              <th class="px-6 py-4 text-left text-xs font-bold text-slate-500 uppercase tracking-wider">No</th>
              <th class="px-6 py-4 text-left text-xs font-bold text-slate-500 uppercase tracking-wider">NISN</th>
              <th class="px-6 py-4 text-left text-xs font-bold text-slate-500 uppercase tracking-wider">Nama Lengkap</th>
              <th class="px-6 py-4 text-left text-xs font-bold text-slate-500 uppercase tracking-wider">Kelas</th>
              <th class="px-6 py-4 text-left text-xs font-bold text-slate-500 uppercase tracking-wider">Ruang</th>
              <th class="px-6 py-4 text-left text-xs font-bold text-slate-500 uppercase tracking-wider">Sesi</th>
            </tr>
          </thead>
          <tbody class="bg-white divide-y divide-slate-50">
            <tr v-if="isLoading">
              <td colspan="7" class="px-6 py-12 text-center text-slate-400">Memuat data siswa...</td>
            </tr>
            <tr v-else-if="filteredSiswas.length === 0">
              <td colspan="7" class="px-6 py-12 text-center text-slate-400">Tidak ada siswa yang ditemukan.</td>
            </tr>
            <tr v-for="(siswa, index) in filteredSiswas" :key="siswa.id" class="hover:bg-slate-50/80 transition-colors">
              <td class="px-6 py-4">
                <input type="checkbox" v-model="selectedSiswaIds" :value="siswa.id" class="w-4 h-4 rounded border-slate-300 text-blue-600 focus:ring-blue-500">
              </td>
              <td class="px-6 py-4 text-sm text-slate-500">{{ index + 1 }}</td>
              <td class="px-6 py-4 text-sm font-bold text-slate-800">{{ siswa.nisn }}</td>
              <td class="px-6 py-4 text-sm text-slate-600 font-medium">{{ siswa.nama_lengkap }}</td>
              <td class="px-6 py-4 text-sm">
                <span class="px-2 py-0.5 bg-blue-50 text-blue-600 rounded text-[10px] font-black border border-blue-100 uppercase">{{ siswa.kelas?.nama_kelas || '-' }}</span>
              </td>
              <td class="px-6 py-4">
                <span :class="siswa.ruang_id ? 'text-emerald-600 bg-emerald-50 px-2.5 py-1 rounded-md text-[10px] font-black border border-emerald-100' : 'text-slate-400 italic text-[10px]'">
                  {{ siswa.ruang?.nama_ruang || 'BELUM DIATUR' }}
                </span>
              </td>
              <td class="px-6 py-4">
                <span :class="siswa.sesi_id ? 'text-indigo-600 bg-indigo-50 px-2.5 py-1 rounded-md text-[10px] font-black border border-indigo-100' : 'text-slate-400 italic text-[10px]'">
                  {{ siswa.sesi?.nama_sesi || 'BELUM DIATUR' }}
                </span>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Floating Action Bar -->
    <transition name="slide-up">
      <div v-if="selectedSiswaIds.length > 0" class="fixed bottom-8 left-1/2 -translate-x-1/2 w-[90%] max-w-4xl bg-slate-900/90 text-white p-6 rounded-[2.5rem] shadow-2xl backdrop-blur-xl z-50 flex flex-col md:flex-row items-center justify-between gap-6 border border-white/10">
        <div class="flex items-center gap-4">
          <div class="w-14 h-14 bg-blue-600 text-white rounded-2xl flex items-center justify-center text-2xl font-black shadow-lg shadow-blue-500/20">
            {{ selectedSiswaIds.length }}
          </div>
          <div>
            <h6 class="font-black text-white text-lg tracking-tight">Siswa Terpilih</h6>
            <p class="text-xs text-slate-400 font-medium">Atur ruang dan sesi secara massal</p>
          </div>
        </div>

        <div class="flex flex-wrap items-center gap-3 w-full md:w-auto">
          <div class="flex-1 md:flex-none">
            <select v-model="bulkRuangId" class="w-full md:w-44 bg-white/10 border border-white/20 text-white rounded-2xl px-4 py-3 text-sm outline-none focus:ring-2 focus:ring-blue-500 transition-all appearance-none cursor-pointer">
              <option value="" class="text-slate-800">Pilih Ruang</option>
              <option v-for="r in ruangs" :key="r.id" :value="r.id" class="text-slate-800">{{ r.nama_ruang }}</option>
            </select>
          </div>
          <div class="flex-1 md:flex-none">
            <select v-model="bulkSesiId" class="w-full md:w-44 bg-white/10 border border-white/20 text-white rounded-2xl px-4 py-3 text-sm outline-none focus:ring-2 focus:ring-blue-500 transition-all appearance-none cursor-pointer">
              <option value="" class="text-slate-800">Pilih Sesi</option>
              <option v-for="s in sesis" :key="s.id" :value="s.id" class="text-slate-800">{{ s.nama_sesi }} ({{ s.waktu_mulai }})</option>
            </select>
          </div>
          <button 
            @click="handleApplyClick" 
            :disabled="isProcessing"
            class="w-full md:w-auto bg-blue-600 text-white hover:bg-blue-700 disabled:bg-slate-700 px-10 py-3 rounded-2xl font-black text-sm uppercase tracking-widest transition-all shadow-xl active:scale-95"
          >
            {{ isProcessing ? 'Memproses...' : 'Terapkan' }}
          </button>
        </div>
      </div>
    </transition>

    <ConfirmModal 
      :show="showConfirm"
      title="Terapkan Plotting Massal"
      :message="`Apakah Anda yakin ingin mengatur ${selectedSiswaIds.length} siswa ke Ruang ${ruangs.find(r => r.id == bulkRuangId)?.nama_ruang} dan Sesi ${sesis.find(s => s.id == bulkSesiId)?.nama_sesi}?`"
      confirmText="Ya, Terapkan"
      :isLoading="isProcessing"
      @close="showConfirm = false"
      @confirm="applyBulkPlot"
    />
  </div>
</template>

<style scoped>
.slide-up-enter-active,
.slide-up-leave-active {
  transition: all 0.4s cubic-bezier(0.16, 1, 0.3, 1);
}

.slide-up-enter-from,
.slide-up-leave-to {
  transform: translate(-50%, 100%);
  opacity: 0;
}

select option {
  background: white;
  color: #1e293b;
}
</style>
