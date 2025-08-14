<template>
    <div class="highlight-reel">
        <Carousel :value="highlightedItems" :numVisible="3" :numScroll="1" :responsiveOptions="responsiveOptions">
            <template #item="slotProps">
                <div class="border border-gray-200 rounded-lg shadow-md p-4 m-2">
                    <img :src="slotProps.data.image.url" :alt="slotProps.data.image.caption" class="w-full h-48 object-cover rounded-t-lg mb-4">
                    <h3 class="text-xl font-semibold mb-2">{{ slotProps.data.title }}</h3>
                    <p class="text-gray-600 mb-4">{{ slotProps.data.description }}</p>
                    <WRTimeRange v-if="slotProps.data.startDateTime" :startDateTime="slotProps.data.startDateTime" :endDateTime="slotProps.data.endDateTime" :separateLines="true" class="mb-2 text-sm" />
                    <a href="#" class="text-green-600 hover:underline">Erfahre mehr</a>
                </div>
            </template>
        </Carousel>
    </div>
</template>

<script setup>
import { storeToRefs } from 'pinia';
import { ref, computed, defineProps } from 'vue';
import Carousel from 'primevue/carousel';
import WRTimeRange from './WRTimeRange.vue';
import { useProjectStore } from '../stores/projects';

const props = defineProps({
    highlights: {
        type: Array,
        default: () => []
    }
});

const projectStore = useProjectStore();
const { projects } = storeToRefs(projectStore);

const highlightedItems = computed(() => {
    const highlightedProjects = projects.value.slice(0, 4); // Get the first 4 projects
    return [...props.highlights, ...highlightedProjects];
});

import { watchEffect } from 'vue';

watchEffect(() => {
  console.log('Highlighted items:', highlightedItems.value);
});


const responsiveOptions = ref([
    {
        breakpoint: '1024px',
        numVisible: 2,
        numScroll: 1
    },
    {
        breakpoint: '768px',
        numVisible: 1,
        numScroll: 1
    }
]);
</script>