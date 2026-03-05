import axios from 'axios'
import { useGoogleAuth } from '@/composables/google_auth'

const { token, signOut } = useGoogleAuth()

// Create axios instance
const apiClient = axios.create({
  baseURL: import.meta.env.VITE_API_URL || 'http://localhost:3000',
  headers: {
    'Content-Type': 'application/json'
  }
})

// Add auth token to requests
apiClient.interceptors.request.use(config => {
  if (token.value) {
    config.headers.Authorization = `Bearer ${token.value}`
  }
  return config
}, error => {
  return Promise.reject(error)
})

// Handle auth errors
apiClient.interceptors.response.use(
  response => response,
  error => {
    if (error.response && error.response.status === 401) {
      signOut()
      if (window.location.pathname !== '/login') {
        window.location.href = '/login'
      }
    }
    return Promise.reject(error)
  }
)

export default apiClient