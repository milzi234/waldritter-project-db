import { ref, computed } from 'vue'
import { decodeCredential } from 'vue3-google-signin'

const user = ref(null)

function isTokenExpired(userData) {
  if (!userData || !userData.expires_at) return true
  // expires_at is in seconds (JWT exp), compare to current time in seconds
  return userData.expires_at <= Math.floor(Date.now() / 1000)
}

// Load saved user on startup
const savedUser = localStorage.getItem('google_user')
if (savedUser) {
  try {
    const parsed = JSON.parse(savedUser)
    if (parsed && parsed.token && !isTokenExpired(parsed)) {
      user.value = parsed
    } else {
      localStorage.removeItem('google_user')
    }
  } catch (e) {
    localStorage.removeItem('google_user')
  }
}

export function useGoogleAuth() {
  const isLoggedIn = computed(() => !!user.value && !isTokenExpired(user.value))
  const token = computed(() => {
    if (!user.value || isTokenExpired(user.value)) return null
    return user.value.token
  })
  const username = computed(() => user.value?.email || null)
  
  const handleCredentialResponse = (response) => {
    const { credential } = response
    const payload = decodeCredential(credential)
    
    // Check domain restriction
    if (payload.hd !== 'waldritter.de') {
      alert('Nur waldritter.de E-Mail-Adressen sind erlaubt')
      return false
    }
    
    // Store user info
    user.value = {
      token: credential,
      email: payload.email,
      name: payload.name,
      picture: payload.picture,
      sub: payload.sub,
      expires_at: payload.exp
    }
    
    localStorage.setItem('google_user', JSON.stringify(user.value))
    return true
  }
  
  const signOut = () => {
    // Sign out from Google
    if (window.google && window.google.accounts) {
      window.google.accounts.id.disableAutoSelect()
    }
    
    // Clear local state
    user.value = null
    localStorage.removeItem('google_user')
  }
  
  return { 
    user,
    token,
    isLoggedIn,
    username,
    handleCredentialResponse,
    signOut
  }
}