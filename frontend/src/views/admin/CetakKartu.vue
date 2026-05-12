<script setup>
import { ref, onMounted } from 'vue'
import { useAuthStore } from '../../store/auth'
import { useAlertStore } from '../../store/alert'

const authStore = useAuthStore()
const alertStore = useAlertStore()

const filters = ref({
  kelas_id: '',
  ruang_id: '',
  sesi_id: ''
})

const masterData = ref({
  kelases: [],
  ruangs: [],
  sesis: []
})

// Konfigurasi Live untuk Header & Footer (Hanya yang bersifat teks dinamis)
const config = ref({
  nama_ujian: '',
  nama_sekolah: '',
  alamat_sekolah: '',
  kota: '',
  tanggal: '',
  logo_sekolah: '',
  logo_kanan: '',
  ttd_kepala_sekolah: '',
  kepala_sekolah_nama: '',
  kepala_sekolah_nip: ''
})

const dataCetak = ref(null)
const isLoading = ref(false)

const fetchMasterData = async () => {
  const headers = { 'Authorization': `Bearer ${authStore.token}` }
  try {
    const [resK, resR, resS] = await Promise.all([
      fetch('/api/admin/kelas', { headers }),
      fetch('/api/admin/ruang', { headers }),
      fetch('/api/admin/sesi', { headers })
    ])
    const [dk, dr, ds] = await Promise.all([resK.json(), resR.json(), resS.json()])
    masterData.value.kelases = dk.data || []
    masterData.value.ruangs = dr.data || []
    masterData.value.sesis = ds.data || []
    
    // Fetch settings awal untuk default config
    const resSet = await fetch('/api/admin/settings', { headers })
    const dSet = await resSet.json()
    if (resSet.ok && dSet.data) {
      const s = dSet.data
      config.value.nama_ujian = s.nama_ujian || ''
      config.value.nama_sekolah = s.nama_sekolah || ''
      config.value.alamat_sekolah = s.alamat_sekolah || ''
      config.value.kota = s.kota || ''
      config.value.logo_sekolah = s.logo_sekolah || ''
      config.value.logo_kanan = s.logo_kanan || ''
      config.value.ttd_kepala_sekolah = s.ttd_kepala_sekolah || ''
      
      if (!config.value.tanggal) {
        config.value.tanggal = new Date().toLocaleDateString('id-ID', { day: 'numeric', month: 'long', year: 'numeric' })
      }
    }
  } catch (e) {
    console.error("Gagal mengambil master data")
  }
}

const loadDataCetak = async () => {
  isLoading.value = true
  const query = new URLSearchParams(filters.value).toString()
  try {
    const res = await fetch(`/api/admin/cetak/kartu?${query}`, {
      headers: { 'Authorization': `Bearer ${authStore.token}` }
    })
    const data = await res.json()
    if (res.ok) {
      dataCetak.value = data.data
      // Update config jika ada data sekolah dari cetak handler (termasuk nama kepsek yang sudah di-join)
      if (data.data.sekolah) {
        const s = data.data.sekolah
        config.value.kepala_sekolah_nama = s.kepala_sekolah_nama || ''
        config.value.kepala_sekolah_nip = s.kepala_sekolah_nip || ''
        
        // Update field lain jika masih kosong
        if (!config.value.nama_ujian) config.value.nama_ujian = s.nama_ujian || ''
        if (!config.value.nama_sekolah) config.value.nama_sekolah = s.nama_sekolah || ''
        if (!config.value.alamat_sekolah) config.value.alamat_sekolah = s.alamat_sekolah || ''
      }
    } else {
      alertStore.showAlert("Gagal memuat data cetak", "error")
    }
  } catch (e) {
    alertStore.showAlert("Kesalahan koneksi", "error")
  } finally {
    isLoading.value = false
  }
}

const printPage = () => {
  window.print()
}

onMounted(() => {
  fetchMasterData()
})
</script>

