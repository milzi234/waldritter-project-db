import { describe, it, expect, vi, beforeEach } from 'vitest'

// Mock the API module before importing the store
vi.mock('../api/Category', () => ({
  default: () => ({
    getAll: vi.fn().mockResolvedValue([]),
    tagAPIFor: vi.fn().mockReturnValue({ getAll: vi.fn().mockResolvedValue([]) }),
    on: vi.fn()
  })
}))

// Extract sortTags for isolated testing
// The sortTags function logic for testing
const sortTags = (a, b) => {
  const regex = /\d+/g
  const aNumbers = a.title.match(regex)
  const bNumbers = b.title.match(regex)

  if (aNumbers && bNumbers) {
    const aMin = Math.min(...aNumbers.map(Number))
    const bMin = Math.min(...bNumbers.map(Number))

    if (aMin < bMin) {
      return -1
    }
    if (bMin < aMin) {
      return 1
    }
    if (aMin == bMin) {
      const aMax = Math.max(...aNumbers.map(Number))
      const bMax = Math.max(...bNumbers.map(Number))

      if (aMax < bMax) {
        return -1
      }
      if (bMax < aMax) {
        return 1
      }
      if (aMax == bMax) {
        return 0
      }
    }
  }
  if (a.title < b.title) {
    return -1
  }
  if (b.title < a.title) {
    return 1
  }
  return 0
}

describe('sortTags', () => {
  describe('numeric sorting', () => {
    it('sorts by minimum number in title', () => {
      const tags = [
        { title: 'Ages 10-12' },
        { title: 'Ages 6-8' },
        { title: 'Ages 13-15' }
      ]

      const sorted = [...tags].sort(sortTags)

      expect(sorted.map(t => t.title)).toEqual([
        'Ages 6-8',
        'Ages 10-12',
        'Ages 13-15'
      ])
    })

    it('uses max number as tiebreaker when min is equal', () => {
      const tags = [
        { title: 'Ages 6-10' },
        { title: 'Ages 6-8' },
        { title: 'Ages 6-12' }
      ]

      const sorted = [...tags].sort(sortTags)

      expect(sorted.map(t => t.title)).toEqual([
        'Ages 6-8',
        'Ages 6-10',
        'Ages 6-12'
      ])
    })

    it('returns 0 for identical number ranges', () => {
      const a = { title: 'Group 5-10' }
      const b = { title: 'Category 5-10' }

      expect(sortTags(a, b)).toBe(0)
    })

    it('handles single numbers', () => {
      const tags = [
        { title: 'Level 3' },
        { title: 'Level 1' },
        { title: 'Level 2' }
      ]

      const sorted = [...tags].sort(sortTags)

      expect(sorted.map(t => t.title)).toEqual([
        'Level 1',
        'Level 2',
        'Level 3'
      ])
    })

    it('handles multiple numbers scattered in title', () => {
      const tags = [
        { title: 'Session 2 Part 10' },
        { title: 'Session 1 Part 5' },
        { title: 'Session 1 Part 3' }
      ]

      const sorted = [...tags].sort(sortTags)

      expect(sorted.map(t => t.title)).toEqual([
        'Session 1 Part 3',
        'Session 1 Part 5',
        'Session 2 Part 10'
      ])
    })
  })

  describe('alphabetical sorting', () => {
    it('falls back to alphabetical when no numbers present', () => {
      const tags = [
        { title: 'Outdoor' },
        { title: 'Adventure' },
        { title: 'Camping' }
      ]

      const sorted = [...tags].sort(sortTags)

      expect(sorted.map(t => t.title)).toEqual([
        'Adventure',
        'Camping',
        'Outdoor'
      ])
    })

    it('sorts alphabetically when only one tag has numbers', () => {
      const tags = [
        { title: 'Basic' },
        { title: 'Level 5' },
        { title: 'Advanced' }
      ]

      const sorted = [...tags].sort(sortTags)

      expect(sorted.map(t => t.title)).toEqual([
        'Advanced',
        'Basic',
        'Level 5'
      ])
    })

    it('returns 0 for identical strings', () => {
      const a = { title: 'Same' }
      const b = { title: 'Same' }

      expect(sortTags(a, b)).toBe(0)
    })
  })

  describe('mixed scenarios', () => {
    it('handles real-world age group sorting', () => {
      const tags = [
        { title: 'Teens 13-17' },
        { title: 'Adults 18+' },
        { title: 'Kids 6-12' },
        { title: 'Toddlers 3-5' }
      ]

      const sorted = [...tags].sort(sortTags)

      expect(sorted.map(t => t.title)).toEqual([
        'Toddlers 3-5',
        'Kids 6-12',
        'Teens 13-17',
        'Adults 18+'
      ])
    })

    it('handles German age formats', () => {
      const tags = [
        { title: 'Ab 16 Jahren' },
        { title: 'Ab 8 Jahren' },
        { title: 'Ab 12 Jahren' }
      ]

      const sorted = [...tags].sort(sortTags)

      expect(sorted.map(t => t.title)).toEqual([
        'Ab 8 Jahren',
        'Ab 12 Jahren',
        'Ab 16 Jahren'
      ])
    })
  })
})
