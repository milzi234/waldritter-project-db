<script setup>
  import { marked } from 'marked'
  import useProjectAPI from '../api/Project';  
  import { useRouter } from 'vue-router'
  import { computed } from 'vue'

  const props = defineProps({
    id: Number,
    title: String,
    description: String,
    image: String,
    date: Array
  })

  const router = useRouter()
  const projectAPI = useProjectAPI();
  const formattedDate = computed(() => {
    if (props.date.length == 0) {
      return ""
    }
    let startDate = new Date(props.date[0]);
    let endDate = new Date(props.date[1]);
    let formattedStartDate = startDate.getDate() + "." + (startDate.getMonth() + 1) + "." + startDate.getFullYear() + " " + startDate.getHours().toString().padStart(2, '0') + ":" + startDate.getMinutes().toString().padStart(2, '0');
    let formattedEndDate = endDate.getDate() + "." + (endDate.getMonth() + 1) + "." + endDate.getFullYear() + " " + endDate.getHours().toString().padStart(2, '0') + ":" + endDate.getMinutes().toString().padStart(2, '0');
    return formattedStartDate + " - " + formattedEndDate;
  })
  const deleteProject = async () => {
    if (!confirm("Wirklich löschen?")) {
      return;
    }
    await projectAPI.delete(props.id);
    router.push({name: 'projects'});
  }
</script>
<template>
  <div class="container">
    <div class="row">
      <div class="col">
        <h2>
          {{ title }} <br/>
        </h2>
      </div>
    </div>
    <div class="row">
      <div class="col">
        <span class="text-muted">{{ formattedDate }}</span>
      </div>
    </div>
    <div class="row">
      <div class="col">
        <img v-if="props.image" :src="props.image" alt="Project image" width="200" height="200" > 
      </div>
    </div>
    <div class="row">
      <div class="col" style="margin-top: 2rem">
        <span class="text-muted" v-html="marked(description)" v-if="description"></span>
      </div>
    </div>
    <div class="row">
      <div class="col">
        <div class="functions" style="margin-top: 2rem">
          <RouterLink :to="{ name: 'edit-project', params: {id: id}}">Bearbeiten</RouterLink>&nbsp;
          <RouterLink :to="{ name: 'tag-project', params: {id: id}}">Tags</RouterLink>&nbsp;
          <RouterLink :to="{ name: 'events', params: {id: id}}">Terminplanung</RouterLink>&nbsp;
          <a href="#" @click.stop.prevent="deleteProject" class="link-danger">Löschen</a>
        </div>
      </div>
    </div>
  </div>
</template>
