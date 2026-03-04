<script setup>
  import { ref } from 'vue'
  import { useRouter } from 'vue-router'
  import useProjectAPI from '../../api/Project';

  const router = useRouter()

  const title = ref('')
  const description = ref('')
  const homepage = ref('')

  const projectAPI = useProjectAPI();

  const createProject = async () => {
    const response = await projectAPI.create({ title: title.value, description: description.value, homepage: homepage.value });
    router.push(`/projects/${response.id}`);
  }

</script>

<template>
<h1 class="text-2xl font-display font-bold text-wald-300 mb-6">Neues Projekt</h1>
<div class="section-panel">
  <form class="space-y-4">
    <div>
      <label for="title">Titel</label>
      <input type="text" id="title" v-model="title">
    </div>
    <div>
      <label for="homepage">Homepage</label>
      <input type="url" id="homepage" v-model="homepage" placeholder="https://example.com">
    </div>
    <div>
      <label for="description">Beschreibung</label>
      <textarea id="description" rows="5" v-model="description"></textarea>
    </div>
    <button type="submit" class="btn-cyber" @click.prevent="createProject">Speichern</button>
  </form>
</div>
</template>
