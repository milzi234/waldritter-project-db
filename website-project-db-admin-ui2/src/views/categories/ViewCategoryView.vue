<script setup>
import { ref, watchEffect } from 'vue'
import { useRoute } from 'vue-router'

import useCategoryAPI from '../../api/Category';
import Category from '../../components/Category.vue'

const route = useRoute()
const category = ref({});

const categoryAPI = useCategoryAPI();

// Fetch categories asynchronously
watchEffect(async () => {
  category.value = await categoryAPI.get(route.params.id);
});

</script>

<template>
  <div class="container">
    <div>
      <Category :title="category.title" :description="category.description" :id="category.id" />
    </div>
  </div>
</template>

