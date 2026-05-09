<script setup>
import { ref, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '../store/auth'
import ConfirmModal from './ConfirmModal.vue'

const authStore = useAuthStore()
const router = useRouter()
const route = useRoute()

const showLogoutConfirm = ref(false)

const isExamPage = computed(() => route.name === 'StudentPengerjaan')

const handleLogout = () => {
  authStore.logout()
  showLogoutConfirm.value = false
  router.push({ name: 'StudentLogin' })
}
</script>

<template>
  <div class="h-screen overflow-y-auto bg-slate-50 font-sans text-slate-900">
    <!-- Simple Header (Sembunyikan jika sedang ujian) -->
    <header v-if="!isExamPage" class="bg-white border-b border-slate-200 sticky top-0 z-50 shadow-sm">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 h-16 flex items-center justify-between">
        <div class="flex items-center gap-3">
          <div class="w-10 h-10 bg-blue-600 rounded-xl flex items-center justify-center text-white shadow-lg shadow-blue-200">
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253"></path></svg>
          </div>
          <h1 class="text-xl font-black text-slate-800 tracking-tight italic">CBT <span class="text-blue-600 font-black">SMP</span></h1>
        </div>

        <div class="flex items-center gap-4">
          <div class="flex flex-col items-end mr-2 text-right">
            <span class="text-[10px] font-black text-slate-400 uppercase tracking-widest leading-none mb-1">Peserta Ujian</span>
            <span class="text-sm font-black text-slate-700 leading-none">{{ authStore.user?.nama || authStore.user?.username }}</span>
          </div>
          <button v-if="!isExamPage" @click="showLogoutConfirm = true" class="p-2 text-slate-400 hover:text-rose-600 hover:bg-rose-50 rounded-xl transition-all" title="Logout">
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1"></path></svg>
          </button>
        </div>
      </div>
    </header>

    <!-- Main Content (Tanpa padding jika sedang ujian) -->
    <main :class="isExamPage ? '' : 'max-w-7xl mx-auto px-4 sm:px-6 py-6 sm:py-10'">
      <router-view v-slot="{ Component }">
        <transition name="fade" mode="out-in">
          <component :is="Component" />
        </transition>
      </router-view>
    </main>

    <!-- Modal Konfirmasi Logout -->
    <ConfirmModal 
      :show="showLogoutConfirm"
      title="Keluar Aplikasi"
      message="Apakah Anda yakin ingin keluar? Sesi ujian Anda mungkin akan terganggu jika Anda keluar sekarang."
      confirmText="Ya, Keluar"
      variant="danger"
      @close="showLogoutConfirm = false"
      @confirm="handleLogout"
    />
  </div>
</template>

<style scoped>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease, transform 0.2s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
  transform: translateY(10px);
}
</style>
