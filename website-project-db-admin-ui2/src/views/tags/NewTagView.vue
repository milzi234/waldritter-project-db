<script setup>
  import { ref } from 'vue'
  import { useRoute, useRouter } from 'vue-router'
  import useCategoryAPI from '../../api/Category';

  const route = useRoute()
  const router = useRouter()

  const title = ref('')
  const description = ref('')

  const categoryAPI = useCategoryAPI();

  const createTag = async () => {
    await categoryAPI.tagAPIFor(route.params.categoryID).create({ title: title.value, description: description.value });
    router.push(`/categories/${route.params.categoryID}`);
  }

</script>

<template>
<h1 class="text-2xl font-display font-bold text-wald-300 mb-6">Neuer Tag</h1>
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
    <button type="submit" class="btn-cyber" @click.prevent="createTag">Speichern</button>
  </form>
</div>
</template>
