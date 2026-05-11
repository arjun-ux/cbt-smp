<script setup>
import { useExamStore } from '../../../store/exam'

const examStore = useExamStore()

const handlePrev = () => {
  examStore.prevQuestion()
}

const handleNext = () => {
  examStore.nextQuestion()
}

const handleRagu = () => {
  if (examStore.currentQuestion) {
    examStore.toggleRagu(examStore.currentQuestion.id)
  }
}
</script>

<template>
  <div class="fixed bottom-0 inset-x-0 z-[70] bg-white/80 backdrop-blur-xl border-t border-slate-200/60 p-4 lg:p-6 shadow-lg">
    <div class="max-w-[1400px] mx-auto flex items-center justify-between gap-2 lg:gap-5">
      <!-- Prev Button -->
      <button 
        @click="handlePrev" 
        :disabled="examStore.currentIdx === 0" 
        class="flex-1 lg:flex-none px-4 lg:px-12 py-3.5 lg:py-4 bg-white border-2 border-slate-100 rounded-2xl text-[10px] lg:text-sm font-black text-slate-500 hover:bg-slate-50 disabled:opacity-20 transition-all uppercase tracking-widest active:scale-90 flex items-center justify-center gap-2 shadow-sm"
      >
        <svg class="w-5 h-5 lg:w-5 lg:h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="3" d="M15 19l-7-7 7-7"></path></svg>
        <span class="hidden md:inline">Sebelumnya</span>
      </button>

      <!-- Ragu Button -->
      <button 
        @click="handleRagu"
        :class="[
          examStore.raguStatus[examStore.currentQuestion?.id] 
          ? 'bg-amber-500 border-amber-500 text-white shadow-md' 
          : 'bg-white border-slate-100 text-amber-600 hover:bg-amber-50/50'
        ]"
        class="flex-[1.5] lg:flex-none px-6 lg:px-14 py-3 lg:py-4 border-2 rounded-xl text-[10px] lg:text-sm font-black transition-all duration-200 uppercase tracking-widest active:scale-95 flex items-center justify-center gap-2 lg:gap-3"
      >
        <div :class="[examStore.raguStatus[examStore.currentQuestion?.id] ? 'bg-white' : 'bg-amber-100']" class="w-4 h-4 lg:w-5 lg:h-5 rounded flex items-center justify-center transition-colors">
          <svg v-if="examStore.raguStatus[examStore.currentQuestion?.id]" class="w-3 h-3 lg:w-4 lg:h-4 text-amber-500" fill="currentColor" viewBox="0 0 20 20"><path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd"></path></svg>
        </div>
        Ragu <span class="hidden sm:inline">-ragu</span>
      </button>

      <!-- Next Button -->
      <button 
        @click="handleNext" 
        :disabled="examStore.currentIdx === examStore.questions.length - 1" 
        class="flex-1 lg:flex-none px-6 lg:px-12 py-3 lg:py-4 bg-indigo-600 text-white rounded-xl text-[10px] lg:text-sm font-black hover:bg-indigo-700 disabled:opacity-30 transition-all uppercase tracking-widest active:scale-95 flex items-center justify-center gap-2 shadow-md shadow-indigo-100"
      >
        <span class="hidden sm:inline">Berikutnya</span>
        <svg class="w-4 h-4 lg:w-5 lg:h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="3" d="M9 5l7 7-7 7"></path></svg>
      </button>
    </div>
  </div>
</template>
