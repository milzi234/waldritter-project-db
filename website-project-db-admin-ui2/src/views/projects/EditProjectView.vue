<script setup>
  import { ref, watchEffect } from 'vue'
  import { useRoute, useRouter } from 'vue-router'
  import useProjectAPI from '../../api/Project';
  import VueDatePicker from '@vuepic/vue-datepicker';
  import '@vuepic/vue-datepicker/dist/main.css'

  const route = useRoute()
  const router = useRouter()

  const title = ref('')
  const description = ref('')
  const date = ref([])
  const simple_event = ref(true) 
  const new_event = ref(true)

  const event_id = ref('')

  const image = ref('')


  const projectAPI = useProjectAPI();
  var file = null;

  watchEffect(async () => {
    const project = await projectAPI.get(route.params.id);
    title.value = project.title
    description.value = project.description
    image.value = project.image
    file = null
  });

  watchEffect(async () => {
    const events = await projectAPI.eventAPIFor(route.params.id).getAll();
    if (events.length > 1) {
      simple_event.value = false
    } else if (events.length == 1) {
      event_id.value = events[0].id
      simple_event.value = true
      var event = events[0]
      if (event.weekly) {
        simple_event.value = false
      } else {
        date.value = [event.start_date, event.end_date]
        new_event.value = false
      }
    } else {
      simple_event.value = true
      new_event.value = true
      date.value = []
    }
  });

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

  const fileChanged = async (event) => {
    file = event.target.files[0];
  }

  const updateProject = async () => {
    const response = await projectAPI.update(route.params.id, { title: title.value, description: description.value });
    const start_date = date.value.length > 0 ? date.value[0] : null
    const end_date = date.value.length > 1 ? date.value[1] : null
    // TODO: Error Handling
    if (file) {
      const formData = new FormData();
      formData.append('image', file);
      await projectAPI.uploadImage(response.id, formData);
    }
    if (simple_event.value && start_date && end_date) {
      if (new_event.value ) {
        await projectAPI.eventAPIFor(response.id).create({ start_date: start_date, end_date: end_date });
      } else {
        await projectAPI.eventAPIFor(response.id).update(event_id.value, { start_date: start_date, end_date: end_date });
      }
    }
    router.push(`/projects/${response.id}`);
  }

</script>

<template>
<h1>Projekt Bearbeiten</h1>
<form>
  <div class="mb-3">
    <label for="title" class="form-label">Titel</label>
    <input type="text" class="form-control" id="title" v-model="title">
  </div>
  <div class="mb-3" v-if="simple_event">
    <label for="date" class="form-label">Wann?</label>
    <VueDatePicker id="date" v-model="date" locale="de" cancelText="abbrechen" selectText="auswählen" :format="format" range></VueDatePicker>
  </div>
  <div class="mb-3">
    <label for="image" class="form-label">Bild</label>
    <input class="form-control" type="file" id="image" @change="fileChanged" accept="image/*">
  </div>
  <div v-if="image" class="mb-3">
    <img :src="image" id="image" alt="Project image" width="200" height="200" >
  </div>
  <div class="mb-3">
    <label for="description" class="form-label">Beschreibung</label>
    <textarea class="form-control" id="description" rows="5" v-model="description"></textarea>
  </div>
  <button type="submit" class="btn btn-primary" @click.prevent="updateProject">Speichern</button>&nbsp;
  <button type="button" class="btn btn-danger" @click.prevent="router.push(`/projects/${route.params.id}`)">Abbrechen</button>
</form>
</template>

