import { computed, ref } from 'vue'

export type PortalLocale = 'en' | 'zh-cn'

const STORAGE_KEY = 'ifs-portal-locale'
const locale = ref<PortalLocale>('zh-cn')

const messages = {
  en: {
    brandName: 'IFS International Logistics', brandTagline: 'FCL and LCL freight services',
    home: 'Home', news: 'Insights', services: 'Services', route: 'China–USA', about: 'About', contact: 'Contact',
    customerCenter: 'Customer Portal', getQuote: 'Get a Freight Quote',
    footerText: 'Freight quotes, shipping routes, cargo planning and customer collaboration.',
    routeLink: 'China to USA Shipping',
  },
  'zh-cn': {
    brandName: 'IFS 国际物流', brandTagline: '海运整柜与拼箱服务',
    home: '首页', news: '新闻资讯', services: '服务能力', route: '中美航线', about: '关于我们', contact: '联系我们',
    customerCenter: '客户中心', getQuote: '获取报价',
    footerText: '提供报价、航线、出货计划与客户协同入口。',
    routeLink: '中国至美国运输',
  },
} as const

export function normalizeLocale(value: unknown): PortalLocale | undefined {
  const text = String(value || '').toLowerCase()
  return text === 'en' || text === 'zh-cn' ? text : undefined
}

export function setPortalLocale(value: PortalLocale) {
  locale.value = value
  document.documentElement.lang = value === 'en' ? 'en' : 'zh-CN'
  localStorage.setItem(STORAGE_KEY, value)
}

export function initialPortalLocale(): PortalLocale {
  const saved = normalizeLocale(localStorage.getItem(STORAGE_KEY))
  if (saved) return saved
  return navigator.language.toLowerCase().startsWith('zh') ? 'zh-cn' : 'en'
}

export function usePortalI18n() {
  const t = (key: keyof typeof messages.en) => messages[locale.value][key]
  const localePath = (path: string, targetLocale = locale.value) => {
    if (path.startsWith('/customer') || path.startsWith('/shipment/share')) return path
    const cleanPath = path.replace(/^\/(en|zh-cn)(?=\/|$)/, '') || '/'
    return `/${targetLocale}${cleanPath === '/' ? '' : cleanPath}`
  }
  return { locale: computed(() => locale.value), t, localePath, setLocale: setPortalLocale }
}
