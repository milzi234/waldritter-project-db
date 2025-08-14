<template>
  <div class="container d-flex flex-column justify-content-center align-items-center vh-100">
    <h1 class="mb-4 text-center fancy-headline">Waldritter Projekt Datenbank</h1>
    <h2 class="mb-5 text-center fancy-subheading">Bitte weisen Sie sich aus</h2>
    <button @click="handleLogin" class="btn btn-primary btn-lg px-5 py-3 shadow-lg rounded-pill">
      <span class="fs-4 fw-bold">Login</span>
    </button>
  </div>
</template>

<script setup>
import { useToken } from '@/composables/auth_token';

const { startAuthFlow } = useToken();

const handleLogin = async () => {
  try {
    await startAuthFlow();
  } catch (error) {
    console.error('Error starting auth flow:', error);
  }
};
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
</style>
