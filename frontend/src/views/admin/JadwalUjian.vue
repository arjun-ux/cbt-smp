<script setup>
import { ref, onMounted, computed } from 'vue'
import { useAuthStore } from '../../store/auth'
import { useAlertStore } from '../../store/alert'
import BaseModal from '../../components/BaseModal.vue'
import ConfirmModal from '../../components/ConfirmModal.vue'

const authStore = useAuthStore()
const alertStore = useAlertStore()

const jadwals = ref([])
const bankSoals = ref([])
const ruangs = ref([])
const sesis = ref([])
const isLoading = ref(false)
const isSaving = ref(false)
const showForm = ref(false)
const isEditMode = ref(false)
const editId = ref(null)

const showDeleteConfirm = ref(false)
const deleteTargetId = ref(null)
const isDeleting = ref(false)

const searchQuery = ref('')

const filteredJadwals = computed(() => {
  let list = [...jadwals.value]
  
  // Filter
  if (searchQuery.value) {
    const q = searchQuery.value.toLowerCase()
    list = list.filter(j => 
      j.bank_soal?.judul_bank_soal?.toLowerCase().includes(q) ||
      j.bank_soal?.mapel?.nama_mapel?.toLowerCase().includes(q)
    )
  }

  // Sort: Terbaru di atas (Date DESC, Time DESC)
  return list.sort((a, b) => {
    const dateA = a.tanggal_ujian + ' ' + a.waktu_mulai
    const dateB = b.tanggal_ujian + ' ' + b.waktu_mulai
    return dateB.localeCompare(dateA)
  })
})

