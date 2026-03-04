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
  const homepage = ref('')
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
    homepage.value = project.homepage || ''
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
    const response = await projectAPI.update(route.params.id, { title: title.value, description: description.value, homepage: homepage.value });
    const start_date = date.value.length > 0 ? date.value[0] : null
    const end_date = date.value.length > 1 ? date.value[1] : null
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
<h1 class="text-2xl font-display font-bold text-wald-300 mb-6">Projekt Bearbeiten</h1>
<div class="section-panel">
  <form class="space-y-4">
    <div>
      <label for="title">Titel</label>
      <input type="text" id="title" v-model="title">
    </div>
    <div v-if="simple_event">
      <label for="date">Wann?</label>
      <VueDatePicker id="date" v-model="date" locale="de" cancelText="abbrechen" selectText="auswählen" :format="format" range dark></VueDatePicker>
    </div>
    <div>
      <label for="image">Bild</label>
      <input type="file" id="image" @change="fileChanged" accept="image/*">
    </div>
    <div v-if="image">
      <img :src="image" alt="Project image" class="w-48 h-48 object-cover rounded border border-wald-500/20">
    </div>
    <div>
      <label for="homepage">Homepage</label>
      <input type="url" id="homepage" v-model="homepage" placeholder="https://example.com">
    </div>
    <div>
      <label for="description">Beschreibung</label>
      <textarea id="description" rows="5" v-model="description"></textarea>
    </div>
    <div class="flex gap-3">
      <button type="submit" class="btn-cyber" @click.prevent="updateProject">Speichern</button>
      <button type="button" class="btn-cyber-danger" @click.prevent="router.push(`/projects/${route.params.id}`)">Abbrechen</button>
    </div>
  </form>
</div>
</template>
