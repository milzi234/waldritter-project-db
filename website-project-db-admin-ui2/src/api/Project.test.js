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
import useProjectAPI from './Project.js'

describe('ProjectAPI', () => {
  let projectAPI

  beforeEach(() => {
    vi.clearAllMocks()
    localStorage.clear()
    projectAPI = useProjectAPI()
    // Clear listeners between tests
    projectAPI.listeners = {}
  })

  describe('getAuthToken', () => {
    it('returns token from localStorage', () => {
      localStorage.setItem('google_user', JSON.stringify({ token: 'test-token' }))

      const token = projectAPI.getAuthToken()

      expect(token).toBe('test-token')
    })

    it('returns null when no user in localStorage', () => {
      const token = projectAPI.getAuthToken()

      expect(token).toBe(null)
    })

    it('returns null on invalid JSON', () => {
      localStorage.setItem('google_user', 'invalid-json')

      const token = projectAPI.getAuthToken()

      expect(token).toBe(null)
    })
  })

  describe('CRUD operations', () => {
    beforeEach(() => {
      localStorage.setItem('google_user', JSON.stringify({ token: 'auth-token' }))
    })

    it('getAll fetches all projects', async () => {
      const projects = [{ id: 1, title: 'Project 1' }, { id: 2, title: 'Project 2' }]
      mockAxios.get.mockResolvedValue({ data: projects })

      const result = await projectAPI.getAll()

      expect(mockAxios.get).toHaveBeenCalledWith(
        'http://localhost:3000/api/v1/projects',
        { headers: { Authorization: 'Bearer auth-token' } }
      )
      expect(result).toEqual(projects)
    })

    it('get fetches single project by id', async () => {
      const project = { id: 1, title: 'Project 1' }
      mockAxios.get.mockResolvedValue({ data: project })

      const result = await projectAPI.get(1)

      expect(mockAxios.get).toHaveBeenCalledWith(
        'http://localhost:3000/api/v1/projects/1',
        { headers: { Authorization: 'Bearer auth-token' } }
      )
      expect(result).toEqual(project)
    })

    it('create posts new project and emits event', async () => {
      const newProject = { title: 'New Project' }
      const createdProject = { id: 1, ...newProject }
      mockAxios.post.mockResolvedValue({ data: createdProject })

      const listener = vi.fn()
      projectAPI.on('created', listener)

      const result = await projectAPI.create(newProject)

      expect(mockAxios.post).toHaveBeenCalledWith(
        'http://localhost:3000/api/v1/projects',
        newProject,
        { headers: { Authorization: 'Bearer auth-token' } }
      )
      expect(result).toEqual(createdProject)
      expect(listener).toHaveBeenCalledWith(createdProject)
    })

    it('update puts project data and emits event', async () => {
      const updateData = { title: 'Updated Project' }
      const updatedProject = { id: 1, ...updateData }
      mockAxios.put.mockResolvedValue({ data: updatedProject })

      const listener = vi.fn()
      projectAPI.on('updated', listener)

      const result = await projectAPI.update(1, updateData)

      expect(mockAxios.put).toHaveBeenCalledWith(
        'http://localhost:3000/api/v1/projects/1',
        updateData,
        { headers: { Authorization: 'Bearer auth-token' } }
      )
      expect(result).toEqual(updatedProject)
      expect(listener).toHaveBeenCalledWith(updatedProject)
    })

    it('delete removes project and emits event', async () => {
      mockAxios.delete.mockResolvedValue({ data: { success: true } })

      const listener = vi.fn()
      projectAPI.on('deleted', listener)

      const result = await projectAPI.delete(1)

      expect(mockAxios.delete).toHaveBeenCalledWith(
        'http://localhost:3000/api/v1/projects/1',
        { headers: { Authorization: 'Bearer auth-token' } }
      )
      expect(listener).toHaveBeenCalled()
    })
  })

  describe('image upload', () => {
    beforeEach(() => {
      localStorage.setItem('google_user', JSON.stringify({ token: 'auth-token' }))
    })

    it('uploads image with multipart form data', async () => {
      const formData = new FormData()
      formData.append('image', new Blob(['test']), 'test.jpg')
      mockAxios.post.mockResolvedValue({ data: { image_url: '/images/test.jpg' } })

      const listener = vi.fn()
      projectAPI.on('image-uploaded', listener)

      const result = await projectAPI.uploadImage(1, formData)

      expect(mockAxios.post).toHaveBeenCalledWith(
        'http://localhost:3000/api/v1/projects/1/upload_image',
        formData,
        {
          headers: {
            'Content-Type': 'multipart/form-data',
            Authorization: 'Bearer auth-token'
          }
        }
      )
      expect(listener).toHaveBeenCalled()
    })
  })

  describe('tags', () => {
    beforeEach(() => {
      localStorage.setItem('google_user', JSON.stringify({ token: 'auth-token' }))
    })

    it('getProjectTags fetches tags for project', async () => {
      const tags = [{ id: 1, title: 'Tag 1' }]
      mockAxios.get.mockResolvedValue({ data: tags })

      const result = await projectAPI.getProjectTags(1)

      expect(mockAxios.get).toHaveBeenCalledWith(
        'http://localhost:3000/api/v1/projects/1/tags',
        { headers: { Authorization: 'Bearer auth-token' } }
      )
      expect(result).toEqual(tags)
    })

    it('setProjectTags updates project tags', async () => {
      mockAxios.post.mockResolvedValue({ data: { success: true } })

      const listener = vi.fn()
      projectAPI.on('tags-updated', listener)

      await projectAPI.setProjectTags(1, [1, 2, 3])

      expect(mockAxios.post).toHaveBeenCalledWith(
        'http://localhost:3000/api/v1/projects/1/tags',
        { tag_ids: [1, 2, 3] },
        { headers: { Authorization: 'Bearer auth-token' } }
      )
      expect(listener).toHaveBeenCalled()
    })
  })

  describe('occurrences', () => {
    beforeEach(() => {
      localStorage.setItem('google_user', JSON.stringify({ token: 'auth-token' }))
    })

    it('getNextOccurrence returns first occurrence when available', async () => {
      const occurrence = { id: 1, date: '2024-01-15' }
      mockAxios.get.mockResolvedValue({ data: [occurrence] })

      const result = await projectAPI.getNextOccurrence(1)

      expect(mockAxios.get).toHaveBeenCalledWith(
        expect.stringMatching(/\/api\/v1\/projects\/1\/occurrences\?start_date=.*&end_date=.*&limit=1/),
        { headers: { Authorization: 'Bearer auth-token' } }
      )
      expect(result).toEqual(occurrence)
    })

    it('getNextOccurrence returns null when no occurrences', async () => {
      mockAxios.get.mockResolvedValue({ data: [] })

      const result = await projectAPI.getNextOccurrence(1)

      expect(result).toBe(null)
    })
  })

  describe('event listener pattern', () => {
    it('registers multiple listeners for same event', () => {
      const listener1 = vi.fn()
      const listener2 = vi.fn()

      projectAPI.on('created', listener1)
      projectAPI.on('created', listener2)

      projectAPI.emit('created', { id: 1 })

      expect(listener1).toHaveBeenCalledWith({ id: 1 })
      expect(listener2).toHaveBeenCalledWith({ id: 1 })
    })

    it('does not error when emitting event with no listeners', () => {
      expect(() => projectAPI.emit('nonexistent', {})).not.toThrow()
    })
  })

  describe('eventAPIFor factory', () => {
    it('returns same EventAPI instance for same projectID', () => {
      const eventAPI1 = projectAPI.eventAPIFor(1)
      const eventAPI2 = projectAPI.eventAPIFor(1)

      expect(eventAPI1).toBe(eventAPI2)
    })

    it('returns different EventAPI instances for different projectIDs', () => {
      const eventAPI1 = projectAPI.eventAPIFor(1)
      const eventAPI2 = projectAPI.eventAPIFor(2)

      expect(eventAPI1).not.toBe(eventAPI2)
    })
  })
})

