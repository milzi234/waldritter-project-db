import { fileURLToPath, URL } from 'node:url'

import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))
    }
  },
  server: {
    https: {
      key: `${process.env.PROJECT_DIR}/local/certs/waldritter.${process.env.PROJECT_TLD}+1-key.pem`,
      cert: `${process.env.PROJECT_DIR}/local/certs/waldritter.${process.env.PROJECT_TLD}+1.pem`
    }
  }
})
