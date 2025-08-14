<script setup>
import { defineProps } from 'vue';

defineProps({
  socialLinks: {
    type: Array,
    required: true,
    validator: (value) => {
      return value.every(link => 
        typeof link.url === 'string' &&
        typeof link.svgIcon === 'string' &&
        typeof link.label === 'string'
      );
    }
  }
});
</script>

<template>
  <ul class="flex space-x-4" aria-label="Social media links">
    <li v-for="link in socialLinks" :key="link.label">
      <a 
        :href="link.url" 
        target="_blank" 
        rel="noopener noreferrer"
        class="text-white hover:text-green-300 transition-colors duration-200"
        :aria-label="`Visit our ${link.label} page`"
      >
        <svg 
          class="w-6 h-6 text-white"
          aria-hidden="true"
          focusable="false"
          viewBox="0 0 24 24"
        >
          <path :d="link.svgIcon" fill="currentColor" />
        </svg>
        <span class="visually-hidden">{{ link.label }}</span>
      </a>
    </li>
  </ul>
</template>

<style scoped>
.visually-hidden {
  position: absolute;
  height: 1px;
  width: 1px;
  overflow: hidden;
  clip: rect(0, 0, 0, 0);
  white-space: nowrap;
}
</style>
