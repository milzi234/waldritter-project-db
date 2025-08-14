<script setup>
import { RouterLink, RouterView } from 'vue-router'
import { useRoute, useRouter } from 'vue-router'
import { computed, ref } from 'vue'
import { useToken } from '@/composables/auth_token'

const route = useRoute()
const router = useRouter()
const isLoginPage = computed(() => route.name === 'login')
const { username, loggedIn, setUser } = useToken()
const showContextMenu = ref(false)

const handleLogout = () => {
  setUser(null)
  // Redirect to login page after logout
  router.push({ name: 'login' })
}

const toggleContextMenu = () => {
  showContextMenu.value = !showContextMenu.value
}

</script>

<template>
  <div class="container">
  <header>
    <div class="wrapper" v-if="!isLoginPage">
      <div class="d-flex justify-content-between align-items-center mb-3">
        <div class="invisible">Placeholder</div>
        <nav class="nav nav-pills justify-content-center">
          <RouterLink to="/projects" class="nav-link">Projekte</RouterLink>
          <RouterLink to="/" class="nav-link">Tags</RouterLink>
        </nav>
        <div v-if="loggedIn" class="d-flex align-items-center position-relative">
          <a href="#" @click.prevent="toggleContextMenu" class="text-muted small">{{ username }}</a>
          <div v-if="showContextMenu" class="position-absolute top-100 end-0 mt-1 bg-white border rounded shadow-sm">
            <button @click="handleLogout" class="btn btn-link btn-sm text-muted p-2">Logout</button>
          </div>
        </div>
      </div>
    </div>
  </header>

  <RouterView />
  </div>
</template>
