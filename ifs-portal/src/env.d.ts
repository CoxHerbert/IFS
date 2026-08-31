/// <reference types="vite/client" />

import 'vue-router'

interface ImportMetaEnv {
  readonly VITE_PORTAL_API_PREFIX?: string
  readonly VITE_AMAP_KEY?: string
  readonly VITE_AMAP_SECURITY_JS_CODE?: string
}

interface ImportMeta {
  readonly env: ImportMetaEnv
}

declare module 'vue-router' {
  interface RouteMeta {
    seoTitle?: string
    seoDescription?: string
    seoType?: 'website' | 'article'
    noIndex?: boolean
    seoEn?: { title: string; description: string }
    seoZh?: { title: string; description: string }
  }
}

declare module '*.vue' {
  import type { DefineComponent } from 'vue'

  const component: DefineComponent<object, object, unknown>
  export default component
}
