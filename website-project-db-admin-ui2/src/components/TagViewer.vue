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
    <div v-for="(tags, categoryID) in selected" :key="categoryID" style="margin-top:2rem">
      <b>{{ categoryStore.find(Number(categoryID)).title }}</b>: <span v-for="tag in tags" :key="tag.id" class="badge bg-primary tag-pill"> {{ tag.title }}</span>     
    </div>
  </div>
</template>

<style>
.tag-pill {
  text-decoration: none;
  color: white;
  margin-right: 0.5rem;
}

.tag-pill:hover {
  color: lightcyan;
}
</style>