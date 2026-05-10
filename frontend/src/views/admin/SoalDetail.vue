<script setup>
import { ref, onMounted, nextTick, watch, computed } from 'vue'
import { useAuthStore } from '../../store/auth'
import { useAlertStore } from '../../store/alert'
import { useRoute, useRouter } from 'vue-router'
import BaseModal from '../../components/BaseModal.vue'
import ConfirmModal from '../../components/ConfirmModal.vue'
import MathHelper from '../../components/MathHelper.vue'

const authStore = useAuthStore()
const alertStore = useAlertStore()
const route = useRoute()
const router = useRouter()
const bankSoalId = route.params.bankSoalId

const showMathHelp = ref(false)

const bankSoalInfo = ref(null)
const soals = ref([])
const isLoading = ref(false)
const showForm = ref(false)
const isEditMode = ref(false)
const editId = ref(null)

// State untuk Modal Konfirmasi Hapus
const showDeleteConfirm = ref(false)
const deleteTargetId = ref(null)
const isDeleting = ref(false)

const form = ref({
  jenis_soal: 'PG',
  pertanyaan: '',
  opsi_a: '',
  opsi_b: '',
  opsi_c: '',
  opsi_d: '',
  kunci_jawaban: 'A',
  bobot_nilai: 1
})

const activeTab = ref('PG')

// State untuk Editor Quill (Manual CDN)
const editorPertanyaan = ref(null)
const editorOpsiA = ref(null)
const editorOpsiB = ref(null)
const editorOpsiC = ref(null)
const editorOpsiD = ref(null)
const editorKunciEssay = ref(null)
const quillInstances = { q: null, a: null, b: null, c: null, d: null, k: null }

const initEditors = () => {
  nextTick(() => {
    if (!window.Quill) return
    
    const config = {
      theme: 'snow',
      modules: {
        toolbar: [
          ['bold', 'italic', 'underline'], 
          [{ 'list': 'ordered'}, { 'list': 'bullet' }], 
          ['image', 'formula']
        ]
      },
      placeholder: 'Tulis di sini...'
    }

    // Inisialisasi Pertanyaan
    if (editorPertanyaan.value && !quillInstances.q) {
      quillInstances.q = new window.Quill(editorPertanyaan.value, config)
      quillInstances.q.root.innerHTML = form.value.pertanyaan || ''
      quillInstances.q.on('text-change', () => { form.value.pertanyaan = quillInstances.q.root.innerHTML })
    }

    // Inisialisasi Opsi (Hanya jika tipe soal PG)
    if (form.value.jenis_soal === 'PG') {
      const refs = { a: editorOpsiA, b: editorOpsiB, c: editorOpsiC, d: editorOpsiD }
      Object.keys(refs).forEach(key => {
        if (refs[key].value && !quillInstances[key]) {
          quillInstances[key] = new window.Quill(refs[key].value, config)
          quillInstances[key].root.innerHTML = form.value['opsi_' + key] || ''
          quillInstances[key].on('text-change', () => { form.value['opsi_' + key] = quillInstances[key].root.innerHTML })
        }
      })
    } else if (form.value.jenis_soal === 'ESSAY') {
      if (editorKunciEssay.value && !quillInstances.k) {
        quillInstances.k = new window.Quill(editorKunciEssay.value, config)
        quillInstances.k.root.innerHTML = form.value.kunci_jawaban || ''
        quillInstances.k.on('text-change', () => { form.value.kunci_jawaban = quillInstances.k.root.innerHTML })
      }
    }

    // Update content if instances exist (safety sync)
    if (quillInstances.q) quillInstances.q.root.innerHTML = form.value.pertanyaan || ''
    if (form.value.jenis_soal === 'PG') {
      ['a','b','c','d'].forEach(k => {
        if (quillInstances[k]) quillInstances[k].root.innerHTML = form.value['opsi_' + k] || ''
      })
    } else if (quillInstances.k) {
      quillInstances.k.root.innerHTML = form.value.kunci_jawaban || ''
    }
  })
}

