<script setup>
import { ref, watchEffect } from 'vue'
import { useRoute } from 'vue-router'

import useProjectAPI from '../../api/Project';
import Project from '../../components/Project.vue'

const route = useRoute()
const project = ref({});
const date = ref([]);

const projectAPI = useProjectAPI();

// Fetch categories asynchronously
watchEffect(async () => {
  project.value = await projectAPI.get(route.params.id);
  projectAPI.getNextOccurrence(route.params.id).then((occurrence) => {
    
    if (occurrence == null) {
      date.value = []
    } else {
      date.value = [occurrence.start_date, occurrence.end_date]
    }
  })
});

</script>

<template>
  <div class="container">
    <div>
      <Project :title="project.title" :description="project.description" :id="project.id" :image="project.image" :date="date" />
    </div>
  </div>
</template>

