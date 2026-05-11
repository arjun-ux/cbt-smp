<script setup>
import { ref, onMounted, computed, watch } from 'vue'
import { useAuthStore } from '../../store/auth'
import { useAlertStore } from '../../store/alert'
import BaseModal from '../../components/BaseModal.vue'
import ConfirmModal from '../../components/ConfirmModal.vue'

const authStore = useAuthStore()
const alertStore = useAlertStore()

const gurus = ref([])
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
 
const showPreview = ref(false)
const previewData = ref([])
const selectedFile = ref(null)

// Pagination State
const currentPage = ref(1)
const itemsPerPage = 10

const form = ref({
  nip: '',
  nama_guru: '',
  password: ''
})

const fetchGurus = async () => {
  isLoading.value = true
  try {
    const res = await fetch('/api/admin/guru', {
      headers: { 'Authorization': `Bearer ${authStore.token}` }
    })
    const data = await res.json()
    if(res.ok) gurus.value = data.data
  } catch (error) {
    alertStore.showAlert("Gagal memuat data guru", "error")
  } finally {
    isLoading.value = false
  }
}

const filteredGurus = computed(() => {
  if (!searchQuery.value) return gurus.value
  const q = searchQuery.value.toLowerCase()
  return gurus.value.filter(g => 
    g.nama_guru.toLowerCase().includes(q) || 
    g.nip.toLowerCase().includes(q)
  )
})

const totalPages = computed(() => Math.ceil(filteredGurus.value.length / itemsPerPage))
const paginatedGurus = computed(() => {
  const start = (currentPage.value - 1) * itemsPerPage
  const end = start + itemsPerPage
  return filteredGurus.value.slice(start, end)
})

watch(searchQuery, () => {
  currentPage.value = 1
})

const openAdd = () => {
  isEditMode.value = false
  editId.value = null
  form.value = { nip: '', nama_guru: '', password: '' }
  showForm.value = true
}

const openEdit = (guru) => {
  isEditMode.value = true
  editId.value = guru.id
  form.value = {
    nip: guru.nip,
    nama_guru: guru.nama_guru,
    password: ''
  }
  showForm.value = true
}

