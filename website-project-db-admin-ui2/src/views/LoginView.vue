<template>
  <div class="flex flex-col items-center justify-center min-h-[80vh]">
    <div class="section-panel max-w-md w-full text-center">
      <h1 class="text-3xl font-display font-bold text-wald-300 mb-2 glitch-text" data-text="Waldritter">Waldritter</h1>
      <p class="text-gray-500 font-mono text-xs tracking-widest uppercase mb-8">Projekt Datenbank</p>

      <GoogleSignInButton
        @success="handleLoginSuccess"
        @error="handleLoginError"
        :client-id="googleClientId"
        prompt-parent-id="prompt"
        :auto-select="false"
        :use-fedcm-for-prompt="true"
      >
        <button class="btn-cyber w-full py-3 text-base">
          Mit Google anmelden
        </button>
      </GoogleSignInButton>

      <div id="prompt"></div>

      <div v-if="errorMessage" class="alert-danger mt-4">
        {{ errorMessage }}
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { GoogleSignInButton } from 'vue3-google-signin'
import { useGoogleAuth } from '@/composables/google_auth'

const router = useRouter()
const { handleCredentialResponse } = useGoogleAuth()
const errorMessage = ref('')

const googleClientId = import.meta.env.VITE_GOOGLE_CLIENT_ID || ''

const handleLoginSuccess = (response) => {
  errorMessage.value = ''

  const success = handleCredentialResponse(response)

  if (success) {
    router.push('/projects')
  } else {
    errorMessage.value = 'Nur waldritter.de E-Mail-Adressen sind erlaubt'
  }
}

const handleLoginError = () => {
  errorMessage.value = 'Anmeldung fehlgeschlagen. Bitte versuchen Sie es erneut.'
  console.error('Google Sign-In error')
}

if (!googleClientId) {
  console.error('Google Client ID not configured. Please set VITE_GOOGLE_CLIENT_ID in your .env file')
  errorMessage.value = 'Google Sign-In ist nicht konfiguriert. Bitte kontaktieren Sie den Administrator.'
}
</script>
