import './assets/main.css'

import { createPinia } from 'pinia'
import PrimeVue from 'primevue/config';
import Aura from '@primevue/themes/aura';
import { createApp, provide, h } from 'vue'
import { ApolloClients } from '@vue/apollo-composable'
import { usePassThrough } from 'primevue/passthrough';

import App from './App.vue'
import router from './router'
import { apolloClientPages, apolloClientProjects } from './apollo'

const app = createApp({
  setup () {
    provide(ApolloClients, {
      default: apolloClientPages,
      pagesAPI: apolloClientPages,
      projectAPI: apolloClientProjects
    })
  },

  render: () => h(App),
})

app.use(PrimeVue, {
  theme: {
    preset: Aura
  },
  unstyled: false,
  pt: usePassThrough(),
});

app.use(createPinia())
app.use(router)
app.mount('#app')