const destroyEditors = () => {
  Object.keys(quillInstances).forEach(k => { quillInstances[k] = null })
}

const filteredSoals = computed(() => {
  return soals.value.filter(s => String(s.jenis_soal).trim().toUpperCase() === activeTab.value)
})

const renderMath = () => {
  nextTick(() => {
    const mathElements = document.querySelectorAll('.math-content')
    if (window.renderMathInElement && mathElements.length > 0) {
      mathElements.forEach(el => {
        try {
          window.renderMathInElement(el, {
            delimiters: [
              { left: '$$', right: '$$', display: true },
              { left: '$', right: '$', display: false }
            ],
            throwOnError: false
          })
        } catch (e) {
          console.warn("KaTeX render error:", e)
        }
      })
    }
  })
}

const apiPrefix = computed(() => authStore.user?.role === 'admin' ? '/api/admin' : '/api/guru')

const fetchData = async () => {
  isLoading.value = true
  const headers = { 'Authorization': `Bearer ${authStore.token}` }
  try {
    const [resBS, resS] = await Promise.all([
      fetch(`${apiPrefix.value}/bank-soal/${bankSoalId}`, { headers }),
      fetch(`${apiPrefix.value}/bank-soal/${bankSoalId}/soal`, { headers })
    ])
    const dbs = await resBS.json()
    const ds = await resS.json()
    bankSoalInfo.value = dbs.data || null
    soals.value = ds.data || []
    renderMath()
  } catch (error) {
    alertStore.showAlert("Gagal memuat data", "error")
  } finally {
    isLoading.value = false
  }
}

const triggerForm = () => {
  destroyEditors()
  isEditMode.value = false
  editId.value = null
  
  // Ambil bobot default dari info bank soal
  const defaultBobot = activeTab.value === 'PG' 
    ? (bankSoalInfo.value?.default_bobot_pg || 1) 
    : (bankSoalInfo.value?.default_bobot_essay || 1)

  form.value = { 
    jenis_soal: activeTab.value, 
    pertanyaan: '', 
    opsi_a: '', 
    opsi_b: '', 
    opsi_c: '', 
    opsi_d: '', 
    kunci_jawaban: 'A', 
    bobot_nilai: defaultBobot 
  }
  showForm.value = true
}

const editSoal = (soal) => {
  destroyEditors()
  isEditMode.value = true
  editId.value = soal.id
  form.value = { ...soal }
  showForm.value = true
}

const saveSoal = async () => {
  isLoading.value = true
  try {
    const url = isEditMode.value ? `${apiPrefix.value}/soal/${editId.value}` : `${apiPrefix.value}/bank-soal/${bankSoalId}/soal`
    const method = isEditMode.value ? 'PUT' : 'POST'
    const res = await fetch(url, {
      method: method,
      headers: { 'Authorization': `Bearer ${authStore.token}`, 'Content-Type': 'application/json' },
      body: JSON.stringify(form.value)
    })
    if (res.ok) {
      alertStore.showAlert("Data berhasil disimpan", "success")
      showForm.value = false
      fetchData()
    } else {
      const data = await res.json()
      alertStore.showAlert(data.error || "Gagal menyimpan", "error")
    }
  } catch (error) {
    alertStore.showAlert("Kesalahan sistem", "error")
  } finally {
    isLoading.value = false
  }
}

const previewSoals = ref([])
const showPreview = ref(false)

