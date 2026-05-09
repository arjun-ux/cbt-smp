<script setup>
import { ref, onMounted, computed } from 'vue'
import { useAuthStore } from '../../store/auth'
import { useAlertStore } from '../../store/alert'
import { useRouter } from 'vue-router'
import BaseModal from '../../components/BaseModal.vue'
import ConfirmModal from '../../components/ConfirmModal.vue'

const authStore = useAuthStore()
const alertStore = useAlertStore()
const router = useRouter()
const bankSoals = ref([])
const mapels = ref([])
const gurus = ref([])
const isLoading = ref(false)
const isSaving = ref(false)
const showForm = ref(false)
const isEditMode = ref(false)
const editId = ref(null)

const showDeleteConfirm = ref(false)
const deleteTargetId = ref(null)
const isDeleting = ref(false)

const form = ref({
  mapel_id: '',
  guru_id: '',
  tingkat_kelas: '',
  judul_bank_soal: '',
  default_bobot_pg: 1,
  default_bobot_essay: 1,
  status: 'Draft'
})

const apiPrefix = computed(() => authStore.user?.role === 'admin' ? '/api/admin' : '/api/guru')
const userRole = computed(() => authStore.user?.role)

const fetchData = async () => {
  isLoading.value = true
  const headers = { 'Authorization': `Bearer ${authStore.token}` }
  try {
    const fetchers = [
      fetch(`${apiPrefix.value}/bank-soal`, { headers }),
      fetch(`${apiPrefix.value}/mapel`, { headers })
    ]
    
    if (userRole.value === 'admin') {
      fetchers.push(fetch('/api/admin/guru', { headers }))
    }

    const responses = await Promise.all(fetchers)
    const [dbs, dm] = await Promise.all([responses[0].json(), responses[1].json()])
    
    bankSoals.value = dbs.data || []
    mapels.value = dm.data || []
    
    if (userRole.value === 'admin') {
      const dg = await responses[2].json()
      gurus.value = dg.data || []
    }
  } catch (error) {
    alertStore.showAlert("Gagal memuat data", "error")
  } finally {
    isLoading.value = false
  }
}

const openAdd = () => {
  isEditMode.value = false
  editId.value = null
  form.value = { 
    mapel_id: '', 
    guru_id: '', 
    tingkat_kelas: '', 
    judul_bank_soal: '', 
    default_bobot_pg: 1,
    default_bobot_essay: 1,
    status: 'Draft' 
  }
  showForm.value = true
}

const openEdit = (bs) => {
  isEditMode.value = true
  editId.value = bs.id
  form.value = {
    mapel_id: bs.mapel_id,
    guru_id: bs.guru_id,
    tingkat_kelas: bs.tingkat_kelas,
    judul_bank_soal: bs.judul_bank_soal,
    default_bobot_pg: bs.default_bobot_pg || 1,
    default_bobot_essay: bs.default_bobot_essay || 1,
    status: bs.status
  }
  showForm.value = true
}

