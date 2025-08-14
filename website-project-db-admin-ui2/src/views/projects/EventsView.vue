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
      disabled = newDate.value.length == 0 // Disable save if no date is selected
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
<h1>Terminplanung - {{ title }}</h1>
<div class="container">
  <div class="row">
    <div class="col">
      <img v-if="image" :src="image" alt="Project image" width="200" height="200" > 
    </div>
  </div>
  <div class="row">
    <div class="col" style="margin-top: 2rem">
      <span class="text-muted">{{ description }}</span>
    </div>
  </div>
  <div class="row">
    <div class="col" style="margin-top: 2rem">
      <h2>Termine</h2>
    </div>
  </div>
  <div v-for="(event, index) in events" :key="event.id">
    <div class="row">
      <div class="col" v-if="dates[index]">
        <VueDatePicker v-model="dates[index]" locale="de" cancelText="abbrechen" selectText="auswählen" :format="format" range></VueDatePicker>
      </div>
    </div>
    <div class="row">
      <div class="col-md-6">
        <a href="#" name="Wiederholungen" @click.prevent="toggleRecurring(event)">Wiederholungen</a> <span class="text-muted">({{ recurrenceLabel(recurrenceTypes[index].recurrence_type, dates[index]) }})</span>
      </div>
      <div class="col-md-6">
        <RouterLink v-if="recurrenceTypes[index].recurrence_type != 'none' && !dirty" :to="{ name: 'edit-occurrences', params: {project_id: route.params.id, id: event.id}}">Ausnahmen</RouterLink>
      </div>
    </div>
    <div class="row" v-if="showRecurring[event.id]">
      <div class="col">
        <input type="radio" id="none" name="repeat" value="none" v-model="recurrenceTypes[index].recurrence_type">
        <label for="none">&nbsp;Keine Wiederholung</label><br>
        <input type="radio" id="weekly" name="repeat" value="weekly" v-model="recurrenceTypes[index].recurrence_type">
        <label for="weekly">&nbsp;{{ weeklyLabel(dates[index]) }}</label><br>
        <input type="radio" id="monthly-date" name="repeat" value="monthly-date" v-model="recurrenceTypes[index].recurrence_type">
        <label for="monthly-date">&nbsp;{{ monthlyDateLabel(dates[index]) }}</label><br>
        <input type="radio" id="monthly-day" name="repeat" value="monthly-day" v-model="recurrenceTypes[index].recurrence_type">
        <label for="monthly-day">&nbsp;{{ monthlyDayLabel(dates[index]) }}</label>
      </div>
    </div>
  </div>
  <div class="row" v-if="newDate">
    <div class="col">
      <VueDatePicker v-model="newDate" locale="de" cancelText="abbrechen" selectText="auswählen" :format="format" range></VueDatePicker>
    </div>
  </div>
  <div class="row">
    <div class="col">
      <div class="functions" style="margin-top: 2rem">
        <button type="button" class="btn btn-success" v-if="dirty" :disabled="saveDisabled"  @click.prevent="save" style="margin-right:1rem">Speichern</button>
        <button type="button" class="btn btn-danger"  v-if="dirty" :disabled="saveDisabled"  @click.prevent="cancel">Abbrechen</button>
        <button type="button" class="btn btn-primary" @click.prevent="createEvent" v-if="!dirty">Neuer Termin</button>
      </div>
      <div class="functions" style="margin-top: 2rem">
        <button type="button" class="btn btn-secondary" @click="$router.go(-1)">Zurück</button>
      </div>
    </div>
  </div>
</div>
</template>

