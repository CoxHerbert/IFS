import { defineConfig, loadEnv } from 'vite';
import { fileURLToPath, URL } from 'node:url';
import createVitePlugins from './vite/plugins';

export default defineConfig(({ mode, command }) => {
    const env = loadEnv(mode, process.cwd(), '');
    const adminApiPrefix = env.VITE_APP_BASE_API || '/admin-api';
    const agentApiPrefix = env.VITE_AGENT_API_PREFIX || '/agent-api';

    return {
        plugins: createVitePlugins(env, command === 'build'),
        resolve: {
            alias: {
                '~': fileURLToPath(new URL('./', import.meta.url)),
                '@': fileURLToPath(new URL('./src', import.meta.url)),
            },
            extensions: ['.mjs', '.js', '.ts', '.jsx', '.tsx', '.json', '.vue'],
        },
        server: {
            port: 8081,
            open: true,
            proxy: {
                [adminApiPrefix]: {
                    target: 'http://127.0.0.1:8080',
                    changeOrigin: true,
                    rewrite: (path) => path.replace(new RegExp(`^${adminApiPrefix}`), ''),
                },
                [agentApiPrefix]: {
                    target: 'http://127.0.0.1:8080',
                    changeOrigin: true,
                },
                '/profile': {
                    target: 'http://127.0.0.1:8080',
                    changeOrigin: true,
                },
            },
        },
    };
});
