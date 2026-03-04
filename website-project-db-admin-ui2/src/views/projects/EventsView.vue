<script setup>
  import VueDatePicker from '@vuepic/vue-datepicker';
  import '@vuepic/vue-datepicker/dist/main.css';

  import { ref, watchEffect } from 'vue'
  import { useRoute } from 'vue-router'
  import useProjectAPI from '../../api/Project';

  const route = useRoute()

  const title = ref('')
  const description = ref('')

  const image = ref('')
  const events = ref([])
  const dates = ref([])
  const recurrenceTypes = ref([])

  const dirty = ref(false)
  const saveDisabled = ref(true)

  const showRecurring = ref({})

  const projectAPI = useProjectAPI();

  watchEffect(async () => {
    events.value = await projectAPI.eventAPIFor(route.params.id).getAll();
    events.value.forEach((event) => {
      dates.value.push([event.start_date, event.end_date])
      recurrenceTypes.value.push({recurrence_type: event.recurrence_type, id: event.id})
    })
  });

  watchEffect(async () => {
    const project = await projectAPI.get(route.params.id);
    title.value = project.title
    description.value = project.description
    image.value = project.image
  });

  const newDate = ref(null)

  const createEvent = () => {
    newDate.value = []
  }

  watchEffect(() => {
    let disabled = false
    let modified = false

    if (newDate.value != null) {
      disabled = newDate.value.length == 0
      modified = true
    } else {
      modified = false
    }
    for (let i = 0; i < events.value.length; i++) {
      let event = events.value[i]
      if (!dates.value[i]) {
        modified = true
        disabled = false
        continue
      } else if (dates.value[i][0] != event.start_date || dates.value[i][1] != event.end_date) {
        modified = true
        disabled = false
      }
      if (recurrenceTypes.value[i].recurrence_type != event.recurrence_type) {
        modified = true
        disabled = false
      }
    }
    dirty.value = modified
    saveDisabled.value = disabled && modified
  })

  const toggleRecurring = (event) => {
    showRecurring.value[event.id] = !showRecurring.value[event.id]
  }

  const recurrenceLabel = (recurrence_type, dates) => {
    if (recurrence_type == "none") {
      return "Keine"
    } else if (recurrence_type == "weekly") {
      return weeklyLabel(dates)
    } else if (recurrence_type == "monthly-date") {
      return monthlyDateLabel(dates)
    } else if (recurrence_type == "monthly-day") {
      return monthlyDayLabel(dates)
    }
  }

  const weeklyLabel = (dates) => {
    const start_date = new Date(dates[0])
    return "jeden " + start_date.toLocaleDateString('de-DE', { weekday: 'long' })
  }

  const monthlyDateLabel = (dates) => {
    const start_date = new Date(dates[0])
    return "am " + start_date.getDate()+ ". des Monats"
  }

  const monthlyDayLabel = (dates) => {
    const start_date = new Date(dates[0])
    const weekDays = ["Sonntag", "Montag", "Dienstag", "Mittwoch", "Donnerstag", "Freitag", "Samstag"];
    let weekDay = weekDays[start_date.getDay()];
    let date = start_date.getDate();
    let nthWeek = Math.floor(date / 7) + 1;
    return `jeden ${nthWeek}. ${weekDay} des Monats`;
  }

  const format = (dates) => {
    if (dates.length == 0) {
      return ""
    }
    let startDate = new Date(dates[0]);
    let endDate = new Date(dates[1]);
    let formattedStartDate = startDate.getDate() + "." + (startDate.getMonth() + 1) + "." + startDate.getFullYear() + " " + startDate.getHours().toString().padStart(2, '0') + ":" + startDate.getMinutes().toString().padStart(2, '0');
    let formattedEndDate = endDate.getDate() + "." + (endDate.getMonth() + 1) + "." + endDate.getFullYear() + " " + endDate.getHours().toString().padStart(2, '0') + ":" + endDate.getMinutes().toString().padStart(2, '0');
    return formattedStartDate + " - " + formattedEndDate;
  }

  const save = async () => {
    let reload = false
    if (newDate.value != null) {
      if (newDate.value.length == 0) {
        return
      }
      await projectAPI.eventAPIFor(route.params.id).create({ start_date: newDate.value[0], end_date: newDate.value[1] });
      newDate.value = null
      dirty.value = false
      reload = true
    } else {
      for (let i = 0; i < events.value.length; i++) {
        let event = events.value[i]
        let update = {}
        let needsUpdate = false
        if (!dates.value[i]) {
          await projectAPI.eventAPIFor(route.params.id).delete(event.id);
          reload = true
        } else if (dates.value[i][0] != event.start_date || dates.value[i][1] != event.end_date) {
          update.start_date = dates.value[i][0]
          update.end_date = dates.value[i][1]
          needsUpdate = true
        }

        if (recurrenceTypes.value[i].recurrence_type != event.recurrence_type) {
          update.recurrence_type = recurrenceTypes.value[i].recurrence_type
          needsUpdate = true
        }
        if (needsUpdate) {
          await projectAPI.eventAPIFor(route.params.id).update(event.id, update);
          reload = true
        }
      }
      dirty.value = false
    }
    if (reload) {
      events.value = await projectAPI.eventAPIFor(route.params.id).getAll();
      dates.value = []
      recurrenceTypes.value = []
      events.value.forEach((event) => {
        dates.value.push([event.start_date, event.end_date])
        recurrenceTypes.value.push({recurrence_type: event.recurrence_type, id: event.id})
      })
    }
  }

  const cancel = () => {
    newDate.value = null
    dirty.value = false
    dates.value = []
    recurrenceTypes.value = []
    events.value.forEach((event) => {
      dates.value.push([event.start_date, event.end_date])
      recurrenceTypes.value.push({recurrence_type: event.recurrence_type, id: event.id})
    })
  }

