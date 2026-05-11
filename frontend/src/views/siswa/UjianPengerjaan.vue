<script setup>
import { ref, onMounted, watch, nextTick } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../../store/auth'
import { useAlertStore } from '../../store/alert'
import { useExamStore } from '../../store/exam'
import { useExamTimer } from '../../composables/useExamTimer'
import { useAntiCheat } from '../../composables/useAntiCheat'

// Components
import ExamHeader from './components/ExamHeader.vue'
import ExamNavigator from './components/ExamNavigator.vue'
import QuestionBox from './components/QuestionBox.vue'
import ExamFooter from './components/ExamFooter.vue'
import ConfirmModal from '../../components/ConfirmModal.vue'

const router = useRouter()
const authStore = useAuthStore()
const alertStore = useAlertStore()
const examStore = useExamStore()
const examContainer = ref(null)

// 1. Logic & Security
useAntiCheat()
useExamTimer((isAuto) => {
  if (isAuto) finishExam(true)
})

const renderMath = () => {
  nextTick(() => {
    if (window.renderMathInElement && examContainer.value) {
      window.renderMathInElement(examContainer.value, {
        delimiters: [
          { left: '$$', right: '$$', display: true },
          { left: '$', right: '$', display: false },
          { left: '\\(', right: '\\)', display: false },
          { left: '\\[', right: '\\]', display: true }
        ],
        throwOnError: false
      })
    }
  })
}

const initializeExam = async () => {
  let session = null
  try {
    session = JSON.parse(sessionStorage.getItem('exam_session'))
  } catch (e) {
    console.error("Session corrupted")
  }

  if (!session) {
    router.push({ name: 'StudentLogin' })
    return
  }
  
  examStore.setExamInfo(session)
  await examStore.fetchQuestions()
  renderMath()
}

const finishExam = async (isAuto = false) => {
  // Tutup navigasi mobile jika sedang terbuka agar tidak menutupi modal konfirmasi
  examStore.isNavModalOpen = false

  if (!isAuto && !examStore.showFinishConfirm) {
    examStore.showFinishConfirm = true
    return
  }

  examStore.isSubmitting = true
  try {
    const res = await fetch(`/api/siswa/submit/${examStore.examInfo.peserta_id}`, {
      method: 'POST',
      headers: { 'Authorization': `Bearer ${authStore.token}` }
    })

    if (res.ok) {
      alertStore.showAlert("Ujian Berhasil Diselesaikan", "success")
      sessionStorage.removeItem('exam_session')
      localStorage.removeItem(`answers_${examStore.examInfo.peserta_id}`)
      if (document.fullscreenElement) document.exitFullscreen()
      router.push({ name: 'StudentDashboard' })
    }
  } catch (e) {
    alertStore.showAlert("Kesalahan koneksi saat mengakhiri ujian", "error")
  } finally {
    examStore.isSubmitting = false
    examStore.showFinishConfirm = false
  }
}

const isDev = import.meta.env.DEV

onMounted(() => {
  initializeExam()
})

watch(() => examStore.currentIdx, () => renderMath())
</script>

