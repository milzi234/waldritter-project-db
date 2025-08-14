<script setup>
import { useProjectStore } from '../stores/projects';

const projectStore = useProjectStore();
const { categories } = storeToRefs(projectStore);

// Create another filter instance
const { filter: anotherFilter, filteredProjects: anotherFilteredProjects, updateFilter: updateAnotherFilter } = projectStore.createFilter();
</script>