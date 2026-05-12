<script setup>
import { ref, onMounted } from 'vue'
import { useAuthStore } from '../../store/auth'
import { useAlertStore } from '../../store/alert'

const authStore = useAuthStore()
const alertStore = useAlertStore()

const settings = ref({
  nama_sekolah: '',
  kota: '',
  alamat_sekolah: '',
  logo_sekolah: '',
  logo_kanan: '',
  nama_ujian: '',
  tahun_ajaran: '',
  semester: 'Ganjil',
  kepala_sekolah_id: '',
  ttd_kepala_sekolah: ''
})

const gurus = ref([])
const isLoading = ref(false)
const isSaving = ref(false)

const fetchGurus = async () => {
  try {
    const res = await fetch('/api/admin/guru', {
      headers: { 'Authorization': `Bearer ${authStore.token}` }
    })
    const data = await res.json()
    if (res.ok) gurus.value = data.data
  } catch (e) {
    console.error("Gagal mengambil data guru")
  }
}

const fetchSettings = async () => {
  isLoading.value = true
  try {
    const res = await fetch('/api/admin/settings', {
      headers: { 'Authorization': `Bearer ${authStore.token}` }
    })
    const data = await res.json()
    if (res.ok) {
      // Gabungkan data dari API ke ref settings kita
      Object.keys(data.data).forEach(key => {
        if (key in settings.value) {
          settings.value[key] = data.data[key]
        }
      })
    }
  } catch (e) {
    alertStore.showAlert("Gagal mengambil pengaturan", "error")
  } finally {
    isLoading.value = false
  }
}

const handleFileUpload = (e, targetKey) => {
  const file = e.target.files[0]
  if (!file) return
  
  if (file.size > 1024 * 1024) { // Max 1MB
    alertStore.showAlert("Ukuran file terlalu besar. Maksimal 1MB.", "warning")
    return
  }

  const reader = new FileReader()
  reader.onload = (event) => {
    settings.value[targetKey] = event.target.result
  }
  reader.readAsDataURL(file)
}

const saveSettings = async () => {
  isSaving.value = true
  try {
    const res = await fetch('/api/admin/settings', {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${authStore.token}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(settings.value)
    })

    if (res.ok) {
      alertStore.showAlert("Pengaturan berhasil disimpan", "success")
    } else {
      alertStore.showAlert("Gagal menyimpan pengaturan", "error")
    }
  } catch (e) {
    alertStore.showAlert("Kesalahan koneksi", "error")
  } finally {
    isSaving.value = false
  }
}

onMounted(() => {
  fetchGurus()
  fetchSettings()
})
</script>

