<script setup>
  const props = defineProps({
    categoryID: Number,
    id: Number,
    title: String,
  })

</script>
<template>
  <RouterLink
    :to='{name: "view-tag", params: {categoryID: categoryID, id: id}}'
    class="badge-cyber no-underline hover:bg-wald-500/30 transition-colors"
  >{{ title }}</RouterLink>
</template>
