<script setup>
  import { ref, watchEffect } from 'vue'
  import { useRoute, useRouter } from 'vue-router'
  import useCategoryAPI from '../../api/Category';

  const route = useRoute()
  const router = useRouter()

  const categoryID = route.params.categoryID
  const id = route.params.id
  const title = ref('')
  const description = ref('')

  const categoryAPI = useCategoryAPI();

  watchEffect(async () => {
    const tag = await categoryAPI.tagAPIFor(categoryID).get(id)
    title.value = tag.title
    description.value = tag.description
  })

  const updateTag = async () => {
    await categoryAPI.tagAPIFor(categoryID).update(id, { title: title.value, description: description.value });
    // TODO: Error Handling
    router.push(`/categories/${categoryID}`);
  }

</script>

<template>
<h1>Tag Bearbeiten</h1>
<form>
  <div class="mb-3">
    <label for="title" class="form-label">Titel</label>
    <input type="text" class="form-control" id="title" v-model="title">
  </div>
  <div class="mb-3">
    <label for="description" class="form-label">Beschreibung</label>
    <textarea class="form-control" id="description" rows="5" v-model="description"></textarea>
  </div>
  <button type="submit" class="btn btn-primary" @click.prevent="updateTag">Speichern</button>
</form>
</template>

