// Global test setup
import { vi } from 'vitest'

// Mock localStorage
const localStorageMock = {
  store: {},
  getItem: vi.fn((key) => localStorageMock.store[key] || null),
  setItem: vi.fn((key, value) => {
    localStorageMock.store[key] = String(value)
  }),
  removeItem: vi.fn((key) => {
    delete localStorageMock.store[key]
  }),
  clear: vi.fn(() => {
    localStorageMock.store = {}
  })
}

Object.defineProperty(global, 'localStorage', {
  value: localStorageMock,
  writable: true
})

// Reset localStorage before each test
beforeEach(() => {
  localStorageMock.clear()
  localStorageMock.getItem.mockClear()
  localStorageMock.setItem.mockClear()
  localStorageMock.removeItem.mockClear()
})

// Mock window.location
const locationMock = {
  href: '',
  pathname: '/',
  assign: vi.fn(),
  replace: vi.fn()
}

Object.defineProperty(global, 'location', {
  value: locationMock,
  writable: true
})

// Mock window.alert
global.alert = vi.fn()

// Mock window.google for Google Sign-In
global.google = {
  accounts: {
    id: {
      disableAutoSelect: vi.fn()
    }
  }
}
