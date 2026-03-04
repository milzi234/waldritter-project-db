<script setup>
  import { ref,watchEffect } from 'vue'
  import { useRoute, useRouter } from 'vue-router'
  import useCategoryAPI from '../../api/Category';

  const router = useRouter()
  const route = useRoute()

  const id = route.params.id

  const title = ref('')
  const description = ref('')

  const categoryAPI = useCategoryAPI();

  watchEffect(async () => {
    const tag = await categoryAPI.tagAPIFor(route.params.categoryID).get(id);
    title.value = tag.title
    description.value = tag.description
  });

  const deleteTag = async () => {
    if (!confirm('Wirklich löschen? Alle Verbindungen zu Projekten gehen verloren!')) {
      return;
    }
    await categoryAPI.tagAPIFor(route.params.categoryID).delete(id);
    router.push({name: 'view-category', params: {id: route.params.categoryID}})
  }
</script>

<template>
<div class="section-panel">
  <h1 class="text-xl font-display font-bold text-wald-300 mb-2">{{ title }}</h1>
  <div class="flex items-center gap-4 mb-4">
    <RouterLink :to="{name: 'edit-tag', params: {categoryID: route.params.categoryID, id: id}}" class="nav-link !px-0">Bearbeiten</RouterLink>
    <a href="#" @click.stop.prevent="deleteTag" class="nav-link !px-0 !text-red-500 hover:!text-red-400">Löschen</a>
  </div>
  <p class="text-gray-400 text-sm">{{ description }}</p>
  <h2 class="text-lg font-display text-wald-400 mt-6 mb-2">Projekte mit diesem Tag</h2>
  <p class="text-gray-500 text-sm font-mono">TODO</p>
</div>
</template>
