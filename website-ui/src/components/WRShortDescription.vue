<script setup>
defineProps({
  title: {
    type: String,
    required: true
  },
  description: {
    type: String,
    required: true
  }
});
</script>

<template>
  <div class="max-w-md">
    <h2 class="text-xl font-heading mb-2">{{ title }}</h2>
    <p class="text-[10px] text-green-300">{{ description }}</p>
  </div>
</template>
