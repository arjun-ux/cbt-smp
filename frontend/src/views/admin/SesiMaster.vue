<script setup>
import { ref, onMounted, computed, watch } from 'vue'
import { useAuthStore } from '../../store/auth'
import { useAlertStore } from '../../store/alert'
import BaseModal from '../../components/BaseModal.vue'
import ConfirmModal from '../../components/ConfirmModal.vue'

const authStore = useAuthStore()
const alertStore = useAlertStore()

const items = ref([])
const searchQuery = ref('')
const isLoading = ref(false)
const isSaving = ref(false)
const showForm = ref(false)
const isEditMode = ref(false)
const editId = ref(null)

const selectedIds = ref([])
const showDeleteConfirm = ref(false)
const deleteTargetId = ref(null)
const isDeleting = ref(false)

// Pagination State
const currentPage = ref(1)
const itemsPerPage = 10

const form = ref({
  nama_sesi: '',
  waktu_mulai: '',
  waktu_selesai: ''
})

const fetchData = async () => {
  isLoading.value = true
  try {
    const res = await fetch('/api/admin/sesi', {
      headers: { 'Authorization': `Bearer ${authStore.token}` }
    })
    const data = await res.json()
    if(res.ok) items.value = data.data
  } catch (error) {
    alertStore.showAlert("Gagal memuat data sesi", "error")
  } finally {
    isLoading.value = false
  }
}

const filteredItems = computed(() => {
  if (!searchQuery.value) return items.value
  const q = searchQuery.value.toLowerCase()
  return items.value.filter(i => 
    i.nama_sesi.toLowerCase().includes(q)
  )
})

const totalPages = computed(() => Math.ceil(filteredItems.value.length / itemsPerPage))
const paginatedItems = computed(() => {
  const start = (currentPage.value - 1) * itemsPerPage
  const end = start + itemsPerPage
  return filteredItems.value.slice(start, end)
})

watch(searchQuery, () => {
  currentPage.value = 1
})

const openAdd = () => {
  isEditMode.value = false
  editId.value = null
  form.value = { nama_sesi: '', waktu_mulai: '', waktu_selesai: '' }
  showForm.value = true
}

const openEdit = (item) => {
  isEditMode.value = true
  editId.value = item.id
  form.value = {
    nama_sesi: item.nama_sesi,
    waktu_mulai: item.waktu_mulai,
    waktu_selesai: item.waktu_selesai
  }
  showForm.value = true
}

