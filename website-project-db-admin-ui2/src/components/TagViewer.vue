<script setup>
  import {useCategoryStore} from '../stores/categories';
  import {computed} from 'vue'

  const props = defineProps({
    selected: Array,
  })

  const categoryStore = useCategoryStore();

  const selected = computed(() => {
    const categories = {}
    categoryStore.categories.forEach(category => {
      categoryStore.tagsFor(category.id).value.sort((a, b) => a.title.localeCompare(b.title)).forEach(tag => {
        if (props.selected.includes(tag.id)) {
          if (categories[category.id] === undefined) {
            categories[category.id] = []
          }
          categories[category.id].push(tag)
        }
      })
    });
    return categories
  });

</script>
<template>
  <div>
    <div v-for="(tags, categoryID) in selected" :key="categoryID" class="mt-4 first:mt-0">
      <span class="text-xs font-mono uppercase tracking-wider text-gray-500">{{ categoryStore.find(Number(categoryID)).title }}:</span>
      <span v-for="tag in tags" :key="tag.id" class="badge-cyber ml-1.5">{{ tag.title }}</span>
    </div>
  </div>
</template>