</script>

<template>
<h1 class="text-2xl font-display font-bold text-wald-300 mb-6">Terminplanung - {{ title }}</h1>
<div class="space-y-6">
  <div v-if="image">
    <img :src="image" alt="Project image" class="w-48 h-48 object-cover rounded border border-wald-500/20">
  </div>

  <p v-if="description" class="text-gray-500 text-sm">{{ description }}</p>

  <div>
    <h2 class="text-lg font-display text-wald-400 mb-4">Termine</h2>

    <div v-for="(event, index) in events" :key="event.id" class="section-panel !p-4 mb-4 space-y-3">
      <div v-if="dates[index]">
        <VueDatePicker v-model="dates[index]" locale="de" cancelText="abbrechen" selectText="auswählen" :format="format" range dark></VueDatePicker>
      </div>
      <div class="flex items-center gap-4">
        <a href="#" @click.prevent="toggleRecurring(event)" class="text-sm">Wiederholungen</a>
        <span class="text-gray-500 text-xs font-mono">({{ recurrenceLabel(recurrenceTypes[index].recurrence_type, dates[index]) }})</span>
        <RouterLink v-if="recurrenceTypes[index].recurrence_type != 'none' && !dirty" :to="{ name: 'edit-occurrences', params: {project_id: route.params.id, id: event.id}}" class="text-sm">Ausnahmen</RouterLink>
      </div>
      <div v-if="showRecurring[event.id]" class="space-y-2 pl-1">
        <label class="flex items-center gap-2 !mb-0 cursor-pointer">
          <input type="radio" :name="'repeat-'+event.id" value="none" v-model="recurrenceTypes[index].recurrence_type">
          <span class="text-sm text-gray-300 normal-case tracking-normal">Keine Wiederholung</span>
        </label>
        <label class="flex items-center gap-2 !mb-0 cursor-pointer">
          <input type="radio" :name="'repeat-'+event.id" value="weekly" v-model="recurrenceTypes[index].recurrence_type">
          <span class="text-sm text-gray-300 normal-case tracking-normal">{{ weeklyLabel(dates[index]) }}</span>
        </label>
        <label class="flex items-center gap-2 !mb-0 cursor-pointer">
          <input type="radio" :name="'repeat-'+event.id" value="monthly-date" v-model="recurrenceTypes[index].recurrence_type">
          <span class="text-sm text-gray-300 normal-case tracking-normal">{{ monthlyDateLabel(dates[index]) }}</span>
        </label>
        <label class="flex items-center gap-2 !mb-0 cursor-pointer">
          <input type="radio" :name="'repeat-'+event.id" value="monthly-day" v-model="recurrenceTypes[index].recurrence_type">
          <span class="text-sm text-gray-300 normal-case tracking-normal">{{ monthlyDayLabel(dates[index]) }}</span>
        </label>
      </div>
    </div>

    <div v-if="newDate" class="section-panel !p-4 mb-4">
      <VueDatePicker v-model="newDate" locale="de" cancelText="abbrechen" selectText="auswählen" :format="format" range dark></VueDatePicker>
    </div>

    <div class="flex gap-3 mt-6">
      <button type="button" class="btn-cyber" v-if="dirty" :disabled="saveDisabled" @click.prevent="save">Speichern</button>
      <button type="button" class="btn-cyber-danger" v-if="dirty" :disabled="saveDisabled" @click.prevent="cancel">Abbrechen</button>
      <button type="button" class="btn-cyber" @click.prevent="createEvent" v-if="!dirty">Neuer Termin</button>
    </div>
    <div class="mt-4">
      <button type="button" class="btn-cyber-secondary" @click="$router.go(-1)">Zurück</button>
    </div>
  </div>
</div>
</template>
