<template>
  <div class="time-range">
    <template v-if="isSameDay && !separateLines">
      <span class="font-heading text-gray-600 mb-1">{{ formatDate(startDateTime) }}</span>
      <span class="text-gray-400"> ( {{ formatTime(startDateTime) }} - {{ formatTime(endDateTime) }} )</span>
    </template>
    <template v-else-if="isSameDay && separateLines">
      <div class="font-heading text-gray-600 mb-1">{{ formatDate(startDateTime) }}</div>
      <div class="text-gray-400">{{ formatTime(startDateTime) }} - {{ formatTime(endDateTime) }}</div>
    </template>
    <template v-else-if="!separateLines">
      <div class="font-heading text-gray-600 mb-1">
        <span>{{ formatDate(startDateTime) }}</span>
        <span class="text-gray-400 ml-2">{{ formatTime(startDateTime) }} Uhr</span>
        <span class="mx-2">-</span>
        <span>{{ formatDate(endDateTime) }}</span>
        <span class="text-gray-400 ml-2">{{ formatTime(endDateTime) }} Uhr</span>
      </div>
    </template>
    <template v-else>
      <div class="font-heading text-gray-600 mb-1">
        <div>{{ formatDate(startDateTime) }} <span class="text-gray-400">{{ formatTime(startDateTime) }} Uhr</span></div>
        <div>{{ formatDate(endDateTime) }} <span class="text-gray-400">{{ formatTime(endDateTime) }} Uhr</span></div>
      </div>
    </template>
  </div>
</template>

<script setup>
import { computed } from 'vue';

const props = defineProps({
  startDateTime: {
    type: String,
    required: true
  },
  endDateTime: {
    type: String,
    required: true
  },
  separateLines: {
    type: Boolean,
    default: false
  }
});

const isSameDay = computed(() => {
  const start = new Date(props.startDateTime);
  const end = new Date(props.endDateTime);
  return start.toDateString() === end.toDateString();
});

const formatDate = (dateString) => {
  const date = new Date(dateString);
  return date.toLocaleDateString('de-DE', {
    weekday: 'long',
    year: 'numeric',
    month: 'long',
    day: 'numeric'
  });
};

const formatTime = (dateString) => {
  const date = new Date(dateString);
  return date.toLocaleTimeString('de-DE', {
    hour: '2-digit',
    minute: '2-digit'
  });
};
</script>
