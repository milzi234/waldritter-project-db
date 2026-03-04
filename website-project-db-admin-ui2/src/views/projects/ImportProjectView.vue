<script setup>
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import VueDatePicker from '@vuepic/vue-datepicker';
import '@vuepic/vue-datepicker/dist/main.css';

import useProjectAPI from '../../api/Project';
import useExtractionAPI from '../../api/Extraction';
import TagChooser from '../../components/TagChooser.vue';
import ScanOverlay from '../../components/ScanOverlay.vue';

const props = defineProps({
  defaultMode: {
    type: String,
    default: 'url'
  }
})

const router = useRouter()
const projectAPI = useProjectAPI();
const extractionAPI = useExtractionAPI();

// Import mode
const importMode = ref(props.defaultMode)
const textInput = ref('')

// URL input
const url = ref('')
const isLoading = ref(false)
const error = ref('')
const hasExtracted = ref(false)

// Extracted data
const title = ref('')
const description = ref('')
const homepage = ref('')
const keywords = ref([])

// Tags
const suggestedTags = ref([])
const selectedTags = ref([])
const showFullTagChooser = ref(false)
const tagReasoning = ref('')

// Events
const events = ref([])
const recurrenceReasoning = ref('')

// Exploration metadata
const pagesExplored = ref(0)
const explorationConfidence = ref(0)
const explorationLog = ref([])
const showExplorationLog = ref(false)

// Image generation
const generatedImage = ref(null)
const isGeneratingImage = ref(false)
const imageError = ref('')
const useImage = ref(true)
const imagePromptUsed = ref('')
const imageVariation = ref(0)

// Recurrence type options
const recurrenceOptions = [
  { value: 'none', label: 'Keine Wiederholung' },
  { value: 'weekly', label: 'Wöchentlich' },
  { value: 'monthly-date', label: 'Monatlich (Datum)' },
  { value: 'monthly-day', label: 'Monatlich (Wochentag)' }
]

const confidenceLevel = computed(() => {
  if (explorationConfidence.value >= 0.7) return 'high'
  if (explorationConfidence.value >= 0.4) return 'medium'
  return 'low'
})

const confidencePercent = computed(() => {
  return Math.round(explorationConfidence.value * 100)
})

const analyze = async () => {
  if (importMode.value === 'url' && !url.value) {
    error.value = 'Bitte eine URL eingeben'
    return
  }
  if (importMode.value === 'text' && !textInput.value.trim()) {
    error.value = 'Bitte eine Beschreibung eingeben'
    return
  }

  isLoading.value = true
  error.value = ''
  hasExtracted.value = false

  try {
    const result = importMode.value === 'url'
      ? await extractionAPI.extract(url.value)
      : await extractionAPI.extractFromText(textInput.value)

    // Populate project fields
    title.value = result.title || ''
    keywords.value = result.keywords || []

    // Set homepage directly
    homepage.value = result.homepage || ''

    // Build description with metadata appended
    let fullDescription = result.description || ''
    const metadata = []
    if (result.location) metadata.push(`**Ort:** ${result.location}`)
    if (result.contact_email) metadata.push(`**Kontakt:** ${result.contact_email}`)
    if (result.keywords && result.keywords.length > 0) {
      const keywordPairs = []
      for (let i = 0; i < result.keywords.length; i += 2) {
        if (i + 1 < result.keywords.length) {
          keywordPairs.push(`${result.keywords[i]}, ${result.keywords[i + 1]}`)
        } else {
          keywordPairs.push(result.keywords[i])
        }
      }
      metadata.push(`**Schlagwörter:**\n${keywordPairs.join('\n')}`)
    }

    if (metadata.length > 0) {
      fullDescription += '\n\n' + metadata.join('\n\n')
    }
    description.value = fullDescription

    // Tags
    suggestedTags.value = result.suggested_tags || []
    selectedTags.value = suggestedTags.value
      .filter(t => t.score >= 0.5)
      .map(t => t.tag_id)
    tagReasoning.value = result.tag_reasoning || ''

    // Events - convert ISO strings to Date objects for VueDatePicker
    events.value = (result.events || []).map(e => ({
      ...e,
      enabled: true,
      dates: e.start_date ? [
        new Date(e.start_date),
        e.end_date ? new Date(e.end_date) : null
      ] : null
    }))
    recurrenceReasoning.value = result.recurrence_reasoning || ''

    // Metadata
    pagesExplored.value = result.pages_explored || 0
    explorationConfidence.value = result.exploration_confidence || 0
    explorationLog.value = result.exploration_log || []

    hasExtracted.value = true

    // Start image generation in parallel (don't await)
    generateProjectImage()
  } catch (e) {
    console.error('Extraction failed:', e)
    error.value = e.response?.data?.error || e.message || 'Extraktion fehlgeschlagen'
  } finally {
    isLoading.value = false
  }
}

