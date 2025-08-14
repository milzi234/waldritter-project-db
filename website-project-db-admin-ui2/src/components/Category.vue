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
  <div v-if="id">
    <h2>
      {{ title }} <br/>
    </h2>
    <span class="text-muted">{{ description }}</span>
    <div class="functions">
      <RouterLink :to="{ name: 'edit-category', params: {id: id}}">Bearbeiten</RouterLink>&nbsp;
      <RouterLink :to="{ name: 'new-tag', params: {categoryID: id}}">Neuer Tag</RouterLink>&nbsp;
      <a href="#" @click.stop.prevent="deleteCategory" class="link-danger">Löschen</a> 
    </div>
    <div style="margin-top:2rem">
      <ul class="list-inline">
        <li class="list-inline-item" v-for="tag in tags" :key="tag.id">
          <TagPill :title="tag.title" :id="tag.id" :categoryID="id"/>
        </li>
      </ul>
    </div>
  </div>
</template>