const saveItem = async () => {
  isSaving.value = true
  try {
    const url = isEditMode.value ? `/api/admin/sesi/${editId.value}` : '/api/admin/sesi'
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
      alertStore.showAlert(isEditMode.value ? "Data diperbarui" : "Sesi berhasil ditambahkan", "success")
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
      const res = await fetch(`/api/admin/sesi/${deleteTargetId.value}`, {
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
        const res = await fetch(`/api/admin/sesi/${id}`, {
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
    
    fetchData()
  } catch (error) {
    alertStore.showAlert("Kesalahan sistem saat menghapus", "error")
  } finally {
    isDeleting.value = false
    showDeleteConfirm.value = false
  }
}

const toggleSelectAll = (e) => {
  if (e.target.checked) {
    selectedIds.value = paginatedItems.value.map(i => i.id)
  } else {
    selectedIds.value = []
  }
}

onMounted(fetchData)
</script>

<template>
  <div class="space-y-6 pb-12">
    <!-- Header -->
    <div class="flex flex-col sm:flex-row justify-between items-start sm:items-center bg-white p-6 rounded-2xl shadow-sm border border-slate-100 gap-4">
      <div>
        <h3 class="text-xl font-bold text-slate-800">Manajemen Sesi</h3>
        <p class="text-sm text-slate-500 mt-1">Atur pembagian waktu pelaksanaan ujian</p>
      </div>
      <button @click="openAdd" class="flex items-center gap-2 bg-blue-600 hover:bg-blue-700 text-white px-6 py-2.5 rounded-xl text-sm font-bold transition-all shadow-lg shadow-blue-200">
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"></path></svg>
        Tambah Sesi
      </button>
    </div>

    <!-- Tabel Utama -->
    <div class="bg-white rounded-2xl shadow-sm border border-slate-100 overflow-hidden">
      <!-- Search & Bulk Actions -->
      <div class="p-4 border-b border-slate-50 bg-slate-50/30 flex flex-col sm:flex-row justify-between items-center gap-4">
        <div class="relative w-full sm:w-80">
          <svg class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"></path></svg>
          <input v-model="searchQuery" type="text" placeholder="Cari nama sesi..." class="block w-full pl-10 pr-4 py-2.5 bg-white border border-slate-200 rounded-xl text-sm focus:ring-2 focus:ring-blue-500 transition-all">
        </div>
        
        <div v-if="selectedIds.length > 0" class="flex items-center gap-3 animate-in fade-in slide-in-from-right-4 duration-300">
          <span class="text-sm font-semibold text-slate-600">{{ selectedIds.length }} item dipilih</span>
          <button @click="confirmDelete()" class="flex items-center gap-2 bg-rose-50 text-rose-600 px-4 py-2 rounded-xl text-sm font-bold hover:bg-rose-100 transition-colors">
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path></svg>
            Hapus Masal
          </button>
        </div>

        <div v-else class="text-xs font-bold text-slate-400 bg-slate-100 px-3 py-1.5 rounded-lg uppercase tracking-wider">
          Total: {{ filteredItems.length }} Sesi
        </div>
      </div>

      <div class="overflow-x-auto">
        <table class="min-w-full divide-y divide-slate-100">
          <thead class="bg-slate-50">
            <tr>
              <th class="px-6 py-4 w-10">
                <input type="checkbox" :checked="selectedIds.length === paginatedItems.length && paginatedItems.length > 0" @change="toggleSelectAll" class="rounded border-slate-300 text-blue-600 focus:ring-blue-500">
              </th>
              <th class="px-6 py-4 text-left text-xs font-bold text-slate-500 uppercase tracking-wider">Nama Sesi</th>
              <th class="px-6 py-4 text-left text-xs font-bold text-slate-500 uppercase tracking-wider">Waktu Mulai</th>
              <th class="px-6 py-4 text-left text-xs font-bold text-slate-500 uppercase tracking-wider">Waktu Selesai</th>
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
            <tr v-else-if="filteredItems.length === 0">
              <td colspan="5" class="px-6 py-12 text-center text-slate-400 italic">Data sesi tidak ditemukan.</td>
            </tr>
            <tr v-for="item in paginatedItems" :key="item.id" class="hover:bg-slate-50/50 transition-colors group">
              <td class="px-6 py-4">
                <input type="checkbox" v-model="selectedIds" :value="item.id" class="rounded border-slate-300 text-blue-600 focus:ring-blue-500">
              </td>
              <td class="px-6 py-4 text-sm font-bold text-slate-800">{{ item.nama_sesi }}</td>
              <td class="px-6 py-4 text-sm text-slate-600 font-medium">{{ item.waktu_mulai }}</td>
              <td class="px-6 py-4 text-sm text-slate-600 font-medium">{{ item.waktu_selesai }}</td>
              <td class="px-6 py-4 text-right">
                <div class="flex justify-end gap-2">
                  <button @click="openEdit(item)" class="p-2 text-slate-400 hover:text-blue-600 hover:bg-blue-50 rounded-xl transition-all" title="Edit">
                    <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"></path></svg>
                  </button>
                  <button @click="confirmDelete(item.id)" class="p-2 text-slate-400 hover:text-rose-600 hover:bg-rose-50 rounded-xl transition-all" title="Hapus">
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
      :title="isEditMode ? 'Edit Sesi' : 'Tambah Sesi Baru'"
      :isLoading="isSaving"
      @close="showForm = false"
      @confirm="saveItem"
    >
      <div class="space-y-4">
        <div>
          <label class="block text-sm font-semibold text-slate-700 mb-1.5">Nama Sesi</label>
          <input v-model="form.nama_sesi" type="text" placeholder="Contoh: Sesi 1, Sesi 2" required class="w-full px-4 py-2.5 bg-slate-50 border border-slate-200 rounded-xl focus:ring-2 focus:ring-blue-500 text-sm">
        </div>
        <div class="grid grid-cols-2 gap-4">
          <div>
            <label class="block text-sm font-semibold text-slate-700 mb-1.5">Waktu Mulai</label>
            <input v-model="form.waktu_mulai" type="time" required class="w-full px-4 py-2.5 bg-slate-50 border border-slate-200 rounded-xl focus:ring-2 focus:ring-blue-500 text-sm">
          </div>
          <div>
            <label class="block text-sm font-semibold text-slate-700 mb-1.5">Waktu Selesai</label>
            <input v-model="form.waktu_selesai" type="time" required class="w-full px-4 py-2.5 bg-slate-50 border border-slate-200 rounded-xl focus:ring-2 focus:ring-blue-500 text-sm">
          </div>
        </div>
      </div>
    </BaseModal>

    <!-- Modal Konfirmasi Hapus -->
    <ConfirmModal 
      :show="showDeleteConfirm"
      title="Hapus Sesi"
      :message="deleteTargetId ? 'Apakah Anda yakin ingin menghapus sesi ini?' : `Apakah Anda yakin ingin menghapus ${selectedIds.length} sesi yang dipilih?`"
      :isLoading="isDeleting"
      @close="showDeleteConfirm = false"
      @confirm="executeDelete"
    />
  </div>
</template>
