<script setup>
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import VueDatePicker from '@vuepic/vue-datepicker';
import '@vuepic/vue-datepicker/dist/main.css';

import useProjectAPI from '../../api/Project';
import useExtractionAPI from '../../api/Extraction';
import TagChooser from '../../components/TagChooser.vue';

const router = useRouter()
const projectAPI = useProjectAPI();
const extractionAPI = useExtractionAPI();

// URL input
const url = ref('')
const isLoading = ref(false)
const error = ref('')
const hasExtracted = ref(false)

// Extracted data
const title = ref('')
const description = ref('')
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

const confidenceClass = computed(() => {
  if (explorationConfidence.value >= 0.7) return 'bg-success'
  if (explorationConfidence.value >= 0.4) return 'bg-warning'
  return 'bg-danger'
})

const confidencePercent = computed(() => {
  return Math.round(explorationConfidence.value * 100)
})

const analyze = async () => {
  if (!url.value) {
    error.value = 'Bitte eine URL eingeben'
    return
  }

  isLoading.value = true
  error.value = ''
  hasExtracted.value = false

  try {
    const result = await extractionAPI.extract(url.value)

    // Populate project fields
    title.value = result.title || ''
    keywords.value = result.keywords || []

    // Build description with metadata appended
    let fullDescription = result.description || ''
    const metadata = []
    if (result.homepage) metadata.push(`**Homepage:** ${result.homepage}`)
    if (result.location) metadata.push(`**Ort:** ${result.location}`)
    if (result.contact_email) metadata.push(`**Kontakt:** ${result.contact_email}`)
    if (result.keywords && result.keywords.length > 0) {
      // Format keywords with each pair on a new line for better markdown rendering
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

  // Increment variation on regenerate to get a different style
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
      description: description.value
    }
    const project = await projectAPI.create(projectData)

    // 2. Set tags
    if (selectedTags.value.length > 0) {
      await projectAPI.setProjectTags(project.id, selectedTags.value)
    }

    // 3. Upload image if generated and user wants to use it
    if (generatedImage.value && useImage.value) {
      try {
        // Convert base64 data URL to blob
        const response = await fetch(generatedImage.value)
        const blob = await response.blob()

        // Create FormData with the image
        const formData = new FormData()
        formData.append('image', blob, 'project-image.png')

        await projectAPI.uploadImage(project.id, formData)
      } catch (imgErr) {
        console.error('Image upload failed:', imgErr)
        // Continue anyway - project is created
      }
    }

    // 4. Create events - convert Date objects to ISO strings for API
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

    // Redirect to project view
    router.push(`/projects/${project.id}`)
  } catch (e) {
    console.error('Save failed:', e)
    error.value = e.response?.data?.error || e.message || 'Speichern fehlgeschlagen'
    isLoading.value = false
  }
}
</script>

