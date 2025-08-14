<script setup>
  const props = defineProps({
    categoryID: Number,
    id: Number,
    title: String,
  })

</script>
<template>
  <span> <RouterLink :to='{name: "view-tag", params: {categoryID: categoryID, id: id}}' class="badge bg-primary tag-pill">{{ title }}</RouterLink> </span>
</template>

<style>
.tag-pill {
  text-decoration: none;
  color: white;
}

.tag-pill:hover {
  color: lightcyan;
}
</style>