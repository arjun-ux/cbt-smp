<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAlertStore } from '../../store/alert'

const router = useRouter()
const alertStore = useAlertStore()
const examData = ref(null)

onMounted(() => {
  const session = sessionStorage.getItem('exam_session')
  if (!session) {
    alertStore.showAlert("Sesi tidak ditemukan, silakan login kembali", "error")
    router.push({ name: 'StudentLogin' })
    return
  }
  examData.value = JSON.parse(session)
})

const startExam = () => {
  router.push({ 
    name: 'StudentPengerjaan', 
    params: { jadwalId: examData.value.jadwal.id } 
  })
}
</script>

<template>
  <div v-if="examData" class="max-w-3xl mx-auto space-y-8">
    <!-- Card Konfirmasi -->
    <div class="bg-white p-8 sm:p-10 rounded-[2.5rem] shadow-xl border border-slate-100">
      <div class="text-center mb-10">
        <h2 class="text-3xl font-black text-slate-800 tracking-tight">Konfirmasi Data Peserta</h2>
        <p class="text-slate-500 mt-2 font-semibold">Pastikan data di bawah ini sudah sesuai sebelum memulai ujian</p>
      </div>

      <div class="grid grid-cols-1 md:grid-cols-2 gap-8 bg-slate-50 p-8 rounded-3xl border border-slate-100">
        <div class="space-y-6">
          <div>
            <label class="text-[10px] font-black text-slate-400 uppercase tracking-widest block mb-1">Nama Lengkap</label>
            <p class="text-lg font-black text-slate-800">{{ examData.siswa.nama_lengkap }}</p>
          </div>
          <div>
            <label class="text-[10px] font-black text-slate-400 uppercase tracking-widest block mb-1">NISN</label>
            <p class="text-lg font-black text-slate-800">{{ examData.siswa.nisn }}</p>
          </div>
          <div>
            <label class="text-[10px] font-black text-slate-400 uppercase tracking-widest block mb-1">Kelas / Ruang</label>
            <p class="text-lg font-black text-slate-800">{{ examData.siswa.kelas?.nama_kelas }} / Ruang {{ examData.siswa.ruang_id }}</p>
          </div>
        </div>

        <div class="space-y-6">
          <div>
            <label class="text-[10px] font-black text-slate-400 uppercase tracking-widest block mb-1">Mata Pelajaran</label>
            <p class="text-lg font-black text-blue-600">{{ examData.jadwal.bank_soal?.mapel?.nama_mapel }}</p>
          </div>
          <div>
            <label class="text-[10px] font-black text-slate-400 uppercase tracking-widest block mb-1">Alokasi Waktu</label>
            <p class="text-lg font-black text-slate-800">{{ examData.jadwal.durasi_menit }} Menit</p>
          </div>
          <div>
            <label class="text-[10px] font-black text-slate-400 uppercase tracking-widest block mb-1">Status Sesi</label>
            <span class="px-3 py-1 bg-green-100 text-green-600 rounded-lg text-xs font-black uppercase tracking-wider border border-green-200">Siap Ujian</span>
          </div>
        </div>
      </div>

      <!-- Aturan Singkat -->
      <div class="mt-10 p-6 bg-amber-50 rounded-2xl border border-amber-100 flex gap-4">
        <div class="w-10 h-10 bg-amber-100 rounded-full flex-shrink-0 flex items-center justify-center text-amber-600">
          <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"></path></svg>
        </div>
        <div>
          <p class="text-sm font-bold text-amber-900 uppercase tracking-widest mb-1">Perhatian!</p>
          <p class="text-sm text-amber-800 leading-relaxed">Jangan menutup browser atau berpindah tab saat ujian berlangsung. Jawaban akan disimpan secara otomatis.</p>
        </div>
      </div>

      <button 
        @click="startExam"
        class="w-full mt-10 bg-blue-600 hover:bg-blue-700 text-white font-black py-5 rounded-3xl shadow-xl shadow-blue-200 transition-all active:scale-[0.98] text-lg tracking-widest"
      >
        MULAI KERJAKAN SOAL
      </button>
    </div>
  </div>
</template>
