<script setup>
defineProps({
  links: {
    type: Array,
    required: true
  }
});
</script>

<template>
  <div class="flex flex-col space-y-2">
    <a v-for="link in links" 
       :key="link.label" 
       :href="link.url"
       class="text-[10px] text-green-200 hover:text-white transition-colors duration-200 text-left">
      {{ link.label }}
    </a>
  </div>
</template>