const generateProjectImage = async (isRegenerate = false) => {
  isGeneratingImage.value = true
  imageError.value = ''
  generatedImage.value = null
  useImage.value = true

  if (isRegenerate) {
    imageVariation.value++
  }

  try {
    const result = await extractionAPI.generateImage(
      title.value,
      description.value,
      keywords.value,
      imageVariation.value
    )
    generatedImage.value = `data:image/png;base64,${result.image_base64}`
    imagePromptUsed.value = result.prompt_used || ''
  } catch (e) {
    console.error('Image generation failed:', e)
    imageError.value = e.response?.data?.error || e.message || 'Bildgenerierung fehlgeschlagen'
  } finally {
    isGeneratingImage.value = false
  }
}

const toggleTag = (tagId) => {
  const index = selectedTags.value.indexOf(tagId)
  if (index >= 0) {
    selectedTags.value.splice(index, 1)
  } else {
    selectedTags.value.push(tagId)
  }
}

const isTagSelected = (tagId) => {
  return selectedTags.value.includes(tagId)
}

const onTagChooserUpdate = (tags) => {
  selectedTags.value = tags
}

const addEvent = () => {
  events.value.push({
    name: '',
    dates: null,
    recurrence_type: 'none',
    recurrence_day: null,
    recurrence_week: null,
    description: '',
    enabled: true
  })
}

const removeEvent = (index) => {
  events.value.splice(index, 1)
}

const formatDateRange = (dates) => {
  if (!dates || dates.length === 0) return ''
  const startDate = new Date(dates[0])
  const endDate = dates[1] ? new Date(dates[1]) : null

  const formatDate = (d) => {
    return d.getDate() + '.' + (d.getMonth() + 1) + '.' + d.getFullYear() + ' ' +
      d.getHours().toString().padStart(2, '0') + ':' + d.getMinutes().toString().padStart(2, '0')
  }

  if (endDate) {
    return formatDate(startDate) + ' - ' + formatDate(endDate)
  }
  return formatDate(startDate)
}

const save = async () => {
  if (!title.value) {
    error.value = 'Titel ist erforderlich'
    return
  }

  isLoading.value = true
  error.value = ''

  try {
    // 1. Create project
    const projectData = {
      title: title.value,
      description: description.value,
      homepage: homepage.value
    }
    const project = await projectAPI.create(projectData)

    // 2. Set tags
    if (selectedTags.value.length > 0) {
      await projectAPI.setProjectTags(project.id, selectedTags.value)
    }

    // 3. Upload image if generated and user wants to use it
    if (generatedImage.value && useImage.value) {
      try {
        const response = await fetch(generatedImage.value)
        const blob = await response.blob()

        const formData = new FormData()
        formData.append('image', blob, 'project-image.png')

        await projectAPI.uploadImage(project.id, formData)
      } catch (imgErr) {
        console.error('Image upload failed:', imgErr)
      }
    }

    // 4. Create events
    const eventAPI = projectAPI.eventAPIFor(project.id)
    for (const event of events.value) {
      if (!event.enabled || !event.dates || !event.dates[0]) continue

      const startDate = event.dates[0] instanceof Date
        ? event.dates[0].toISOString()
        : event.dates[0]
      const endDate = event.dates[1]
        ? (event.dates[1] instanceof Date ? event.dates[1].toISOString() : event.dates[1])
        : startDate

      await eventAPI.create({
        name: event.name || title.value,
        start_date: startDate,
        end_date: endDate,
        recurrence_type: event.recurrence_type || 'none',
        description: event.description || ''
      })
    }

    router.push(`/projects/${project.id}`)
  } catch (e) {
    console.error('Save failed:', e)
    error.value = e.response?.data?.error || e.message || 'Speichern fehlgeschlagen'
    isLoading.value = false
  }
}
</script>

