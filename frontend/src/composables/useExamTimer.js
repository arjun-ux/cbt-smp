import { onMounted, onUnmounted } from 'vue'
import { useExamStore } from '../store/exam'

export function useExamTimer(onFinish) {
  const examStore = useExamStore()
  let timerInterval = null

  const startTimer = () => {
    timerInterval = setInterval(() => {
      if (examStore.timeLeft > 0) {
        examStore.decrementTimer()
        
        // Auto-sync every 15 seconds
        if (examStore.timeLeft % 15 === 0) {
          examStore.syncData()
        }
      } else {
        stopTimer()
        if (onFinish) onFinish(true)
      }
    }, 1000)
  }

  const stopTimer = () => {
    if (timerInterval) {
      clearInterval(timerInterval)
      timerInterval = null
    }
  }

  onMounted(() => {
    startTimer()
  })

  onUnmounted(() => {
    stopTimer()
  })

  return {
    startTimer,
    stopTimer
  }
}
