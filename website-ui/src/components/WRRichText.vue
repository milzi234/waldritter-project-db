<script setup>
import { ref, computed } from 'vue';
import { marked } from 'marked';

const props = defineProps({
  markdown: {
    type: String,
    required: true
  }
});

const htmlContent = computed(() => marked(props.markdown));
</script>

<template>
  <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-16 sm:py-20 lg:py-24">
    <div class="prose prose-lg prose-emerald mx-auto" v-html="htmlContent"></div>
  </div>
</template>