<script setup>
import { useAlertStore } from '../store/alert'

const alertStore = useAlertStore()
</script>

<template>
  <Transition name="slide-fade">
    <div 
      v-if="alertStore.show" 
      class="fixed top-6 right-6 z-[300] flex items-center gap-4 px-6 py-4 rounded-2xl shadow-2xl border min-w-[320px] max-w-md animate-in slide-in-from-right duration-300"
      :class="{
        'bg-emerald-50 border-emerald-100 text-emerald-800': alertStore.type === 'success',
        'bg-rose-50 border-rose-100 text-rose-800': alertStore.type === 'error',
        'bg-amber-50 border-amber-100 text-amber-800': alertStore.type === 'warning',
        'bg-blue-50 border-blue-100 text-blue-800': alertStore.type === 'info',
      }"
    >
      <!-- Icon Success -->
      <div v-if="alertStore.type === 'success'" class="w-10 h-10 rounded-full bg-emerald-500/10 flex items-center justify-center flex-shrink-0">
        <svg class="w-6 h-6 text-emerald-600" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"></path></svg>
      </div>
      
      <!-- Icon Error -->
      <div v-if="alertStore.type === 'error'" class="w-10 h-10 rounded-full bg-rose-500/10 flex items-center justify-center flex-shrink-0">
        <svg class="w-6 h-6 text-rose-600" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path></svg>
      </div>

      <!-- Icon Warning -->
      <div v-if="alertStore.type === 'warning'" class="w-10 h-10 rounded-full bg-amber-500/10 flex items-center justify-center flex-shrink-0">
        <svg class="w-6 h-6 text-amber-600" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"></path></svg>
      </div>

      <div class="flex-1">
        <h4 class="text-xs font-black uppercase tracking-widest opacity-50 mb-0.5">
          {{ alertStore.type }}
        </h4>
        <p class="text-sm font-bold leading-tight">{{ alertStore.message }}</p>
      </div>

      <button @click="alertStore.closeAlert" class="text-current opacity-40 hover:opacity-100 transition-opacity">
        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path></svg>
      </button>
    </div>
  </Transition>
</template>

<style scoped>
.slide-fade-enter-active {
  transition: all 0.3s ease-out;
}

.slide-fade-leave-active {
  transition: all 0.2s cubic-bezier(1, 0.5, 0.8, 1);
}

.slide-fade-enter-from,
.slide-fade-leave-to {
  transform: translateX(20px);
  opacity: 0;
}
</style>
