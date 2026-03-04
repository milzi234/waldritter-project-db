<script setup>
import Category from '../../components/Category.vue'
import {useCategoryStore} from '../../stores/categories';

const categoryStore = useCategoryStore();

</script>

<template>
  <div>
    <div class="flex items-center justify-between mb-6">
      <h1 class="text-2xl font-display font-bold text-wald-300">Kategorien</h1>
      <RouterLink :to="{name: 'new-category'}" class="btn-cyber btn-cyber-sm">Neue Kategorie</RouterLink>
    </div>
    <div class="space-y-6">
      <div v-for="category in categoryStore.categories" :key="category.id">
        <Category :title="category.title" :description="category.description" :id="category.id" />
      </div>
    </div>
  </div>
</template>
