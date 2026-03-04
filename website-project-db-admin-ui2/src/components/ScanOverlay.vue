<script setup>
import { ref, onMounted, onUnmounted, computed } from 'vue'

const props = defineProps({
  url: { type: String, default: '' }
})

// Cycling status messages
const statusMessages = [
  'Verbindung wird hergestellt...',
  'DNS-Auflösung...',
  'TLS-Handshake...',
  'Seiteninhalte werden geladen...',
  'HTML-Struktur wird analysiert...',
  'Navigation wird kartografiert...',
  'Relevante Links werden priorisiert...',
  'Unterseiten werden erkundet...',
  'Textinhalte werden extrahiert...',
  'Semantische Analyse läuft...',
  'Projekt-Metadaten werden erkannt...',
  'Schlüsselwörter werden identifiziert...',
  'Terminstrukturen werden gesucht...',
  'Wiederholungsmuster werden analysiert...',
  'Tags werden zugeordnet...',
  'Konfidenz wird berechnet...',
  'Ergebnisse werden zusammengestellt...',
]

const currentMessageIndex = ref(0)
const dataLines = ref([])
const ringProgress = ref(0)

const currentMessage = computed(() => statusMessages[currentMessageIndex.value])

// Generate fake data stream lines
const generateDataLine = () => {
  const hex = () => Math.floor(Math.random() * 256).toString(16).padStart(2, '0')
  const types = [
    () => `0x${hex()}${hex()}${hex()}${hex()}  ${hex()} ${hex()} ${hex()} ${hex()} ${hex()} ${hex()} ${hex()} ${hex()}`,
    () => `GET /${['page', 'api', 'data', 'content', 'assets', 'img'][Math.floor(Math.random() * 6)]}/${hex()}${hex()} HTTP/1.1`,
    () => `<${['div', 'section', 'article', 'main', 'header', 'nav', 'p', 'h2', 'a'][Math.floor(Math.random() * 9)]} class="${hex()}${hex()}">`,
    () => `>>> node_${hex()} :: depth=${Math.floor(Math.random() * 8)} children=${Math.floor(Math.random() * 24)}`,
    () => `[${String(Math.random().toFixed(4)).padStart(6)}] match: "${['projekt', 'termin', 'kontakt', 'über uns', 'angebot', 'team', 'events'][Math.floor(Math.random() * 7)]}"`,
    () => `sha256:${hex()}${hex()}${hex()}${hex()}${hex()}${hex()}${hex()}${hex()}`,
  ]
  return types[Math.floor(Math.random() * types.length)]()
}

let messageInterval, dataInterval, progressInterval

onMounted(() => {
  // Cycle status messages
  messageInterval = setInterval(() => {
    currentMessageIndex.value = (currentMessageIndex.value + 1) % statusMessages.length
  }, 2800)

  // Add data lines
  dataInterval = setInterval(() => {
    dataLines.value.push(generateDataLine())
    if (dataLines.value.length > 12) {
      dataLines.value.shift()
    }
  }, 300)

  // Animate ring progress (loops slowly)
  progressInterval = setInterval(() => {
    ringProgress.value = (ringProgress.value + 0.4) % 100
  }, 50)
})

onUnmounted(() => {
  clearInterval(messageInterval)
  clearInterval(dataInterval)
  clearInterval(progressInterval)
})

const ringDasharray = computed(() => {
  const circumference = 2 * Math.PI * 70
  const filled = (ringProgress.value / 100) * circumference
  return `${filled} ${circumference - filled}`
})
</script>

