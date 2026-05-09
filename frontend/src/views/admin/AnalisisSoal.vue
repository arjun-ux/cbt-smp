<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '../../store/auth'
import { useAlertStore } from '../../store/alert'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()
const alertStore = useAlertStore()

const jadwalId = route.params.jadwalId
const isLoading = ref(true)
const analisisData = ref([])
const selectedSoal = ref(null)
const showDetail = ref(false)

const fetchData = async () => {
  isLoading.value = true
  try {
    const rolePrefix = authStore.user?.role === 'admin' ? 'admin' : 'guru'
    const res = await fetch(`/api/${rolePrefix}/analisis/${jadwalId}`, {
      headers: { 'Authorization': `Bearer ${authStore.token}` }
    })
    const data = await res.json()
    if (res.ok) {
      analisisData.value = data.data || data || []
    } else {
      alertStore.showAlert(data.error || "Gagal memuat analisis", "error")
    }
  } catch (e) {
    alertStore.showAlert("Kesalahan koneksi", "error")
  } finally {
    isLoading.value = false
  }
}

const openDetail = (soal) => {
  selectedSoal.value = soal
  showDetail.value = true
}

const getKesulitanClass = (kesulitan) => {
  if (kesulitan === 'Mudah') return 'bg-emerald-50 text-emerald-600 border-emerald-100'
  if (kesulitan === 'Sukar') return 'bg-rose-50 text-rose-600 border-rose-100'
  return 'bg-amber-50 text-amber-600 border-amber-100'
}

onMounted(fetchData)
</script>

