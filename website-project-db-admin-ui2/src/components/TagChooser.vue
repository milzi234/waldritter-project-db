<script setup>
  import {useCategoryStore} from '../stores/categories';
  import {ref, watchEffect} from 'vue'
  
  defineEmits(['update'])

  const props = defineProps({
    selected: Array,
  })
  
  const categoryStore = useCategoryStore();
  const selected = ref([])
  watchEffect(() => {
    selected.value = props.selected
  })

  const chunk = (arr, size) => {
    return arr.reduce((chunks, el, i) => {
      if (i % size === 0) {
        chunks.push([el])
      } else {
        chunks[chunks.length - 1].push(el)
      }
      return chunks
    }, [])
  }

</script>
<template>
  <form>
    <div v-for="category in categoryStore.categories" :key="category.id" style="margin-top:2rem">
      <div class="row">
        <div class="col">
          <h2 class="category-header">{{ category.title }}</h2>          
        </div>
      </div>
      <div class="row" v-for="(tagChunk, index) in chunk(categoryStore.tagsFor(category.id).value, 2)" :key="`chunk-${index}`">
        <div class="col" v-for="tag in tagChunk">
          <div class="form-check">
            <input class="form-check-input" type="checkbox" :value="tag.id" :id="`tag-${tag.id}`" v-model="selected" @change="$emit('update', selected)">
            <label class="form-check-label" :for="`tag-${ tag.id }`">
              {{ tag.title }}
            </label>
          </div>
        </div>
      </div>
    </div>
  </form>
</template>

<style>
  .category-header {
    font-size: 1.2rem;
  }
</style>
