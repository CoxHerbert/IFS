import type { RouteLocationNormalizedLoaded } from 'vue-router'
import { normalizeLocale } from '@/i18n'

const SITE_NAME = '中美货运转运平台'
const DEFAULT_TITLE = '中美货运转运平台 | 海运、FBA头程、清关与派送'
const DEFAULT_DESCRIPTION = '提供中美线海运、FBA头程、清关与派送服务，在线获取运输方案、物流咨询和报价。'

interface SeoOptions {
  title?: string
  description?: string
  image?: string
  type?: 'website' | 'article'
  noIndex?: boolean
  canonicalPath?: string
}

function setMeta(selector: string, attribute: 'name' | 'property', key: string, content: string) {
  let element = document.head.querySelector<HTMLMetaElement>(selector)
  if (!element) {
    element = document.createElement('meta')
    element.setAttribute(attribute, key)
    document.head.appendChild(element)
  }
  element.content = content
}

function setCanonical(url: string) {
  let link = document.head.querySelector<HTMLLinkElement>('link[rel="canonical"]')
  if (!link) {
    link = document.createElement('link')
    link.rel = 'canonical'
    document.head.appendChild(link)
  }
  link.href = url
}

function setAlternate(language: string, url: string) {
  let link = document.head.querySelector<HTMLLinkElement>(`link[rel="alternate"][hreflang="${language}"]`)
  if (!link) {
    link = document.createElement('link')
    link.rel = 'alternate'
    link.hreflang = language
    document.head.appendChild(link)
  }
  link.href = url
}

export function updateSeo(options: SeoOptions = {}) {
  const title = options.title || DEFAULT_TITLE
  const description = options.description || DEFAULT_DESCRIPTION
  const canonical = new URL(options.canonicalPath || window.location.pathname, window.location.origin).href
  const image = new URL(options.image || '/logo.svg', window.location.origin).href

  document.title = title
  setMeta('meta[name="description"]', 'name', 'description', description)
  setMeta('meta[name="robots"]', 'name', 'robots', options.noIndex ? 'noindex,nofollow' : 'index,follow,max-image-preview:large,max-snippet:-1,max-video-preview:-1')
  setMeta('meta[property="og:site_name"]', 'property', 'og:site_name', SITE_NAME)
  setMeta('meta[property="og:type"]', 'property', 'og:type', options.type || 'website')
  setMeta('meta[property="og:title"]', 'property', 'og:title', title)
  setMeta('meta[property="og:description"]', 'property', 'og:description', description)
  setMeta('meta[property="og:url"]', 'property', 'og:url', canonical)
  setMeta('meta[property="og:image"]', 'property', 'og:image', image)
  setMeta('meta[name="twitter:card"]', 'name', 'twitter:card', 'summary_large_image')
  setMeta('meta[name="twitter:title"]', 'name', 'twitter:title', title)
  setMeta('meta[name="twitter:description"]', 'name', 'twitter:description', description)
  setMeta('meta[name="twitter:image"]', 'name', 'twitter:image', image)
  setCanonical(canonical)
}

export function updateRouteSeo(route: RouteLocationNormalizedLoaded) {
  document.head.querySelectorAll('script[data-seo-schema]').forEach((element) => element.remove())
  const activeLocale = normalizeLocale(route.params.locale) || 'zh-cn'
  const cleanPath = route.path.replace(/^\/(en|zh-cn)(?=\/|$)/, '') || '/'
  const localizedMeta = activeLocale === 'en' ? route.meta.seoEn : route.meta.seoZh
  updateSeo({
    title: (localizedMeta as { title?: string } | undefined)?.title || route.meta.seoTitle as string | undefined,
    description: (localizedMeta as { description?: string } | undefined)?.description || route.meta.seoDescription as string | undefined,
    type: route.meta.seoType as SeoOptions['type'],
    noIndex: Boolean(route.meta.noIndex || route.path.startsWith('/customer')),
  })
  const base = window.location.origin
  setAlternate('en', new URL(`/en${cleanPath === '/' ? '' : cleanPath}`, base).href)
  setAlternate('zh-CN', new URL(`/zh-cn${cleanPath === '/' ? '' : cleanPath}`, base).href)
  setAlternate('x-default', new URL(`/en${cleanPath === '/' ? '' : cleanPath}`, base).href)
}

export function setStructuredData(id: string, data: Record<string, unknown>) {
  let script = document.head.querySelector<HTMLScriptElement>(`script[data-seo-schema="${id}"]`)
  if (!script) {
    script = document.createElement('script')
    script.type = 'application/ld+json'
    script.dataset.seoSchema = id
    document.head.appendChild(script)
  }
  script.textContent = JSON.stringify(data)
}
