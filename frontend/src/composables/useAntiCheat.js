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

  onMounted(() => {
    toggleFullscreen()
    window.addEventListener('blur', handleBlur)
  })

  onUnmounted(() => {
    window.removeEventListener('blur', handleBlur)
  })

  return {
    toggleFullscreen,
    handleSecurityBreach
  }
}
