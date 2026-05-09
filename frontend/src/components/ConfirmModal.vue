<script setup>
import BaseModal from './BaseModal.vue'

defineProps({
  show: Boolean,
  title: { type: String, default: 'Konfirmasi Tindakan' },
  message: String,
  isLoading: Boolean,
  confirmText: { type: String, default: 'Ya, Lanjutkan' },
  variant: { type: String, default: 'danger' } // danger, primary, success
})

const emit = defineEmits(['close', 'confirm'])
</script>

<template>
  <BaseModal 
    :show="show" 
    :title="title" 
    :isLoading="isLoading" 
    :confirmText="confirmText"
    :variant="variant"
    size="max-w-md"
    @close="emit('close')" 
    @confirm="emit('confirm')"
  >
    <div class="flex items-start gap-5">
      <div 
        class="w-14 h-14 rounded-2xl flex items-center justify-center flex-shrink-0"
        :class="{
          'bg-rose-50 text-rose-600': variant === 'danger',
          'bg-blue-50 text-blue-600': variant === 'primary',
          'bg-emerald-50 text-emerald-600': variant === 'success',
        }"
      >
        <svg v-if="variant === 'danger'" class="w-8 h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path></svg>
        <svg v-else-if="variant === 'primary'" class="w-8 h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1"></path></svg>
        <svg v-else class="w-8 h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"></path></svg>
      </div>
      <div class="flex-1">
        <p class="text-slate-600 text-sm leading-relaxed font-medium">{{ message }}</p>
      </div>
    </div>
  </BaseModal>
</template>
