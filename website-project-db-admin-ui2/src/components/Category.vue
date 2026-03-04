<script setup>
  import TagPill from '../components/TagPill.vue'
  import useCategoryAPI from '../api/Category';
  import { ref, watchEffect } from 'vue'
  import { useRouter } from 'vue-router'

  const props = defineProps({
    id: Number,
    title: String,
    description: String
  })

  const tags = ref([]);

  const router = useRouter()
  const categoryAPI = useCategoryAPI();

  watchEffect(async () => {
    if (!props.id) {
      return;
    }
    const tagAPI = categoryAPI.tagAPIFor(props.id);
    tags.value = await tagAPI.getAll();
  });

  const deleteCategory = async () => {
    if (!confirm('Wirklich löschen? Alle Tags in dieser Kategorie werden auch gelöscht!')) {
      return;
    }
    await categoryAPI.delete(props.id);
    router.push({name: 'categories'})
  }

</script>
<template>
  <div v-if="id" class="section-panel">
    <h2 class="text-lg font-display font-bold text-wald-300 mb-1">{{ title }}</h2>
    <p class="text-gray-500 text-sm mb-3">{{ description }}</p>
    <div class="flex items-center gap-4 mb-4">
      <RouterLink :to="{ name: 'edit-category', params: {id: id}}" class="nav-link !px-0">Bearbeiten</RouterLink>
      <RouterLink :to="{ name: 'new-tag', params: {categoryID: id}}" class="nav-link !px-0">Neuer Tag</RouterLink>
      <a href="#" @click.stop.prevent="deleteCategory" class="nav-link !px-0 !text-red-500 hover:!text-red-400">Löschen</a>
    </div>
    <div class="flex flex-wrap gap-2">
      <TagPill v-for="tag in tags" :key="tag.id" :title="tag.title" :id="tag.id" :categoryID="id"/>
    </div>
  </div>
</template>