const saveBankSoal = async () => {
  isSaving.value = true
  try {
    const url = isEditMode.value ? `${apiPrefix.value}/bank-soal/${editId.value}` : `${apiPrefix.value}/bank-soal`
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
      alertStore.showAlert(isEditMode.value ? "Data diperbarui" : "Wadah Bank Soal dibuat", "success")
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

const confirmDelete = (id) => {
  deleteTargetId.value = id
  showDeleteConfirm.value = true
}

const executeDelete = async () => {
  isDeleting.value = true
  try {
    const res = await fetch(`${apiPrefix.value}/bank-soal/${deleteTargetId.value}`, {
      method: 'DELETE',
      headers: { 'Authorization': `Bearer ${authStore.token}` }
    })
    if(res.ok) {
      alertStore.showAlert("Bank Soal berhasil dihapus", "success")
      showDeleteConfirm.value = false
      fetchData()
    } else {
      const data = await res.json()
      alertStore.showAlert(data.error || "Gagal menghapus Bank Soal", "error")
      showDeleteConfirm.value = false
    }
  } catch (error) {
    alertStore.showAlert("Kesalahan sistem saat menghapus data", "error")
  } finally {
    isDeleting.value = false
  }
}

const manageQuestions = (id) => {
  const routeName = userRole.value === 'admin' ? 'AdminSoalDetail' : 'GuruSoalDetail'
  router.push({ name: routeName, params: { bankSoalId: id } })
}

onMounted(fetchData)
</script>

<template>
  <div class="space-y-6 pb-20">
    <!-- Header -->
    <div class="flex flex-col sm:flex-row justify-between items-start sm:items-center bg-white p-6 rounded-2xl shadow-sm border border-slate-100 gap-4">
      <div>
        <h3 class="text-xl font-bold text-slate-800">Manajemen Bank Soal</h3>
        <p class="text-sm text-slate-500 mt-1">Buat wadah ujian untuk dikelola oleh Guru</p>
      </div>
      <button @click="openAdd" class="flex items-center gap-2 bg-blue-600 hover:bg-blue-700 text-white px-6 py-2.5 rounded-xl text-sm font-bold transition-all shadow-lg shadow-blue-200">
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"></path></svg>
        Buat Wadah Baru
      </button>
    </div>

    <!-- Table Utama -->
    <div class="bg-white rounded-2xl shadow-sm border border-slate-100 overflow-hidden">
      <div class="overflow-x-auto">
        <table class="min-w-full divide-y divide-slate-100">
          <thead class="bg-slate-50">
            <tr>
              <th class="px-6 py-4 text-left text-xs font-bold text-slate-500 uppercase tracking-wider">No</th>
              <th class="px-6 py-4 text-left text-xs font-bold text-slate-500 uppercase tracking-wider">Info Bank Soal</th>
              <th class="px-6 py-4 text-left text-xs font-bold text-slate-500 uppercase tracking-wider">Pengampu</th>
              <th class="px-6 py-4 text-left text-xs font-bold text-slate-500 uppercase tracking-wider">Status</th>
              <th class="px-6 py-4 text-right text-xs font-bold text-slate-500 uppercase tracking-wider">Aksi</th>
            </tr>
          </thead>
          <tbody class="bg-white divide-y divide-slate-50">
            <tr v-if="isLoading">
              <td colspan="5" class="px-6 py-12 text-center text-slate-400">
                <div class="flex flex-col items-center gap-2">
                  <svg class="w-8 h-8 animate-spin text-blue-500" fill="none" viewBox="0 0 24 24"><circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle><path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path></svg>
                  <span class="text-sm font-medium">Memuat data bank soal...</span>
                </div>
              </td>
            </tr>
            <tr v-else-if="bankSoals.length === 0">
              <td colspan="5" class="px-6 py-12 text-center text-slate-400 font-medium italic">Belum ada wadah bank soal.</td>
            </tr>
            <tr v-for="(bs, index) in bankSoals" :key="bs.id" class="hover:bg-slate-50/50 transition-colors group">
              <td class="px-6 py-4 text-sm text-slate-500 font-bold">{{ index + 1 }}</td>
              <td class="px-6 py-4">
                <div class="flex flex-col">
                  <span class="text-sm font-black text-slate-800 group-hover:text-blue-600 transition-colors leading-tight">{{ bs.judul_bank_soal }}</span>
                  <div class="flex items-center gap-2 mt-1.5">
                    <span class="px-2 py-0.5 bg-blue-50 text-blue-600 text-[10px] font-black rounded border border-blue-100 uppercase">{{ bs.mapel?.nama_mapel }}</span>
                    <span class="text-[10px] text-slate-400 font-bold uppercase tracking-widest">Kelas {{ bs.tingkat_kelas }}</span>
                  </div>
                </div>
              </td>
              <td class="px-6 py-4">
                <div class="flex items-center gap-3">
                  <div class="w-9 h-9 rounded-xl bg-blue-600 flex items-center justify-center text-xs font-black text-white shadow-lg shadow-blue-200 ring-2 ring-white">
                    {{ bs.guru?.nama_guru?.charAt(0) || 'A' }}
                  </div>
                  <span class="text-sm font-bold text-slate-700">{{ bs.guru?.nama_guru || 'Administrator' }}</span>
                </div>
              </td>
              <td class="px-6 py-4">
                <span v-if="bs.status === 'Aktif'" class="px-2.5 py-1 text-[10px] font-black uppercase rounded-md bg-emerald-100 text-emerald-700 border border-emerald-200">Aktif</span>
                <span v-else class="px-2.5 py-1 text-[10px] font-black uppercase rounded-md bg-amber-100 text-amber-700 border border-amber-200">Draft</span>
              </td>
              <td class="px-6 py-4 text-right">
                <div class="flex justify-end gap-2">
                  <button @click="manageQuestions(bs.id)" class="px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-xl text-xs font-black uppercase tracking-wider transition-all shadow-md shadow-blue-100 active:scale-95">Kelola Soal</button>
                  <button @click="openEdit(bs)" class="p-2 text-slate-400 hover:text-blue-600 hover:bg-blue-50 rounded-xl transition-all" title="Edit">
                    <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"></path></svg>
                  </button>
                  <button @click="confirmDelete(bs.id)" class="p-2 text-slate-400 hover:text-rose-600 hover:bg-rose-50 rounded-xl transition-all" title="Hapus">
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
      :title="isEditMode ? 'Edit Wadah Bank Soal' : 'Buat Wadah Baru'"
      :isLoading="isSaving"
      @close="showForm = false"
      @confirm="saveBankSoal"
    >
      <div class="space-y-4">
        <div>
          <label class="block text-sm font-semibold text-slate-700 mb-1.5">Judul Bank Soal</label>
          <input v-model="form.judul_bank_soal" type="text" required placeholder="Contoh: PAS Ganjil Matematika Kelas 7" class="w-full px-4 py-2.5 bg-slate-50 border border-slate-200 rounded-xl focus:ring-2 focus:ring-blue-500 text-sm">
        </div>
        
        <div class="grid grid-cols-2 gap-4">
          <div>
            <label class="block text-sm font-semibold text-slate-700 mb-1.5">Tingkat Kelas</label>
            <select v-model="form.tingkat_kelas" required class="w-full px-4 py-2.5 bg-slate-50 border border-slate-200 rounded-xl focus:ring-2 focus:ring-blue-500 text-sm">
              <option value="">-- Pilih --</option>
              <option value="7">Kelas 7</option>
              <option value="8">Kelas 8</option>
              <option value="9">Kelas 9</option>
            </select>
          </div>
          <div>
            <label class="block text-sm font-semibold text-slate-700 mb-1.5">Mata Pelajaran</label>
            <select v-model="form.mapel_id" required class="w-full px-4 py-2.5 bg-slate-50 border border-slate-200 rounded-xl focus:ring-2 focus:ring-blue-500 text-sm">
              <option value="">-- Pilih --</option>
              <option v-for="m in mapels" :key="m.id" :value="m.id">{{ m.nama_mapel }}</option>
            </select>
          </div>
        </div>

        <div v-if="userRole === 'admin'">
          <label class="block text-sm font-semibold text-slate-700 mb-1.5">Guru Pengampu</label>
          <select v-model="form.guru_id" class="w-full px-4 py-2.5 bg-slate-50 border border-slate-200 rounded-xl focus:ring-2 focus:ring-blue-500 text-sm">
            <option value="">-- Tanpa Penugasan (Admin) --</option>
            <option v-for="g in gurus" :key="g.id" :value="g.id">{{ g.nama_guru }}</option>
          </select>
        </div>

        <div class="grid grid-cols-2 gap-4">
          <div>
            <label class="block text-sm font-semibold text-slate-700 mb-1.5">Bobot PG Standar</label>
            <input v-model.number="form.default_bobot_pg" type="number" step="0.1" required class="w-full px-4 py-2.5 bg-slate-50 border border-slate-200 rounded-xl focus:ring-2 focus:ring-blue-500 text-sm">
          </div>
          <div>
            <label class="block text-sm font-semibold text-slate-700 mb-1.5">Bobot Essay Standar</label>
            <input v-model.number="form.default_bobot_essay" type="number" step="0.1" required class="w-full px-4 py-2.5 bg-slate-50 border border-slate-200 rounded-xl focus:ring-2 focus:ring-blue-500 text-sm">
          </div>
        </div>

        <div v-if="isEditMode" class="p-4 bg-slate-50 border border-slate-100 rounded-2xl">
          <label class="block text-sm font-bold text-slate-700 mb-3">Status Publikasi</label>
          <div class="flex gap-6">
            <label class="flex items-center gap-2 cursor-pointer group">
              <input type="radio" v-model="form.status" value="Draft" class="w-5 h-5 text-blue-600 focus:ring-blue-500 border-slate-300">
              <span class="text-sm font-bold text-slate-600 group-hover:text-blue-600 transition-colors">Draft</span>
            </label>
            <label class="flex items-center gap-2 cursor-pointer group">
              <input type="radio" v-model="form.status" value="Aktif" class="w-5 h-5 text-blue-600 focus:ring-blue-500 border-slate-300">
              <span class="text-sm font-bold text-slate-600 group-hover:text-blue-600 transition-colors">Aktif</span>
            </label>
          </div>
        </div>
      </div>
    </BaseModal>

    <!-- Delete Confirmation -->
    <ConfirmModal 
      :show="showDeleteConfirm"
      title="Hapus Bank Soal"
      message="Menghapus wadah ini akan menghapus seluruh soal yang ada di dalamnya secara permanen. Lanjutkan?"
      :isLoading="isDeleting"
      @close="showDeleteConfirm = false"
      @confirm="executeDelete"
    />
  </div>
</template>
