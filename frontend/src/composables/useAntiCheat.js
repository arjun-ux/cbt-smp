import { onMounted, onUnmounted } from 'vue'
import { useAuthStore } from '../store/auth'
import { useExamStore } from '../store/exam'

export function useAntiCheat() {
  const authStore = useAuthStore()
  const examStore = useExamStore()

  const handleSecurityBreach = async (reason) => {
    if (!examStore.examInfo || examStore.isTerblokir) return
    
    try {
      const res = await fetch('/api/siswa/log', {
        method: 'POST',
        headers: { 
          'Authorization': `Bearer ${authStore.token}`,
          'Content-Type': 'application/json' 
        },
        body: JSON.stringify({
          peserta_ujian_id: examStore.examInfo.peserta_id,
          keterangan: reason
        })
      })
      
      if (res.ok) {
        examStore.isTerblokir = true
      }
    } catch (e) {
      console.error("Failed to log security breach:", e)
    }
  }

  const toggleFullscreen = () => {
    if (!document.fullscreenElement) {
      document.documentElement.requestFullscreen().catch(err => {
        console.log(`Error: ${err.message}`)
      })
    }
  }

  const handleBlur = () => {
    handleSecurityBreach("Pindah Tab / Keluar Layar")
  }

  const preventDefaults = (e) => {
    e.preventDefault()
    return false
  }

  const handleKeydown = (e) => {
    // Blokir F12, Ctrl+Shift+I, Ctrl+Shift+J, Ctrl+U
    if (
      e.keyCode === 123 || 
      (e.ctrlKey && e.shiftKey && (e.keyCode === 73 || e.keyCode === 74)) || 
      (e.ctrlKey && e.keyCode === 85)
    ) {
      e.preventDefault()
      return false
    }
  }

  onMounted(() => {
    toggleFullscreen()
    window.addEventListener('blur', handleBlur)
    
    // Keamanan Konten - Hanya aktif jika BUKAN mode Development
    if (!import.meta.env.DEV) {
      document.addEventListener('contextmenu', preventDefaults)
      document.addEventListener('copy', preventDefaults)
      document.addEventListener('paste', preventDefaults)
      document.addEventListener('selectstart', preventDefaults)
      document.addEventListener('keydown', handleKeydown)
    } else {
      console.log("Anti-Cheat: Mode Development terdeteksi, proteksi konten dinonaktifkan untuk debugging.")
    }
  })

  onUnmounted(() => {
    window.removeEventListener('blur', handleBlur)
    if (!import.meta.env.DEV) {
      document.removeEventListener('contextmenu', preventDefaults)
      document.removeEventListener('copy', preventDefaults)
      document.removeEventListener('paste', preventDefaults)
      document.removeEventListener('selectstart', preventDefaults)
      document.removeEventListener('keydown', handleKeydown)
    }
  })

  return {
    toggleFullscreen,
    handleSecurityBreach
  }
}
