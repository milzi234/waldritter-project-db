import { ref, computed } from 'vue'
import { UserManager } from 'oidc-client-ts'

const oidcUserManager = new UserManager({
  authority: 'https://auth.waldritter.cisco.local:8443',
  client_id: 'waldritter-admin',
  redirect_uri: 'https://admin.waldritter.cisco.local:5173/login/callback'
})

// Eagerly fetch the OpenID configuration
oidcUserManager.metadataService.getMetadata().catch(error => {
  console.error('Failed to fetch OpenID configuration:', error)
})

export function useToken() {
  const storedUser = JSON.parse(localStorage.getItem('user'))
  const user = ref(null)

  if (storedUser && storedUser.expires_at) {
    const expiresAt = new Date(storedUser.expires_at * 1000) // Convert to milliseconds
    if (expiresAt > new Date()) {
      user.value = storedUser
    } else {
      localStorage.removeItem('user')
    }
  }

  const token = computed(() => user.value?.id_token || null)

  const loggedIn = computed(() => !!token.value)
  const username = computed(() => user.value?.profile?.aud || null)

  const setUser = (newUser) => {
    user.value = newUser
    if (newUser) {
      localStorage.setItem('user', JSON.stringify(newUser))
    } else {
      localStorage.removeItem('user')
    }
  }

  async function startAuthFlow(url_state) {
    await oidcUserManager.signinRedirect({ url_state })
  }
  
  async function handleCallback() {
    const userData = await oidcUserManager.signinCallback()
    console.log('userData', userData)
    setUser(userData)
    return userData.url_state
  }

  return {
    token,
    user,
    loggedIn,
    username,
    setUser,
    startAuthFlow,
    handleCallback
  }
}