const handleImportWord = async (event) => {
  const file = event.target.files[0]
  if (!file) return
  isLoading.value = true
  const reader = new FileReader()
  reader.onload = async (e) => {
    try {
      const arrayBuffer = e.target.result
      const result = await window.mammoth.convertToHtml({ arrayBuffer })
      const html = result.value
      const tempDiv = document.createElement('div')
      tempDiv.innerHTML = html
      // Hanya ambil baris (tr) yang merupakan bagian dari tabel utama (bukan nested table)
      const rows = Array.from(tempDiv.querySelectorAll('tr')).filter(tr => {
        const parentTable = tr.closest('table')
        return parentTable && !parentTable.parentElement.closest('table')
      })
      if (rows.length < 2) throw new Error("Tabel tidak ditemukan")
      let parsed = []
      let current = null
      rows.forEach(row => {
        const cells = Array.from(row.cells)
        if (cells.length < 3) return
        const no = cells[0].innerText.trim()
        if (no && /^\d+$/.test(no)) {
          if (current) parsed.push(current)
          const isEssay = cells[2].innerText.trim() === '2' || cells[2].innerText.trim() === '5'
          current = { 
            pertanyaan: cells[1].innerHTML.trim(), 
            jenis_soal: isEssay ? 'ESSAY' : 'PG',
            opsi_a: '', opsi_b: '', opsi_c: '', opsi_d: '', kunci_jawaban: '', bobot_nilai: null
          }
          if (isEssay && cells[cells.length-1]) current.kunci_jawaban = cells[cells.length-1].innerHTML.trim()
        }
        if (current && current.jenis_soal === 'PG') {
          cells.forEach((cell, idx) => {
            const txt = cell.innerText.trim().toUpperCase()
            if (['A','B','C','D'].includes(txt) && txt.length === 1) {
              if (cells[idx+1]) current['opsi_'+txt.toLowerCase()] = cells[idx+1].innerHTML.trim()
              if ((cells[idx+2] || cells[5])?.innerText.trim().toLowerCase() === 'v') current.kunci_jawaban = txt
            }
          })
        }
      })
      if (current) parsed.push(current)
      previewSoals.value = parsed
      showPreview.value = true
    } catch (err) {
      alertStore.showAlert("Gagal impor: " + err.message, "error")
    } finally {
      isLoading.value = false
      event.target.value = ''
    }
  }
  reader.readAsArrayBuffer(file)
}

const confirmDeleteSoal = (id) => {
  deleteTargetId.value = id
  showDeleteConfirm.value = true
}

const executeDeleteSoal = async () => {
  isDeleting.value = true
  try {
    const res = await fetch(`${apiPrefix.value}/soal/${deleteTargetId.value}`, {
      method: 'DELETE',
      headers: { 'Authorization': `Bearer ${authStore.token}` }
    })
    if (res.ok) {
      showDeleteConfirm.value = false
      alertStore.showAlert("Soal dihapus", "success")
      fetchData()
    } else {
      const data = await res.json()
      alertStore.showAlert(data.error || "Gagal menghapus soal", "error")
      showDeleteConfirm.value = false
    }
  } catch (error) {
    alertStore.showAlert("Kesalahan sistem saat menghapus", "error")
  } finally {
    isDeleting.value = false
  }
}

const saveImportedSoals = async () => {
  // Filter hanya soal yang memiliki pertanyaan (menghindari baris kosong dari Word)
  const validSoals = previewSoals.value.filter(s => s.pertanyaan && s.pertanyaan.trim() !== '')
  
  if (validSoals.length === 0) {
    alertStore.showAlert("Tidak ada soal valid yang ditemukan dalam dokumen", "error")
    return
  }

  isLoading.value = true
  try {
    // 1. Hapus semua soal lama
    const clearRes = await fetch(`${apiPrefix.value}/bank-soal/${bankSoalId}/soal/clear`, {
      method: 'DELETE',
      headers: { 'Authorization': `Bearer ${authStore.token}` }
    })

    if (!clearRes.ok) {
      const data = await clearRes.json()
      throw new Error(data.error || "Gagal membersihkan soal lama")
    }

    // 2. Simpan soal baru
    let count = 0
    const defPG = bankSoalInfo.value?.default_bobot_pg || 1
    const defEssay = bankSoalInfo.value?.default_bobot_essay || 1

    for (const s of validSoals) {
      // Tentukan bobot default berdasarkan jenis soal yang di-parse
      const isEssay = String(s.jenis_soal).toUpperCase() === 'ESSAY'
      const defaultBobot = isEssay ? defEssay : defPG

      const payload = {
        ...s,
        bobot_nilai: s.bobot_nilai || defaultBobot
      }

      const res = await fetch(`${apiPrefix.value}/bank-soal/${bankSoalId}/soal`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json', 'Authorization': `Bearer ${authStore.token}` },
        body: JSON.stringify(payload)
      })
      if (res.ok) count++
    }
    alertStore.showAlert(`Berhasil mengganti soal! ${count} soal baru tersimpan.`, "success")
    showPreview.value = false
    fetchData()
  } catch (err) {
    alertStore.showAlert("Gagal simpan impor: " + err.message, "error")
  } finally {
    isLoading.value = false
  }
}

