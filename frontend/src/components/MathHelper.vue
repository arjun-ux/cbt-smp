<script setup>
const mathTips = [
  { label: 'Pecahan', code: '\\frac{a}{b}', preview: '\\frac{1}{2}' },
  { label: 'Pangkat', code: 'x^{2}', preview: 'x^{2}' },
  { label: 'Akar', code: '\\sqrt{x}', preview: '\\sqrt{25}' },
  { label: 'Perkalian', code: '\\times', preview: '5 \\times 5' },
  { label: 'Pembagian', code: '\\div', preview: '10 \\div 2' },
  { label: 'Kurang Lebih', code: '\\pm', preview: '\\pm 5' },
  { label: 'Derajat', code: '90^{\\circ}', preview: '90^{\\circ}' },
  { label: 'Garis Bawah', code: '\\underline{x}', preview: '\\underline{Teks}' },
]

const copyToClipboard = (text) => {
  navigator.clipboard.writeText('$' + text + '$')
  // Kita bisa tambah toast kecil di sini jika mau
}
</script>

<template>
  <div class="bg-slate-900 text-white rounded-2xl p-6 shadow-xl border border-slate-700">
    <div class="flex items-center gap-3 mb-6">
      <div class="w-10 h-10 bg-blue-500 rounded-xl flex items-center justify-center shadow-lg shadow-blue-500/20">
        <span class="font-black text-lg">Σ</span>
      </div>
      <div>
        <h4 class="font-bold text-sm">Panduan Rumus Matematika</h4>
        <p class="text-[10px] text-slate-400 uppercase tracking-widest font-bold">Gunakan tanda $ di awal dan akhir rumus</p>
      </div>
    </div>

    <div class="grid grid-cols-2 gap-3">
      <div 
        v-for="tip in mathTips" 
        :key="tip.code"
        @click="copyToClipboard(tip.code)"
        class="group p-3 bg-slate-800/50 border border-slate-700 rounded-xl hover:border-blue-500 transition-all cursor-pointer relative overflow-hidden"
      >
        <div class="flex flex-col gap-1">
          <span class="text-[10px] font-bold text-slate-500 uppercase tracking-tighter">{{ tip.label }}</span>
          <code class="text-xs text-blue-400 font-mono">/{{ tip.code }}</code>
        </div>
        <div class="mt-2 text-lg font-serif text-slate-200">
          <!-- Kita gunakan gambar statis atau simbol jika KaTeX belum render di sini -->
          <span v-if="tip.label === 'Pecahan'">½</span>
          <span v-else-if="tip.label === 'Pangkat'">x²</span>
          <span v-else-if="tip.label === 'Akar'">√</span>
          <span v-else-if="tip.label === 'Perkalian'">×</span>
          <span v-else-if="tip.label === 'Pembagian'">÷</span>
          <span v-else-if="tip.label === 'Kurang Lebih'">±</span>
          <span v-else-if="tip.label === 'Derajat'">°</span>
          <span v-else-if="tip.label === 'Garis Bawah'"><u>x</u></span>
        </div>
        <div class="absolute inset-y-0 right-0 w-8 bg-blue-500 flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity">
          <svg class="w-4 h-4 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 012-2v-8a2 2 0 01-2-2h-8a2 2 0 01-2 2v8a2 2 0 012 2z"></path></svg>
        </div>
      </div>
    </div>

    <div class="mt-6 p-4 bg-blue-500/10 border border-blue-500/20 rounded-xl">
      <p class="text-[10px] text-blue-300 leading-relaxed font-medium italic">
        * Tip: Klik pada kotak untuk menyalin kode rumus. Tempelkan (Paste) di editor dan pastikan rumus diapit tanda dolar. Contoh: <span class="text-white not-italic">$x^2$</span>
      </p>
    </div>
  </div>
</template>
