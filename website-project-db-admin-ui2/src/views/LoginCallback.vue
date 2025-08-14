<template>
  <div>Verarbeite Login...</div>
</template>

<script setup>
import { useToken } from '@/composables/auth_token'
import { useRouter } from 'vue-router'

const router = useRouter()
const { handleCallback } = useToken()

const processCallback = async () => {
  try {
    await handleCallback()
    router.push('/projects')
  } catch (error) {
    console.error('Error during callback:', error)
    router.push('/login')
  }
}

processCallback()
</script>
