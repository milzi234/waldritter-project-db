<script setup>
  import { ref, watchEffect } from 'vue'
  import { useRoute, useRouter } from 'vue-router'
  import useCategoryAPI from '../../api/Category';

  const route = useRoute()
  const router = useRouter()

  const title = ref('')
  const description = ref('')

  const categoryAPI = useCategoryAPI();

  watchEffect(async () => {
    const category = await categoryAPI.get(route.params.id);
    title.value = category.title
    description.value = category.value
  });

  const updateCategory = async () => {
    const response = await categoryAPI.update(route.params.id, { title: title.value, description: description.value });
    // TODO: Error Handling
    router.push(`/categories/${response.id}`);
  }

</script>

<template>
<h1>Kategorie Bearbeiten</h1>
<form>
  <div class="mb-3">
    <label for="title" class="form-label">Titel</label>
    <input type="text" class="form-control" id="title" v-model="title">
  </div>
  <div class="mb-3">
    <label for="description" class="form-label">Beschreibung</label>
    <textarea class="form-control" id="description" rows="5" v-model="description"></textarea>
  </div>
  <button type="submit" class="btn btn-primary" @click.prevent="updateCategory">Speichern</button>
</form>
</template>

