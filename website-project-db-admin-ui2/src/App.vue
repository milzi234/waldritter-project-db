<script setup>
import { RouterLink, RouterView } from 'vue-router'
import { useRoute, useRouter } from 'vue-router'
import { computed, ref } from 'vue'
import { useGoogleAuth } from '@/composables/google_auth'

const route = useRoute()
const router = useRouter()
const isLoginPage = computed(() => route.name === 'login')
const { username, isLoggedIn, signOut } = useGoogleAuth()
const showContextMenu = ref(false)

const handleLogout = () => {
  signOut()
  router.push({ name: 'login' })
}

const toggleContextMenu = () => {
  showContextMenu.value = !showContextMenu.value
}

</script>

<template>
  <div class="min-h-screen">
    <header v-if="!isLoginPage" class="sticky top-0 z-50 border-b border-wald-500/20 bg-black/80 backdrop-blur-md">
      <div class="max-w-6xl mx-auto px-4 py-3 flex items-center justify-between">
        <div class="flex items-center gap-2">
          <span class="text-wald-400 font-mono text-xs tracking-widest uppercase">SYS::ADMIN</span>
          <span class="w-1.5 h-1.5 rounded-full bg-wald-400 animate-pulse"></span>
        </div>
        <nav class="flex items-center gap-1">
          <RouterLink to="/projects" class="nav-link">Projekte</RouterLink>
          <RouterLink to="/" class="nav-link">Tags</RouterLink>
        </nav>
        <div v-if="isLoggedIn" class="relative">
          <a href="#" @click.prevent="toggleContextMenu" class="text-gray-500 text-xs font-mono hover:text-gray-300 transition-colors">{{ username }}</a>
          <div v-if="showContextMenu" class="absolute top-full right-0 mt-1 section-panel !p-0 min-w-[120px]">
            <button @click="handleLogout" class="w-full text-left px-3 py-2 text-sm text-gray-400 hover:text-red-400 font-mono transition-colors">Logout</button>
          </div>
        </div>
      </div>
    </header>

    <main class="max-w-4xl mx-auto px-4 py-8">
      <RouterView />
    </main>
  </div>
</template>
