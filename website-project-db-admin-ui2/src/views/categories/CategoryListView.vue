<script setup>
import Category from '../../components/Category.vue'
import {useCategoryStore} from '../../stores/categories';

// Declare categories as a ref
const categoryStore = useCategoryStore();

</script>

<template>
  <div class="container">
    <h1>Kategorien</h1>
    <RouterLink :to="{name: 'new-category'}">Neue Kategorie</RouterLink>
    <div v-for="category in categoryStore.categories" :key="category.id" style="margin-top:2rem">
      <Category :title="category.title" :description="category.description" :id="category.id" />
    </div>
  </div>
</template>

