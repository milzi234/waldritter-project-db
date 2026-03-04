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
    router.push(`/categories/${categoryID}`);
  }

</script>

<template>
<h1 class="text-2xl font-display font-bold text-wald-300 mb-6">Tag Bearbeiten</h1>
<div class="section-panel">
  <form class="space-y-4">
    <div>
      <label for="title">Titel</label>
      <input type="text" id="title" v-model="title">
    </div>
    <div>
      <label for="description">Beschreibung</label>
      <textarea id="description" rows="5" v-model="description"></textarea>
    </div>
    <button type="submit" class="btn-cyber" @click.prevent="updateTag">Speichern</button>
  </form>
</div>
</template>
