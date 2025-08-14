<script setup>
import { RouterView } from 'vue-router'
import WRMainMenuBar from '@/components/WRMainMenuBar.vue'
import WRMainFooter from '@/components/WRMainFooter.vue'
</script>

<template>
  <WRMainMenuBar />
  <RouterView />
  <WRMainFooter />
</template>
