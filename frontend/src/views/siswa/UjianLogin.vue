<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../../store/auth'
import { useAlertStore } from '../../store/alert'

const authStore = useAuthStore()
const alertStore = useAlertStore()
const router = useRouter()

const nisn = ref('')
const password = ref('')
const isLoading = ref(false)

const handleLogin = async () => {
  if (!nisn.value || !password.value) {
    alertStore.showAlert("NISN dan Password wajib diisi", "warning")
    return
  }

  isLoading.value = true
  try {
    const res = await fetch('/api/auth/login', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ username: nisn.value, password: password.value })
    })

    const data = await res.json()
    if (res.ok) {
      authStore.setAuth(data.data.token, data.data.user)
      alertStore.showAlert("Login Berhasil", "success")
      router.push({ name: 'StudentDashboard' })
    } else {
      alertStore.showAlert(data.error || "Gagal masuk. Cek NISN & Password.", "error")
    }
  } catch (error) {
    alertStore.showAlert("Kesalahan koneksi ke server", "error")
  } finally {
    isLoading.value = false
  }
}
</script>

<template>
  <div class="min-h-screen bg-slate-50 flex flex-col items-center justify-center p-6 relative font-sans">
    <!-- Abstract Background -->
    <div class="absolute inset-0 overflow-hidden pointer-events-none">
      <div class="absolute -top-[10%] -left-[10%] w-[40%] h-[40%] bg-blue-100 blur-[120px] rounded-full"></div>
      <div class="absolute -bottom-[10%] -right-[10%] w-[40%] h-[40%] bg-cyan-100 blur-[120px] rounded-full"></div>
    </div>

    <div class="w-full max-w-md z-10">
      <!-- Logo/Brand -->
      <div class="text-center mb-10">
        <div class="inline-flex items-center justify-center w-20 h-20 bg-gradient-to-tr from-blue-600 to-cyan-500 rounded-3xl shadow-2xl shadow-blue-500/20 mb-6 rotate-3">
          <svg class="w-10 h-10 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253"></path></svg>
        </div>
        <h1 class="text-4xl font-black text-slate-800 tracking-tight italic">CBT <span class="text-blue-600 not-italic uppercase">System</span></h1>
        <p class="text-slate-400 mt-2 font-bold tracking-widest uppercase text-xs">Portal Peserta Ujian v1.0</p>
      </div>

      <!-- Login Card -->
      <div class="bg-white/80 backdrop-blur-xl p-8 rounded-[2.5rem] border border-white shadow-[0_20px_50px_rgba(0,0,0,0.05)]">
        <form @submit.prevent="handleLogin" class="space-y-6">
          <div>
            <label class="block text-slate-400 text-[10px] font-black uppercase tracking-widest mb-2 ml-1">Nomor Induk Siswa Nasional</label>
            <div class="relative group">
              <div class="absolute inset-y-0 left-0 pl-4 flex items-center pointer-events-none">
                <svg class="w-5 h-5 text-slate-300 group-focus-within:text-blue-500 transition-colors" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 6H5a2 2 0 00-2 2v9a2 2 0 002 2h14a2 2 0 002-2V8a2 2 0 00-2-2h-5m-4 0V5a2 2 0 114 0v1m-4 0a2 2 0 104 0m-5 8a2 2 0 100-4 2 2 0 000 4zm5 3a3 3 0 01-3 3H9a3 3 0 01-3-3v-1h10v1z"></path></svg>
              </div>
              <input 
                v-model="nisn"
                type="text" 
                required
                placeholder="Masukkan NISN Anda"
                class="block w-full pl-12 pr-4 py-4 bg-slate-50 border border-slate-100 rounded-2xl text-slate-800 placeholder-slate-300 focus:ring-4 focus:ring-blue-500/10 focus:border-blue-500 focus:bg-white transition-all font-bold"
              >
            </div>
          </div>

          <div>
            <label class="block text-slate-400 text-[10px] font-black uppercase tracking-widest mb-2 ml-1">Kata Sandi (Password)</label>
            <div class="relative group">
              <div class="absolute inset-y-0 left-0 pl-4 flex items-center pointer-events-none">
                <svg class="w-5 h-5 text-slate-300 group-focus-within:text-cyan-500 transition-colors" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z"></path></svg>
              </div>
              <input 
                v-model="password"
                type="password" 
                required
                placeholder="••••••••"
                class="block w-full pl-12 pr-4 py-4 bg-slate-50 border border-slate-100 rounded-2xl text-slate-800 placeholder-slate-300 focus:ring-4 focus:ring-cyan-500/10 focus:border-cyan-500 focus:bg-white transition-all font-bold"
              >
            </div>
          </div>

          <button 
            type="submit" 
            :disabled="isLoading"
            class="w-full bg-blue-600 hover:bg-blue-700 text-white font-black py-4 rounded-2xl shadow-xl shadow-blue-200 transform transition-all active:scale-[0.98] disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-3"
          >
            <span v-if="isLoading" class="w-5 h-5 border-2 border-white/30 border-t-white rounded-full animate-spin"></span>
            {{ isLoading ? 'VERIFIKASI...' : 'MASUK UJIAN' }}
          </button>
        </form>
      </div>

      <!-- Footer Info -->
      <div class="mt-10 text-center space-y-2">
        <p class="text-slate-400 text-[10px] font-bold uppercase tracking-[0.3em]">Computer Based Test System</p>
        <p class="text-slate-300 text-[9px] font-medium italic">Pastikan Anda menggunakan browser versi terbaru dan koneksi stabil.</p>
      </div>
    </div>
  </div>
</template>
