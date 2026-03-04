<script setup>
  const props = defineProps({
    categoryID: Number,
    id: Number,
    title: String,
  })

</script>
<template>
  <RouterLink
    :to='{name: "view-project", params: {id: id}}'
    class="block section-panel !p-4 hover:glow-border transition-all duration-300 no-underline"
  >
    <span class="text-gray-200 font-body">{{ title }}</span>
  </RouterLink>
</template>
