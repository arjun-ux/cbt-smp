<script setup>
import { ref, onMounted, computed, watch } from 'vue'
import { useAuthStore } from '../../store/auth'
import { useAlertStore } from '../../store/alert'
import BaseModal from '../../components/BaseModal.vue'
import ConfirmModal from '../../components/ConfirmModal.vue'

const authStore = useAuthStore()
const alertStore = useAlertStore()

const siswas = ref([])
const searchQuery = ref('')
const isLoading = ref(false)
const isSaving = ref(false)
const showForm = ref(false)
const isEditMode = ref(false)
const editId = ref(null)
const fileInput = ref(null)
const isUploading = ref(false)

const selectedIds = ref([])
const showDeleteConfirm = ref(false)
const deleteTargetId = ref(null)
const isDeleting = ref(false)

// Pagination State
const currentPage = ref(1)
const itemsPerPage = 10

// State untuk Preview
const showPreview = ref(false)
const previewData = ref([])
const selectedFile = ref(null)

// Master Data untuk Dropdown
const kelases = ref([])
const ruangs = ref([])
const sesis = ref([])

const form = ref({
  nisn: '',
  nama_lengkap: '',
  password: '',
  kelas_id: null,
  ruang_id: null,
  sesi_id: null
})

const fetchSiswas = async () => {
  isLoading.value = true
  try {
    const res = await fetch('/api/admin/siswa', {
      headers: { 'Authorization': `Bearer ${authStore.token}` }
    })
    const data = await res.json()
    if(res.ok) siswas.value = data.data
  } catch (error) {
    alertStore.showAlert("Gagal memuat data siswa", "error")
  } finally {
    isLoading.value = false
  }
}

const filteredSiswas = computed(() => {
  if (!searchQuery.value) return siswas.value
  const q = searchQuery.value.toLowerCase()
  return siswas.value.filter(s => 
    s.nama_lengkap.toLowerCase().includes(q) || 
    s.nisn.toLowerCase().includes(q) ||
    s.kelas?.nama_kelas?.toLowerCase().includes(q)
  )
})

const totalPages = computed(() => Math.ceil(filteredSiswas.value.length / itemsPerPage))
const paginatedSiswas = computed(() => {
  const start = (currentPage.value - 1) * itemsPerPage
  const end = start + itemsPerPage
  return filteredSiswas.value.slice(start, end)
})

watch(searchQuery, () => {
  currentPage.value = 1
})

const fetchMasterData = async () => {
  const headers = { 'Authorization': `Bearer ${authStore.token}` }
  try {
    const [resK, resR, resS] = await Promise.all([
      fetch('/api/admin/kelas', { headers }),
      fetch('/api/admin/ruang', { headers }),
      fetch('/api/admin/sesi', { headers })
    ])
    const [dk, dr, ds] = await Promise.all([resK.json(), resR.json(), resS.json()])
    kelases.value = dk.data || []
    ruangs.value = dr.data || []
    sesis.value = ds.data || []
  } catch (e) {
    console.error("Gagal memuat master data dropdown")
  }
}

const openAdd = () => {
  isEditMode.value = false
  editId.value = null
  form.value = { nisn: '', nama_lengkap: '', password: '', kelas_id: null, ruang_id: null, sesi_id: null }
  showForm.value = true
}

const openEdit = (siswa) => {
  isEditMode.value = true
  editId.value = siswa.id
  form.value = {
    nisn: siswa.nisn,
    nama_lengkap: siswa.nama_lengkap,
    password: '',
    kelas_id: siswa.kelas_id,
    ruang_id: siswa.ruang_id,
    sesi_id: siswa.sesi_id
  }
  showForm.value = true
}

