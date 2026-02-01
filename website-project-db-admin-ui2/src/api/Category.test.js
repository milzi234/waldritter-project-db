import { describe, it, expect, vi, beforeEach } from 'vitest'

// Mock axios - use hoisted to ensure mock is defined before vi.mock runs
const mockAxios = vi.hoisted(() => ({
  get: vi.fn(),
  post: vi.fn(),
  put: vi.fn(),
  delete: vi.fn()
}))

vi.mock('axios', () => ({
  default: mockAxios
}))

// Import after mocking
import useCategoryAPI from './Category.js'

describe('CategoryAPI', () => {
  let categoryAPI

  beforeEach(() => {
    vi.clearAllMocks()
    localStorage.clear()
    categoryAPI = useCategoryAPI()
    // Clear listeners between tests
    categoryAPI.listeners = {}
  })

  describe('getAuthToken', () => {
    it('returns token from localStorage', () => {
      localStorage.setItem('google_user', JSON.stringify({ token: 'test-token' }))

      const token = categoryAPI.getAuthToken()

      expect(token).toBe('test-token')
    })

    it('returns null when no user in localStorage', () => {
      const token = categoryAPI.getAuthToken()

      expect(token).toBe(null)
    })

    it('returns null on invalid JSON', () => {
      localStorage.setItem('google_user', 'invalid-json')

      const token = categoryAPI.getAuthToken()

      expect(token).toBe(null)
    })
  })

  describe('CRUD operations', () => {
    beforeEach(() => {
      localStorage.setItem('google_user', JSON.stringify({ token: 'auth-token' }))
    })

    it('getAll fetches all categories', async () => {
      const categories = [{ id: 1, title: 'Category 1' }, { id: 2, title: 'Category 2' }]
      mockAxios.get.mockResolvedValue({ data: categories })

      const result = await categoryAPI.getAll()

      expect(mockAxios.get).toHaveBeenCalledWith(
        'http://localhost:3000/api/v1/categories',
        { headers: { Authorization: 'Bearer auth-token' } }
      )
      expect(result).toEqual(categories)
    })

    it('get fetches single category by id', async () => {
      const category = { id: 1, title: 'Category 1' }
      mockAxios.get.mockResolvedValue({ data: category })

      const result = await categoryAPI.get(1)

      expect(mockAxios.get).toHaveBeenCalledWith(
        'http://localhost:3000/api/v1/categories/1',
        { headers: { Authorization: 'Bearer auth-token' } }
      )
      expect(result).toEqual(category)
    })

    it('create posts new category and emits event', async () => {
      const newCategory = { title: 'New Category' }
      const createdCategory = { id: 1, ...newCategory }
      mockAxios.post.mockResolvedValue({ data: createdCategory })

      const listener = vi.fn()
      categoryAPI.on('created', listener)

      const result = await categoryAPI.create(newCategory)

      expect(mockAxios.post).toHaveBeenCalledWith(
        'http://localhost:3000/api/v1/categories',
        newCategory,
        { headers: { Authorization: 'Bearer auth-token' } }
      )
      expect(result).toEqual(createdCategory)
      expect(listener).toHaveBeenCalledWith(createdCategory)
    })

    it('update puts category data and emits event', async () => {
      const updateData = { title: 'Updated Category' }
      const updatedCategory = { id: 1, ...updateData }
      mockAxios.put.mockResolvedValue({ data: updatedCategory })

      const listener = vi.fn()
      categoryAPI.on('updated', listener)

      const result = await categoryAPI.update(1, updateData)

      expect(mockAxios.put).toHaveBeenCalledWith(
        'http://localhost:3000/api/v1/categories/1',
        updateData,
        { headers: { Authorization: 'Bearer auth-token' } }
      )
      expect(result).toEqual(updatedCategory)
      expect(listener).toHaveBeenCalledWith(updatedCategory)
    })

    it('delete removes category and emits event', async () => {
      mockAxios.delete.mockResolvedValue({ data: { success: true } })

      const listener = vi.fn()
      categoryAPI.on('deleted', listener)

      const result = await categoryAPI.delete(1)

      expect(mockAxios.delete).toHaveBeenCalledWith(
        'http://localhost:3000/api/v1/categories/1',
        { headers: { Authorization: 'Bearer auth-token' } }
      )
      expect(listener).toHaveBeenCalled()
    })
  })

  describe('event listener pattern', () => {
    it('registers multiple listeners for same event', () => {
      const listener1 = vi.fn()
      const listener2 = vi.fn()

      categoryAPI.on('created', listener1)
      categoryAPI.on('created', listener2)

      categoryAPI.emit('created', { id: 1 })

      expect(listener1).toHaveBeenCalledWith({ id: 1 })
      expect(listener2).toHaveBeenCalledWith({ id: 1 })
    })

    it('does not error when emitting event with no listeners', () => {
      expect(() => categoryAPI.emit('nonexistent', {})).not.toThrow()
    })
  })

  describe('tagAPIFor factory', () => {
    it('returns same TagAPI instance for same categoryID', () => {
      const tagAPI1 = categoryAPI.tagAPIFor(1)
      const tagAPI2 = categoryAPI.tagAPIFor(1)

      expect(tagAPI1).toBe(tagAPI2)
    })

    it('returns different TagAPI instances for different categoryIDs', () => {
      const tagAPI1 = categoryAPI.tagAPIFor(1)
      const tagAPI2 = categoryAPI.tagAPIFor(2)

      expect(tagAPI1).not.toBe(tagAPI2)
    })
  })
})