describe('EventAPI', () => {
  let projectAPI
  let eventAPI

  beforeEach(() => {
    vi.clearAllMocks()
    localStorage.clear()
    localStorage.setItem('google_user', JSON.stringify({ token: 'auth-token' }))
    projectAPI = useProjectAPI()
    projectAPI.listeners = {}
    eventAPI = projectAPI.eventAPIFor(1)
  })

  describe('CRUD operations', () => {
    it('getAll fetches all events for project', async () => {
      const events = [{ id: 1, title: 'Event 1' }]
      mockAxios.get.mockResolvedValue({ data: events })

      const result = await eventAPI.getAll()

      expect(mockAxios.get).toHaveBeenCalledWith(
        'http://localhost:3000/api/v1/projects/1/events',
        { headers: { Authorization: 'Bearer auth-token' } }
      )
      expect(result).toEqual(events)
    })

    it('create emits event through parent projectAPI', async () => {
      mockAxios.post.mockResolvedValue({ data: { id: 1 } })

      const listener = vi.fn()
      projectAPI.on('event-created', listener)

      await eventAPI.create({ title: 'New Event' })

      expect(listener).toHaveBeenCalled()
    })

    it('update emits event through parent projectAPI', async () => {
      mockAxios.put.mockResolvedValue({ data: { id: 1 } })

      const listener = vi.fn()
      projectAPI.on('event-updated', listener)

      await eventAPI.update(1, { title: 'Updated Event' })

      expect(listener).toHaveBeenCalled()
    })

    it('delete emits event through parent projectAPI', async () => {
      mockAxios.delete.mockResolvedValue({ data: { success: true } })

      const listener = vi.fn()
      projectAPI.on('event-deleted', listener)

      await eventAPI.delete(1)

      expect(listener).toHaveBeenCalled()
    })
  })

  describe('occurrences and exceptions', () => {
    it('getOccurrences fetches occurrences within date range', async () => {
      const occurrences = [{ id: 1, date: '2024-01-15' }]
      mockAxios.get.mockResolvedValue({ data: occurrences })

      const startDate = new Date('2024-01-01')
      const endDate = new Date('2024-12-31')

      const result = await eventAPI.getOccurrences(1, startDate, endDate)

      expect(mockAxios.get).toHaveBeenCalledWith(
        expect.stringContaining('/events/1/occurrences?start_date='),
        { headers: { Authorization: 'Bearer auth-token' } }
      )
      expect(result).toEqual(occurrences)
    })

    it('createException emits exception-created event', async () => {
      mockAxios.delete.mockResolvedValue({ data: { success: true } })

      const listener = vi.fn()
      projectAPI.on('exception-created', listener)

      await eventAPI.createException(1, 'occ-123')

      expect(mockAxios.delete).toHaveBeenCalledWith(
        'http://localhost:3000/api/v1/projects/1/events/1/occurrences/occ-123',
        { headers: { Authorization: 'Bearer auth-token' } }
      )
      expect(listener).toHaveBeenCalled()
    })

    it('deleteException emits exception-deleted event', async () => {
      mockAxios.delete.mockResolvedValue({ data: { success: true } })

      const listener = vi.fn()
      projectAPI.on('exception-deleted', listener)

      await eventAPI.deleteException(1, 'exc-456')

      expect(listener).toHaveBeenCalled()
    })
  })
})
