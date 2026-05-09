<script setup>
import { ref, onMounted, computed } from 'vue'
import { useAuthStore } from '../../store/auth'
import { useAlertStore } from '../../store/alert'

const authStore = useAuthStore()
const alertStore = useAlertStore()

const jadwals = ref([])
const gururs = ref([])
const isLoading = ref(false)
const isUpdating = ref(null) // ID jadwal yang sedang diupdate

const apiPrefix = computed(() => authStore.user?.role === 'admin' ? '/api/admin' : '/api/guru')

const fetchData = async () => {
  isLoading.value = true
  const headers = { 'Authorization': `Bearer ${authStore.token}` }
  const isGuru = authStore.user?.role === 'guru'
  try {
    const [resJ, resG] = await Promise.all([
      fetch(isGuru ? '/api/guru/pengawas-jadwal' : '/api/admin/jadwal', { headers }),
      fetch(isGuru ? '/api/guru/guru' : '/api/admin/guru', { headers })
    ])
    
    if (resJ.ok && resG.ok) {
      const dj = await resJ.json()
      const dg = await resG.json()
      jadwals.value = dj.data || []
      gururs.value = dg.data || []
    }
  } catch (error) {
    alertStore.showAlert("Gagal memuat data", "error")
  } finally {
    isLoading.value = false
  }
}

