<script setup>
defineProps({
  show: Boolean,
  title: String,
  isLoading: Boolean,
  confirmText: { type: String, default: 'Simpan Data' },
  cancelText: { type: String, default: 'Batal' },
  size: { type: String, default: 'max-w-2xl' },
  variant: { type: String, default: 'primary' } // primary, danger, success
})

const emit = defineEmits(['close', 'confirm'])
</script>

<template>
  <Transition name="modal-fade">
    <div v-if="show" class="fixed inset-0 z-[100] flex items-center justify-center p-4 bg-slate-900/60 backdrop-blur-md">
      <div 
        class="bg-white w-full rounded-3xl shadow-2xl flex flex-col overflow-hidden animate-in zoom-in duration-200 max-h-[90vh]"
        :class="size"
      >
        <!-- Header -->
        <div class="p-6 border-b border-slate-100 flex justify-between items-center bg-slate-50/50">
          <div>
            <h3 class="text-xl font-bold text-slate-800">{{ title }}</h3>
          </div>
          <button @click="emit('close')" class="text-slate-400 hover:text-slate-600 p-2 transition-colors rounded-full hover:bg-slate-100">
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path></svg>
          </button>
        </div>

        <!-- Body -->
        <div class="flex-1 overflow-y-auto p-6">
          <slot></slot>
        </div>

        <!-- Footer -->
        <div class="p-6 border-t border-slate-100 bg-white flex justify-end gap-3">
          <button 
            @click="emit('close')" 
            class="px-6 py-2.5 text-sm font-bold text-slate-500 hover:bg-slate-100 rounded-xl transition-all"
          >
            {{ cancelText }}
          </button>
          <button 
            @click="emit('confirm')" 
            :disabled="isLoading"
            class="px-8 py-2.5 text-white rounded-xl text-sm font-bold shadow-lg flex items-center gap-2 transition-all disabled:opacity-50 disabled:cursor-not-allowed"
            :class="{
              'bg-blue-600 hover:bg-blue-700 shadow-blue-200': variant === 'primary',
              'bg-rose-600 hover:bg-rose-700 shadow-rose-200': variant === 'danger',
              'bg-emerald-600 hover:bg-emerald-700 shadow-emerald-200': variant === 'success'
            }"
          >
            <svg v-if="isLoading" class="animate-spin h-4 w-4 text-white" fill="none" viewBox="0 0 24 24"><circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle><path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path></svg>
            <span>{{ isLoading ? 'Memproses...' : confirmText }}</span>
          </button>
        </div>
      </div>
    </div>
  </Transition>
</template>

<style scoped>
.modal-fade-enter-active,
.modal-fade-leave-active {
  transition: opacity 0.3s ease;
}

.modal-fade-enter-from,
.modal-fade-leave-to {
  opacity: 0;
}
</style>
