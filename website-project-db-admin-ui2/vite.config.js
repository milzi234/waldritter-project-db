import { fileURLToPath, URL } from 'node:url'
import { readFileSync } from 'node:fs'
import { resolve } from 'node:path'

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
      key: readFileSync(resolve(__dirname, 'localhost+2-key.pem')),
      cert: readFileSync(resolve(__dirname, 'localhost+2.pem'))
    },
    port: 5173
  }
})
