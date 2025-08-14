<template>
  <div class="project-list space-y-6">
    <div v-for="project in displayedProjects" :key="project.id" class="project-item bg-white text-black shadow-sm rounded-lg overflow-hidden">
      <div class="flex flex-col md:flex-row">
        <div class="w-full md:w-1/3">
          <img 
            :src="project.image.url" 
            :alt="project.image.caption" 
            class="w-full h-48 md:h-full object-cover"
          />
        </div>
        <div class="w-full md:w-2/3 p-4">
          <h2 class="text-lg text-gray-600 font-heading mb-2">{{ project.title }}</h2>
          <WRTimeRange :startDateTime="project.startDateTime" :endDateTime="project.endDateTime" class="mb-2 text-sm" />
          <p class="text-sm mb-3">{{ project.description }}</p>
          <div class="text-xs mb-2">
            <span v-for="(category, index) in project.categories" :key="category.id" class="mr-2">
              <span class="font-medium">{{ category.title }}:</span> 
              {{ category.tags.map(tag => tag.title).join(', ') }}
              <span v-if="index < project.categories.length - 1">|</span>
            </span>
          </div>
          <div v-if="project.occurrences.length > 0" class="mt-3">
            <a href="#" @click.prevent="toggleOccurrences(project.id)" 
               class="text-green-600 hover:text-green-800 text-sm">
              {{ occurrencesExpanded[project.id] ? 'Ausblenden' : 'Weitere Termine' }}
            </a>
            <div v-if="occurrencesExpanded[project.id]">
              <ul class="list-disc list-inside pl-3 mt-2">
                <li v-for="(occurrence, index) in paginatedOccurrences(project.id)" :key="index" class="text-xs list-none">
                  <WRTimeRange :startDateTime="occurrence.startDateTime" :endDateTime="occurrence.endDateTime" />
                </li>
              </ul>
              <Paginator v-if="project.occurrences.length > itemsPerPage" 
                         :rows="itemsPerPage" 
                         :totalRecords="project.occurrences.length"
                         @page="onPageChange($event, project.id)"
                         class="mt-3" />
            </div>
          </div>
        </div>
      </div>
    </div>
    <Paginator v-if="projects.length > itemsPerPage" 
               :rows="itemsPerPage" 
               :totalRecords="projects.length"
               @page="onMainPageChange"
               class="mt-6" />
  </div>
</template>

<script setup>
import { ref, computed } from 'vue';
import Button from 'primevue/button';
import Paginator from 'primevue/paginator';
import WRTimeRange from './WRTimeRange.vue';

const props = defineProps({
  projects: {
    type: Array,
    required: true
  }
});

const itemsPerPage = 5;
const currentPage = ref(1);
const occurrencesExpanded = ref({});
const occurrencesPages = ref({});

const displayedProjects = computed(() => {
  const start = (currentPage.value - 1) * itemsPerPage;
  const end = start + itemsPerPage;
  return props.projects.slice(start, end);
});

const toggleOccurrences = (projectId) => {
  occurrencesExpanded.value[projectId] = !occurrencesExpanded.value[projectId];
  if (!occurrencesPages.value[projectId]) {
    occurrencesPages.value[projectId] = 1;
  }
};

const paginatedOccurrences = (projectId) => {
  const project = props.projects.find(p => p.id === projectId);
  if (!project) return [];
  
  const page = occurrencesPages.value[projectId] || 1;
  const start = (page - 1) * itemsPerPage;
  const end = start + itemsPerPage;
  
  return project.occurrences.slice(start, end);
};

const onPageChange = (event, projectId) => {
  occurrencesPages.value[projectId] = event.page + 1;
};

const onMainPageChange = (event) => {
  currentPage.value = event.page + 1;
};
</script>