<template>
  <div>
    <!-- Analysis overlay -->
    <ScanOverlay v-if="isLoading" :url="importMode === 'url' ? url : ''" />

    <h1 class="text-2xl font-display font-bold text-wald-300 mb-6">{{ importMode === 'url' ? 'Von Website importieren' : 'Aus Textbeschreibung' }}</h1>

    <div class="section-panel mb-6">
      <!-- URL input -->
      <div v-if="importMode === 'url'">
        <label for="url">Website-URL</label>
        <div class="flex gap-2">
          <input
            type="url"
            id="url"
            v-model="url"
            placeholder="https://example.com"
            :disabled="isLoading"
            @keyup.enter="analyze"
            class="flex-1"
          >
          <button
            class="btn-cyber whitespace-nowrap"
            type="button"
            @click="analyze"
            :disabled="isLoading || !url"
          >
            <svg v-if="isLoading" class="animate-spin -ml-1 mr-2 h-4 w-4 inline-block" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24"><circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle><path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path></svg>
            {{ isLoading ? 'Analysiere...' : 'Analysieren' }}
          </button>
        </div>
      </div>

      <!-- Text input -->
      <div v-else>
        <label for="textInput">Projektbeschreibung</label>
        <textarea
          id="textInput"
          rows="6"
          v-model="textInput"
          placeholder="Beschreibe das Projekt, seine Aktivitäten, Termine, und andere relevante Informationen..."
          :disabled="isLoading"
        ></textarea>
        <button
          class="btn-cyber mt-3"
          type="button"
          @click="analyze"
          :disabled="isLoading || !textInput.trim()"
        >
          <svg v-if="isLoading" class="animate-spin -ml-1 mr-2 h-4 w-4 inline-block" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24"><circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle><path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path></svg>
          {{ isLoading ? 'Analysiere...' : 'Analysieren' }}
        </button>
      </div>

      <div v-if="error" class="alert-danger mt-4">
        {{ error }}
      </div>
    </div>

    <!-- Extraction Results -->
    <div v-if="hasExtracted" class="space-y-6">
      <!-- Confidence Indicator (only for URL mode) -->
      <div v-if="importMode === 'url'" class="section-panel">
        <div class="flex justify-between items-center mb-2">
          <span class="text-sm text-gray-400 font-mono">Extraktions-Konfidenz</span>
          <span class="text-sm font-mono text-wald-300">{{ confidencePercent }}%</span>
        </div>
        <div class="progress-cyber">
          <div
            class="progress-cyber-bar"
            :class="confidenceLevel"
            :style="{ width: confidencePercent + '%' }"
          ></div>
        </div>
        <div class="mt-2 text-gray-500 text-xs font-mono">
          {{ pagesExplored }} Seite(n) analysiert
          <a href="#" @click.prevent="showExplorationLog = !showExplorationLog" class="ml-2">
            {{ showExplorationLog ? 'Log ausblenden' : 'Log anzeigen' }}
          </a>
        </div>
        <div v-if="showExplorationLog" class="mt-2">
          <pre class="bg-black/80 border border-wald-500/10 p-3 text-xs text-gray-500 font-mono rounded max-h-48 overflow-y-auto">{{ explorationLog.join('\n') }}</pre>
        </div>
      </div>

      <!-- Project Section -->
      <div class="section-panel">
        <h3 class="text-sm font-mono uppercase tracking-wider text-wald-400 mb-4">Projekt-Informationen</h3>
        <div class="space-y-4">
          <div>
            <label for="title">Titel</label>
            <input type="text" id="title" v-model="title">
          </div>
          <div>
            <label for="description">Beschreibung</label>
            <textarea id="description" rows="8" v-model="description"></textarea>
          </div>
          <div>
            <label for="homepage">Homepage</label>
            <input type="url" id="homepage" v-model="homepage" placeholder="https://...">
          </div>
          <div v-if="keywords.length > 0">
            <label class="text-gray-600">Extrahierte Schlüsselwörter</label>
            <div class="flex flex-wrap gap-1.5">
              <span v-for="keyword in keywords" :key="keyword" class="badge-cyber">
                {{ keyword }}
              </span>
            </div>
          </div>
        </div>
      </div>

      <!-- Image Preview -->
      <div class="section-panel">
        <div class="flex justify-between items-center mb-4">
          <h3 class="text-sm font-mono uppercase tracking-wider text-wald-400">Projektbild</h3>
          <div v-if="generatedImage && useImage" class="flex gap-2">
            <button
              class="btn-cyber-outline btn-cyber-sm"
              @click="generateProjectImage(true)"
              :disabled="isGeneratingImage"
            >
              Neu generieren
            </button>
            <button
              class="btn-cyber-danger btn-cyber-sm"
              @click="useImage = false"
            >
              Kein Bild
            </button>
          </div>
        </div>
        <div class="text-center">
          <div v-if="isGeneratingImage" class="py-8">
            <svg class="animate-spin h-8 w-8 text-wald-400 mx-auto" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24"><circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle><path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path></svg>
            <p class="mt-3 text-gray-500 text-sm font-mono">Bild wird generiert...</p>
          </div>

          <div v-else-if="generatedImage && useImage">
            <img :src="generatedImage" class="max-w-[256px] rounded border border-wald-500/20 mx-auto">
            <p v-if="imagePromptUsed" class="mt-2 text-gray-600 text-xs">{{ imagePromptUsed }}</p>
          </div>

          <div v-else-if="generatedImage && !useImage" class="py-6">
            <p class="text-gray-500 text-sm">Kein Bild wird verwendet</p>
            <button class="btn-cyber-outline btn-cyber-sm mt-2" @click="useImage = true">
              Generiertes Bild doch verwenden
            </button>
          </div>

          <div v-else-if="imageError" class="py-6 text-red-400 text-sm">
            {{ imageError }}
            <button class="text-wald-400 hover:text-wald-300 ml-2 underline" @click="generateProjectImage(true)">Erneut versuchen</button>
          </div>

          <div v-else class="py-6">
            <p class="text-gray-500 text-sm font-mono">Bildgenerierung wird gestartet...</p>
          </div>
        </div>
      </div>

      <!-- Tags Section -->
      <div class="section-panel">
        <div class="flex justify-between items-center mb-4">
          <h3 class="text-sm font-mono uppercase tracking-wider text-wald-400">Tags</h3>
          <button
            class="btn-cyber-outline btn-cyber-sm"
            @click="showFullTagChooser = !showFullTagChooser"
          >
            {{ showFullTagChooser ? 'Vorschläge anzeigen' : 'Alle Tags anzeigen' }}
          </button>
        </div>
        <div v-if="!showFullTagChooser">
          <div v-if="tagReasoning" class="text-gray-500 text-xs mb-3">
            <span class="text-wald-400 font-mono">KI-Begründung:</span> {{ tagReasoning }}
          </div>
          <div v-if="suggestedTags.length > 0">
            <div class="flex flex-wrap gap-2">
              <button
                v-for="tag in suggestedTags"
                :key="tag.tag_id"
                class="px-3 py-1 text-xs font-mono uppercase tracking-wider rounded-full border transition-all"
                :class="isTagSelected(tag.tag_id)
                  ? 'border-wald-500 bg-wald-500/20 text-wald-300'
                  : 'border-gray-600 text-gray-400 hover:border-gray-400'"
                @click="toggleTag(tag.tag_id)"
              >
                {{ tag.tag_name }}
                <span class="ml-1 opacity-60">{{ Math.round(tag.score * 100) }}%</span>
              </button>
            </div>
          </div>
          <div v-else class="text-gray-500 text-sm">
            Keine Tag-Vorschläge. Klicke "Alle Tags anzeigen" um manuell auszuwählen.
          </div>
        </div>
        <div v-else>
          <TagChooser :selected="selectedTags" @update="onTagChooserUpdate" />
        </div>
      </div>

      <!-- Events Section -->
      <div class="section-panel">
        <div class="flex justify-between items-center mb-4">
          <h3 class="text-sm font-mono uppercase tracking-wider text-wald-400">Termine</h3>
          <button class="btn-cyber-outline btn-cyber-sm" @click="addEvent">
            + Termin hinzufügen
          </button>
        </div>
        <div v-if="recurrenceReasoning" class="text-gray-500 text-xs mb-3">
          <span class="text-wald-400 font-mono">Wiederholungserkennung:</span> {{ recurrenceReasoning }}
        </div>

        <div v-if="events.length === 0" class="text-gray-500 text-sm">
          Keine Termine erkannt. Klicke "Termin hinzufügen" um manuell einen zu erstellen.
        </div>

        <div v-for="(event, index) in events" :key="index" class="section-panel !p-4 mb-3">
          <div class="flex gap-3">
            <div class="pt-1">
              <input
                type="checkbox"
                :id="'event-enabled-' + index"
                v-model="event.enabled"
              >
            </div>
            <div class="flex-1 space-y-3">
              <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
                <div>
                  <label class="!text-[10px]">Name</label>
                  <input
                    type="text"
                    v-model="event.name"
                    :placeholder="title"
                    :disabled="!event.enabled"
                  >
                </div>
                <div>
                  <label class="!text-[10px]">Wiederholung</label>
                  <select
                    v-model="event.recurrence_type"
                    :disabled="!event.enabled"
                  >
                    <option v-for="opt in recurrenceOptions" :key="opt.value" :value="opt.value">
                      {{ opt.label }}
                    </option>
                  </select>
                </div>
              </div>
              <div class="grid grid-cols-1 md:grid-cols-3 gap-3">
                <div class="md:col-span-2">
                  <label class="!text-[10px]">Zeitraum</label>
                  <VueDatePicker
                    v-model="event.dates"
                    locale="de"
                    cancelText="abbrechen"
                    selectText="auswählen"
                    :format="formatDateRange"
                    range
                    :disabled="!event.enabled"
                    dark
                  />
                </div>
                <div>
                  <label class="!text-[10px]">Beschreibung</label>
                  <input
                    type="text"
                    v-model="event.description"
                    :disabled="!event.enabled"
                  >
                </div>
              </div>
            </div>
            <div>
              <button
                class="btn-cyber-danger btn-cyber-sm"
                @click="removeEvent(index)"
                title="Entfernen"
              >
                &times;
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- Save Button -->
      <div class="flex gap-3 pb-8">
        <button
          class="btn-cyber"
          @click="save"
          :disabled="isLoading || !title"
        >
          <svg v-if="isLoading" class="animate-spin -ml-1 mr-2 h-4 w-4 inline-block" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24"><circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle><path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path></svg>
          Projekt erstellen
        </button>
        <button class="btn-cyber-secondary" @click="router.push('/projects')">
          Abbrechen
        </button>
      </div>
    </div>

    <!-- Back button when no extraction -->
    <div v-else-if="!isLoading" class="pb-8">
      <RouterLink to="/projects" class="btn-cyber-secondary inline-block">Zurück zur Liste</RouterLink>
    </div>
  </div>
</template>