const saveSiswa = async () => {
  isSaving.value = true
  try {
    const url = isEditMode.value ? `/api/admin/siswa/${editId.value}` : '/api/admin/siswa'
    const method = isEditMode.value ? 'PUT' : 'POST'

    const res = await fetch(url, {
      method: method,
      headers: { 
        'Authorization': `Bearer ${authStore.token}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(form.value)
    })
    
    if(res.ok) {
      alertStore.showAlert(isEditMode.value ? "Data Siswa diperbarui" : "Siswa berhasil ditambahkan", "success")
      showForm.value = false
      fetchSiswas()
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

const confirmDelete = (id = null) => {
  deleteTargetId.value = id
  showDeleteConfirm.value = true
}

const executeDelete = async () => {
  isDeleting.value = true
  let successCount = 0
  let failCount = 0
  let lastErrorMessage = ""

  try {
    if (deleteTargetId.value) {
      const res = await fetch(`/api/admin/siswa/${deleteTargetId.value}`, {
        method: 'DELETE',
        headers: { 'Authorization': `Bearer ${authStore.token}` }
      })
      if (res.ok) {
        successCount++
      } else {
        const data = await res.json()
        lastErrorMessage = data.error || "Gagal menghapus data"
        failCount++
      }
    } else {
      for (const id of selectedIds.value) {
        const res = await fetch(`/api/admin/siswa/${id}`, {
          method: 'DELETE',
          headers: { 'Authorization': `Bearer ${authStore.token}` }
        })
        if (res.ok) {
          successCount++
        } else {
          failCount++
        }
      }
      selectedIds.value = []
    }

    if (failCount > 0) {
      if (successCount > 0) {
        alertStore.showAlert(`${successCount} data dihapus, ${failCount} gagal (terikat data lain)`, "warning")
      } else {
        alertStore.showAlert(lastErrorMessage || "Gagal menghapus: Data sedang digunakan/terikat", "error")
      }
    } else if (successCount > 0) {
      alertStore.showAlert("Data berhasil dihapus", "success")
    }
    
    fetchSiswas()
  } catch (error) {
    alertStore.showAlert("Kesalahan sistem saat menghapus", "error")
  } finally {
    isDeleting.value = false
    showDeleteConfirm.value = false
  }
}

const toggleStatus = async (item) => {
  try {
    const res = await fetch(`/api/admin/siswa/${item.id}/status`, {
      method: 'PATCH',
      headers: { 'Authorization': `Bearer ${authStore.token}` }
    })
    const data = await res.json()
    if (res.ok) {
      alertStore.showAlert(data.message, "success")
      fetchSiswas()
    } else {
      alertStore.showAlert(data.error || "Gagal mengubah status", "error")
    }
  } catch (error) {
    alertStore.showAlert("Kesalahan sistem saat mengubah status", "error")
  }
}

const toggleSelectAll = (e) => {
  if (e.target.checked) {
    selectedIds.value = paginatedSiswas.value.map(i => i.id)
  } else {
    selectedIds.value = []
  }
}

const triggerFileInput = () => {
  fileInput.value.click()
}

const handleFileChange = (event) => {
  const file = event.target.files[0]
  if (!file) return

  if (!file.name.endsWith('.csv')) {
    alertStore.showAlert("File harus berformat CSV", "error")
    event.target.value = ''
    return
  }

  selectedFile.value = file
  const reader = new FileReader()
  reader.onload = (e) => {
    const text = e.target.result
    const rows = text.split('\n')
    const parsedData = []
    
    for (let i = 1; i < rows.length; i++) {
      const cols = rows[i].split(',')
      if (cols.length >= 3) {
        parsedData.push({
          nisn: cols[0]?.trim(),
          nama: cols[1]?.trim(),
          password: cols[2]?.trim(),
          kelas: cols[3]?.trim() || '-'
        })
      }
    }
    
    if (parsedData.length === 0) {
      alertStore.showAlert("File CSV kosong atau format salah", "error")
      return
    }

    previewData.value = parsedData
    showPreview.value = true
  }
  reader.readAsText(file)
}

const confirmImport = async () => {
  isUploading.value = true
  const formData = new FormData()
  formData.append('file', selectedFile.value)

  try {
    const res = await fetch('/api/admin/siswa/import', {
      method: 'POST',
      headers: { 'Authorization': `Bearer ${authStore.token}` },
      body: formData
    })
    
    if(res.ok) {
      alertStore.showAlert("Import data siswa berhasil!", "success")
      showPreview.value = false
      fetchSiswas()
    } else {
      const data = await res.json()
      alertStore.showAlert(data.error || "Gagal melakukan import", "error")
    }
  } catch (error) {
    alertStore.showAlert("Gagal mengunggah file", "error")
  } finally {
    isUploading.value = false
    selectedFile.value = null
    if(fileInput.value) fileInput.value.value = ''
  }
}

const downloadTemplate = () => {
  const csvContent = "NISN,NamaLengkap,Password,Kelas\n2021001,Ahmad Siswa,Siswa123,7A\n2021002,Siti Siswi,Siswa456,8B"
  const blob = new Blob([csvContent], { type: 'text/csv;charset=utf-8;' })
  const link = document.createElement("a")
  const url = URL.createObjectURL(blob)
  link.setAttribute("href", url)
  link.setAttribute("download", "Template_Import_Siswa.csv")
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
}

onMounted(() => {
  fetchSiswas()
  fetchMasterData()
})
</script>

<template>
  <div class="space-y-6 pb-12">
    <!-- Header -->
    <div class="flex flex-col sm:flex-row justify-between items-start sm:items-center bg-white p-6 rounded-2xl shadow-sm border border-slate-100 gap-4">
      <div>
        <h3 class="text-xl font-bold text-slate-800">Manajemen Siswa</h3>
        <p class="text-sm text-slate-500 mt-1">Kelola data peserta dan akun login siswa</p>
      </div>
      <div class="flex items-center gap-3 w-full sm:w-auto">
        <input type="file" ref="fileInput" @change="handleFileChange" accept=".csv" class="hidden">
        
        <button @click="downloadTemplate" class="flex-1 sm:flex-none flex items-center justify-center gap-2 bg-emerald-50 hover:bg-emerald-100 text-emerald-700 px-4 py-2.5 rounded-xl text-sm font-bold transition-all border border-emerald-200">
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 10v6m0 0l-3-3m3 3l3-3m2 8H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"></path></svg>
          Template
        </button>

        <button @click="triggerFileInput" class="flex-1 sm:flex-none flex items-center justify-center gap-2 bg-slate-100 hover:bg-slate-200 text-slate-700 px-4 py-2.5 rounded-xl text-sm font-bold transition-all">
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-8l-4-4m0 0L8 8m4-4v12"></path></svg>
          Import
        </button>
        <button @click="openAdd" class="flex-1 sm:flex-none flex items-center justify-center gap-2 bg-blue-600 hover:bg-blue-700 text-white px-6 py-2.5 rounded-xl text-sm font-bold transition-all shadow-lg shadow-blue-200">
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"></path></svg>
          Tambah Siswa
        </button>
      </div>
    </div>

    <!-- Modal Preview Import -->
    <div v-if="showPreview" class="fixed inset-0 bg-slate-900/60 backdrop-blur-sm z-[60] flex items-center justify-center p-4">
      <div class="bg-white rounded-[2rem] shadow-2xl w-full max-w-4xl max-h-[85vh] flex flex-col overflow-hidden animate-in fade-in zoom-in duration-300">
        <div class="p-8 border-b border-slate-100 flex justify-between items-center bg-slate-50/50">
          <div>
            <h4 class="text-xl font-black text-slate-800">Pratinjau Data Impor</h4>
            <p class="text-sm text-slate-500">Periksa kembali data sebelum disimpan ke sistem</p>
          </div>
          <button @click="showPreview = false" class="p-2 bg-white hover:bg-rose-50 text-slate-400 hover:text-rose-500 rounded-xl transition-all shadow-sm">
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path></svg>
          </button>
        </div>
        <div class="flex-1 overflow-auto p-8">
          <div class="border border-slate-100 rounded-2xl overflow-hidden shadow-sm">
            <table class="min-w-full divide-y divide-slate-100">
              <thead class="bg-slate-50">
                <tr>
                  <th class="px-6 py-4 text-left text-xs font-black text-slate-400 uppercase tracking-widest">NISN</th>
                  <th class="px-6 py-4 text-left text-xs font-black text-slate-400 uppercase tracking-widest">Nama Lengkap</th>
                  <th class="px-6 py-4 text-left text-xs font-black text-slate-400 uppercase tracking-widest">Kelas</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-slate-100 bg-white">
                <tr v-for="(row, i) in previewData" :key="i" class="hover:bg-slate-50 transition-colors">
                  <td class="px-6 py-4 text-sm font-bold text-slate-700">{{ row.nisn }}</td>
                  <td class="px-6 py-4 text-sm text-slate-600">{{ row.nama }}</td>
                  <td class="px-6 py-4 text-sm">
                    <span class="px-3 py-1 bg-blue-50 text-blue-600 font-black rounded-lg text-xs border border-blue-100">{{ row.kelas }}</span>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
        <div class="p-8 border-t border-slate-100 flex justify-end gap-3 bg-slate-50/50">
          <button @click="showPreview = false" class="px-8 py-3 rounded-2xl text-sm font-bold text-slate-500 hover:bg-slate-200 transition-all">Batal</button>
          <button @click="confirmImport" :disabled="isUploading" class="px-10 py-3 rounded-2xl text-sm font-black bg-blue-600 text-white hover:bg-blue-700 shadow-xl shadow-blue-200 transition-all disabled:opacity-50">
            {{ isUploading ? 'Mengimpor...' : 'Konfirmasi & Impor' }}
          </button>
        </div>
      </div>
    </div>

    <!-- Tabel Utama -->
    <div class="bg-white rounded-2xl shadow-sm border border-slate-100 overflow-hidden">
      <!-- Search & Bulk Actions -->
      <div class="p-4 border-b border-slate-50 bg-slate-50/30 flex flex-col sm:flex-row justify-between items-center gap-4">
        <div class="relative w-full sm:w-80">
          <svg class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"></path></svg>
          <input v-model="searchQuery" type="text" placeholder="Cari NISN, nama, atau kelas..." class="block w-full pl-10 pr-4 py-2.5 bg-white border border-slate-200 rounded-xl text-sm focus:ring-2 focus:ring-blue-500 transition-all">
        </div>
        
        <div v-if="selectedIds.length > 0" class="flex items-center gap-3 animate-in fade-in slide-in-from-right-4 duration-300">
          <span class="text-sm font-semibold text-slate-600">{{ selectedIds.length }} item dipilih</span>
          <button @click="confirmDelete()" class="flex items-center gap-2 bg-rose-50 text-rose-600 px-4 py-2 rounded-xl text-sm font-bold hover:bg-rose-100 transition-colors">
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path></svg>
            Hapus Masal
          </button>
        </div>

        <div v-else class="text-xs font-bold text-slate-400 bg-slate-100 px-3 py-1.5 rounded-lg uppercase tracking-wider">
          Total: {{ filteredSiswas.length }} Siswa
        </div>
      </div>

      <div class="overflow-x-auto">
        <table class="min-w-full divide-y divide-slate-100">
          <thead class="bg-slate-50">
            <tr>
              <th class="px-6 py-4 w-10">
                <input type="checkbox" :checked="selectedIds.length === paginatedSiswas.length && paginatedSiswas.length > 0" @change="toggleSelectAll" class="rounded border-slate-300 text-blue-600 focus:ring-blue-500">
              </th>
              <th class="px-6 py-4 text-left text-xs font-bold text-slate-500 uppercase tracking-wider">NISN / Username</th>
              <th class="px-6 py-4 text-left text-xs font-bold text-slate-500 uppercase tracking-wider">Nama Lengkap</th>
              <th class="px-6 py-4 text-left text-xs font-bold text-slate-500 uppercase tracking-wider">Penempatan</th>
              <th class="px-6 py-4 text-left text-xs font-bold text-slate-500 uppercase tracking-wider">Status</th>
              <th class="px-6 py-4 text-right text-xs font-bold text-slate-500 uppercase tracking-wider">Aksi</th>
            </tr>
          </thead>
          <tbody class="bg-white divide-y divide-slate-50">
            <tr v-if="isLoading">
              <td colspan="6" class="px-6 py-12 text-center text-slate-400">
                <div class="flex flex-col items-center gap-2">
                  <svg class="w-8 h-8 animate-spin text-blue-500" fill="none" viewBox="0 0 24 24"><circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle><path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path></svg>
                  <span class="text-sm font-medium">Memuat data...</span>
                </div>
              </td>
            </tr>
            <tr v-else-if="filteredSiswas.length === 0">
              <td colspan="6" class="px-6 py-12 text-center text-slate-400 italic">Data siswa tidak ditemukan.</td>
            </tr>
            <tr v-for="siswa in paginatedSiswas" :key="siswa.id" class="hover:bg-slate-50/50 transition-colors group">
              <td class="px-6 py-4">
                <input type="checkbox" v-model="selectedIds" :value="siswa.id" class="rounded border-slate-300 text-blue-600 focus:ring-blue-500">
              </td>
              <td class="px-6 py-4 text-sm font-bold text-slate-800">{{ siswa.nisn }}</td>
              <td class="px-6 py-4 text-sm text-slate-600 font-medium">{{ siswa.nama_lengkap }}</td>
              <td class="px-6 py-4 text-xs font-semibold text-slate-500">
                <div class="flex flex-wrap gap-2">
                  <span class="px-2 py-0.5 bg-blue-50 text-blue-600 rounded-md border border-blue-100" title="Kelas">K: {{ siswa.kelas?.nama_kelas || '-' }}</span>
                  <span class="px-2 py-0.5 bg-indigo-50 text-indigo-600 rounded-md border border-indigo-100" title="Ruang">R: {{ siswa.ruang?.nama_ruang || '-' }}</span>
                  <span class="px-2 py-0.5 bg-slate-100 text-slate-600 rounded-md border border-slate-200" title="Sesi">S: {{ siswa.sesi?.nama_sesi || '-' }}</span>
                </div>
              </td>
              <td class="px-6 py-4 text-sm cursor-pointer" @click="toggleStatus(siswa)" title="Klik untuk mengubah status">
                <span v-if="siswa.user?.is_active" class="px-2.5 py-1 text-[10px] font-black uppercase rounded-md bg-emerald-100 text-emerald-700 border border-emerald-200 hover:bg-emerald-200 transition-colors">Aktif</span>
                <span v-else class="px-2.5 py-1 text-[10px] font-black uppercase rounded-md bg-rose-100 text-rose-700 border border-rose-200 hover:bg-rose-200 transition-colors">Nonaktif</span>
              </td>
              <td class="px-6 py-4 text-right">
                <div class="flex justify-end gap-2">
                  <button @click="openEdit(siswa)" class="p-2 text-slate-400 hover:text-blue-600 hover:bg-blue-50 rounded-xl transition-all" title="Edit">
                    <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"></path></svg>
                  </button>
                  <button @click="confirmDelete(siswa.id)" class="p-2 text-slate-400 hover:text-rose-600 hover:bg-rose-50 rounded-xl transition-all" title="Hapus">
                    <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path></svg>
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- Pagination Controls -->
      <div v-if="totalPages > 1" class="p-4 border-t border-slate-50 bg-slate-50/30 flex justify-between items-center">
        <button 
          @click="currentPage--" 
          :disabled="currentPage === 1"
          class="flex items-center gap-2 px-4 py-2 text-sm font-bold text-slate-600 hover:bg-slate-100 rounded-xl disabled:opacity-30 transition-all"
        >
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7"></path></svg>
          Sebelumnya
        </button>
        <div class="hidden sm:flex items-center gap-1">
          <button 
            v-for="p in totalPages" 
            :key="p" 
            @click="currentPage = p"
            :class="[
              'w-9 h-9 rounded-xl text-xs font-bold transition-all',
              currentPage === p ? 'bg-blue-600 text-white shadow-lg shadow-blue-200' : 'text-slate-600 hover:bg-slate-200'
            ]"
          >
            {{ p }}
          </button>
        </div>
        <button 
          @click="currentPage++" 
          :disabled="currentPage === totalPages"
          class="flex items-center gap-2 px-4 py-2 text-sm font-bold text-slate-600 hover:bg-slate-100 rounded-xl disabled:opacity-30 transition-all"
        >
          Selanjutnya
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"></path></svg>
        </button>
      </div>
    </div>

    <!-- Modal Form -->
    <BaseModal 
      :show="showForm" 
      :title="isEditMode ? 'Edit Data Siswa' : 'Tambah Siswa Baru'"
      :isLoading="isSaving"
      @close="showForm = false"
      @confirm="saveSiswa"
    >
      <div class="space-y-4">
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <label class="block text-sm font-semibold text-slate-700 mb-1.5">NISN / Username Login</label>
            <input v-model="form.nisn" type="text" placeholder="Contoh: 2021001" required class="w-full px-4 py-2.5 bg-slate-50 border border-slate-200 rounded-xl focus:ring-2 focus:ring-blue-500 text-sm">
          </div>
          <div>
            <label class="block text-sm font-semibold text-slate-700 mb-1.5">
              Password 
              <span v-if="isEditMode" class="text-[10px] font-normal text-slate-400 italic">(Kosongkan jika tidak diubah)</span>
            </label>
            <input v-model="form.password" type="password" :required="!isEditMode" placeholder="••••••••" class="w-full px-4 py-2.5 bg-slate-50 border border-slate-200 rounded-xl focus:ring-2 focus:ring-blue-500 text-sm">
          </div>
        </div>
        <div>
          <label class="block text-sm font-semibold text-slate-700 mb-1.5">Nama Lengkap Siswa</label>
          <input v-model="form.nama_lengkap" type="text" placeholder="Masukkan nama lengkap siswa" required class="w-full px-4 py-2.5 bg-slate-50 border border-slate-200 rounded-xl focus:ring-2 focus:ring-blue-500 text-sm">
        </div>
        <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
          <div>
            <label class="block text-sm font-semibold text-slate-700 mb-1.5">Kelas</label>
            <select v-model="form.kelas_id" class="w-full px-4 py-2.5 bg-slate-50 border border-slate-200 rounded-xl focus:ring-2 focus:ring-blue-500 text-sm">
              <option :value="null">-- Pilih --</option>
              <option v-for="k in kelases" :key="k.id" :value="k.id">{{ k.nama_kelas }}</option>
            </select>
          </div>
          <div>
            <label class="block text-sm font-semibold text-slate-700 mb-1.5">Ruang</label>
            <select v-model="form.ruang_id" class="w-full px-4 py-2.5 bg-slate-50 border border-slate-200 rounded-xl focus:ring-2 focus:ring-blue-500 text-sm">
              <option :value="null">-- Pilih --</option>
              <option v-for="r in ruangs" :key="r.id" :value="r.id">{{ r.nama_ruang }}</option>
            </select>
          </div>
          <div>
            <label class="block text-sm font-semibold text-slate-700 mb-1.5">Sesi</label>
            <select v-model="form.sesi_id" class="w-full px-4 py-2.5 bg-slate-50 border border-slate-200 rounded-xl focus:ring-2 focus:ring-blue-500 text-sm">
              <option :value="null">-- Pilih --</option>
              <option v-for="s in sesis" :key="s.id" :value="s.id">{{ s.nama_sesi }}</option>
            </select>
          </div>
        </div>
      </div>
    </BaseModal>

    <!-- Modal Konfirmasi Hapus -->
    <ConfirmModal 
      :show="showDeleteConfirm"
      title="Hapus Data Siswa"
      :message="deleteTargetId ? 'Apakah Anda yakin ingin menghapus siswa ini? Akun login terkait juga akan dihapus.' : `Apakah Anda yakin ingin menghapus ${selectedIds.length} siswa yang dipilih?`"
      :isLoading="isDeleting"
      @close="showDeleteConfirm = false"
      @confirm="executeDelete"
    />
  </div>
</template>
