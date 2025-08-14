import { createRouter, createWebHistory } from 'vue-router'
import { useQuery, provideApolloClient } from '@vue/apollo-composable'
import { watch } from 'vue'
import gql from 'graphql-tag'
import WRPage from '../components/WRPage.vue'
import { apolloClientPages } from '../apollo'

const ROUTING_QUERY = gql`
  query Routing {
    pages {
      documentId
      url
    }
  }
`

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: []
})

const { result, loading } = provideApolloClient(apolloClientPages)(() => 
  useQuery(ROUTING_QUERY)
)

const updateRoutes = () => {
  if (result.value && result.value.pages) {
    const newRoutes = result.value.pages.map(page => ({
      path: page.url,
      name: page.url === '/' ? '$home' : page.url.slice(1),
      component: WRPage,
      props: { documentId: page.documentId }
    }))

    newRoutes.forEach(route => router.addRoute(route))
    router.replace(router.currentRoute.value.fullPath)
  }
}

router.isReady().then(() => {
  if (!loading.value) {
    updateRoutes()
  } else {
    const unwatch = watch(loading, (newValue) => {
      if (!newValue) {
        updateRoutes()
        unwatch()
      }
    })
  }
})

export default router
