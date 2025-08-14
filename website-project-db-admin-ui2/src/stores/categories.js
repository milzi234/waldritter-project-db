import { ref } from 'vue'
import { defineStore } from 'pinia'
import useCategoryAPI from '../api/Category'


const categoryAPI = useCategoryAPI()
export const useCategoryStore = defineStore('categories', () => {
  const categories = ref([])
  const tagsPerCategory = {}
  const tags = ref([])

  const reload = async () => {
    categories.value = await categoryAPI.getAll()
    categories.value.sort((a, b) => {
      if (a.title < b.title) {
        return -1
      }
      if (b.title < a.title) {
        return 1
      }
      return 0
    })
    tags.value = []
    Object.keys(tagsPerCategory).forEach((categoryID) => {
      if (!categories.value.find(category => category.id === categoryID)) {
        tagsPerCategory[categoryID].value = []
      }
    })
    categories.value.forEach(async (category) => {
      const tagAPI = categoryAPI.tagAPIFor(category.id)
      tagsPerCategory[category.id] = tagsPerCategory[category.id] || ref([])
      tagsPerCategory[category.id].value = await tagAPI.getAll()
      tags.value.push(...tagsPerCategory[category.id].value)
    })
  }

  const find = (id) => {
    const found = categories.value.find((category) => category.id === id)
    return found
  }

  const tagsFor = (categoryID) => {
    tagsPerCategory[categoryID] = tagsPerCategory[categoryID] || ref([])
    return tagsPerCategory[categoryID]
  }

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

  categoryAPI.on('created', reload)
  categoryAPI.on('updated', reload)
  categoryAPI.on('tag-created', reload)
  categoryAPI.on('tag-deleted', reload)
  categoryAPI.on('deleted', reload)

  reload()

  return { categories, tags, tagsFor, sortTags, reload, find }
})
