<script setup>
import { ref, computed } from 'vue'
import {useProjectStore} from '../../stores/projects';
import ProjectPill from '../../components/ProjectPill.vue'

const projectStore = useProjectStore();

const searchQuery = ref('')
const filteredProjects = computed(() => {
  if (!searchQuery.value.trim()) return projectStore.projects
  const q = searchQuery.value.toLowerCase()
  return projectStore.projects.filter(p => p.title?.toLowerCase().includes(q))
})

</script>

<template>
  <div>
    <div class="flex items-center justify-between mb-6">
      <h1 class="text-2xl font-display font-bold text-wald-300">Projekte</h1>
      <div class="flex items-center gap-3">
        <RouterLink :to="{name: 'new-project'}" class="btn-cyber btn-cyber-sm">Neu</RouterLink>
        <RouterLink :to="{name: 'import-project'}" class="btn-cyber-outline btn-cyber-sm">URL Import</RouterLink>
        <RouterLink :to="{name: 'import-text'}" class="btn-cyber-outline btn-cyber-sm">Text Import</RouterLink>
      </div>
    </div>

    <div class="mb-6">
      <div class="flex items-center gap-2">
        <span class="text-wald-500 font-mono text-xs tracking-wider">FILTER::</span>
        <input
          type="text"
          v-model="searchQuery"
          placeholder="Projekt suchen..."
          class="flex-1"
        >
        <span v-if="searchQuery.trim()" class="badge-cyber">{{ filteredProjects.length }}</span>
      </div>
    </div>

    <div class="space-y-3">
      <div v-for="project in filteredProjects" :key="project.id">
        <ProjectPill :title="project.title" :id="project.id" />
      </div>
    </div>
  </div>
</template>