onMounted(() => {
  fetchData()
})

watch(activeTab, () => {
  renderMath()
})

watch(showForm, (val) => {
  if (val) {
    initEditors()
    renderMath()
  }
})
</script>

<template>
  <div>
    <!-- Loading Overlay Global -->
    <div v-if="isLoading" class="fixed inset-0 z-[9999] flex items-center justify-center bg-slate-900/40 backdrop-blur-sm">
      <div class="bg-white p-8 rounded-3xl shadow-2xl flex flex-col items-center gap-4">
        <div class="relative">
          <div class="w-12 h-12 border-4 border-slate-100 rounded-full"></div>
          <div class="absolute top-0 w-12 h-12 border-4 border-blue-600 border-t-transparent rounded-full animate-spin"></div>
        </div>
        <div class="text-center">
          <p class="font-bold text-slate-800">Sedang Memproses...</p>
          <p class="text-xs text-slate-500 mt-1">Sistem sedang membaca dokumen Anda</p>
        </div>
      </div>
    </div>
    <!-- Konten Utama -->
    <div class="space-y-6 pb-4">
      <!-- Header -->
      <div class="bg-white p-6 rounded-2xl shadow-sm border border-slate-100 flex flex-col md:flex-row justify-between items-center gap-4">
        <div class="flex items-center gap-4">
          <button @click="router.back()" class="p-2 hover:bg-slate-100 rounded-xl text-slate-500 transition-colors">
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 19l-7-7m0 0l7-7m-7 7h18"></path></svg>
          </button>
          <div>
            <h3 class="text-xl font-bold text-slate-800">{{ bankSoalInfo?.judul_bank_soal || 'Memuat...' }}</h3>
            <p class="text-sm text-slate-500">{{ bankSoalInfo?.mapel?.nama_mapel }} • Kelas {{ bankSoalInfo?.tingkat_kelas }} • {{ soals.length }} Soal</p>
          </div>
        </div>
        <div class="flex items-center gap-3">
          <input type="file" id="importWord" class="hidden" accept=".docx" @change="handleImportWord">
          <label for="importWord" class="cursor-pointer bg-blue-600 hover:bg-blue-700 text-white px-6 py-2.5 rounded-xl text-sm font-bold border border-blue-600 flex items-center gap-2 shadow-lg shadow-blue-100">
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"></path></svg>
            Import dari Word (.docx)
          </label>
          <button @click="triggerForm" class="bg-blue-600 hover:bg-blue-700 text-white px-6 py-2.5 rounded-xl text-sm font-black uppercase tracking-wider shadow-lg shadow-blue-100 transition-all active:scale-95">{{ showForm ? 'Batal' : '+ Tambah Manual' }}</button>
        </div>
      </div>

      <!-- Tabs -->
      <div class="flex border-b border-slate-200 gap-8">
        <button 
          @click="activeTab = 'PG'"
          :class="activeTab === 'PG' ? 'border-blue-600 text-blue-600' : 'border-transparent text-slate-400 hover:text-slate-600'"
          class="pb-4 border-b-2 font-bold text-sm transition-all flex items-center gap-2"
        >
          Pilihan Ganda
          <span class="bg-slate-100 text-slate-500 px-2 py-0.5 rounded-full text-[10px]">{{ soals.filter(s => String(s.jenis_soal).trim().toUpperCase() === 'PG').length }}</span>
        </button>
        <button 
          @click="activeTab = 'ESSAY'"
          :class="activeTab === 'ESSAY' ? 'border-blue-600 text-blue-600' : 'border-transparent text-slate-400 hover:text-slate-600'"
          class="pb-4 border-b-2 font-bold text-sm transition-all flex items-center gap-2"
        >
          Essay / Uraian
          <span class="bg-slate-100 text-slate-500 px-2 py-0.5 rounded-full text-[10px]">{{ soals.filter(s => String(s.jenis_soal).trim().toUpperCase() === 'ESSAY').length }}</span>
        </button>
      </div>

      <!-- List -->
      <div class="space-y-4">
        <div v-if="filteredSoals.length === 0" class="bg-white p-12 rounded-2xl border border-dashed border-slate-200 text-center text-slate-400">
          Belum ada soal di kategori ini.
        </div>
        <div v-for="(soal, index) in filteredSoals" :key="soal.id" class="bg-white rounded-2xl shadow-sm border border-slate-100 overflow-hidden hover:border-blue-200 transition-all">
          <div class="p-6">
            <div class="flex gap-4">
              <div class="w-9 h-9 bg-blue-600 text-white rounded-xl flex items-center justify-center font-black text-sm shadow-lg shadow-blue-100">{{ index + 1 }}</div>
              <div class="flex-1">
                <div class="text-slate-800 font-medium leading-relaxed math-content ql-editor !p-0 !min-h-0" v-html="soal.pertanyaan"></div>
                <div v-if="String(soal.jenis_soal).trim().toUpperCase() === 'PG'" class="grid grid-cols-1 md:grid-cols-2 gap-4 mt-1">
                  <div v-for="o in ['a','b','c','d']" :key="o" 
                    :class="String(soal.kunci_jawaban).trim().toUpperCase() === o.toUpperCase() ? 'text-emerald-600 font-bold bg-emerald-50 border-emerald-100' : 'text-slate-500 border-slate-50'" 
                    class="text-sm px-4 py-2 border rounded-xl flex items-center gap-2 math-content">
                    <span class="w-5 h-5 flex-shrink-0 flex items-center justify-center border border-current rounded text-[10px] uppercase font-black">{{ o }}</span>
                    <div class="ql-editor !p-0" v-html="soal['opsi_'+o]"></div>
                  </div>
                </div>
                <div v-else class="bg-slate-50 p-4 rounded-xl border border-slate-100">
                  <p class="text-[10px] font-black text-slate-400 uppercase tracking-widest mb-1">Referensi Jawaban:</p>
                  <div class="text-sm text-slate-600 italic math-content ql-editor !p-0" v-html="soal.kunci_jawaban"></div>
                </div>
              </div>
              <div class="flex flex-col gap-2">
                <button @click="editSoal(soal)" class="p-2 text-slate-400 hover:text-blue-600 hover:bg-blue-50 rounded-xl transition-all">
                  <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"></path></svg>
                </button>
                <button @click="confirmDeleteSoal(soal.id)" class="p-2 text-slate-400 hover:text-rose-600 hover:bg-rose-50 rounded-xl transition-all">
                  <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path></svg>
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Modals -->
    <BaseModal 
      :show="showForm" 
      :title="isEditMode ? 'Edit Soal' : 'Tambah Soal'"
      :isLoading="isLoading"
      size="max-w-4xl"
      @close="showForm = false"
      @confirm="saveSoal"
    >
      <div class="grid grid-cols-1 md:grid-cols-2 gap-6 relative">
        <div v-if="showMathHelp" class="md:col-span-2 animate-in slide-in-from-top duration-300">
           <MathHelper />
        </div>
        <div class="md:col-span-2">
          <div class="flex justify-between items-center mb-2">
            <label class="text-sm font-semibold text-slate-700">Pertanyaan</label>
            <button @click="showMathHelp = !showMathHelp" class="text-[10px] font-black uppercase tracking-widest px-3 py-1 rounded-lg transition-all" :class="showMathHelp ? 'bg-rose-50 text-rose-600 hover:bg-rose-100' : 'bg-blue-50 text-blue-600 hover:bg-blue-100'">
              {{ showMathHelp ? 'Tutup Bantuan' : 'Bantuan Rumus (Σ)' }}
            </button>
          </div>
          <div class="bg-white border border-slate-200 rounded-2xl overflow-hidden shadow-sm">
            <div ref="editorPertanyaan" class="min-h-[200px] ql-editor-custom"></div>
          </div>
        </div>
        <template v-if="form.jenis_soal === 'PG'">
          <div v-for="o in ['a','b','c','d']" :key="o" class="space-y-2 p-4 bg-slate-50 border border-slate-100 rounded-2xl">
            <label class="block text-xs font-black uppercase tracking-widest text-slate-400">Opsi {{ o }}</label>
            <div class="bg-white border border-slate-200 rounded-xl overflow-hidden">
              <div :ref="(el) => { if (o === 'a') editorOpsiA = el; if (o === 'b') editorOpsiB = el; if (o === 'c') editorOpsiC = el; if (o === 'd') editorOpsiD = el; }" class="min-h-[100px] ql-editor-custom"></div>
            </div>
          </div>
          <div class="md:col-span-2 p-4 bg-blue-50 border border-blue-100 rounded-2xl">
            <label class="block text-sm font-bold text-blue-800 mb-2">Kunci Jawaban</label>
            <select v-model="form.kunci_jawaban" class="block w-full px-4 py-3 bg-white border border-blue-200 rounded-xl focus:ring-2 focus:ring-blue-500 font-bold text-blue-900 transition-all">
              <option value="">-- Pilih Kunci --</option>
              <option value="A">Opsi A</option>
              <option value="B">Opsi B</option>
              <option value="C">Opsi C</option>
              <option value="D">Opsi D</option>
            </select>
          </div>
        </template>
        <div v-else class="md:col-span-2 space-y-4">
          <div class="p-4 bg-blue-50 border border-blue-100 rounded-2xl">
            <label class="block text-sm font-bold text-blue-800 mb-2">Kunci/Referensi Jawaban Essay</label>
            <div class="bg-white border border-blue-200 rounded-xl overflow-hidden">
              <div ref="editorKunciEssay" class="min-h-[150px] ql-editor-custom"></div>
            </div>
          </div>
        </div>
        <div>
          <label class="block text-sm font-semibold text-slate-700 mb-2">Bobot Nilai</label>
          <input v-model.number="form.bobot_nilai" type="number" step="0.1" min="0" class="block w-full px-4 py-3 bg-slate-50 border border-slate-200 rounded-xl focus:ring-2 focus:ring-blue-500 transition-all font-bold">
        </div>
      </div>
    </BaseModal>

    <BaseModal 
      :show="showPreview" 
      title="Pratinjau Import Word"
      confirmText="Simpan Semua Soal"
      :isLoading="isLoading"
      size="max-w-5xl"
      @close="showPreview = false"
      @confirm="saveImportedSoals"
    >
      <div class="space-y-6">
        <div class="p-4 bg-amber-50 border border-amber-200 rounded-2xl flex items-start gap-4 shadow-sm">
          <div class="w-10 h-10 bg-amber-500 text-white rounded-xl flex items-center justify-center font-black flex-shrink-0 shadow-lg shadow-amber-200">
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"></path></svg>
          </div>
          <div>
            <p class="text-xs text-amber-800 font-bold uppercase tracking-wider mb-1">Peringatan Penting!</p>
            <p class="text-sm text-amber-700 font-medium leading-relaxed">Melanjutkan proses ini akan <span class="font-black underline decoration-amber-500">MENGHAPUS SELURUH SOAL LAMA</span> yang ada di wadah ini dan menggantinya dengan soal baru dari dokumen.</p>
          </div>
        </div>
        <div class="p-4 bg-blue-50 border border-blue-100 rounded-2xl flex items-center gap-4">
          <div class="w-8 h-8 bg-blue-600 text-white rounded-lg flex items-center justify-center font-black">i</div>
          <p class="text-sm text-blue-800 font-medium italic">Ditemukan {{ previewSoals.length }} soal dalam dokumen.</p>
        </div>
        <div v-for="(ps, i) in previewSoals" :key="i" class="bg-white p-6 rounded-2xl border border-slate-200 shadow-sm space-y-4 hover:border-blue-300 transition-all group">
          <div class="flex justify-between items-center pb-4 border-b border-slate-50">
            <div class="flex items-center gap-3">
              <span class="w-9 h-9 bg-blue-600 text-white text-xs font-black rounded-xl flex items-center justify-center italic shadow-lg shadow-blue-100">#{{ i + 1 }}</span>
              <span :class="ps.jenis_soal === 'PG' ? 'bg-blue-100 text-blue-600' : 'bg-amber-100 text-amber-600'" class="px-3 py-1 text-[10px] font-black rounded-lg uppercase tracking-tighter">{{ ps.jenis_soal }}</span>
            </div>
            <button @click="previewSoals.splice(i, 1)" class="text-rose-400 hover:text-rose-600 opacity-0 group-hover:opacity-100 transition-opacity">
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path></svg>
            </button>
          </div>
          <div class="prose prose-sm max-w-none math-content" v-html="ps.pertanyaan"></div>
          <div v-if="ps.jenis_soal === 'PG'" class="grid grid-cols-1 md:grid-cols-2 gap-3 mt-4">
            <div v-for="o in ['a','b','c','d']" :key="o" :class="ps.kunci_jawaban === o.toUpperCase() ? 'bg-emerald-50 border-emerald-200 ring-1 ring-emerald-500' : 'bg-slate-50 border-slate-100'" class="p-3 border rounded-xl flex items-start gap-3">
              <span class="w-5 h-5 flex-shrink-0 flex items-center justify-center border border-slate-300 rounded text-[10px] uppercase font-black bg-white shadow-sm">{{ o }}</span>
              <div class="text-xs leading-tight" v-html="ps['opsi_'+o]"></div>
            </div>
          </div>
          <div v-else class="mt-4 p-4 bg-emerald-50 border border-emerald-100 rounded-xl">
            <p class="text-[10px] font-black text-emerald-600 uppercase tracking-widest mb-1">Kunci Jawaban:</p>
            <div class="text-sm italic" v-html="ps.kunci_jawaban"></div>
          </div>
        </div>
      </div>
    </BaseModal>

    <ConfirmModal 
      :show="showDeleteConfirm"
      title="Hapus Soal"
      message="Apakah Anda yakin ingin menghapus soal ini? Tindakan ini tidak dapat dibatalkan."
      :isLoading="isDeleting"
      @close="showDeleteConfirm = false"
      @confirm="executeDeleteSoal"
    />
  </div>
</template>

<style scoped>
.ql-editor-custom :deep(.ql-editor) {
  min-height: 100px !important;
  padding: 12px 16px !important;
  font-size: 14px !important;
  line-height: 1.6 !important;
}

:deep(.ql-container.ql-snow) {
  border: none !important;
}

:deep(.ql-toolbar.ql-snow) {
  border: none !important;
  border-bottom: 1px solid #f1f5f9 !important;
  background: #f8fafc !important;
}

:deep(.ql-editor) {
  min-height: 0 !important;
  height: auto !important;
  padding: 0 !important;
  overflow-y: hidden !important;
}

:deep(.ql-editor p) {
  margin-bottom: 0 !important;
}

:deep(.ql-editor img) {
  max-width: 100%;
  height: auto;
  border-radius: 12px;
  margin: 8px 0;
  box-shadow: 0 4px 6px -1px rgb(0 0 0 / 0.1);
  display: block;
}

:deep(.math-content img) {
  max-width: 100%;
  height: auto;
  border-radius: 8px;
}
</style>
