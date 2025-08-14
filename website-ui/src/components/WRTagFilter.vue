<template>
    <div class="mb-4">
        <Button @click="toggleDrawer" :label="buttonLabel" />
        
        <Drawer v-model:visible="drawerVisible" :style="{ width: '40rem' }" :dismissable="true" :modal="true">
            <template #header>
                <h3 class="text-xl font-semibold">Filter</h3>
            </template>
            
            <div class="space-y-6">
                <div v-for="category in categories" :key="category.id" class="category">
                    <h4 class="text-lg font-medium text-gray-600 mb-2">{{ category.title }}</h4>
                    <div class="grid grid-cols-2 gap-4 sm:grid-cols-3">
                        <div v-for="tag in category.tags" :key="tag.id" class="flex items-center">
                            <Checkbox :inputId="'tag-' + tag.id" :value="tag.id" v-model="selectedTagIds" :binary="false" class="mr-2" @change="emitUpdate" />
                            <label :for="'tag-' + tag.id" class="text-sm leading-tight break-words">{{ tag.title }}</label>
                        </div>
                    </div>
                </div>
            </div>
        </Drawer>
    </div>
</template>

<script setup>
import { ref, watch, computed } from 'vue';
import Button from 'primevue/button';
import Drawer from 'primevue/drawer';
import Checkbox from 'primevue/checkbox';

const props = defineProps({
    categories: {
        type: Array,
        required: true
    },
    selectedTags: {
        type: Array,
        default: () => []
    }
});

const emit = defineEmits(['update']);

const drawerVisible = ref(false);
const selectedTagIds = ref([...props.selectedTags]);

const buttonLabel = computed(() => {
    return selectedTagIds.value.length > 0 ? `Filter (${selectedTagIds.value.length})` : 'Filter';
});

const toggleDrawer = () => {
    drawerVisible.value = !drawerVisible.value;
};

const emitUpdate = () => {
    emit('update', selectedTagIds.value);
};

watch(() => props.selectedTags, (newTags) => {
    selectedTagIds.value = [...newTags];
}, { deep: true });
</script>
