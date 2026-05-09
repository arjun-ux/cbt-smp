import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useAlertStore = defineStore('alert', () => {
  const show = ref(false)
  const message = ref('')
  const type = ref('success') // success, error, warning, info

  const showAlert = (msg, t = 'success') => {
    message.value = msg
    type.value = t
    show.value = true
    
    // Otomatis tutup (Error lebih lama agar bisa dibaca)
    const timeout = t === 'error' || t === 'warning' ? 8000 : 5000
    setTimeout(() => {
      show.value = false
    }, timeout)
  }

  const closeAlert = () => {
    show.value = false
  }

  return { show, message, type, showAlert, closeAlert }
})
