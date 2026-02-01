import { describe, it, expect, vi, beforeEach } from 'vitest'

// Mock vue3-google-signin
vi.mock('vue3-google-signin', () => ({
  decodeCredential: vi.fn((credential) => {
    // Return a mock decoded JWT payload
    const payloads = {
      'valid-token-waldritter': {
        hd: 'waldritter.de',
        email: 'user@waldritter.de',
        name: 'Test User',
        picture: 'https://example.com/pic.jpg',
        sub: '12345',
        exp: Math.floor(Date.now() / 1000) + 3600 // 1 hour from now
      },
      'valid-token-other-domain': {
        hd: 'gmail.com',
        email: 'user@gmail.com',
        name: 'External User',
        picture: 'https://example.com/pic2.jpg',
        sub: '67890',
        exp: Math.floor(Date.now() / 1000) + 3600
      },
      'expired-token': {
        hd: 'waldritter.de',
        email: 'expired@waldritter.de',
        name: 'Expired User',
        picture: 'https://example.com/pic3.jpg',
        sub: '11111',
        exp: Math.floor(Date.now() / 1000) - 3600 // 1 hour ago
      }
    }
    return payloads[credential] || null
  })
}))

// Reset module state between tests by re-importing
let useGoogleAuth

describe('useGoogleAuth', () => {
  beforeEach(async () => {
    // Clear all mocks
    vi.clearAllMocks()

    // Reset localStorage
    localStorage.clear()

    // Reset the module to get fresh state
    vi.resetModules()

    // Re-import to get fresh state
    const module = await import('./google_auth.js')
    useGoogleAuth = module.useGoogleAuth
  })

  describe('initial state', () => {
    it('starts logged out with no saved user', () => {
      const { isLoggedIn, token, username } = useGoogleAuth()

      expect(isLoggedIn.value).toBe(false)
      expect(token.value).toBe(null)
      expect(username.value).toBe(null)
    })

    it('restores user from localStorage on load', async () => {
      const savedUser = {
        token: 'saved-token',
        email: 'saved@waldritter.de',
        name: 'Saved User'
      }
      localStorage.setItem('google_user', JSON.stringify(savedUser))

      // Re-import to trigger load
      vi.resetModules()
      const module = await import('./google_auth.js')
      const { isLoggedIn, token, username } = module.useGoogleAuth()

      expect(isLoggedIn.value).toBe(true)
      expect(token.value).toBe('saved-token')
      expect(username.value).toBe('saved@waldritter.de')
    })

    it('clears invalid localStorage data', async () => {
      localStorage.setItem('google_user', 'invalid-json')

      vi.resetModules()
      const module = await import('./google_auth.js')
      const { isLoggedIn } = module.useGoogleAuth()

      expect(isLoggedIn.value).toBe(false)
      expect(localStorage.getItem('google_user')).toBe(null)
    })

    it('clears localStorage when token is missing', async () => {
      localStorage.setItem('google_user', JSON.stringify({ email: 'no-token@test.de' }))

      vi.resetModules()
      const module = await import('./google_auth.js')
      const { isLoggedIn } = module.useGoogleAuth()

      expect(isLoggedIn.value).toBe(false)
      expect(localStorage.getItem('google_user')).toBe(null)
    })
  })

  describe('handleCredentialResponse', () => {
    it('accepts waldritter.de domain', () => {
      const { handleCredentialResponse, isLoggedIn, user } = useGoogleAuth()

      const result = handleCredentialResponse({ credential: 'valid-token-waldritter' })

      expect(result).toBe(true)
      expect(isLoggedIn.value).toBe(true)
      expect(user.value.email).toBe('user@waldritter.de')
      expect(user.value.token).toBe('valid-token-waldritter')
    })

    it('rejects non-waldritter.de domains', () => {
      const { handleCredentialResponse, isLoggedIn } = useGoogleAuth()

      const result = handleCredentialResponse({ credential: 'valid-token-other-domain' })

      expect(result).toBe(false)
      expect(isLoggedIn.value).toBe(false)
      expect(alert).toHaveBeenCalledWith('Nur waldritter.de E-Mail-Adressen sind erlaubt')
    })

    it('stores user data in localStorage on successful login', () => {
      const { handleCredentialResponse } = useGoogleAuth()

      handleCredentialResponse({ credential: 'valid-token-waldritter' })

      const stored = JSON.parse(localStorage.getItem('google_user'))
      expect(stored.email).toBe('user@waldritter.de')
      expect(stored.token).toBe('valid-token-waldritter')
      expect(stored.name).toBe('Test User')
      expect(stored.picture).toBe('https://example.com/pic.jpg')
      expect(stored.sub).toBe('12345')
    })

    it('stores expiration time from token', () => {
      const { handleCredentialResponse } = useGoogleAuth()

      handleCredentialResponse({ credential: 'valid-token-waldritter' })

      const stored = JSON.parse(localStorage.getItem('google_user'))
      expect(stored.expires_at).toBeDefined()
      expect(typeof stored.expires_at).toBe('number')
    })
  })

  describe('signOut', () => {
    it('clears user state', () => {
      const { handleCredentialResponse, signOut, isLoggedIn, user } = useGoogleAuth()

      handleCredentialResponse({ credential: 'valid-token-waldritter' })
      expect(isLoggedIn.value).toBe(true)

      signOut()

      expect(isLoggedIn.value).toBe(false)
      expect(user.value).toBe(null)
    })

    it('removes user from localStorage', () => {
      const { handleCredentialResponse, signOut } = useGoogleAuth()

      handleCredentialResponse({ credential: 'valid-token-waldritter' })
      expect(localStorage.getItem('google_user')).not.toBe(null)

      signOut()

      expect(localStorage.getItem('google_user')).toBe(null)
    })

    it('calls Google disableAutoSelect', () => {
      const { signOut } = useGoogleAuth()

      signOut()

      expect(window.google.accounts.id.disableAutoSelect).toHaveBeenCalled()
    })

    it('handles missing Google object gracefully', () => {
      const originalGoogle = window.google
      window.google = undefined

      const { signOut } = useGoogleAuth()

      expect(() => signOut()).not.toThrow()

      window.google = originalGoogle
    })
  })

  describe('computed properties', () => {
    it('token returns current token value', () => {
      const { handleCredentialResponse, token } = useGoogleAuth()

      expect(token.value).toBe(null)

      handleCredentialResponse({ credential: 'valid-token-waldritter' })

      expect(token.value).toBe('valid-token-waldritter')
    })

    it('username returns email', () => {
      const { handleCredentialResponse, username } = useGoogleAuth()

      expect(username.value).toBe(null)

      handleCredentialResponse({ credential: 'valid-token-waldritter' })

      expect(username.value).toBe('user@waldritter.de')
    })
  })
})
