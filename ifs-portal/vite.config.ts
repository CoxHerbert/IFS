import { defineConfig, loadEnv } from 'vite'
import vue from '@vitejs/plugin-vue'
import { fileURLToPath, URL } from 'node:url'

export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, process.cwd(), '')
  const portalApiPrefix = env.VITE_PORTAL_API_PREFIX || '/portal-api'
  const agentApiPrefix = env.VITE_AGENT_API_PREFIX || '/agent-api'

  return {
    plugins: [vue()],
    resolve: {
      alias: {
        '@': fileURLToPath(new URL('./src', import.meta.url))
      }
    },
    server: {
      host: '0.0.0.0',
      port: 5173,
      open: true,
      proxy: {
        [portalApiPrefix]: {
          target: 'http://localhost:8080',
          changeOrigin: true,
          rewrite: (path) => path.replace(new RegExp(`^${portalApiPrefix}`), '/portal')
        },
        [agentApiPrefix]: {
          target: 'http://localhost:8080',
          changeOrigin: true
        },
        '/profile': {
          target: 'http://localhost:8080',
          changeOrigin: true
        }
      },
    }
  }
})
