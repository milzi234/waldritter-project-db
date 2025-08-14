<template>
  <div class="container d-flex flex-column justify-content-center align-items-center vh-100">
    <h1 class="mb-4 text-center fancy-headline">Waldritter Projekt Datenbank</h1>
    
    <!-- Google Sign-In Button -->
    <GoogleSignInButton
      @success="handleLoginSuccess"
      @error="handleLoginError"
      :client-id="googleClientId"
      prompt-parent-id="prompt"
      :auto-select="false"
      :use-fedcm-for-prompt="true"
    >
      <button class="btn btn-primary btn-lg px-5 py-3 shadow-lg rounded-pill">
        <span class="fs-4 fw-bold">Mit Google anmelden</span>
      </button>
    </GoogleSignInButton>
    
    <!-- Google One Tap prompt container -->
    <div id="prompt"></div>
    
    <!-- Error message -->
    <div v-if="errorMessage" class="alert alert-danger mt-3" role="alert">
      {{ errorMessage }}
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

// Get Google Client ID from environment
const googleClientId = import.meta.env.VITE_GOOGLE_CLIENT_ID || ''

const handleLoginSuccess = (response) => {
  errorMessage.value = ''
  
  const success = handleCredentialResponse(response)
  
  if (success) {
    // Navigate to projects page after successful login
    router.push('/projects')
  } else {
    errorMessage.value = 'Nur waldritter.de E-Mail-Adressen sind erlaubt'
  }
}

const handleLoginError = () => {
  errorMessage.value = 'Anmeldung fehlgeschlagen. Bitte versuchen Sie es erneut.'
  console.error('Google Sign-In error')
}

// Check if client ID is configured
if (!googleClientId) {
  console.error('Google Client ID not configured. Please set VITE_GOOGLE_CLIENT_ID in your .env file')
  errorMessage.value = 'Google Sign-In ist nicht konfiguriert. Bitte kontaktieren Sie den Administrator.'
}
</script>

<style scoped>
.btn-primary {
  background: linear-gradient(45deg, #4CAF50, #2E7D32);
  border: none;
  transition: all 0.3s ease;
}

.btn-primary:hover {
  transform: scale(1.05);
  box-shadow: 0 0 15px rgba(76, 175, 80, 0.5) !important;
}

.fancy-headline {
  font-size: 2.5rem;
  font-weight: bold;
  color: #4CAF50;
  text-shadow: 2px 2px 4px rgba(0, 0, 0, 0.1);
}

.fancy-subheading {
  font-size: 1.5rem;
  color: #2E7D32;
  font-style: italic;
}

.alert {
  max-width: 400px;
}
</style>