import { createRouter, createWebHistory } from 'vue-router'
import { useGoogleAuth } from '@/composables/google_auth'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/login',
      name: 'login',
      component: () => import('../views/LoginView.vue')
    },
    {
      path: '/projects',
      name: 'projects',
      component: () => import('../views/projects/ProjectListView.vue')
    },
    {
      path: '/project/new',
      name: 'new-project',
      component: () => import('../views/projects/NewProjectView.vue')
    },
    {
      path: '/project/import',
      name: 'import-project',
      component: () => import('../views/projects/ImportProjectView.vue')
    },
    {
      path: '/project/import-text',
      name: 'import-text',
      component: () => import('../views/projects/ImportProjectView.vue'),
      props: { defaultMode: 'text' }
    },
    {
      path: '/projects/:id',
      name: 'view-project',
      component: () => import('../views/projects/ViewProjectView.vue')
    },
    {
      path: '/projects/:id/edit',
      name: 'edit-project',
      component: () => import('../views/projects/EditProjectView.vue')
    },
    {
      path: '/projects/:id/tag',
      name: 'tag-project',
      component: () => import('../views/projects/TagProjectView.vue')
    },
    {
      path: '/projects/:id/events',
      name: 'events',
      component: () => import('../views/projects/EventsView.vue')
    },
    {
      path: '/projects/:project_id/events/:id/occurrences',
      name: 'edit-occurrences',
      component: () => import('../views/events/EditOccurrencesView.vue')
    },
    {
      path: '/',
      name: 'categories',
      component: () => import('../views/categories/CategoryListView.vue')
    },
    {
      path: '/categories/new',
      name: 'new-category',
      component: () => import('../views/categories/NewCategoryView.vue')
    },
    {
      path: '/categories/:id',
      name: 'view-category',
      component: () => import('../views/categories/ViewCategoryView.vue')
    },
    {
      path: '/categories/:id/edit',
      name: 'edit-category',
      component: () => import('../views/categories/EditCategoryView.vue')
    },
    {
      path: '/categories/:categoryID/tags/new',
      name: 'new-tag',
      component: () => import('../views/tags/NewTagView.vue')
    },
    {
      path: '/categories/:categoryID/tags/:id',
      name: 'view-tag',
      component: () => import('../views/tags/ViewTagView.vue')
    },
    {
      path: '/categories/:categoryID/tags/:id/edit',
      name: 'edit-tag',
      component: () => import('../views/tags/EditTagView.vue')
    }
  ]
})

router.beforeEach((to, from, next) => {
  const { isLoggedIn } = useGoogleAuth()
  
  if (to.name !== 'login' && !isLoggedIn.value) {
    next({ name: 'login' })
  } else {
    next()
  }
})

export default router
