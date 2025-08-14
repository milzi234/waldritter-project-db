<script setup>
  import { ref,watchEffect } from 'vue'
  import { useRoute, useRouter } from 'vue-router'
  import useCategoryAPI from '../../api/Category';

  const router = useRouter()
  const route = useRoute()

  const id = route.params.id

  const title = ref('')
  const description = ref('')
//  const projects = ref([])
  
  const categoryAPI = useCategoryAPI();

  watchEffect(async () => {
    const tag = await categoryAPI.tagAPIFor(route.params.categoryID).get(id);
    title.value = tag.title
    description.value = tag.description
  });

  // watchEffect(async () => {
  //   projects.value = await tagAPI.getProjects(id)
  // });

  const deleteTag = async () => {
    if (!confirm('Wirklich löschen? Alle Verbindungen zu Projekten gehen verloren!')) {
      return;
    }
    await categoryAPI.tagAPIFor(route.params.categoryID).delete(id);
    router.push({name: 'view-category', params: {id: route.params.categoryID}})
  }
</script>

<template>
<h1>{{ title }}</h1>
<div>
  <RouterLink :to="{name: 'edit-tag', params: {categoryID: route.params.categoryID, id: id}}">Bearbeiten</RouterLink> &nbsp;
  <a href="#" @click.stop.prevent="deleteTag" class="link-danger">Löschen</a>
</div>
<div>{{ description }}</div>
<h2>Projekte mit diesem Tag</h2>
TODO
<!-- <div style="margin-top:2rem">
  <ul class="list-unstyled">
    <li v-for="project in projects" :key="project.id">
      <ProjectPill :title="project.title" :description="project.description" :id="project.id" />
    </li>
  </ul>
</div> -->
</template>