const saveGuru = async () => {
  isSaving.value = true
  try {
    const url = isEditMode.value ? `/api/admin/guru/${editId.value}` : '/api/admin/guru'
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
      alertStore.showAlert(isEditMode.value ? "Data Guru diperbarui" : "Guru berhasil ditambahkan", "success")
      showForm.value = false
      fetchGurus()
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
      const res = await fetch(`/api/admin/guru/${deleteTargetId.value}`, {
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
        const res = await fetch(`/api/admin/guru/${id}`, {
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
    
    fetchGurus()
  } catch (error) {
    alertStore.showAlert("Kesalahan sistem saat menghapus", "error")
  } finally {
    isDeleting.value = false
    showDeleteConfirm.value = false
  }
}

const toggleStatus = async (item) => {
  try {
    const res = await fetch(`/api/admin/guru/${item.id}/status`, {
      method: 'PATCH',
      headers: { 'Authorization': `Bearer ${authStore.token}` }
    })
    const data = await res.json()
    if (res.ok) {
      alertStore.showAlert(data.message, "success")
      fetchGurus()
    } else {
      alertStore.showAlert(data.error || "Gagal mengubah status", "error")
    }
  } catch (error) {
    alertStore.showAlert("Kesalahan sistem saat mengubah status", "error")
  }
}

const toggleSelectAll = (e) => {
  if (e.target.checked) {
    selectedIds.value = paginatedGurus.value.map(i => i.id)
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
          nip: cols[0]?.trim(),
          nama_guru: cols[1]?.trim(),
          password: cols[2]?.trim()
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

const commitImport = async () => {
  isUploading.value = true
  const formData = new FormData()
  formData.append('file', selectedFile.value)

  try {
    const res = await fetch('/api/admin/guru/import', {
      method: 'POST',
      headers: { 'Authorization': `Bearer ${authStore.token}` },
      body: formData
    })
    
    const data = await res.json()
    if(res.ok) {
      alertStore.showAlert(data.message || "Import guru berhasil", "success")
      showPreview.value = false
      fetchGurus()
    } else {
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
  const csvContent = "NIP,NamaGuru,Password\n198001012023011001,Drs. Budi Santoso,guru123\n198505052023021002,Siti Aminah S.Pd,guru456"
  const blob = new Blob([csvContent], { type: 'text/csv;charset=utf-8;' })
  const link = document.createElement("a")
  const url = URL.createObjectURL(blob)
  link.setAttribute("href", url)
  link.setAttribute("download", "Template_Import_Guru.csv")
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
}

onMounted(fetchGurus)
</script>

<template>
  <div class="space-y-6 pb-12">
    <!-- Header -->
    <div class="flex flex-col sm:flex-row justify-between items-start sm:items-center bg-white p-6 rounded-2xl shadow-sm border border-slate-100 gap-4">
      <div>
        <h3 class="text-xl font-bold text-slate-800">Manajemen Guru</h3>
        <p class="text-sm text-slate-500 mt-1">Kelola data login dan biodata tenaga pengajar</p>
      </div>
      <div class="flex items-center gap-3 w-full sm:w-auto">
        <input type="file" ref="fileInput" @change="handleFileChange" accept=".csv" class="hidden">
        
        <button @click="downloadTemplate" class="flex-1 sm:flex-none flex items-center justify-center gap-2 bg-emerald-50 hover:bg-emerald-100 text-emerald-700 px-4 py-2.5 rounded-xl text-sm font-bold transition-all border border-emerald-200">
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 10v6m0 0l-3-3m3 3l3-3m2 8H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"></path></svg>
          Template
        </button>

        <button @click="triggerFileInput" :disabled="isUploading" class="flex-1 sm:flex-none flex items-center justify-center gap-2 bg-slate-100 hover:bg-slate-200 text-slate-700 px-4 py-2.5 rounded-xl text-sm font-bold transition-all disabled:opacity-50">
          <svg v-if="!isUploading" class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-8l-4-4m0 0L8 8m4-4v12"></path></svg>
          <svg v-else class="w-4 h-4 animate-spin" fill="none" viewBox="0 0 24 24"><circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle><path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path></svg>
          {{ isUploading ? 'Mengimpor...' : 'Import CSV' }}
        </button>
        <button @click="openAdd" class="flex-1 sm:flex-none flex items-center justify-center gap-2 bg-blue-600 hover:bg-blue-700 text-white px-6 py-2.5 rounded-xl text-sm font-bold transition-all shadow-lg shadow-blue-200">
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"></path></svg>
          Tambah Guru
        </button>
      </div>
    </div>

    <!-- Tabel Utama -->
    <div class="bg-white rounded-2xl shadow-sm border border-slate-100 overflow-hidden">
      <!-- Search & Bulk Actions -->
      <div class="p-4 border-b border-slate-50 bg-slate-50/30 flex flex-col sm:flex-row justify-between items-center gap-4">
        <div class="relative w-full sm:w-80">
          <svg class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"></path></svg>
          <input v-model="searchQuery" type="text" placeholder="Cari NIP atau nama guru..." class="block w-full pl-10 pr-4 py-2.5 bg-white border border-slate-200 rounded-xl text-sm focus:ring-2 focus:ring-blue-500 transition-all">
        </div>
        
        <div v-if="selectedIds.length > 0" class="flex items-center gap-3 animate-in fade-in slide-in-from-right-4 duration-300">
          <span class="text-sm font-semibold text-slate-600">{{ selectedIds.length }} item dipilih</span>
          <button @click="confirmDelete()" class="flex items-center gap-2 bg-rose-50 text-rose-600 px-4 py-2 rounded-xl text-sm font-bold hover:bg-rose-100 transition-colors">
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path></svg>
            Hapus Masal
          </button>
        </div>

        <div v-else class="text-xs font-bold text-slate-400 bg-slate-100 px-3 py-1.5 rounded-lg uppercase tracking-wider">
          Total: {{ filteredGurus.length }} Guru
        </div>
      </div>

      <div class="overflow-x-auto">
        <table class="min-w-full divide-y divide-slate-100">
          <thead class="bg-slate-50">
            <tr>
              <th class="px-6 py-4 w-10">
                <input type="checkbox" :checked="selectedIds.length === paginatedGurus.length && paginatedGurus.length > 0" @change="toggleSelectAll" class="rounded border-slate-300 text-blue-600 focus:ring-blue-500">
              </th>
              <th class="px-6 py-4 text-left text-xs font-bold text-slate-500 uppercase tracking-wider">NIP / Username</th>
              <th class="px-6 py-4 text-left text-xs font-bold text-slate-500 uppercase tracking-wider">Nama Lengkap</th>
              <th class="px-6 py-4 text-left text-xs font-bold text-slate-500 uppercase tracking-wider">Status</th>
              <th class="px-6 py-4 text-right text-xs font-bold text-slate-500 uppercase tracking-wider">Aksi</th>
            </tr>
          </thead>
          <tbody class="bg-white divide-y divide-slate-50">
            <tr v-if="isLoading">
              <td colspan="5" class="px-6 py-12 text-center text-slate-400">
                <div class="flex flex-col items-center gap-2">
                  <svg class="w-8 h-8 animate-spin text-blue-500" fill="none" viewBox="0 0 24 24"><circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle><path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path></svg>
                  <span class="text-sm font-medium">Memuat data...</span>
                </div>
              </td>
            </tr>
            <tr v-else-if="filteredGurus.length === 0">
              <td colspan="5" class="px-6 py-12 text-center text-slate-400 italic">Data guru tidak ditemukan.</td>
            </tr>
            <tr v-for="guru in paginatedGurus" :key="guru.id" class="hover:bg-slate-50/50 transition-colors group">
              <td class="px-6 py-4">
                <input type="checkbox" v-model="selectedIds" :value="guru.id" class="rounded border-slate-300 text-blue-600 focus:ring-blue-500">
              </td>
              <td class="px-6 py-4 text-sm font-bold text-slate-800">{{ guru.nip }}</td>
              <td class="px-6 py-4 text-sm text-slate-600 font-medium">{{ guru.nama_guru }}</td>
              <td class="px-6 py-4 text-sm cursor-pointer" @click="toggleStatus(guru)" title="Klik untuk mengubah status">
                <span v-if="guru.user?.is_active" class="px-2.5 py-1 text-[10px] font-black uppercase rounded-md bg-emerald-100 text-emerald-700 border border-emerald-200 hover:bg-emerald-200 transition-colors">Aktif</span>
                <span v-else class="px-2.5 py-1 text-[10px] font-black uppercase rounded-md bg-rose-100 text-rose-700 border border-rose-200 hover:bg-rose-200 transition-colors">Nonaktif</span>
              </td>
              <td class="px-6 py-4 text-right">
                <div class="flex justify-end gap-2">
                  <button @click="openEdit(guru)" class="p-2 text-slate-400 hover:text-blue-600 hover:bg-blue-50 rounded-xl transition-all" title="Edit">
                    <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"></path></svg>
                  </button>
                  <button @click="confirmDelete(guru.id)" class="p-2 text-slate-400 hover:text-rose-600 hover:bg-rose-50 rounded-xl transition-all" title="Hapus">
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
      :title="isEditMode ? 'Edit Data Guru' : 'Tambah Guru Baru'"
      :isLoading="isSaving"
      @close="showForm = false"
      @confirm="saveGuru"
    >
      <div class="space-y-4">
        <div>
          <label class="block text-sm font-semibold text-slate-700 mb-1.5">NIP / Username Login</label>
          <input v-model="form.nip" type="text" placeholder="Masukkan NIP atau username" required class="w-full px-4 py-2.5 bg-slate-50 border border-slate-200 rounded-xl focus:ring-2 focus:ring-blue-500 text-sm">
        </div>
        <div>
          <label class="block text-sm font-semibold text-slate-700 mb-1.5">Nama Lengkap & Gelar</label>
          <input v-model="form.nama_guru" type="text" placeholder="Masukkan nama lengkap guru" required class="w-full px-4 py-2.5 bg-slate-50 border border-slate-200 rounded-xl focus:ring-2 focus:ring-blue-500 text-sm">
        </div>
        <div>
          <label class="block text-sm font-semibold text-slate-700 mb-1.5">
            Password 
            <span v-if="isEditMode" class="text-[10px] font-normal text-slate-400 italic">(Kosongkan jika tidak ingin diubah)</span>
          </label>
          <input v-model="form.password" type="password" :required="!isEditMode" placeholder="••••••••" class="w-full px-4 py-2.5 bg-slate-50 border border-slate-200 rounded-xl focus:ring-2 focus:ring-blue-500 text-sm font-mono">
        </div>
      </div>
    </BaseModal>

    <!-- Modal Preview Import -->
    <div v-if="showPreview" class="fixed inset-0 bg-slate-900/60 backdrop-blur-sm z-[60] flex items-center justify-center p-4">
      <div class="bg-white rounded-[2rem] shadow-2xl w-full max-w-4xl max-h-[85vh] flex flex-col overflow-hidden animate-in fade-in zoom-in duration-300">
        <div class="p-8 border-b border-slate-100 flex justify-between items-center bg-slate-50/50">
          <div>
            <h4 class="text-xl font-black text-slate-800">Pratinjau Data Impor</h4>
            <p class="text-sm text-slate-500">Periksa kembali data guru sebelum disimpan ke sistem</p>
          </div>
          <button @click="showPreview = false" class="p-2 bg-white hover:bg-rose-50 text-slate-400 hover:text-rose-500 rounded-xl transition-all shadow-sm">
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path></svg>
          </button>
        </div>
        <div class="flex-1 overflow-auto p-8">
          <div class="border border-slate-100 rounded-2xl overflow-hidden shadow-sm">
            <table class="min-w-full divide-y divide-slate-100">
              <thead class="bg-slate-50 sticky top-0">
                <tr>
                  <th class="px-6 py-4 text-left text-xs font-black text-slate-400 uppercase tracking-widest">NIP</th>
                  <th class="px-6 py-4 text-left text-xs font-black text-slate-400 uppercase tracking-widest">Nama Guru</th>
                  <th class="px-6 py-4 text-left text-xs font-black text-slate-400 uppercase tracking-widest">Password</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-slate-100 bg-white">
                <tr v-for="(item, idx) in previewData" :key="idx" class="hover:bg-slate-50 transition-colors">
                  <td class="px-6 py-4 text-sm font-bold text-slate-700">{{ item.nip }}</td>
                  <td class="px-6 py-4 text-sm text-slate-600 font-medium">{{ item.nama_guru }}</td>
                  <td class="px-6 py-4 text-sm font-mono text-slate-400">{{ item.password }}</td>
                </tr>
              </tbody>
            </table>
          </div>
          <p class="mt-4 text-[10px] text-slate-400 font-medium italic">* Pastikan NIP tidak ada yang sama dengan data yang sudah terdaftar.</p>
        </div>
        <div class="p-8 border-t border-slate-100 flex justify-end gap-3 bg-slate-50/50">
          <button @click="showPreview = false" class="px-8 py-3 rounded-2xl text-sm font-bold text-slate-500 hover:bg-slate-200 transition-all">Batal</button>
          <button @click="commitImport" :disabled="isUploading" class="px-10 py-3 rounded-2xl text-sm font-black bg-blue-600 text-white hover:bg-blue-700 shadow-xl shadow-blue-200 transition-all disabled:opacity-50">
            {{ isUploading ? 'Mengimpor...' : 'Konfirmasi & Impor' }}
          </button>
        </div>
      </div>
    </div>

    <!-- Modal Konfirmasi Hapus -->
    <ConfirmModal 
      :show="showDeleteConfirm"
      title="Hapus Data Guru"
      :message="deleteTargetId ? 'Apakah Anda yakin ingin menghapus data guru ini? Akun login terkait juga akan dihapus.' : `Apakah Anda yakin ingin menghapus ${selectedIds.length} data guru yang dipilih?`"
      :isLoading="isDeleting"
      @close="showDeleteConfirm = false"
      @confirm="executeDelete"
    />
  </div>
</template>
