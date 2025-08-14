<script setup>
  const props = defineProps({
    categoryID: Number,
    id: Number,
    title: String,
  })

</script>
<template>
  <span> <RouterLink :to='{name: "view-project", params: {id: id}}'>{{ title }}</RouterLink> </span>
</template>