const formatDate = (dateStr) => {
  if (!dateStr) return '-'
  
  // Jika formatnya YYYY-MM-DD (string murni), hindari pergeseran zona waktu
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

const form = ref({
  bank_soal_id: '',
  tanggal_ujian: '',
  waktu_mulai: '',
  durasi_menit: 90,
  ruang_id: '',
  sesi_id: '',
  acak_soal: false,
  acak_jawaban: false,
  status: 'Belum Mulai'
})

const apiPrefix = computed(() => authStore.user?.role === 'admin' ? '/api/admin' : '/api/guru')

const fetchData = async () => {
  isLoading.value = true
  const headers = { 'Authorization': `Bearer ${authStore.token}` }
  
  try {
    const [resJ, resBS, resR, resS] = await Promise.all([
      fetch(`${apiPrefix.value}/jadwal`, { headers }),
      fetch(`${apiPrefix.value}/bank-soal`, { headers }),
      fetch(`${apiPrefix.value}/ruang`, { headers }),
      fetch(`${apiPrefix.value}/sesi`, { headers })
    ])
    
    if (!resJ.ok || !resBS.ok) {
      throw new Error(`Gagal memuat data utama`)
    }

    const dj = await resJ.json()
    const dbs = await resBS.json()
    jadwals.value = dj.data || []
    bankSoals.value = dbs.data || []

    if (resR.ok) {
      const dr = await resR.json()
      ruangs.value = dr.data || []
    }
    if (resS.ok) {
      const ds = await resS.json()
      sesis.value = ds.data || []
    }
  } catch (error) {
    console.error("Fetch Error:", error)
    alertStore.showAlert("Gagal memuat data: " + error.message, "error")
  } finally {
    isLoading.value = false
  }
}

const triggerForm = () => {
  isEditMode.value = false
  editId.value = null
  form.value = {
    bank_soal_id: '',
    tanggal_ujian: new Date().toISOString().split('T')[0],
    waktu_mulai: '07:30',
    durasi_menit: 90,
    ruang_id: '',
    sesi_id: '',
    acak_soal: false,
    acak_jawaban: false,
    status: 'Belum Mulai'
  }
  showForm.value = true
}

const editJadwal = (j) => {
  isEditMode.value = true
  editId.value = j.id
  
  // Pastikan format tanggal hanya YYYY-MM-DD untuk input date
  const tgl = j.tanggal_ujian && j.tanggal_ujian.includes('T') 
    ? j.tanggal_ujian.split('T')[0] 
    : j.tanggal_ujian

  form.value = {
    bank_soal_id: j.bank_soal_id,
    tanggal_ujian: tgl,
    waktu_mulai: j.waktu_mulai,
    durasi_menit: j.durasi_menit,
    ruang_id: j.ruang_id || '',
    sesi_id: j.sesi_id || '',
    acak_soal: j.acak_soal,
    acak_jawaban: j.acak_jawaban,
    status: j.status
  }
  showForm.value = true
}

const saveJadwal = async () => {
  isSaving.value = true
  try {
    const url = isEditMode.value ? `${apiPrefix.value}/jadwal/${editId.value}` : `${apiPrefix.value}/jadwal`
    const method = isEditMode.value ? 'PUT' : 'POST'
    const res = await fetch(url, {
      method: method,
      headers: { 'Authorization': `Bearer ${authStore.token}`, 'Content-Type': 'application/json' },
      body: JSON.stringify(form.value)
    })
    if (res.ok) {
      alertStore.showAlert(isEditMode.value ? "Jadwal diperbarui" : "Jadwal berhasil dibuat", "success")
      showForm.value = false
      fetchData()
    } else {
      const data = await res.json()
      alertStore.showAlert(data.error || "Gagal menyimpan", "error")
    }
  } catch (error) {
    alertStore.showAlert("Kesalahan sistem", "error")
  } finally {
    isSaving.value = false
  }
}

const refreshToken = async (id) => {
  try {
    const res = await fetch(`${apiPrefix.value}/jadwal/${id}/token`, {
      method: 'PATCH',
      headers: { 'Authorization': `Bearer ${authStore.token}` }
    })
    if (res.ok) {
      const data = await res.json()
      // Menggunakan data.data.token karena dibungkus oleh helper SendSuccess di backend
      alertStore.showAlert("Token baru: " + data.data.token, "success")
      fetchData()
    }
  } catch (error) {
    alertStore.showAlert("Gagal refresh token", "error")
  }
}

const confirmDelete = (id) => {
  deleteTargetId.value = id
  showDeleteConfirm.value = true
}

const executeDelete = async () => {
  isDeleting.value = true
  try {
    const res = await fetch(`${apiPrefix.value}/jadwal/${deleteTargetId.value}`, {
      method: 'DELETE',
      headers: { 'Authorization': `Bearer ${authStore.token}` }
    })
    if (res.ok) {
      alertStore.showAlert("Jadwal dihapus", "success")
      showDeleteConfirm.value = false
      fetchData()
    } else {
      const data = await res.json()
      alertStore.showAlert(data.error || "Gagal menghapus jadwal", "error")
      showDeleteConfirm.value = false
    }
  } catch (error) {
    alertStore.showAlert("Kesalahan sistem saat menghapus", "error")
  } finally {
    isDeleting.value = false
  }
}

onMounted(fetchData)
</script>

<template>
  <div class="space-y-6 pb-20">
    <!-- Header -->
    <div class="flex flex-col sm:flex-row justify-between items-start sm:items-center bg-white p-6 rounded-2xl shadow-sm border border-slate-100 gap-4">
      <div class="flex-1">
        <h3 class="text-xl font-bold text-slate-800">Penjadwalan Ujian</h3>
        <p class="text-sm text-slate-500 mt-1">Atur jadwal pelaksanaan ujian dan generate token</p>
      </div>
      <div class="flex flex-col sm:flex-row items-center gap-3 w-full sm:w-auto">
        <div class="relative w-full sm:w-64">
          <svg class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"></path></svg>
          <input 
            v-model="searchQuery"
            type="text" 
            placeholder="Cari ujian/mapel..." 
            class="w-full pl-10 pr-4 py-2.5 bg-slate-50 border border-slate-200 rounded-xl text-sm focus:ring-2 focus:ring-blue-500 transition-all"
          >
        </div>
        <button 
          v-if="authStore.user?.role === 'admin'"
          @click="triggerForm" 
          class="w-full sm:w-auto flex items-center justify-center gap-2 bg-blue-600 hover:bg-blue-700 text-white px-6 py-2.5 rounded-xl text-sm font-bold transition-all shadow-lg shadow-blue-200"
        >
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"></path></svg>
          <span class="whitespace-nowrap">Tambah Jadwal</span>
        </button>
      </div>
    </div>

    <!-- Table -->
    <div class="bg-white rounded-2xl shadow-sm border border-slate-100 overflow-hidden">
      <div class="overflow-x-auto">
        <table class="min-w-full divide-y divide-slate-100">
          <thead class="bg-slate-50">
            <tr>
              <th class="px-6 py-4 text-left text-xs font-bold text-slate-500 uppercase tracking-wider">Info Ujian</th>
              <th class="px-6 py-4 text-left text-xs font-bold text-slate-500 uppercase tracking-wider">Waktu & Durasi</th>
              <th class="px-6 py-4 text-center text-xs font-bold text-slate-500 uppercase tracking-wider">Token</th>
              <th class="px-6 py-4 text-left text-xs font-bold text-slate-500 uppercase tracking-wider">Status</th>
              <th class="px-6 py-4 text-right text-xs font-bold text-slate-500 uppercase tracking-wider">Aksi</th>
            </tr>
          </thead>
          <tbody class="bg-white divide-y divide-slate-50">
            <tr v-if="isLoading">
              <td colspan="5" class="px-6 py-12 text-center text-slate-400">Memuat data...</td>
            </tr>
            <tr v-else-if="jadwals.length === 0">
              <td colspan="5" class="px-6 py-12 text-center text-slate-400 italic">Belum ada jadwal ujian.</td>
            </tr>
            <tr v-for="j in filteredJadwals" :key="j.id" class="hover:bg-slate-50/80 transition-colors group">
              <td class="px-6 py-4">
                <div class="flex flex-col">
                  <span class="text-sm font-bold text-slate-800">{{ j.bank_soal?.judul_bank_soal }}</span>
                  <div class="flex items-center gap-2 mt-1">
                    <span class="px-2 py-0.5 bg-blue-50 text-blue-600 text-[10px] font-bold rounded uppercase">{{ j.bank_soal?.mapel?.nama_mapel }}</span>
                    <span class="text-[10px] text-slate-400 font-bold uppercase tracking-wider">Kelas {{ j.bank_soal?.tingkat_kelas }}</span>
                  </div>
                </div>
              </td>
              <td class="px-6 py-4">
                <div class="flex flex-col">
                  <div class="flex items-center gap-1.5 text-sm font-semibold text-slate-700">
                    <svg class="w-3.5 h-3.5 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"></path></svg>
                    {{ formatDate(j.tanggal_ujian) }}
                  </div>
                  <div class="flex items-center gap-1.5 text-xs text-slate-500 mt-1">
                    <svg class="w-3.5 h-3.5 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"></path></svg>
                    {{ j.waktu_mulai }} ({{ j.durasi_menit }} Menit)
                  </div>
                  <div v-if="j.ruang || j.sesi" class="flex items-center gap-2 mt-1.5">
                    <span v-if="j.ruang" class="px-1.5 py-0.5 bg-slate-100 text-slate-600 text-[10px] font-bold rounded border border-slate-200">
                      {{ j.ruang.nama_ruang }}
                    </span>
                    <span v-if="j.sesi" class="px-1.5 py-0.5 bg-slate-100 text-slate-600 text-[10px] font-bold rounded border border-slate-200">
                      {{ j.sesi.nama_sesi }}
                    </span>
                  </div>
                </div>
              </td>
              <td class="px-6 py-4 text-center">
                <div class="inline-flex items-center gap-2 bg-slate-100 px-3 py-1.5 rounded-lg border border-slate-200 group-hover:bg-blue-50 group-hover:border-blue-200 transition-colors">
                  <span class="text-sm font-black text-slate-800 tracking-widest font-mono group-hover:text-blue-600">{{ j.token_ujian }}</span>
                  <button 
                    @click="refreshToken(j.id)" 
                    class="p-1 hover:bg-white rounded transition-colors text-slate-400 hover:text-blue-600" 
                    title="Refresh Token"
                  >
                    <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"></path></svg>
                  </button>
                </div>
              </td>
              <td class="px-6 py-4">
                <span :class="{
                  'bg-amber-50 text-amber-600': j.status === 'Belum Mulai',
                  'bg-blue-50 text-blue-600': j.status === 'Berlangsung',
                  'bg-emerald-50 text-emerald-600': j.status === 'Selesai'
                }" class="px-2.5 py-1 rounded-lg text-[10px] font-black uppercase tracking-wider border">
                  {{ j.status }}
                </span>
              </td>
              <td class="px-6 py-4 text-right">
                <div class="flex justify-end gap-2">
                  <router-link 
                    :to="{ name: authStore.user?.role === 'admin' ? 'AdminMonitor' : 'GuruMonitor', params: { jadwalId: j.id } }"
                    class="text-slate-400 hover:text-blue-600 hover:bg-blue-50 p-2 rounded-xl transition-all"
                    title="Monitor Ujian"
                  >
                    <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z"></path></svg>
                  </router-link>
                  <!-- Tombol Rekap Nilai: Hanya untuk Admin (Guru akses via menu Laporan) -->
                  <router-link 
                    v-if="authStore.user?.role === 'admin'"
                    :to="{ name: 'AdminRekap', params: { jadwalId: j.id } }"
                    class="text-slate-400 hover:text-emerald-600 hover:bg-emerald-50 p-2 rounded-xl transition-all"
                    title="Rekap Nilai"
                  >
                    <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-3 7h3m-3 4h3m-6-4h.01M9 16h.01"></path></svg>
                  </router-link>


                  <button 
                    v-if="authStore.user?.role === 'admin'"
                    @click="editJadwal(j)" 
                    class="text-slate-400 hover:text-blue-600 hover:bg-blue-50 p-2 rounded-xl transition-all"
                  >
                    <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"></path></svg>
                  </button>
                  <button 
                    v-if="authStore.user?.role === 'admin'"
                    @click="confirmDelete(j.id)" 
                    :disabled="j.status === 'Berlangsung'"
                    class="text-slate-400 hover:text-rose-600 hover:bg-rose-50 p-2 rounded-xl transition-all disabled:opacity-20 disabled:cursor-not-allowed"
                  >
                    <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path></svg>
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Modal Form -->
    <BaseModal 
      :show="showForm" 
      :title="isEditMode ? 'Edit Jadwal Ujian' : 'Tambah Jadwal Baru'"
      :isLoading="isSaving"
      @close="showForm = false"
      @confirm="saveJadwal"
    >
      <div class="space-y-5">
        <div>
          <label class="block text-sm font-semibold text-slate-700 mb-1.5">Pilih Bank Soal</label>
          <select v-model="form.bank_soal_id" required class="block w-full px-4 py-3 bg-slate-50 border border-slate-200 rounded-xl focus:ring-2 focus:ring-blue-500 text-sm font-medium">
            <option value="">-- Pilih Bank Soal --</option>
            <option v-for="bs in bankSoals" :key="bs.id" :value="bs.id">
              [{{ bs.tingkat_kelas }}] {{ bs.judul_bank_soal }} - {{ bs.mapel?.nama_mapel }}
            </option>
          </select>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-2 gap-5">
          <div>
            <label class="block text-sm font-semibold text-slate-700 mb-1.5">Tanggal Ujian</label>
            <input v-model="form.tanggal_ujian" type="date" required class="block w-full px-4 py-3 bg-slate-50 border border-slate-200 rounded-xl focus:ring-2 focus:ring-blue-500 text-sm font-medium">
          </div>
          <div>
            <label class="block text-sm font-semibold text-slate-700 mb-1.5">Waktu Mulai</label>
            <input v-model="form.waktu_mulai" type="time" required class="block w-full px-4 py-3 bg-slate-50 border border-slate-200 rounded-xl focus:ring-2 focus:ring-blue-500 text-sm font-medium">
          </div>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-2 gap-5">
          <div>
            <label class="block text-sm font-semibold text-slate-700 mb-1.5">Lokasi / Ruang</label>
            <select v-model="form.ruang_id" class="block w-full px-4 py-3 bg-slate-50 border border-slate-200 rounded-xl focus:ring-2 focus:ring-blue-500 text-sm font-medium">
              <option value="">-- Tanpa Ruang --</option>
              <option v-for="r in ruangs" :key="r.id" :value="r.id">
                {{ r.nama_ruang }} ({{ r.kode_ruang }})
              </option>
            </select>
          </div>
          <div>
            <label class="block text-sm font-semibold text-slate-700 mb-1.5">Sesi Ujian</label>
            <select v-model="form.sesi_id" class="block w-full px-4 py-3 bg-slate-50 border border-slate-200 rounded-xl focus:ring-2 focus:ring-blue-500 text-sm font-medium">
              <option value="">-- Tanpa Sesi --</option>
              <option v-for="s in sesis" :key="s.id" :value="s.id">
                {{ s.nama_sesi }} ({{ s.waktu_mulai }} - {{ s.waktu_selesai }})
              </option>
            </select>
          </div>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-2 gap-5">
          <div>
            <label class="block text-sm font-semibold text-slate-700 mb-1.5">Durasi (Menit)</label>
            <input v-model.number="form.durasi_menit" type="number" required class="block w-full px-4 py-3 bg-slate-50 border border-slate-200 rounded-xl focus:ring-2 focus:ring-blue-500 text-sm font-medium">
          </div>
          <div class="flex flex-col gap-4 h-full pt-6">
            <!-- Toggle Acak Soal -->
            <label class="flex items-center justify-between p-3 bg-slate-50 rounded-2xl border border-slate-200 cursor-pointer hover:border-blue-300 transition-all group">
              <span class="text-sm font-bold text-slate-700 group-hover:text-blue-700 transition-colors">Acak Urutan Soal</span>
              <div class="relative w-11 h-6 bg-slate-300 rounded-full transition-colors peer-checked:bg-blue-600 overflow-hidden">
                <input v-model="form.acak_soal" type="checkbox" class="sr-only peer">
                <div class="absolute left-1 top-1 bg-white w-4 h-4 rounded-full transition-all peer-checked:translate-x-5 shadow-sm"></div>
                <div class="w-full h-full rounded-full peer-checked:bg-blue-600 transition-colors"></div>
              </div>
            </label>

            <!-- Toggle Acak Jawaban -->
            <label class="flex items-center justify-between p-3 bg-slate-50 rounded-2xl border border-slate-200 cursor-pointer hover:border-emerald-300 transition-all group">
              <span class="text-sm font-bold text-slate-700 group-hover:text-emerald-700 transition-colors">Acak Pilihan Jawaban</span>
              <div class="relative w-11 h-6 bg-slate-300 rounded-full transition-colors peer-checked:bg-emerald-600 overflow-hidden">
                <input v-model="form.acak_jawaban" type="checkbox" class="sr-only peer">
                <div class="absolute left-1 top-1 bg-white w-4 h-4 rounded-full transition-all peer-checked:translate-x-5 shadow-sm"></div>
                <div class="w-full h-full rounded-full peer-checked:bg-emerald-600 transition-colors"></div>
              </div>
            </label>
          </div>
        </div>

        <div v-if="isEditMode" class="p-4 bg-slate-50 border border-slate-100 rounded-2xl">
          <label class="block text-sm font-bold text-slate-700 mb-3">Status Ujian</label>
          <div class="grid grid-cols-3 gap-3">
            <button 
              v-for="s in ['Belum Mulai', 'Berlangsung', 'Selesai']" 
              :key="s"
              type="button"
              @click="form.status = s"
              :class="form.status === s ? 'bg-blue-600 text-white border-blue-600' : 'bg-white text-slate-600 border-slate-200 hover:border-blue-400'"
              class="px-3 py-2 rounded-xl text-[11px] font-black uppercase tracking-wider border transition-all"
            >
              {{ s }}
            </button>
          </div>
          <p v-if="form.status === 'Berlangsung'" class="text-[10px] text-amber-600 font-bold mt-3 bg-amber-50 p-2 rounded-lg border border-amber-100">
            ⚠ PERINGATAN: Jadwal sedang berlangsung. Mengubah data selain status dapat mengganggu jalannya ujian.
          </p>
        </div>
      </div>
    </BaseModal>

    <!-- Delete Confirmation -->
    <ConfirmModal 
      :show="showDeleteConfirm"
      title="Hapus Jadwal"
      message="Apakah Anda yakin ingin menghapus jadwal ujian ini? Data nilai siswa (jika ada) mungkin akan terpengaruh."
      :isLoading="isDeleting"
      @close="showDeleteConfirm = false"
      @confirm="executeDelete"
    />
  </div>
</template>