<template>
  <div class="space-y-6 pb-20">
    <!-- Header -->
    <div class="bg-white p-8 rounded-[2.5rem] shadow-sm border border-slate-100 flex flex-col md:flex-row md:items-center justify-between gap-6">
      <div>
        <button @click="router.back()" class="group flex items-center gap-2 text-slate-400 hover:text-blue-600 transition-colors mb-2">
          <svg class="w-4 h-4 group-hover:-translate-x-1 transition-transform" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7"></path></svg>
          <span class="text-xs font-black uppercase tracking-widest">Kembali</span>
        </button>
        <h1 class="text-3xl font-black text-slate-800 tracking-tight">Analisis Butir Soal</h1>
        <p class="text-slate-500 font-medium mt-1">Evaluasi kualitas soal berdasarkan performa jawaban siswa</p>
      </div>

      <div class="flex items-center gap-3">
        <div class="p-4 bg-slate-50 rounded-3xl border border-slate-100 flex items-center gap-4">
          <div class="w-10 h-10 bg-blue-600 text-white rounded-2xl flex items-center justify-center shadow-lg shadow-blue-200">
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z"></path></svg>
          </div>
          <div>
            <p class="text-[10px] font-black text-slate-400 uppercase tracking-widest leading-none mb-1">Total Soal</p>
            <p class="text-xl font-black text-slate-800 leading-none">{{ analisisData.length }}</p>
          </div>
        </div>
      </div>
    </div>

    <!-- Main Content -->
    <div class="bg-white rounded-[2.5rem] shadow-xl shadow-slate-200/50 border border-slate-100 overflow-hidden">
      <div class="overflow-x-auto">
        <table class="w-full text-left border-collapse">
          <thead>
            <tr class="bg-slate-50/50 border-b border-slate-100">
              <th class="px-8 py-6 text-[10px] font-black text-slate-400 uppercase tracking-widest">No</th>
              <th class="px-8 py-6 text-[10px] font-black text-slate-400 uppercase tracking-widest">Pertanyaan</th>
              <th class="px-8 py-6 text-[10px] font-black text-slate-400 uppercase tracking-widest text-center">Benar</th>
              <th class="px-8 py-6 text-[10px] font-black text-slate-400 uppercase tracking-widest text-center">Salah</th>
              <th class="px-8 py-6 text-[10px] font-black text-slate-400 uppercase tracking-widest text-center">Kosong</th>
              <th class="px-8 py-6 text-[10px] font-black text-slate-400 uppercase tracking-widest text-center">Kesulitan</th>
              <th class="px-8 py-6 text-[10px] font-black text-slate-400 uppercase tracking-widest text-right">Aksi</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-50">
            <tr v-if="isLoading">
              <td colspan="7" class="px-8 py-20 text-center">
                <div class="flex flex-col items-center gap-4">
                  <div class="w-12 h-12 border-4 border-blue-600 border-t-transparent rounded-full animate-spin"></div>
                  <p class="text-sm font-black text-slate-400 uppercase tracking-widest">Menganalisis Data...</p>
                </div>
              </td>
            </tr>
            <tr v-else-if="analisisData.length === 0">
              <td colspan="7" class="px-8 py-20 text-center text-slate-400 font-bold">Belum ada data pengerjaan untuk dianalisis.</td>
            </tr>
            <tr v-for="(soal, idx) in analisisData" :key="soal.soal_id" class="hover:bg-slate-50/50 transition-colors group">
              <td class="px-8 py-5 font-black text-slate-400">{{ idx + 1 }}</td>
              <td class="px-8 py-5">
                <div class="max-w-md">
                  <div class="text-sm font-bold text-slate-700 line-clamp-2" v-html="soal.pertanyaan"></div>
                  <span class="text-[10px] font-black px-2 py-0.5 bg-slate-100 text-slate-400 rounded-md uppercase tracking-widest mt-1 inline-block">
                    {{ soal.jenis_soal }}
                  </span>
                </div>
              </td>
              <td class="px-8 py-5 text-center">
                <span class="text-sm font-black text-emerald-600 bg-emerald-50 px-3 py-1 rounded-full">{{ soal.benar }}</span>
              </td>
              <td class="px-8 py-5 text-center">
                <span class="text-sm font-black text-rose-600 bg-rose-50 px-3 py-1 rounded-full">{{ soal.salah }}</span>
              </td>
              <td class="px-8 py-5 text-center">
                <span class="text-sm font-black text-slate-400 bg-slate-50 px-3 py-1 rounded-full">{{ soal.kosong }}</span>
              </td>
              <td class="px-8 py-5 text-center">
                <div class="flex flex-col items-center gap-1">
                  <span :class="getKesulitanClass(soal.kesulitan)" class="px-3 py-1 border rounded-full text-[10px] font-black uppercase tracking-widest">
                    {{ soal.kesulitan }}
                  </span>
                  <span class="text-[10px] font-bold text-slate-400">{{ soal.persentase.toFixed(1) }}%</span>
                </div>
              </td>
              <td class="px-8 py-5 text-right">
                <button 
                  @click="openDetail(soal)"
                  class="p-3 text-blue-600 hover:bg-blue-600 hover:text-white rounded-2xl transition-all active:scale-90 border border-transparent hover:border-blue-100 shadow-sm"
                >
                  <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z"></path></svg>
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Detail Modal -->
    <div v-if="showDetail && selectedSoal" class="fixed inset-0 z-[100] flex items-center justify-center p-6">
      <div class="absolute inset-0 bg-slate-900/60 backdrop-blur-sm" @click="showDetail = false"></div>
      <div class="bg-white w-full max-w-2xl rounded-[3rem] shadow-2xl relative z-10 overflow-hidden animate-slide-up">
        <div class="p-8 border-b border-slate-100 flex items-center justify-between bg-slate-50/50">
          <div>
            <h3 class="text-xl font-black text-slate-800 tracking-tight">Detail Analisis Soal</h3>
            <p class="text-xs font-bold text-slate-400 uppercase tracking-widest mt-1">ID Soal: #{{ selectedSoal.soal_id }}</p>
          </div>
          <button @click="showDetail = false" class="p-3 bg-white text-slate-400 hover:text-rose-600 rounded-2xl transition-colors shadow-sm">
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path></svg>
          </button>
        </div>

        <div class="p-10 space-y-8 max-h-[70vh] overflow-y-auto custom-scrollbar">
          <!-- Pertanyaan -->
          <div class="space-y-3">
            <h4 class="text-[10px] font-black text-slate-400 uppercase tracking-widest">Pertanyaan</h4>
            <div class="p-6 bg-slate-50 rounded-3xl border border-slate-100 text-slate-700 font-medium leading-relaxed prose prose-slate max-w-none" v-html="selectedSoal.pertanyaan"></div>
          </div>

          <!-- Sebaran Jawaban (Khusus PG) -->
          <div v-if="selectedSoal.jenis_soal === 'PG'" class="space-y-6">
            <h4 class="text-[10px] font-black text-slate-400 uppercase tracking-widest">Sebaran Jawaban (Efektivitas Pengecoh)</h4>
            <div class="grid grid-cols-1 gap-4">
              <div v-for="(count, opt) in selectedSoal.sebaran" :key="opt" class="space-y-2">
                <div class="flex justify-between items-center text-sm">
                  <span class="font-black text-slate-700 uppercase">Pilihan {{ opt }}</span>
                  <span class="font-bold text-slate-500">{{ count }} Siswa</span>
                </div>
                <div class="w-full h-3 bg-slate-100 rounded-full overflow-hidden">
                  <div 
                    class="h-full bg-blue-500 rounded-full transition-all duration-1000"
                    :style="{ width: `${(count / (selectedSoal.benar + selectedSoal.salah + selectedSoal.kosong)) * 100}%` }"
                  ></div>
                </div>
              </div>
              <div v-if="!selectedSoal.sebaran || Object.keys(selectedSoal.sebaran).length === 0" class="text-center py-6 bg-slate-50 rounded-3xl border border-dashed border-slate-200">
                <p class="text-sm font-bold text-slate-400 uppercase tracking-widest">Belum ada sebaran data</p>
              </div>
            </div>
          </div>
          
          <div v-else class="p-10 text-center bg-blue-50 rounded-[2.5rem] border border-blue-100">
            <svg class="w-12 h-12 text-blue-400 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"></path></svg>
            <p class="text-blue-700 font-bold">Analisis sebaran jawaban hanya tersedia untuk tipe soal Pilihan Ganda.</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
@keyframes slide-up {
  from { transform: translateY(20px); opacity: 0; }
  to { transform: translateY(0); opacity: 1; }
}
.animate-slide-up {
  animation: slide-up 0.4s cubic-bezier(0.4, 0, 0.2, 1);
}
.custom-scrollbar::-webkit-scrollbar {
  width: 6px;
}
.custom-scrollbar::-webkit-scrollbar-thumb {
  background: #e2e8f0;
  border-radius: 10px;
}
</style>