<template>
  <div ref="examContainer" :class="{ 'no-select': !isDev }" class="min-h-screen bg-slate-50 font-sans text-slate-900 relative">
    
    <!-- BLOCKED OVERLAY -->
    <div v-if="examStore.isTerblokir" class="fixed inset-0 z-[99999] bg-slate-900/95 backdrop-blur-2xl flex items-center justify-center p-6 text-center">
      <div class="max-w-md w-full bg-white rounded-[3rem] p-12 shadow-2xl border border-slate-100">
        <div class="w-24 h-24 bg-red-50 text-red-600 rounded-full flex items-center justify-center mx-auto mb-8 animate-bounce">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-12 w-12" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m0 0v2m0-2h2m-2 0H10m11-3V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2h14a2 2 0 002-2v-4z" />
          </svg>
        </div>
        <h2 class="text-4xl font-black text-slate-900 mb-4 tracking-tight">Ujian Terkunci!</h2>
        <p class="text-slate-500 mb-10 leading-relaxed font-bold text-lg">Anda terdeteksi keluar dari layar ujian.</p>
        <button @click="router.push({ name: 'StudentDashboard' })" class="w-full py-5 bg-slate-900 text-white font-black rounded-2xl transition-all hover:bg-slate-800 shadow-xl shadow-slate-200">
          KEMBALI KE DASHBOARD
        </button>
      </div>
    </div>

    <ExamHeader 
      @open-nav="examStore.isNavModalOpen = true" 
      @finish="finishExam(false)" 
    />

    <div v-if="examStore.isLoading" class="flex items-center justify-center h-[60vh]">
      <div class="animate-spin w-10 h-10 border-4 border-blue-600 border-t-transparent rounded-full"></div>
    </div>
    
    <div v-else class="max-w-[1400px] mx-auto p-3 lg:p-8 pb-32 lg:pb-40">
      <div class="grid grid-cols-1 lg:grid-cols-12 gap-8">
        <div class="lg:col-span-8 space-y-6">
          

          <QuestionBox :key="examStore.currentIdx" :soal="examStore.currentQuestion" />
        </div>

        <div class="hidden lg:block lg:col-span-4 space-y-4">
          <!-- INFORMASI PESERTA (COMPACT) -->
          <div class="bg-white p-5 rounded-2xl border border-slate-100 shadow-sm flex items-center gap-4 group">
            <div class="w-14 h-14 bg-gradient-to-br from-indigo-500 to-blue-600 text-white rounded-xl flex items-center justify-center text-xl font-black shadow-lg shadow-indigo-50 group-hover:rotate-2 transition-transform">
              {{ authStore.user?.nama?.charAt(0) || '?' }}
            </div>
            <div>
              <p class="text-[8px] font-black text-slate-400 uppercase tracking-widest mb-0.5">Peserta</p>
              <h3 class="font-bold text-slate-700 text-base leading-tight">{{ authStore.user?.nama }}</h3>
              <p class="text-[10px] font-bold text-slate-400 mt-1 uppercase">{{ authStore.user?.username }} • {{ authStore.user?.kelas_nama || 'Kelas' }}</p>
            </div>
          </div>

          <ExamNavigator @finish="finishExam(false)" />
        </div>
      </div>
    </div>

    <ExamFooter />

    <!-- Mobile Navigation Modal -->
    <div v-if="examStore.isNavModalOpen" class="fixed inset-0 z-[100] bg-slate-900/60 backdrop-blur-sm lg:hidden flex items-end">
      <div class="w-full bg-white rounded-t-[3rem] p-8 shadow-2xl animate-slide-up">
        <div class="flex items-center justify-between mb-8">
          <h4 class="text-xl font-black text-slate-800 tracking-tight">Daftar Soal</h4>
          <button @click="examStore.isNavModalOpen = false" class="p-2 bg-slate-100 rounded-xl text-slate-400">
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path></svg>
          </button>
        </div>
        <ExamNavigator :is-mobile="true" @close="examStore.isNavModalOpen = false" @finish="finishExam(false)" />
      </div>
    </div>

    <ConfirmModal 
      :show="examStore.showFinishConfirm"
      title="Akhiri Ujian?"
      message="Apakah Anda yakin ingin mengakhiri ujian ini? Pastikan semua jawaban sudah benar. Anda tidak dapat kembali setelah menekan tombol Selesai."
      confirmText="Ya, Selesai"
      :isLoading="examStore.isSubmitting"
      variant="primary"
      @close="examStore.showFinishConfirm = false"
      @confirm="finishExam(true)"
    />
  </div>
</template>

<style scoped>
@keyframes slide-up {
  from { transform: translateY(100%); }
  to { transform: translateY(0); }
}
.animate-slide-up {
  animation: slide-up 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}
</style>
