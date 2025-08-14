<template>
    <div class="p-6">
        <WRTagFilter :categories="categories" :selectedTags="filter" @update="updateFilter" />
        <WRProjectList :projects="filteredProjects" />
    </div>
</template>

<script setup>
import { storeToRefs } from 'pinia';
import { useProjectStore } from '../stores/projects';
import WRTagFilter from './WRTagFilter.vue';
import WRProjectList from './WRProjectList.vue';

const projectStore = useProjectStore();
const { categories } = storeToRefs(projectStore);
const { createFilter } = projectStore;

const { filter, updateFilter, filteredProjects } = createFilter();
</script>