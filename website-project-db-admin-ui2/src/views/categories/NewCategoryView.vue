<script setup>
  import { ref } from 'vue'
  import { useRouter } from 'vue-router'
  import useCategoryAPI from '../../api/Category';

  const router = useRouter()

  const title = ref('')
  const description = ref('')

  const categoryAPI = useCategoryAPI();

  const createCategory = async () => {
    const response = await categoryAPI.create({ title: title.value, description: description.value });
    // TODO: Error Handling
    router.push(`/categories/${response.id}`);
  }

</script>

<template>
<h1>Neue Kategorie</h1>
<form>
  <div class="mb-3">
    <label for="title" class="form-label">Titel</label>
    <input type="text" class="form-control" id="title" v-model="title">
  </div>
  <div class="mb-3">
    <label for="description" class="form-label">Beschreibung</label>
    <textarea class="form-control" id="description" rows="5" v-model="description"></textarea>
  </div>
  <button type="submit" class="btn btn-primary" @click.prevent="createCategory">Speichern</button>
</form>
</template>

