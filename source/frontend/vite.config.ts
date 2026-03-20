import { fileURLToPath, URL } from 'node:url'
import { defineConfig, loadEnv } from 'vite'
import vue from '@vitejs/plugin-vue'
import vuetify from 'vite-plugin-vuetify'

// https://vite.dev/config/
export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, process.cwd(), '')
  const apiGatewayUrl = env.VITE_API_GATEWAY_URL || 'http://localhost:8080'

  return {
    plugins: [vue(), vuetify({ autoImport: true })],
    server: {
      proxy: {
        '/api': {
          target: apiGatewayUrl,
          changeOrigin: true,
        },
        '/swagger': {
          target: apiGatewayUrl,
          changeOrigin: true,
        },
      },
    },
    resolve: {
      alias: {
        '~': fileURLToPath(new URL('./src', import.meta.url)),
        '@shared': fileURLToPath(new URL('./src/shared', import.meta.url)),
        '~/features': fileURLToPath(new URL('./src/features', import.meta.url)),
      },
    },
  }
})
