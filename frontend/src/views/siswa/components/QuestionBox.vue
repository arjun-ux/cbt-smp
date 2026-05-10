<script setup>
import { useExamStore } from '../../../store/exam'

const examStore = useExamStore()

defineProps({
  soal: Object
})
</script>

<template>
  <div class="bg-white rounded-2xl shadow-sm border border-slate-100 overflow-hidden min-h-[400px] flex flex-col transition-all duration-500">
    <!-- Header Box -->
    <div class="px-5 py-4 bg-slate-50/50 border-b border-slate-100/80 flex justify-between items-center">
      <div class="flex items-center gap-2">
        <div class="w-1.5 h-1.5 bg-indigo-500 rounded-full animate-pulse"></div>
        <h1 class="text-[9px] lg:text-[10px] font-black text-slate-400 uppercase tracking-widest">Pertanyaan {{ examStore.currentIdx + 1 }} / {{ examStore.questions.length }}</h1>
      </div>
      <div v-if="examStore.answers[soal?.id]" class="flex items-center gap-2 px-3 py-1 bg-emerald-50 text-emerald-600 rounded-xl border border-emerald-100">
        <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="3" d="M5 13l4 4L19 7"></path></svg>
        <span class="text-[8px] font-black uppercase tracking-widest">Saved</span>
      </div>
    </div>

    <!-- Content Box -->
    <div class="p-6 lg:p-10 flex-1">
      <!-- Question Text -->
      <div class="prose prose-slate max-w-none">
        <div class="text-base lg:text-xl font-bold text-slate-800 leading-relaxed mb-8 lg:mb-10" v-html="soal?.pertanyaan"></div>
      </div>
      
      <!-- PG Options -->
      <div v-if="soal?.jenis_soal === 'PG'" class="grid grid-cols-1 gap-3">
        <button 
          v-for="(opt, idx) in soal.displayOptions" 
          :key="opt.key"
          @click="examStore.setAnswer(soal.id, opt.key)"
          :class="[
            examStore.answers[soal?.id] === opt.key 
            ? 'bg-indigo-600 border-indigo-600 text-white shadow-lg shadow-indigo-50' 
            : 'bg-white border-slate-100 text-slate-600 hover:border-indigo-200 hover:bg-slate-50'
          ]"
          class="w-full text-left p-4 lg:p-5 border-2 rounded-2xl transition-all duration-200 flex items-start gap-4 group relative overflow-hidden"
        >
          <!-- Letter Bubble -->
          <span :class="[
            examStore.answers[soal?.id] === opt.key 
            ? 'bg-white/20 text-white' 
            : 'bg-slate-100 text-slate-400 group-hover:bg-indigo-50 group-hover:text-indigo-600'
          ]" class="w-8 h-8 lg:w-10 lg:h-10 rounded-xl flex items-center justify-center font-black uppercase text-xs lg:text-sm flex-shrink-0 transition-all shadow-sm">
            {{ String.fromCharCode(65 + idx) }}
          </span>

          <!-- Option Text -->
          <span class="font-bold pt-1.5 lg:pt-2 text-sm lg:text-lg leading-relaxed" v-html="opt.text"></span>
        </button>
      </div>
      
      <!-- Essay Input -->
      <div v-else class="space-y-4">
        <div class="flex items-center gap-2 ml-1">
          <div class="w-1 h-3 bg-indigo-500 rounded-full"></div>
          <label class="text-[9px] font-black text-slate-400 uppercase tracking-widest block">Jawaban Essay</label>
        </div>
        <textarea 
          :value="examStore.answers[soal?.id]"
          @input="e => examStore.setAnswer(soal.id, e.target.value)"
          placeholder="Tuliskan jawaban lengkap Anda di sini..."
          class="w-full h-48 lg:h-64 p-6 lg:p-8 bg-slate-50/50 border-2 border-slate-100 rounded-2xl focus:border-indigo-500 transition-all outline-none font-bold text-slate-700 text-base lg:text-lg placeholder-slate-300"
        ></textarea>
      </div>
    </div>
  </div>
</template>
