import { ref } from 'vue'
import { defineStore } from 'pinia'
import useProjectAPI from '../api/Project'


const projectAPI = useProjectAPI()
export const useProjectStore = defineStore('projects', () => {
  const projects = ref([])

  const reload = async () => {
    projects.value = await projectAPI.getAll()
  }

  projectAPI.on('created', reload)
  projectAPI.on('updated', reload)
  projectAPI.on('deleted', reload)

  reload()

  return { projects, reload }
})
