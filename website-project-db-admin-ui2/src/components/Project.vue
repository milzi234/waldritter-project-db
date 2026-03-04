<script setup>
  import { marked } from 'marked'
  import useProjectAPI from '../api/Project';
  import { useRouter } from 'vue-router'
  import { computed } from 'vue'

  const props = defineProps({
    id: Number,
    title: String,
    description: String,
    homepage: String,
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
  <div class="section-panel glow-border">
    <h2 class="text-xl font-display font-bold text-wald-300 mb-2">{{ title }}</h2>

    <p v-if="formattedDate" class="text-gray-500 text-sm font-mono mb-4">{{ formattedDate }}</p>

    <img v-if="props.image" :src="props.image" alt="Project image" class="w-48 h-48 object-cover rounded mb-4 border border-wald-500/20">

    <div v-if="homepage" class="mb-4">
      <a :href="homepage" target="_blank" rel="noopener noreferrer" class="btn-cyber-outline btn-cyber-sm inline-block">Homepage besuchen</a>
    </div>

    <div v-if="description" class="text-gray-400 text-sm leading-relaxed mb-6 prose-invert" v-html="marked(description)"></div>

    <div class="flex items-center gap-4 pt-4 border-t border-wald-500/10">
      <RouterLink :to="{ name: 'edit-project', params: {id: id}}" class="nav-link !px-0">Bearbeiten</RouterLink>
      <RouterLink :to="{ name: 'tag-project', params: {id: id}}" class="nav-link !px-0">Tags</RouterLink>
      <RouterLink :to="{ name: 'events', params: {id: id}}" class="nav-link !px-0">Terminplanung</RouterLink>
      <a href="#" @click.stop.prevent="deleteProject" class="nav-link !px-0 !text-red-500 hover:!text-red-400">Löschen</a>
    </div>
  </div>
</template>
