<script setup>
import {useProjectStore} from '../../stores/projects';
import ProjectPill from '../../components/ProjectPill.vue'

const projectStore = useProjectStore();

</script>

<template>
  <div class="container">
    <h1>Projekte</h1>
    <RouterLink :to="{name: 'new-project'}">Neues Projekt</RouterLink>
    <div v-for="project in projectStore.projects" :key="project.id" style="margin-top:2rem">
      <ProjectPill :title="project.title" :id="project.id" />
    </div>
  </div>
</template>