<template>
  <div class="space-y-6 pb-12">
    <!-- Header -->
    <div class="flex flex-col sm:flex-row justify-between items-start sm:items-center bg-white p-6 rounded-2xl shadow-sm border border-slate-100 gap-4">
      <div>
        <h3 class="text-xl font-bold text-slate-800">Pengaturan Identitas Sekolah</h3>
        <p class="text-sm text-slate-500 mt-1">Kelola profil sekolah dan konfigurasi administrasi ujian</p>
      </div>
      <div class="flex items-center gap-3 w-full sm:w-auto">
        <button 
          @click="saveSettings" 
          :disabled="isSaving"
          class="flex-1 sm:flex-none flex items-center justify-center gap-2 bg-blue-600 hover:bg-blue-700 text-white px-8 py-2.5 rounded-xl text-sm font-bold transition-all shadow-lg shadow-blue-200 disabled:opacity-50"
        >
          <svg v-if="!isSaving" class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7H5a2 2 0 00-2 2v9a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-3m-1 4l-3 3m0 0l-3-3m3 3V4"></path></svg>
          <div v-else class="animate-spin w-4 h-4 border-2 border-white border-t-transparent rounded-full"></div>
          {{ isSaving ? 'Menyimpan...' : 'Simpan Perubahan' }}
        </button>
      </div>
    </div>

    <div v-if="isLoading" class="flex items-center justify-center h-64 bg-white rounded-2xl border border-slate-100">
      <div class="animate-spin w-8 h-8 border-4 border-blue-600 border-t-transparent rounded-full"></div>
    </div>

    <div v-else class="grid grid-cols-1 lg:grid-cols-3 gap-6">
      
      <!-- Panel Kiri: Media -->
      <div class="lg:col-span-1 space-y-6">
        <!-- Logo Sekolah -->
        <div class="bg-white p-6 rounded-2xl shadow-sm border border-slate-100">
          <div class="flex items-center gap-3 mb-6">
            <div class="w-8 h-8 bg-blue-50 text-blue-600 rounded-lg flex items-center justify-center">
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z"></path></svg>
            </div>
            <h4 class="font-bold text-slate-800">Logo Sekolah</h4>
          </div>
          
          <div class="aspect-square w-full bg-slate-50 rounded-xl border-2 border-dashed border-slate-200 flex items-center justify-center overflow-hidden mb-4 relative group">
            <img v-if="settings.logo_sekolah" :src="settings.logo_sekolah" class="w-full h-full object-contain p-4" />
            <div v-else class="text-center">
              <svg class="w-10 h-10 text-slate-300 mx-auto mb-2" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z"></path></svg>
              <p class="text-[10px] font-bold text-slate-400 uppercase tracking-widest">Belum ada logo</p>
            </div>
          </div>

          <label class="block w-full">
            <span class="sr-only">Upload Logo</span>
            <input type="file" @change="handleFileUpload($event, 'logo_sekolah')" accept="image/*" class="block w-full text-xs text-slate-500 file:mr-4 file:py-2 file:px-4 file:rounded-xl file:border-0 file:text-xs file:font-bold file:bg-blue-50 file:text-blue-700 hover:file:bg-blue-100 transition-all cursor-pointer" />
          </label>
        </div>

        <!-- Logo Kanan -->
        <div class="bg-white p-6 rounded-2xl shadow-sm border border-slate-100">
          <div class="flex items-center gap-3 mb-6">
            <div class="w-8 h-8 bg-amber-50 text-amber-600 rounded-lg flex items-center justify-center">
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z"></path></svg>
            </div>
            <h4 class="font-bold text-slate-800">Logo Kanan (Opsional)</h4>
          </div>
          
          <div class="aspect-square w-full bg-slate-50 rounded-xl border-2 border-dashed border-slate-200 flex items-center justify-center overflow-hidden mb-4 relative group">
            <img v-if="settings.logo_kanan" :src="settings.logo_kanan" class="w-full h-full object-contain p-4" />
            <div v-else class="text-center">
              <svg class="w-10 h-10 text-slate-300 mx-auto mb-2" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z"></path></svg>
              <p class="text-[10px] font-bold text-slate-400 uppercase tracking-widest">Kosong</p>
            </div>
          </div>

          <label class="block w-full">
            <span class="sr-only">Upload Logo Kanan</span>
            <input type="file" @change="handleFileUpload($event, 'logo_kanan')" accept="image/*" class="block w-full text-xs text-slate-500 file:mr-4 file:py-2 file:px-4 file:rounded-xl file:border-0 file:text-xs file:font-bold file:bg-amber-50 file:text-amber-700 hover:file:bg-amber-100 transition-all cursor-pointer" />
          </label>
        </div>

        <!-- Tanda Tangan -->
        <div class="bg-white p-6 rounded-2xl shadow-sm border border-slate-100">
          <div class="flex items-center gap-3 mb-6">
            <div class="w-8 h-8 bg-indigo-50 text-indigo-600 rounded-lg flex items-center justify-center">
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z"></path></svg>
            </div>
            <h4 class="font-bold text-slate-800">Tanda Tangan Kepsek</h4>
          </div>

          <div class="h-32 w-full bg-slate-50 rounded-xl border-2 border-dashed border-slate-200 flex items-center justify-center overflow-hidden mb-4">
            <img v-if="settings.ttd_kepala_sekolah" :src="settings.ttd_kepala_sekolah" class="h-full object-contain p-4" />
            <div v-else class="text-center">
              <p class="text-[10px] font-bold text-slate-400 uppercase tracking-widest">Belum ada TTD</p>
            </div>
          </div>

          <label class="block w-full">
            <span class="sr-only">Upload TTD</span>
            <input type="file" @change="handleFileUpload($event, 'ttd_kepala_sekolah')" accept="image/*" class="block w-full text-xs text-slate-500 file:mr-4 file:py-2 file:px-4 file:rounded-xl file:border-0 file:text-xs file:font-bold file:bg-indigo-50 file:text-indigo-700 hover:file:bg-indigo-100 transition-all cursor-pointer" />
          </label>
        </div>
      </div>

      <!-- Panel Kanan: Detail Form -->
      <div class="lg:col-span-2 space-y-6">
        
        <!-- Form Identitas -->
        <div class="bg-white p-8 rounded-2xl shadow-sm border border-slate-100">
          <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
            
            <div class="md:col-span-2">
              <h4 class="text-sm font-bold text-slate-800 mb-4 pb-2 border-b border-slate-50">Profil Lembaga</h4>
            </div>

            <div class="md:col-span-1 space-y-1.5">
              <label class="text-xs font-bold text-slate-600 ml-1">Nama Sekolah</label>
              <input v-model="settings.nama_sekolah" type="text" class="w-full px-4 py-2.5 bg-slate-50 border border-slate-100 rounded-xl focus:bg-white focus:ring-2 focus:ring-blue-500/20 focus:border-blue-500 outline-none transition-all text-sm font-medium text-slate-700" placeholder="Contoh: SMP Negeri 1 Jakarta">
            </div>

            <div class="md:col-span-1 space-y-1.5">
              <label class="text-xs font-bold text-slate-600 ml-1">Kabupaten/Kota</label>
              <input v-model="settings.kota" type="text" class="w-full px-4 py-2.5 bg-slate-50 border border-slate-100 rounded-xl focus:bg-white focus:ring-2 focus:ring-blue-500/20 focus:border-blue-500 outline-none transition-all text-sm font-medium text-slate-700" placeholder="Contoh: Jakarta">
            </div>

            <div class="md:col-span-2 space-y-1.5">
              <label class="text-xs font-bold text-slate-600 ml-1">Alamat Sekolah</label>
              <textarea v-model="settings.alamat_sekolah" rows="2" class="w-full px-4 py-2.5 bg-slate-50 border border-slate-100 rounded-xl focus:bg-white focus:ring-2 focus:ring-blue-500/20 focus:border-blue-500 outline-none transition-all text-sm font-medium text-slate-700" placeholder="Alamat lengkap sekolah"></textarea>
            </div>

            <div class="md:col-span-2 mt-4">
              <h4 class="text-sm font-bold text-slate-800 mb-4 pb-2 border-b border-slate-50">Konfigurasi Ujian Aktif</h4>
            </div>

            <div class="space-y-1.5">
              <label class="text-xs font-bold text-slate-600 ml-1">Nama Jenis Ujian</label>
              <input v-model="settings.nama_ujian" type="text" class="w-full px-4 py-2.5 bg-slate-50 border border-slate-100 rounded-xl focus:bg-white focus:ring-2 focus:ring-blue-500/20 focus:border-blue-500 outline-none transition-all text-sm font-medium text-slate-700" placeholder="Contoh: Penilaian Akhir Semester">
            </div>

            <div class="space-y-1.5">
              <label class="text-xs font-bold text-slate-600 ml-1">Tahun Pelajaran</label>
              <input v-model="settings.tahun_ajaran" type="text" class="w-full px-4 py-2.5 bg-slate-50 border border-slate-100 rounded-xl focus:bg-white focus:ring-2 focus:ring-blue-500/20 focus:border-blue-500 outline-none transition-all text-sm font-medium text-slate-700" placeholder="Contoh: 2024/2025">
            </div>

            <div class="space-y-1.5">
              <label class="text-xs font-bold text-slate-600 ml-1">Semester</label>
              <select v-model="settings.semester" class="w-full px-4 py-2.5 bg-slate-50 border border-slate-100 rounded-xl focus:bg-white focus:ring-2 focus:ring-blue-500/20 focus:border-blue-500 outline-none transition-all text-sm font-medium text-slate-700 cursor-pointer">
                <option value="Ganjil">Semester Ganjil</option>
                <option value="Genap">Semester Genap</option>
              </select>
            </div>

            <div class="md:col-span-2 mt-4">
              <h4 class="text-sm font-bold text-slate-800 mb-4 pb-2 border-b border-slate-50">Pejabat Berwenang</h4>
            </div>

            <div class="md:col-span-2 space-y-1.5">
              <label class="text-xs font-bold text-slate-600 ml-1">Kepala Sekolah</label>
              <select v-model="settings.kepala_sekolah_id" class="w-full px-4 py-2.5 bg-slate-50 border border-slate-100 rounded-xl focus:bg-white focus:ring-2 focus:ring-blue-500/20 focus:border-blue-500 outline-none transition-all text-sm font-medium text-slate-700 cursor-pointer">
                <option value="">-- Pilih Kepala Sekolah dari Data Guru --</option>
                <option v-for="g in gurus" :key="g.id" :value="g.id.toString()">
                  {{ g.nama_guru }} (NIP. {{ g.nip || '-' }})
                </option>
              </select>
              <div class="flex items-center gap-2 mt-2 px-3 py-2 bg-amber-50 rounded-lg border border-amber-100">
                <svg class="w-4 h-4 text-amber-600" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path></svg>
                <p class="text-[11px] font-medium text-amber-700">Nama dan NIP akan otomatis terupdate jika data Guru diubah.</p>
              </div>
            </div>

          </div>
        </div>
      </div>

    </div>
  </div>
</template>
