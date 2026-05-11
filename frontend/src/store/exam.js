import { defineStore } from 'pinia'
import { useAuthStore } from './auth'
import { useAlertStore } from './alert'

export const useExamStore = defineStore('exam', {
  state: () => ({
    examInfo: null,
    questions: [],
    answers: {},
    raguStatus: {},
    timeLeft: 0,
    currentIdx: 0,
    attemptId: 0,
    isTerblokir: false,
    isLoading: true,
    isSubmitting: false,
    
    // Sync states
    syncStatus: 'synced', // synced, pending, syncing, offline
    lastSyncTime: null,
    dirtyQuestions: new Set(),
  }),

  getters: {
    currentQuestion: (state) => state.questions[state.currentIdx] || null,
    isAllAnswered: (state) => {
      if (state.questions.length === 0) return false
      return state.questions.every(s => state.answers[s.id] !== undefined && state.answers[s.id] !== "")
    },
    progressPercentage: (state) => {
      if (state.questions.length === 0) return 0
      const answered = state.questions.filter(s => state.answers[s.id] !== undefined && state.answers[s.id] !== "").length
      return Math.round((answered / state.questions.length) * 100)
    }
  },

  actions: {
    setExamInfo(info) {
      // Selalu bersihkan state lama sebelum memulai sesi baru
      this.clearStore()
      this.examInfo = info
    },

    clearStore() {
      this.questions = []
      this.answers = {}
      this.raguStatus = {}
      this.timeLeft = 0
      this.currentIdx = 0
      this.isTerblokir = false
      this.isLoading = false
      this.isSubmitting = false
      this.syncStatus = 'synced'
      this.dirtyQuestions = new Set()
      this.examInfo = null
      this.attemptId = 0
    },

    async fetchQuestions() {
      const authStore = useAuthStore()
      const alertStore = useAlertStore()
      
      this.isLoading = true
      try {
        const res = await fetch(`/api/siswa/soal/${this.examInfo.jadwal.id}`, {
          headers: { 'Authorization': `Bearer ${authStore.token}` }
        })
        const d = await res.json()
        
        if (res.ok) {
          let rawQuestions = (d.data.items || []).map(s => {
            const ops = [
              { key: 'A', text: s.opsi_a },
              { key: 'B', text: s.opsi_b },
              { key: 'C', text: s.opsi_c },
              { key: 'D', text: s.opsi_d },
            ]
            
            // 1. Acak Jawaban jika aktif (Gunakan seed gabungan Peserta + Soal agar tetap stabil)
            const displayOptions = d.data.acak_jawaban 
              ? this.shuffleArray(ops, this.examInfo.peserta_id + s.id) 
              : ops
            
            return { ...s, displayOptions }
          })
          
          // 2. Acak Soal jika aktif (Gunakan seed Peserta agar urutan soal tetap sama bagi siswa tersebut)
          if (d.data.acak_soal) {
            this.questions = this.shuffleArray(rawQuestions, this.examInfo.peserta_id)
          } else {
            this.questions = rawQuestions
          }
          
          this.timeLeft = d.data.sisa_waktu
          this.isTerblokir = d.data.is_terblokir
          
          // Jika status adalah 'Belum Mengerjakan' (hasil Reset Total dari Admin),
          // Hapus sisa-sisa memori dan LocalStorage lama agar tidak mencemari ujian baru.
          if (d.data.status_ujian === 'Belum Mengerjakan' || (this.attemptId > 0 && d.data.attempt_id > this.attemptId)) {
            console.log("!!! RESET DETECTED !!!", "Server Attempt:", d.data.attempt_id, "Local Attempt:", this.attemptId)
            this.clearLocalStorage()
            this.answers = {}
            this.raguStatus = {}
            this.currentIdx = 0
          }
          this.attemptId = d.data.attempt_id

          // Load existing answers from DB
          if (d.data.existing_answers) {
            Object.keys(d.data.existing_answers).forEach(sId => {
              const item = d.data.existing_answers[sId]
              this.answers[sId] = item.jawaban
              this.raguStatus[sId] = item.ragu_ragu
            })
          }

          // Fallback from LocalStorage
          this.loadFromLocalStorage()
        }
      } catch (error) {
        alertStore.showAlert("Gagal memuat soal. Periksa koneksi internet.", "error")
      } finally {
        this.isLoading = false
      }
    },

    setAnswer(soalId, ans) {
      if (this.answers[soalId] !== ans) {
        this.answers[soalId] = ans
        this.dirtyQuestions.add(soalId)
        this.syncStatus = 'pending'
        this.saveToLocalStorage()
      }
    },

    toggleRagu(soalId) {
      this.raguStatus[soalId] = !this.raguStatus[soalId]
      this.dirtyQuestions.add(soalId)
      this.syncStatus = 'pending'
      this.saveToLocalStorage()
    },

    nextQuestion() {
      if (this.currentIdx < this.questions.length - 1) {
        this.currentIdx++
      }
    },

    prevQuestion() {
      if (this.currentIdx > 0) {
        this.currentIdx--
      }
    },

    jumpToQuestion(index) {
      if (index >= 0 && index < this.questions.length) {
        this.currentIdx = index
      }
    },

    decrementTimer() {
      if (this.timeLeft > 0) {
        this.timeLeft--
      }
    },

    async syncData(force = false) {
      const authStore = useAuthStore()
      
      // Tetap jalankan sync jika:
      // 1. Ada pertanyaan baru (dirty)
      // 2. Dipaksa (force)
      // 3. SEDANG TERBLOKIR (Penting: agar bisa polling status "Lepas Blokir")
      const shouldSync = this.dirtyQuestions.size > 0 || force || this.isTerblokir
      
      if (!this.examInfo || !shouldSync || this.syncStatus === 'syncing') return
      
      this.syncStatus = 'syncing'
      
      const itemsToSync = []
      this.dirtyQuestions.forEach(sId => {
        itemsToSync.push({
          soal_id: sId,
          jawaban_teks: this.answers[sId] || "",
          ragu_ragu: !!this.raguStatus[sId]
        })
      })

      const payload = {
        peserta_ujian_id: this.examInfo.peserta_id,
        items: itemsToSync,
        sisa_waktu: this.timeLeft,
        attempt_id: this.attemptId
      }

      try {
        const res = await fetch('/api/siswa/sync', {
          method: 'POST',
          headers: { 
            'Authorization': `Bearer ${authStore.token}`,
            'Content-Type': 'application/json' 
          },
          body: JSON.stringify(payload)
        })

        const data = await res.json()

        if (res.status === 401) {
          authStore.logout()
          return 'unauthorized'
        }

        if (res.status === 403) {
          this.isTerblokir = true
          this.syncStatus = 'offline'
          return 'blocked'
        }

        if (res.ok) {
          // JIKA BERHASIL: Berarti sudah tidak terblokir
          this.isTerblokir = false 
          this.dirtyQuestions.clear()
          this.syncStatus = 'synced'
          this.lastSyncTime = new Date().toLocaleTimeString('id-ID')
          
          // Deteksi Reset Total dari hasil Sync (jika status returned)
          if (data.data && data.data.status === 'OUTDATED_SESSION') {
            console.error("Sesi Kedaluwarsa! Membersihkan data...")
            this.clearLocalStorage()
            window.location.reload() // Paksa refresh untuk ambil status baru
            return 'outdated'
          }

          if (data.data.sisa_waktu !== undefined && Math.abs(data.data.sisa_waktu - this.timeLeft) > 5) {
            this.timeLeft = data.data.sisa_waktu
          }
          return 'success'
        } else {
          throw new Error("Sync failed")
        }
      } catch (e) {
        this.syncStatus = 'offline'
        return 'error'
      }
    },

    saveToLocalStorage() {
      if (this.examInfo) {
        const dataToSave = {
          answers: this.answers,
          raguStatus: this.raguStatus
        }
        localStorage.setItem(`answers_${this.examInfo.peserta_id}`, JSON.stringify(dataToSave))
      }
    },

    clearLocalStorage() {
      if (this.examInfo) {
        localStorage.removeItem(`answers_${this.examInfo.peserta_id}`)
      }
    },

    loadFromLocalStorage() {
      const savedLocal = localStorage.getItem(`answers_${this.examInfo.peserta_id}`)
      if (savedLocal) {
        try {
          const localData = JSON.parse(savedLocal)
          // Support old format and new format with raguStatus
          const answers = localData.answers || localData
          const ragu = localData.raguStatus || {}

          Object.keys(answers).forEach(sId => {
            if (this.answers[sId] === undefined || this.answers[sId] === "") {
              this.answers[sId] = answers[sId]
              this.dirtyQuestions.add(parseInt(sId))
              this.syncStatus = 'pending'
            }
          })

          Object.keys(ragu).forEach(sId => {
            if (this.raguStatus[sId] === undefined) {
              this.raguStatus[sId] = ragu[sId]
            }
          })
        } catch (e) {
          console.warn("LocalStorage data corrupted")
        }
      }
    },

    // Generator angka acak yang deterministik (Seeded PRNG)
    seededRandom(seed) {
      const x = Math.sin(seed) * 10000;
      return x - Math.floor(x);
    },

    shuffleArray(array, seed) {
      const newArr = [...array]
      let m = newArr.length, t, i;
      let s = seed;
      
      while (m) {
        i = Math.floor(this.seededRandom(s++) * m--);
        t = newArr[m];
        newArr[m] = newArr[i];
        newArr[i] = t;
      }
      return newArr
    }
  }
})
