<script setup>
import Image from 'primevue/image';

const props = defineProps({
  name: {
    type: String,
    required: true
  },
  email: {
    type: String,
    required: true
  },
  socials: {
    type: Array,
    required: true
  },
  imageSrc: {
    type: String,
    required: true
  },
  imageAlt: {
    type: String,
    default: 'Waldritter Kontakt'
  }
});
</script>

<template>
  <div class="grid grid-cols-1 md:grid-cols-2 gap-8 p-4 md:p-8">
    <div class="flex items-center justify-center">
      <Image :src="imageSrc" :alt="imageAlt" class="rounded-lg shadow-lg" style="max-width: 200px; max-height: 200px; object-fit: cover;" />
    </div>
    <div class="flex flex-col justify-center">
      <h2 class="text-2xl font-bold mb-4">{{ name }}</h2>
      <div class="space-y-4">
        <a :href="`mailto:${email}`" class="text-emerald-700 hover:text-emerald-900 hover:cursor-pointer">{{ email }}</a>
        <div v-if="socials && socials.length > 0">
          <ul class="list-none pl-0 mt-2">
            <li v-for="social in socials" :key="social.name" class="inline-block mr-4">
              <a :href="social.url" target="_blank" rel="noopener noreferrer" class="text-emerald-700 hover:text-emerald-900 hover:cursor-pointer">
                {{ social.name }}
              </a>
            </li>
          </ul>
        </div>
      </div>
    </div>
  </div>
</template>