<template>
  <div class="space-y-6 pb-20">
    <!-- Main Header -->
    <div class="flex flex-col sm:flex-row justify-between items-start sm:items-center bg-white p-6 rounded-2xl shadow-sm border border-slate-100 gap-4 no-print">
      <div>
        <h3 class="text-xl font-bold text-slate-800">Manajemen Cetak Kartu</h3>
        <p class="text-sm text-slate-500 mt-1">Atur desain dan filter peserta ujian</p>
      </div>
      <button 
        v-if="dataCetak"
        @click="printPage"
        class="flex items-center gap-2 bg-slate-900 text-white px-8 py-3 rounded-xl text-sm font-bold transition-all hover:bg-blue-600 shadow-xl shadow-slate-200"
      >
        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 17h2a2 2 0 002-2v-4a2 2 0 00-2-2H5a2 2 0 00-2 2v4a2 2 0 002 2h2m2 4h6a2 2 0 002-2v-4a2 2 0 00-2-2H9a2 2 0 00-2 2v4a2 2 0 002 2zm8-12V5a2 2 0 00-2-2H9a2 2 0 00-2 2v4h10z"></path></svg>
        CETAK SEKARANG
      </button>
    </div>

    <!-- Layout Grid (Top Section) -->
    <div class="grid grid-cols-1 lg:grid-cols-12 gap-8 items-start no-print">
      
      <!-- SISI KIRI: Konfigurasi -->
      <div class="lg:col-span-4 space-y-6">
        <div class="bg-white p-6 rounded-2xl shadow-sm border border-slate-100">
          <div class="flex items-center gap-3 mb-6">
            <div class="w-10 h-10 bg-blue-50 text-blue-600 rounded-xl flex items-center justify-center">
              <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"></path></svg>
            </div>
            <div>
              <h4 class="font-bold text-slate-800">Konfigurasi Kartu</h4>
              <p class="text-[11px] text-slate-500 uppercase tracking-wider font-medium">Header & Footer Settings</p>
            </div>
          </div>

          <div class="space-y-5">
            <div class="space-y-1.5">
              <label class="text-[11px] font-bold text-slate-500 uppercase ml-1">Judul Ujian (Baris 1)</label>
              <textarea v-model="config.nama_ujian" rows="1" class="w-full px-4 py-3 bg-slate-50 border border-slate-100 rounded-xl text-sm focus:ring-2 focus:ring-blue-500/10 focus:bg-white outline-none transition-all"></textarea>
            </div>
            <div class="space-y-1.5">
              <label class="text-[11px] font-bold text-slate-500 uppercase ml-1">Nama Instansi (Baris 2)</label>
              <textarea v-model="config.nama_sekolah" rows="1" class="w-full px-4 py-3 bg-slate-50 border border-slate-100 rounded-xl text-sm font-bold focus:ring-2 focus:ring-blue-500/10 focus:bg-white outline-none transition-all"></textarea>
            </div>
            <div class="space-y-1.5">
              <label class="text-[11px] font-bold text-slate-500 uppercase ml-1">Alamat/Lainnya (Baris 3)</label>
              <textarea v-model="config.alamat_sekolah" rows="2" class="w-full px-4 py-3 bg-slate-50 border border-slate-100 rounded-xl text-xs focus:ring-2 focus:ring-blue-500/10 focus:bg-white outline-none transition-all"></textarea>
            </div>
            <div class="grid grid-cols-2 gap-4 pt-2">
              <div class="space-y-1.5">
                <label class="text-[11px] font-bold text-slate-500 uppercase ml-1">Kota</label>
                <input v-model="config.kota" type="text" class="w-full px-4 py-2.5 bg-slate-50 border border-slate-100 rounded-xl text-sm focus:ring-2 focus:ring-blue-500/10 focus:bg-white outline-none transition-all">
              </div>
              <div class="space-y-1.5">
                <label class="text-[11px] font-bold text-slate-500 uppercase ml-1">Tanggal</label>
                <input v-model="config.tanggal" type="text" class="w-full px-4 py-2.5 bg-slate-50 border border-slate-100 rounded-xl text-sm focus:ring-2 focus:ring-blue-500/10 focus:bg-white outline-none transition-all">
              </div>
            </div>
          </div>

          <div class="mt-8 p-4 bg-blue-50 rounded-2xl border border-blue-100">
            <div class="flex gap-3">
              <svg class="w-5 h-5 text-blue-600 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path></svg>
              <p class="text-[11px] text-blue-700 leading-relaxed font-medium">
                Perubahan di panel ini bersifat langsung pada pratinjau. Data Kepala Sekolah diambil otomatis dari sistem.
              </p>
            </div>
          </div>
        </div>
      </div>

      <!-- SISI KANAN: Preview & Filter -->
      <div class="lg:col-span-8 space-y-8">
        <!-- 1. Live Preview Card -->
        <div class="bg-white p-8 rounded-3xl shadow-sm border border-slate-100 flex flex-col items-center">
          <h4 class="text-sm font-bold text-slate-800 mb-6 self-start flex items-center gap-2">
            <span class="w-2 h-6 bg-blue-600 rounded-full"></span>
            Pratinjau Desain
          </h4>
          
          <div class="exam-card relative border border-black p-2 bg-white w-[100mm] h-[70mm] flex flex-col overflow-hidden font-serif shadow-2xl">
              <div class="grid grid-cols-[50px_1fr_50px] gap-2 items-center border-b-2 border-black pb-1 mb-1.5">
                <div class="w-[50px] h-[50px] flex items-center justify-center">
                  <img v-if="config.logo_sekolah" :src="config.logo_sekolah" class="w-full h-full object-contain" />
                </div>
                <div class="text-center leading-tight">
                  <h4 class="text-[10px] font-black uppercase tracking-tight text-black">{{ config.nama_ujian || 'JUDUL UJIAN' }}</h4>
                  <h2 class="text-[13px] font-black uppercase text-black leading-none tracking-tighter">{{ config.nama_sekolah || 'NAMA SEKOLAH' }}</h2>
                  <p class="text-[8.5px] font-normal italic leading-none text-slate-800">{{ config.alamat_sekolah || 'Alamat Lengkap Sekolah' }}</p>
                </div>
                <div class="w-[50px] h-[50px] flex items-center justify-center">
                  <img v-if="config.logo_kanan" :src="config.logo_kanan" class="w-full h-full object-contain" />
                </div>
              </div>
              <div class="px-4 space-y-0 flex-1">
                <div v-for="label in ['Nomor Peserta', 'Nama Siswa', 'NIS / NISN', 'Kelas / Sesi', 'Ruang Ujian', 'Username', 'Password']" :key="label" class="grid grid-cols-[90px_10px_1fr] items-baseline text-[10.5px] leading-tight">
                  <span>{{ label }}</span>
                  <span class="text-center">:</span>
                  <span class="font-bold">Contoh Data</span>
                </div>
              </div>
              <div class="mt-2 flex justify-between items-end px-4 pb-1">
                <div class="w-[50px] h-[65px] border border-black flex items-center justify-center p-0.5 bg-white mb-1 uppercase text-[7px] text-slate-300 italic">Foto</div>
                <div class="text-center text-[8.5px] leading-tight min-w-[140px]">
                  <p class="mb-0.5">{{ config.kota || 'Kota' }}, {{ config.tanggal }}</p>
                  <p class="font-bold">Kepala Sekolah,</p>
                  <div class="h-8 py-0.5 flex justify-center items-center"><img v-if="config.ttd_kepala_sekolah" :src="config.ttd_kepala_sekolah" class="h-full object-contain" /></div>
                  <p class="font-bold underline uppercase">{{ config.kepala_sekolah_nama || 'Nama Kepala Sekolah' }}</p>
                  <p class="text-[7.5px]">NIP. {{ config.kepala_sekolah_nip || '-' }}</p>
                </div>
              </div>
          </div>
        </div>

        <!-- 2. Opsi Cetak & Filter -->
        <div class="bg-white p-8 rounded-3xl shadow-sm border border-slate-100">
          <h4 class="text-sm font-bold text-slate-800 mb-6 flex items-center gap-2">
            <span class="w-2 h-6 bg-emerald-500 rounded-full"></span>
            Pilih Data Peserta
          </h4>
          <div class="grid grid-cols-1 md:grid-cols-4 gap-6 items-end">
            <div class="space-y-1.5">
              <label class="text-[11px] font-bold text-slate-500 uppercase ml-1">Kelas / Rombel</label>
              <select v-model="filters.kelas_id" class="w-full px-4 py-2.5 bg-slate-50 border border-slate-100 rounded-xl focus:bg-white focus:border-blue-500 outline-none text-sm font-medium">
                <option value="">Semua Kelas</option>
                <option v-for="k in masterData.kelases" :key="k.id" :value="k.id">{{ k.nama_kelas }}</option>
              </select>
            </div>
            <div class="space-y-1.5">
              <label class="text-[11px] font-bold text-slate-500 uppercase ml-1">Ruang</label>
              <select v-model="filters.ruang_id" class="w-full px-4 py-2.5 bg-slate-50 border border-slate-100 rounded-xl focus:bg-white focus:border-blue-500 outline-none text-sm font-medium">
                <option value="">Semua Ruang</option>
                <option v-for="r in masterData.ruangs" :key="r.id" :value="r.id">{{ r.nama_ruang }}</option>
              </select>
            </div>
            <div class="space-y-1.5">
              <label class="text-[11px] font-bold text-slate-500 uppercase ml-1">Sesi</label>
              <select v-model="filters.sesi_id" class="w-full px-4 py-2.5 bg-slate-50 border border-slate-100 rounded-xl focus:bg-white focus:border-blue-500 outline-none text-sm font-medium">
                <option value="">Semua Sesi</option>
                <option v-for="s in masterData.sesis" :key="s.id" :value="s.id">{{ s.nama_sesi }}</option>
              </select>
            </div>
            <button 
              @click="loadDataCetak" 
              :disabled="isLoading"
              class="bg-blue-600 hover:bg-blue-700 text-white h-[42px] px-6 rounded-xl text-sm font-bold transition-all shadow-lg shadow-blue-100 disabled:opacity-50"
            >
              <span v-if="!isLoading" class="flex items-center justify-center gap-2">
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 002-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10"></path></svg>
                TAMPILKAN DATA
              </span>
              <span v-else class="animate-pulse">Loading...</span>
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- 3. Hasil Generate (Full Width Below Grid) -->
    <div v-if="dataCetak" class="space-y-6 pt-10">
      <div class="flex items-center justify-between px-4 no-print">
          <h4 class="text-sm font-bold text-slate-800 flex items-center gap-2">
          <span class="w-2 h-6 bg-slate-800 rounded-full"></span>
          Hasil Generate ({{ dataCetak.siswa.length }} Kartu)
        </h4>
        <span class="text-xs text-slate-400 font-medium tracking-wide">READY FOR A4 PRINTING</span>
      </div>

      <div class="bg-white p-4 sm:p-12 rounded-[40px] border border-slate-100 shadow-sm print:shadow-none print:border-0 print:p-0">
        <div class="grid grid-cols-1 md:grid-cols-2 print:grid-cols-2 gap-6 print:gap-x-0 print:gap-y-2 w-full max-w-[210mm] mx-auto print-area">
          <div v-for="s in dataCetak.siswa" :key="s.id" class="exam-card relative border border-black p-2 bg-white print:w-[100mm] print:h-[70mm] flex flex-col overflow-hidden font-serif">
            <div class="grid grid-cols-[50px_1fr_50px] gap-2 items-center border-b-2 border-black pb-1 mb-1.5">
              <div class="w-[50px] h-[50px] flex items-center justify-center">
                <img v-if="config.logo_sekolah" :src="config.logo_sekolah" class="w-full h-full object-contain" />
              </div>
              <div class="text-center leading-tight">
                <h4 class="text-[10px] font-black uppercase tracking-tight text-black">{{ config.nama_ujian }}</h4>
                <h2 class="text-[13px] font-black uppercase text-black leading-none tracking-tighter">{{ config.nama_sekolah }}</h2>
                <p class="text-[8.5px] font-normal italic leading-none text-slate-800">{{ config.alamat_sekolah }}</p>
              </div>
              <div class="w-[50px] h-[50px] flex items-center justify-center">
                <img v-if="config.logo_kanan" :src="config.logo_kanan" class="w-full h-full object-contain" />
              </div>
            </div>
            <div class="px-4 space-y-0 flex-1">
              <div class="grid grid-cols-[90px_10px_1fr] items-baseline text-[10.5px] leading-tight">
                <span>Nomor Peserta</span><span class="text-center">:</span><span class="text-black uppercase">{{ s.nomor_peserta || '-' }}</span>
              </div>
              <div class="grid grid-cols-[90px_10px_1fr] items-baseline text-[10.5px] leading-tight">
                <span>Nama Siswa</span><span class="text-center">:</span><span class="font-bold text-black uppercase truncate">{{ s.nama_lengkap }}</span>
              </div>
              <div class="grid grid-cols-[90px_10px_1fr] items-baseline text-[10.5px] leading-tight">
                <span>NIS / NISN</span><span class="text-center">:</span><span class="text-black">{{ s.nisn }}</span>
              </div>
              <div class="grid grid-cols-[90px_10px_1fr] items-baseline text-[10.5px] leading-tight">
                <span>Kelas / Sesi</span><span class="text-center">:</span><span class="text-black">{{ s.kelas }} / {{ s.sesi }}</span>
              </div>
              <div class="grid grid-cols-[90px_10px_1fr] items-baseline text-[10.5px] leading-tight">
                <span>Ruang Ujian</span><span class="text-center">:</span><span class="text-black">{{ s.ruang }}</span>
              </div>
              <div class="grid grid-cols-[90px_10px_1fr] items-baseline text-[10.5px] leading-tight">
                <span>Username</span><span class="text-center">:</span><span class="font-bold text-black">{{ s.username }}</span>
              </div>
              <div class="grid grid-cols-[90px_10px_1fr] items-baseline text-[10.5px] leading-tight">
                <span>Password</span><span class="text-center">:</span><span class="font-bold text-black">{{ s.password || '********' }}</span>
              </div>
            </div>
            <div class="mt-2 flex justify-between items-end px-4 pb-1">
              <div class="w-[50px] h-[65px] border border-black flex items-center justify-center p-0.5 bg-white mb-1 uppercase text-[7px] text-slate-300 italic">Foto</div>
              <div class="text-center text-[8.5px] leading-tight min-w-[140px]">
                <p class="text-black mb-0.5">{{ config.kota || '............' }}, {{ config.tanggal }}</p>
                <p class="font-bold">Kepala Sekolah,</p>
                <div class="h-8 flex justify-center items-center relative py-0.5"><img v-if="config.ttd_kepala_sekolah" :src="config.ttd_kepala_sekolah" class="h-full object-contain" /></div>
                <p class="font-bold underline uppercase">{{ config.kepala_sekolah_nama }}</p>
                <p class="text-[7.5px]">NIP. {{ config.kepala_sekolah_nip || '-' }}</p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
@media print {
  @page {
    margin: 5mm;
    size: A4 portrait;
  }
  
  /* Sembunyikan secara brutal semua elemen UI */
  nav, header, aside, .navbar, .sidebar, .no-print, button, .breadcrumb, .lg-col-span-4, .topbar, #sidebar, .main-header {
    display: none !important;
    opacity: 0 !important;
    visibility: hidden !important;
  }

  body {
    margin: 0;
    padding: 0 !important;
    background: white !important;
    -webkit-print-color-adjust: exact;
  }

  .print-area {
    padding: 0 !important;
    margin: 0 !important;
    width: 202mm !important;
    position: absolute !important;
    top: 0 !important;
    left: 50% !important;
    transform: translateX(-50%) !important;
    z-index: 99999 !important;
    background: white !important;
    border: none !important;
    display: grid !important;
  }

  .exam-card {
    break-inside: avoid;
    page-break-inside: avoid;
    border: 1px solid black !important;
  }

  .print-area * {
    -webkit-print-color-adjust: exact !important;
    print-color-adjust: exact !important;
  }
}

.exam-card {
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}
</style>