const updatePengawas = async (jadwalId, pengawasId) => {
  if (authStore.user?.role !== 'admin') return // Hanya admin yang bisa ubah
  isUpdating.value = jadwalId
  try {
    const res = await fetch(`${apiPrefix.value}/jadwal/${jadwalId}/pengawas`, {
      method: 'POST',
      headers: { 
        'Authorization': `Bearer ${authStore.token}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({ pengawas_id: pengawasId ? parseInt(pengawasId) : null })
    })
    
    if (res.ok) {
      alertStore.showAlert("Pengawas berhasil diplot", "success")
    } else {
      alertStore.showAlert("Gagal mengupdate pengawas", "error")
    }
  } catch (error) {
    alertStore.showAlert("Kesalahan sistem", "error")
  } finally {
    isUpdating.value = null
  }
}

const formatDate = (dateStr) => {
  if (!dateStr) return '-'
  const d = new Date(dateStr)
  return d.toLocaleDateString('id-ID', { day: '2-digit', month: '2-digit', year: 'numeric' })
}

const printPlotting = () => {
  const printWindow = window.open('', '_blank')
  let rows = ''
  jadwals.value.forEach((j, index) => {
    const pengawas = gururs.value.find(g => g.id === j.pengawas_id)?.nama_guru || '-'
    rows += `
      <tr>
        <td style="border: 1px solid #000; padding: 8px; text-align: center;">${index + 1}</td>
        <td style="border: 1px solid #000; padding: 8px;">${formatDate(j.tanggal_ujian)}</td>
        <td style="border: 1px solid #000; padding: 8px; text-align: center;">${j.ruang?.nama_ruang || '-'}</td>
        <td style="border: 1px solid #000; padding: 8px; text-align: center;">${j.sesi?.nama_sesi || j.sesi_id}</td>
        <td style="border: 1px solid #000; padding: 8px;">${j.bank_soal?.judul_bank_soal}</td>
        <td style="border: 1px solid #000; padding: 8px; font-weight: bold;">${pengawas}</td>
      </tr>
    `
  })

  const html = `
    <html>
      <head>
        <title>Daftar Pengawas Ujian</title>
        <style>
          body { font-family: sans-serif; padding: 20px; }
          table { width: 100%; border-collapse: collapse; margin-top: 20px; }
          th { border: 1px solid #000; padding: 10px; background: #eee; }
          h2 { text-align: center; margin-bottom: 5px; }
          p { text-align: center; margin-top: 0; font-size: 14px; }
        </style>
      </head>
      <body>
        <h2>DAFTAR PLOTTING PENGAWAS UJIAN</h2>
        <p>CBT SMP APP - Tanggal Cetak: ${new Date().toLocaleString('id-ID')}</p>
        <table>
          <thead>
            <tr>
              <th>No</th>
              <th>Tanggal</th>
              <th>Ruang</th>
              <th>Sesi</th>
              <th>Mata Pelajaran</th>
              <th>Nama Pengawas</th>
            </tr>
          </thead>
          <tbody>${rows}</tbody>
        </table>
        <div style="margin-top: 40px; text-align: right; padding-right: 50px;">
          <p>Panitia Ujian,</p>
          <br><br><br>
          <p>( ____________________ )</p>
        </div>
        <script>window.print(); window.close();<\/script>
      </body>
    </html>
  `
  printWindow.document.write(html)
  printWindow.document.close()
}

onMounted(fetchData)
</script>

<template>
  <div class="space-y-6 pb-12">
    <!-- Header -->
    <div class="flex flex-col sm:flex-row justify-between items-start sm:items-center bg-white p-6 rounded-2xl shadow-sm border border-slate-100 gap-4">
      <div class="flex items-center gap-4">
        <div>
          <h3 class="text-xl font-bold text-slate-800">Plotting Pengawas</h3>
          <p class="text-sm text-slate-500 mt-0.5">Tugaskan guru untuk mengawas ujian di setiap ruang dan sesi</p>
        </div>
      </div>
      <button @click="printPlotting" class="flex items-center gap-2 bg-emerald-600 hover:bg-emerald-700 text-white px-6 py-2.5 rounded-xl text-sm font-bold transition-all shadow-lg shadow-emerald-100">
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 17h2a2 2 0 002-2v-4a2 2 0 00-2-2H5a2 2 0 00-2 2v4a2 2 0 002 2h2m2 4h6a2 2 0 002-2v-4a2 2 0 00-2-2H9a2 2 0 00-2 2v4a2 2 0 002 2zm8-12V5a2 2 0 00-2-2H9a2 2 0 00-2 2v4h10z"></path></svg>
        Cetak Daftar
      </button>
    </div>

    <!-- Table Section -->
    <div class="bg-white rounded-2xl shadow-sm border border-slate-100 overflow-hidden">
      <div class="p-4 border-b border-slate-50 bg-slate-50/30 flex justify-between items-center">
        <div class="text-xs font-bold text-slate-400 uppercase tracking-widest">Daftar Jadwal Aktif</div>
        <div class="text-xs font-bold text-blue-500 bg-blue-50 px-3 py-1 rounded-lg">{{ jadwals.length }} Jadwal</div>
      </div>

      <div class="overflow-x-auto">
        <table class="min-w-full divide-y divide-slate-100">
          <thead class="bg-slate-50">
            <tr>
              <th class="px-6 py-4 text-left text-xs font-bold text-slate-500 uppercase tracking-wider w-16">No</th>
              <th class="px-6 py-4 text-left text-xs font-bold text-slate-500 uppercase tracking-wider">Jadwal & Ruang</th>
              <th class="px-6 py-4 text-left text-xs font-bold text-slate-500 uppercase tracking-wider">Mata Pelajaran</th>
              <th class="px-6 py-4 text-left text-xs font-bold text-slate-500 uppercase tracking-wider">Sesi</th>
              <th class="px-6 py-4 text-left text-xs font-bold text-slate-500 uppercase tracking-wider w-72">Pengawas</th>
            </tr>
          </thead>
          <tbody class="bg-white divide-y divide-slate-50">
            <tr v-if="isLoading">
              <td colspan="5" class="px-6 py-12 text-center text-slate-400">
                <div class="flex flex-col items-center gap-2 animate-pulse">
                  <div class="w-8 h-8 border-4 border-blue-500 border-t-transparent rounded-full animate-spin"></div>
                  <span class="text-sm font-medium">Memuat data...</span>
                </div>
              </td>
            </tr>
            <tr v-else-if="jadwals.length === 0">
              <td colspan="5" class="px-6 py-12 text-center text-slate-400 italic font-medium">Belum ada jadwal ujian yang dibuat.</td>
            </tr>
            <tr v-for="(j, index) in jadwals" :key="j.id" class="hover:bg-slate-50/50 transition-colors group">
              <td class="px-6 py-5 text-sm font-bold text-slate-400">{{ index + 1 }}</td>
              <td class="px-6 py-5">
                <div class="flex flex-col">
                  <span class="font-bold text-slate-800">{{ formatDate(j.tanggal_ujian) }}</span>
                  <span class="text-[10px] font-black text-blue-500 uppercase tracking-tighter">RUANG: {{ j.ruang?.nama_ruang || '-' }}</span>
                </div>
              </td>
              <td class="px-6 py-5">
                <div class="flex flex-col">
                  <span class="text-sm font-bold text-slate-700 leading-tight">{{ j.bank_soal?.judul_bank_soal }}</span>
                  <div class="flex items-center gap-2 mt-1">
                    <span :class="{
                      'bg-amber-100 text-amber-700 border-amber-200': j.status === 'Belum Mulai',
                      'bg-blue-100 text-blue-700 border-blue-200': j.status === 'Berlangsung',
                      'bg-emerald-100 text-emerald-700 border-emerald-200': j.status === 'Selesai' || j.status === 'Diarsipkan'
                    }" class="px-2 py-0.5 border rounded text-[8px] font-black uppercase">
                      {{ j.status }}
                    </span>
                  </div>
                </div>
              </td>
              <td class="px-6 py-5">
                <span class="px-2.5 py-1 bg-slate-100 text-slate-600 rounded-md text-[10px] font-bold border border-slate-200">
                  {{ j.sesi?.nama_sesi || 'SESI ' + j.sesi_id }}
                </span>
              </td>
              <td class="px-6 py-5">
                <div class="relative">
                  <select 
                    v-model="j.pengawas_id"
                    @change="updatePengawas(j.id, $event.target.value)"
                    :disabled="isUpdating === j.id || j.status === 'Berlangsung' || j.status === 'Diarsipkan'"
                    class="w-full pl-4 pr-10 py-2.5 bg-slate-50 border border-slate-200 rounded-xl text-sm font-bold focus:ring-2 focus:ring-blue-500 transition-all outline-none appearance-none disabled:opacity-50"
                  >
                    <option :value="null">-- Pilih Pengawas --</option>
                    <option v-for="g in gururs" :key="g.id" :value="g.id">{{ g.nama_guru }}</option>
                  </select>
                  <div class="absolute right-3 top-1/2 -translate-y-1/2 pointer-events-none">
                    <svg v-if="isUpdating !== j.id" class="w-4 h-4 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7"></path></svg>
                    <div v-else class="w-4 h-4 border-2 border-blue-500 border-t-transparent rounded-full animate-spin"></div>
                  </div>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<style scoped>
/* Custom scrollbar for better look */
.overflow-x-auto::-webkit-scrollbar {
  height: 6px;
}
.overflow-x-auto::-webkit-scrollbar-thumb {
  background: #e2e8f0;
  border-radius: 10px;
}
</style>
