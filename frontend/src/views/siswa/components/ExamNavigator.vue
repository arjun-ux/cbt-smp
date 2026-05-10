<script setup>
import { useExamStore } from '../../../store/exam'

const examStore = useExamStore()

const props = defineProps({
  isMobile: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['close', 'finish'])

const handleJump = (idx) => {
  examStore.jumpToQuestion(idx)
  if (props.isMobile) emit('close')
}
</script>

<template>
  <div :class="isMobile ? '' : 'bg-white p-6 rounded-2xl shadow-sm border border-slate-100 sticky top-24'">
    <div v-if="!isMobile" class="flex items-center gap-2 mb-6">
      <div class="w-1 h-3 bg-indigo-500 rounded-full"></div>
      <h4 class="text-[9px] font-black text-slate-400 uppercase tracking-widest">Navigasi Soal</h4>
    </div>

    <div class="grid grid-cols-5 gap-2" :class="isMobile ? 'max-h-[50vh] overflow-y-auto pb-6' : ''">
      <button 
        v-for="(s, idx) in examStore.questions" 
        :key="s.id"
        @click="handleJump(idx)"
        :class="[
          examStore.currentIdx === idx 
            ? 'border-indigo-600 ring-2 ring-indigo-500/10' 
            : 'border-transparent',
          
          examStore.raguStatus[s.id] 
            ? 'bg-amber-500 text-white shadow-sm' 
            : (examStore.answers[s.id] 
                ? 'bg-indigo-600 text-white shadow-sm' 
                : 'bg-slate-50 text-slate-400 hover:bg-slate-100')
        ]"
        class="w-full aspect-square border-2 rounded-lg text-xs font-black transition-all duration-200 flex items-center justify-center relative overflow-hidden"
      >
        {{ idx + 1 }}
      </button>
    </div>
    
    <div v-if="!isMobile" class="mt-8 pt-6 border-t border-slate-100">
      <button 
        @click="$emit('finish')"
        :disabled="!examStore.isAllAnswered"
        class="w-full py-4 bg-slate-900 text-white hover:bg-rose-600 disabled:bg-slate-100 disabled:text-slate-300 rounded-xl font-black text-[10px] uppercase tracking-widest transition-all duration-200 shadow-lg shadow-slate-100 flex items-center justify-center gap-2"
      >
        <span>Selesai Ujian</span>
        <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="3" d="M13 7l5 5m0 0l-5 5m5-5H6"></path></svg>
      </button>
    </div>
  </div>
</template>

<style scoped>
@keyframes shimmer {
  100% { transform: translateX(100%); }
}
.animate-shimmer {
  animation: shimmer 1.5s infinite;
}
</style>
