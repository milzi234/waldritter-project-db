import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'

// Mock google_auth before importing api_client
const mockToken = { value: null }
const mockSignOut = vi.fn()

vi.mock('@/composables/google_auth', () => ({
  useGoogleAuth: () => ({
    token: mockToken,
    signOut: mockSignOut
  })
}))

// Mock axios
const mockRequestInterceptors = []
const mockResponseInterceptors = []

vi.mock('axios', () => ({
  default: {
    create: vi.fn(() => ({
      interceptors: {
        request: {
          use: vi.fn((onFulfilled, onRejected) => {
            mockRequestInterceptors.push({ onFulfilled, onRejected })
          })
        },
        response: {
          use: vi.fn((onFulfilled, onRejected) => {
            mockResponseInterceptors.push({ onFulfilled, onRejected })
          })
        }
      }
    }))
  }
}))

describe('api_client', () => {
  let originalLocation

  beforeEach(() => {
    vi.clearAllMocks()
    mockRequestInterceptors.length = 0
    mockResponseInterceptors.length = 0
    mockToken.value = null

    // Save original location
    originalLocation = window.location

    // Reset modules to re-register interceptors
    vi.resetModules()
  })

  afterEach(() => {
    // Restore location
    window.location = originalLocation
  })

  describe('axios instance creation', () => {
    it('creates axios instance with correct baseURL', async () => {
      const axios = await import('axios')
      await import('./api_client.js')

      expect(axios.default.create).toHaveBeenCalledWith({
        baseURL: expect.any(String),
        headers: {
          'Content-Type': 'application/json'
        }
      })
    })
  })

  describe('request interceptor', () => {
    it('adds Authorization header when token exists', async () => {
      await import('./api_client.js')

      mockToken.value = 'test-jwt-token'

      const requestInterceptor = mockRequestInterceptors[0]
      const config = { headers: {} }

      const result = requestInterceptor.onFulfilled(config)

      expect(result.headers.Authorization).toBe('Bearer test-jwt-token')
    })

    it('does not add Authorization header when token is null', async () => {
      await import('./api_client.js')

      mockToken.value = null

      const requestInterceptor = mockRequestInterceptors[0]
      const config = { headers: {} }

      const result = requestInterceptor.onFulfilled(config)

      expect(result.headers.Authorization).toBeUndefined()
    })

    it('rejects on request error', async () => {
      await import('./api_client.js')

      const requestInterceptor = mockRequestInterceptors[0]
      const error = new Error('Request failed')

      await expect(requestInterceptor.onRejected(error)).rejects.toThrow('Request failed')
    })
  })

  describe('response interceptor', () => {
    it('passes through successful responses', async () => {
      await import('./api_client.js')

      const responseInterceptor = mockResponseInterceptors[0]
      const response = { data: { success: true }, status: 200 }

      const result = responseInterceptor.onFulfilled(response)

      expect(result).toBe(response)
    })

    it('calls signOut on 401 error', async () => {
      await import('./api_client.js')

      const responseInterceptor = mockResponseInterceptors[0]
      const error = {
        response: { status: 401 }
      }

      // Mock location
      delete window.location
      window.location = { pathname: '/projects', href: '' }

      await expect(responseInterceptor.onRejected(error)).rejects.toBe(error)
      expect(mockSignOut).toHaveBeenCalled()
    })

    it('redirects to login on 401 when not already on login page', async () => {
      await import('./api_client.js')

      const responseInterceptor = mockResponseInterceptors[0]
      const error = {
        response: { status: 401 }
      }

      delete window.location
      window.location = { pathname: '/projects', href: '' }

      try {
        await responseInterceptor.onRejected(error)
      } catch (e) {
        // Expected to reject
      }

      expect(window.location.href).toBe('/login')
    })

    it('does not redirect when already on login page', async () => {
      await import('./api_client.js')

      const responseInterceptor = mockResponseInterceptors[0]
      const error = {
        response: { status: 401 }
      }

      delete window.location
      window.location = { pathname: '/login', href: '/login' }

      try {
        await responseInterceptor.onRejected(error)
      } catch (e) {
        // Expected to reject
      }

      expect(window.location.href).toBe('/login') // unchanged
    })

    it('does not call signOut on non-401 errors', async () => {
      await import('./api_client.js')

      const responseInterceptor = mockResponseInterceptors[0]
      const error = {
        response: { status: 500 }
      }

      try {
        await responseInterceptor.onRejected(error)
      } catch (e) {
        // Expected to reject
      }

      expect(mockSignOut).not.toHaveBeenCalled()
    })

    it('handles errors without response object', async () => {
      await import('./api_client.js')

      const responseInterceptor = mockResponseInterceptors[0]
      const error = new Error('Network error')

      await expect(responseInterceptor.onRejected(error)).rejects.toThrow('Network error')
      expect(mockSignOut).not.toHaveBeenCalled()
    })
  })
})