describe('TagAPI', () => {
  let categoryAPI
  let tagAPI

  beforeEach(() => {
    vi.clearAllMocks()
    localStorage.clear()
    localStorage.setItem('google_user', JSON.stringify({ token: 'auth-token' }))
    categoryAPI = useCategoryAPI()
    categoryAPI.listeners = {}
    tagAPI = categoryAPI.tagAPIFor(1)
  })

  describe('CRUD operations', () => {
    it('getAll fetches all tags for category', async () => {
      const tags = [{ id: 1, title: 'Tag 1' }]
      mockAxios.get.mockResolvedValue({ data: tags })

      const result = await tagAPI.getAll()

      expect(mockAxios.get).toHaveBeenCalledWith(
        'http://localhost:3000/api/v1/categories/1/tags',
        { headers: { Authorization: 'Bearer auth-token' } }
      )
      expect(result).toEqual(tags)
    })

    it('get fetches single tag by id', async () => {
      const tag = { id: 1, title: 'Tag 1' }
      mockAxios.get.mockResolvedValue({ data: tag })

      const result = await tagAPI.get(1)

      expect(mockAxios.get).toHaveBeenCalledWith(
        'http://localhost:3000/api/v1/categories/1/tags/1',
        { headers: { Authorization: 'Bearer auth-token' } }
      )
      expect(result).toEqual(tag)
    })

    it('create emits tag-created through parent categoryAPI', async () => {
      mockAxios.post.mockResolvedValue({ data: { id: 1 } })

      const listener = vi.fn()
      categoryAPI.on('tag-created', listener)

      await tagAPI.create({ title: 'New Tag' })

      expect(listener).toHaveBeenCalled()
    })

    it('update emits tag-updated through parent categoryAPI', async () => {
      mockAxios.put.mockResolvedValue({ data: { id: 1 } })

      const listener = vi.fn()
      categoryAPI.on('tag-updated', listener)

      await tagAPI.update(1, { title: 'Updated Tag' })

      expect(listener).toHaveBeenCalled()
    })

    it('delete emits tag-deleted through parent categoryAPI', async () => {
      mockAxios.delete.mockResolvedValue({ data: { success: true } })

      const listener = vi.fn()
      categoryAPI.on('tag-deleted', listener)

      await tagAPI.delete(1)

      expect(listener).toHaveBeenCalled()
    })
  })

  describe('getAuthToken', () => {
    it('returns token from localStorage', () => {
      const token = tagAPI.getAuthToken()

      expect(token).toBe('auth-token')
    })

    it('returns null when no user in localStorage', () => {
      localStorage.clear()

      const token = tagAPI.getAuthToken()

      expect(token).toBe(null)
    })
  })
})