<template>
  <div class="container">
    <h1>Projekt importieren</h1>

    <!-- URL Input Section -->
    <div class="card mb-4">
      <div class="card-body">
        <div class="mb-3">
          <label for="url" class="form-label">Website-URL</label>
          <div class="input-group">
            <input
              type="url"
              class="form-control"
              id="url"
              v-model="url"
              placeholder="https://example.com"
              :disabled="isLoading"
              @keyup.enter="analyze"
            >
            <button
              class="btn btn-primary"
              type="button"
              @click="analyze"
              :disabled="isLoading || !url"
            >
              <span v-if="isLoading" class="spinner-border spinner-border-sm me-2" role="status"></span>
              {{ isLoading ? 'Analysiere...' : 'Analysieren' }}
            </button>
          </div>
        </div>

        <div v-if="error" class="alert alert-danger" role="alert">
          {{ error }}
        </div>
      </div>
    </div>

    <!-- Extraction Results -->
    <div v-if="hasExtracted">
      <!-- Confidence Indicator -->
      <div class="card mb-4">
        <div class="card-body">
          <div class="d-flex justify-content-between align-items-center mb-2">
            <span>Extraktions-Konfidenz</span>
            <span>{{ confidencePercent }}%</span>
          </div>
          <div class="progress" style="height: 10px;">
            <div
              class="progress-bar"
              :class="confidenceClass"
              role="progressbar"
              :style="{ width: confidencePercent + '%' }"
            ></div>
          </div>
          <div class="mt-2 text-muted small">
            {{ pagesExplored }} Seite(n) analysiert
            <a href="#" @click.prevent="showExplorationLog = !showExplorationLog" class="ms-2">
              {{ showExplorationLog ? 'Log ausblenden' : 'Log anzeigen' }}
            </a>
          </div>
          <div v-if="showExplorationLog" class="mt-2">
            <pre class="bg-light p-2 small" style="max-height: 200px; overflow-y: auto;">{{ explorationLog.join('\n') }}</pre>
          </div>
        </div>
      </div>

      <!-- Project Section -->
      <div class="card mb-4">
        <div class="card-header">
          <h5 class="mb-0">Projekt-Informationen</h5>
        </div>
        <div class="card-body">
          <div class="mb-3">
            <label for="title" class="form-label">Titel</label>
            <input type="text" class="form-control" id="title" v-model="title">
          </div>
          <div class="mb-3">
            <label for="description" class="form-label">Beschreibung</label>
            <textarea class="form-control" id="description" rows="8" v-model="description"></textarea>
          </div>
          <div v-if="keywords.length > 0" class="mb-3">
            <label class="form-label text-muted">Extrahierte Schlüsselwörter</label>
            <div>
              <span v-for="keyword in keywords" :key="keyword" class="badge bg-secondary me-1">
                {{ keyword }}
              </span>
            </div>
          </div>
        </div>
      </div>

      <!-- Image Preview Card -->
      <div class="card mb-4">
        <div class="card-header d-flex justify-content-between align-items-center">
          <h5 class="mb-0">Projektbild</h5>
          <div v-if="generatedImage && useImage" class="btn-group btn-group-sm">
            <button
              class="btn btn-outline-secondary"
              @click="generateProjectImage(true)"
              :disabled="isGeneratingImage"
            >
              Neu generieren
            </button>
            <button
              class="btn btn-outline-danger"
              @click="useImage = false"
            >
              Kein Bild
            </button>
          </div>
        </div>
        <div class="card-body text-center">
          <!-- Generating state -->
          <div v-if="isGeneratingImage" class="py-4">
            <div class="spinner-border text-primary" role="status"></div>
            <p class="mt-2 text-muted">Bild wird generiert...</p>
          </div>

          <!-- Image generated and user wants to use it -->
          <div v-else-if="generatedImage && useImage">
            <img :src="generatedImage" class="img-fluid rounded" style="max-width: 256px;">
            <p v-if="imagePromptUsed" class="mt-2 text-muted small">{{ imagePromptUsed }}</p>
          </div>

          <!-- User chose not to use image -->
          <div v-else-if="generatedImage && !useImage" class="py-3">
            <p class="text-muted">Kein Bild wird verwendet</p>
            <button class="btn btn-sm btn-outline-primary" @click="useImage = true">
              Generiertes Bild doch verwenden
            </button>
          </div>

          <!-- Error state -->
          <div v-else-if="imageError" class="text-danger py-3">
            {{ imageError }}
            <button class="btn btn-sm btn-link" @click="generateProjectImage(true)">Erneut versuchen</button>
          </div>

          <!-- Not yet generated -->
          <div v-else class="py-3">
            <p class="text-muted">Bildgenerierung wird gestartet...</p>
          </div>
        </div>
      </div>

      <!-- Tags Section -->
      <div class="card mb-4">
        <div class="card-header d-flex justify-content-between align-items-center">
          <h5 class="mb-0">Tags</h5>
          <button
            class="btn btn-sm btn-outline-secondary"
            @click="showFullTagChooser = !showFullTagChooser"
          >
            {{ showFullTagChooser ? 'Vorschläge anzeigen' : 'Alle Tags anzeigen' }}
          </button>
        </div>
        <div class="card-body">
          <div v-if="!showFullTagChooser">
            <div v-if="tagReasoning" class="text-muted small mb-3">
              <strong>KI-Begründung:</strong> {{ tagReasoning }}
            </div>
            <div v-if="suggestedTags.length > 0">
              <div class="d-flex flex-wrap gap-2">
                <button
                  v-for="tag in suggestedTags"
                  :key="tag.tag_id"
                  class="btn btn-sm"
                  :class="isTagSelected(tag.tag_id) ? 'btn-primary' : 'btn-outline-secondary'"
                  @click="toggleTag(tag.tag_id)"
                >
                  {{ tag.tag_name }}
                  <span class="badge bg-light text-dark ms-1">{{ Math.round(tag.score * 100) }}%</span>
                </button>
              </div>
            </div>
            <div v-else class="text-muted">
              Keine Tag-Vorschläge. Klicke "Alle Tags anzeigen" um manuell auszuwählen.
            </div>
          </div>
          <div v-else>
            <TagChooser :selected="selectedTags" @update="onTagChooserUpdate" />
          </div>
        </div>
      </div>

      <!-- Events Section -->
      <div class="card mb-4">
        <div class="card-header d-flex justify-content-between align-items-center">
          <h5 class="mb-0">Termine</h5>
          <button class="btn btn-sm btn-outline-primary" @click="addEvent">
            + Termin hinzufügen
          </button>
        </div>
        <div class="card-body">
          <div v-if="recurrenceReasoning" class="text-muted small mb-3">
            <strong>Wiederholungserkennung:</strong> {{ recurrenceReasoning }}
          </div>

          <div v-if="events.length === 0" class="text-muted">
            Keine Termine erkannt. Klicke "Termin hinzufügen" um manuell einen zu erstellen.
          </div>

          <div v-for="(event, index) in events" :key="index" class="card mb-3">
            <div class="card-body">
              <div class="row align-items-start">
                <div class="col-auto">
                  <div class="form-check mt-2">
                    <input
                      class="form-check-input"
                      type="checkbox"
                      :id="'event-enabled-' + index"
                      v-model="event.enabled"
                    >
                  </div>
                </div>
                <div class="col">
                  <div class="row">
                    <div class="col-md-6 mb-2">
                      <label class="form-label small">Name</label>
                      <input
                        type="text"
                        class="form-control form-control-sm"
                        v-model="event.name"
                        :placeholder="title"
                        :disabled="!event.enabled"
                      >
                    </div>
                    <div class="col-md-6 mb-2">
                      <label class="form-label small">Wiederholung</label>
                      <select
                        class="form-select form-select-sm"
                        v-model="event.recurrence_type"
                        :disabled="!event.enabled"
                      >
                        <option v-for="opt in recurrenceOptions" :key="opt.value" :value="opt.value">
                          {{ opt.label }}
                        </option>
                      </select>
                    </div>
                  </div>
                  <div class="row">
                    <div class="col-md-8 mb-2">
                      <label class="form-label small">Zeitraum</label>
                      <VueDatePicker
                        v-model="event.dates"
                        locale="de"
                        cancelText="abbrechen"
                        selectText="auswählen"
                        :format="formatDateRange"
                        range
                        :disabled="!event.enabled"
                      />
                    </div>
                    <div class="col-md-4 mb-2">
                      <label class="form-label small">Beschreibung</label>
                      <input
                        type="text"
                        class="form-control form-control-sm"
                        v-model="event.description"
                        :disabled="!event.enabled"
                      >
                    </div>
                  </div>
                </div>
                <div class="col-auto">
                  <button
                    class="btn btn-sm btn-outline-danger"
                    @click="removeEvent(index)"
                    title="Entfernen"
                  >
                    &times;
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Save Button -->
      <div class="d-flex gap-2 mb-4">
        <button
          class="btn btn-success btn-lg"
          @click="save"
          :disabled="isLoading || !title"
        >
          <span v-if="isLoading" class="spinner-border spinner-border-sm me-2" role="status"></span>
          Projekt erstellen
        </button>
        <button class="btn btn-secondary btn-lg" @click="router.push('/projects')">
          Abbrechen
        </button>
      </div>
    </div>

    <!-- Back button when no extraction -->
    <div v-else-if="!isLoading" class="mb-4">
      <RouterLink to="/projects" class="btn btn-secondary">Zurück zur Liste</RouterLink>
    </div>
  </div>
</template>
