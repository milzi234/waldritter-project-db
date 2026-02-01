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

    if (metadata.length > 0) {
      fullDescription += '\n\n' + metadata.join('\n')
    }
    description.value = fullDescription

    // Tags
    suggestedTags.value = result.suggested_tags || []
    selectedTags.value = suggestedTags.value
      .filter(t => t.score >= 0.5)
      .map(t => t.tag_id)
    tagReasoning.value = result.tag_reasoning || ''

    // Events
    events.value = (result.events || []).map(e => ({
      ...e,
      enabled: true,
      dates: e.start_date ? [e.start_date, e.end_date] : null
    }))
    recurrenceReasoning.value = result.recurrence_reasoning || ''

    // Metadata
    pagesExplored.value = result.pages_explored || 0
    explorationConfidence.value = result.exploration_confidence || 0
    explorationLog.value = result.exploration_log || []

    hasExtracted.value = true
  } catch (e) {
    console.error('Extraction failed:', e)
    error.value = e.response?.data?.error || e.message || 'Extraktion fehlgeschlagen'
  } finally {
    isLoading.value = false
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

    // 3. Create events
    const eventAPI = projectAPI.eventAPIFor(project.id)
    for (const event of events.value) {
      if (!event.enabled || !event.dates || !event.dates[0]) continue

      await eventAPI.create({
        name: event.name || title.value,
        start_date: event.dates[0],
        end_date: event.dates[1] || event.dates[0],
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
