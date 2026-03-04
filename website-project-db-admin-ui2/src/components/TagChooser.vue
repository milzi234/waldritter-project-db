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
    <div v-for="category in categoryStore.categories" :key="category.id" class="mt-6 first:mt-0">
      <h2 class="text-sm font-mono uppercase tracking-wider text-wald-400 mb-2">{{ category.title }}</h2>
      <div class="grid grid-cols-2 gap-x-4 gap-y-1">
        <div v-for="tag in categoryStore.tagsFor(category.id).value" :key="tag.id">
          <label class="flex items-center gap-2 !mb-0 cursor-pointer">
            <input type="checkbox" :value="tag.id" :id="`tag-${tag.id}`" v-model="selected" @change="$emit('update', selected)">
            <span class="text-sm text-gray-300 normal-case tracking-normal">{{ tag.title }}</span>
          </label>
        </div>
      </div>
    </div>
  </form>
</template>
