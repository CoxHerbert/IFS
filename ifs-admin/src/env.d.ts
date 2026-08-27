/// <reference types="vite/client" />

interface ImportMetaEnv {
  readonly VITE_APP_BASE_API?: string
  readonly VITE_AGENT_API_PREFIX?: string
  readonly VITE_PORTAL_BASE_URL?: string
}

interface ImportMeta {
  readonly env: ImportMetaEnv
}

declare module '*.vue' {
  import type { DefineComponent } from 'vue'

  const component: DefineComponent<Record<string, unknown>, Record<string, unknown>, unknown>
  export default component
}