<template>
  <div class="fixed inset-0 z-40 flex items-center justify-center bg-black/85 backdrop-blur-sm">
    <!-- Hex grid background pattern -->
    <div class="absolute inset-0 overflow-hidden opacity-20">
      <div class="hex-grid"></div>
    </div>

    <!-- Main content -->
    <div class="relative z-10 flex flex-col items-center max-w-2xl w-full px-6">

      <!-- Scan ring -->
      <div class="relative w-48 h-48 mb-8">
        <!-- Outer rotating ring -->
        <svg class="absolute inset-0 w-full h-full animate-[spin_8s_linear_infinite]" viewBox="0 0 160 160">
          <circle cx="80" cy="80" r="70" fill="none" stroke="rgba(27, 174, 112, 0.1)" stroke-width="1" />
          <circle cx="80" cy="80" r="70" fill="none" stroke="url(#ring-gradient)" stroke-width="2"
            stroke-linecap="round" :stroke-dasharray="ringDasharray" transform="rotate(-90 80 80)" />
          <defs>
            <linearGradient id="ring-gradient" x1="0%" y1="0%" x2="100%" y2="100%">
              <stop offset="0%" stop-color="#00FF88" />
              <stop offset="50%" stop-color="#1BAE70" />
              <stop offset="100%" stop-color="#FFD700" />
            </linearGradient>
          </defs>
        </svg>

        <!-- Inner counter-rotating ring -->
        <svg class="absolute inset-3 w-[calc(100%-1.5rem)] h-[calc(100%-1.5rem)] animate-[spin_12s_linear_infinite_reverse]" viewBox="0 0 160 160">
          <circle cx="80" cy="80" r="70" fill="none" stroke="rgba(212, 160, 23, 0.15)" stroke-width="1"
            stroke-dasharray="8 12" />
        </svg>

        <!-- Pulsing center -->
        <div class="absolute inset-0 flex items-center justify-center">
          <div class="w-16 h-16 rounded-full bg-wald-500/10 animate-ping"></div>
        </div>
        <div class="absolute inset-0 flex items-center justify-center">
          <div class="w-10 h-10 rounded-full border border-wald-400/40 flex items-center justify-center">
            <div class="w-3 h-3 rounded-full bg-wald-400 shadow-[0_0_12px_rgba(0,255,136,0.6)]"></div>
          </div>
        </div>

        <!-- Corner markers -->
        <div class="absolute top-0 left-0 w-4 h-4 border-t border-l border-wald-500/40"></div>
        <div class="absolute top-0 right-0 w-4 h-4 border-t border-r border-wald-500/40"></div>
        <div class="absolute bottom-0 left-0 w-4 h-4 border-b border-l border-wald-500/40"></div>
        <div class="absolute bottom-0 right-0 w-4 h-4 border-b border-r border-wald-500/40"></div>
      </div>

      <!-- URL target display -->
      <div v-if="url" class="mb-4 px-4 py-1.5 border border-wald-500/20 rounded bg-black/60">
        <span class="text-[10px] font-mono text-gray-600 uppercase tracking-widest">Target :: </span>
        <span class="text-xs font-mono text-wald-300">{{ url }}</span>
      </div>

      <!-- Status message -->
      <p class="text-sm font-mono text-wald-300 mb-6 tracking-wider transition-opacity duration-300">
        {{ currentMessage }}
      </p>

      <!-- Data stream -->
      <div class="w-full max-w-lg h-48 overflow-hidden rounded border border-wald-500/10 bg-black/60 p-3">
        <div class="flex items-center gap-2 mb-2 pb-2 border-b border-wald-500/10">
          <div class="w-1.5 h-1.5 rounded-full bg-wald-400 animate-pulse"></div>
          <span class="text-[10px] font-mono text-gray-600 uppercase tracking-widest">Live Data Stream</span>
        </div>
        <div class="space-y-0.5">
          <div
            v-for="(line, i) in dataLines"
            :key="i"
            class="text-[11px] font-mono leading-relaxed truncate"
            :class="i === dataLines.length - 1 ? 'text-wald-400' : 'text-gray-600'"
            :style="{ opacity: 0.3 + (i / dataLines.length) * 0.7 }"
          >
            {{ line }}
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.hex-grid {
  width: 200%;
  height: 200%;
  position: absolute;
  top: -50%;
  left: -50%;
  background-image:
    linear-gradient(30deg, rgba(27, 174, 112, 0.08) 12%, transparent 12.5%, transparent 87%, rgba(27, 174, 112, 0.08) 87.5%, rgba(27, 174, 112, 0.08)),
    linear-gradient(150deg, rgba(27, 174, 112, 0.08) 12%, transparent 12.5%, transparent 87%, rgba(27, 174, 112, 0.08) 87.5%, rgba(27, 174, 112, 0.08)),
    linear-gradient(30deg, rgba(27, 174, 112, 0.08) 12%, transparent 12.5%, transparent 87%, rgba(27, 174, 112, 0.08) 87.5%, rgba(27, 174, 112, 0.08)),
    linear-gradient(150deg, rgba(27, 174, 112, 0.08) 12%, transparent 12.5%, transparent 87%, rgba(27, 174, 112, 0.08) 87.5%, rgba(27, 174, 112, 0.08)),
    linear-gradient(60deg, rgba(212, 160, 23, 0.04) 25%, transparent 25.5%, transparent 75%, rgba(212, 160, 23, 0.04) 75%, rgba(212, 160, 23, 0.04)),
    linear-gradient(60deg, rgba(212, 160, 23, 0.04) 25%, transparent 25.5%, transparent 75%, rgba(212, 160, 23, 0.04) 75%, rgba(212, 160, 23, 0.04));
  background-size: 80px 140px;
  background-position: 0 0, 0 0, 40px 70px, 40px 70px, 0 0, 40px 70px;
  animation: hex-drift 20s linear infinite;
}

@keyframes hex-drift {
  0% { transform: translate(0, 0) rotate(0deg); }
  100% { transform: translate(-40px, -70px) rotate(3deg); }
}
</style>
