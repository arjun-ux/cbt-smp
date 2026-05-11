<script setup>
import { computed } from 'vue'
import { useExamStore } from '../../../store/exam'

const examStore = useExamStore()

const formatTime = (seconds) => {
  const h = Math.floor(seconds / 3600)
  const m = Math.floor((seconds % 3600) / 60)
  const s = seconds % 60
  return `${h.toString().padStart(2, '0')}:${m.toString().padStart(2, '0')}:${s.toString().padStart(2, '0')}`
}

defineEmits(['open-nav', 'finish'])
</script>

<template>
  <div class="sticky top-0 z-[60] bg-white/80 backdrop-blur-xl border-b border-slate-200/60 px-4 py-3 lg:px-10 lg:py-4 shadow-[0_10px_30px_rgba(0,0,0,0.02)]">
    <div class="max-w-[1400px] mx-auto flex items-center justify-between">
      <!-- Left: Subject & Progress -->
      <div class="flex items-center gap-4 lg:gap-8">
        <div class="hidden md:block">
          <p class="text-[9px] font-black text-slate-400 uppercase tracking-[0.2em] mb-0.5">Ujian Berlangsung</p>
          <h2 class="text-sm lg:text-base font-black text-slate-800 tracking-tight truncate max-w-[200px]">
            {{ examStore.examInfo?.jadwal?.nama_ujian || 'Ujian Sekolah' }}
          </h2>
        </div>
        <div class="h-8 w-[1px] bg-slate-200 hidden md:block"></div>
        <div class="flex items-center gap-3">
          <div class="bg-indigo-600 text-white px-4 py-1.5 rounded-2xl font-black text-[10px] lg:text-xs uppercase tracking-[0.15em] shadow-lg shadow-indigo-100">
            {{ examStore.currentIdx + 1 }} / {{ examStore.questions.length }}
          </div>
        </div>
      </div>
      
      <!-- Right: Timer & Actions -->
      <div class="flex items-center gap-3 lg:gap-6">
        <!-- SYNC INDICATOR (NEW: Mobile Friendly) -->
        <div 
          :title="examStore.syncStatus"
          class="w-8 h-8 lg:w-10 lg:h-10 rounded-xl flex items-center justify-center transition-all duration-500 border"
          :class="[
            examStore.syncStatus === 'synced' ? 'bg-emerald-50 border-emerald-100 text-emerald-500' : 
            examStore.syncStatus === 'pending' ? 'bg-amber-50 border-amber-100 text-amber-500' :
            examStore.syncStatus === 'offline' ? 'bg-rose-50 border-rose-100 text-rose-500' : 'bg-indigo-50 border-indigo-100 text-indigo-500'
          ]"
        >
          <div class="w-2 h-2 rounded-full bg-current" :class="{ 'animate-ping': examStore.syncStatus !== 'synced' }"></div>
        </div>

        <div :class="examStore.timeLeft < 300 ? 'bg-rose-50 text-rose-600 border-rose-100 animate-pulse' : 'bg-slate-50 text-slate-700 border-slate-100'" class="flex items-center gap-3 px-3 py-2 lg:px-5 lg:py-2.5 rounded-xl lg:rounded-2xl border transition-all duration-500">
          <svg class="w-4 h-4 hidden xs:block" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"></path></svg>
          <span class="text-xs lg:text-sm font-black font-mono tracking-widest">{{ formatTime(examStore.timeLeft) }}</span>
        </div>

        
        <button @click="$emit('open-nav')" class="lg:hidden flex items-center justify-center w-10 h-10 bg-slate-100 text-slate-600 rounded-xl hover:bg-slate-200 transition-colors">
          <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16"></path></svg>
        </button>
      </div>
    </div>
  </div>
</template>
