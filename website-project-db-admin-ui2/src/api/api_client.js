import axios from '@bundled-es-modules/axios'
import { useGoogleAuth } from '@/composables/google_auth'

// Create axios instance
const apiClient = axios.create({
  baseURL: import.meta.env.VITE_API_URL || 'http://localhost:3000',
  headers: {
    'Content-Type': 'application/json'
  }
})

// Add auth token to requests
apiClient.interceptors.request.use(config => {
  const { token } = useGoogleAuth()
  
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
      // Token expired or invalid - sign out
      const { signOut } = useGoogleAuth()
      signOut()
      
      // Redirect to login
      if (window.location.pathname !== '/login') {
        window.location.href = '/login'
      }
    }
    return Promise.reject(error)
  }
)

export default apiClient