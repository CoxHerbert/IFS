export type ThemeName = 'lavender' | 'mint' | 'ocean' | 'rose' | 'peach' | 'amber' | 'forest' | 'teal' | 'slate' | 'graphite'

export const themes: { name: ThemeName; label: string; primary: string; soft: string; page: string; gradientStart: string; gradientEnd: string }[] = [
  { name: 'lavender', label: '薰衣草紫', primary: '#6556b2', soft: '#efedf8', page: '#f6f5fa', gradientStart: '#332e5d', gradientEnd: '#6b5cb5' },
  { name: 'mint', label: '薄荷绿', primary: '#27866b', soft: '#e9f5ef', page: '#f3f8f5', gradientStart: '#174d42', gradientEnd: '#48a88b' },
  { name: 'ocean', label: '海盐蓝', primary: '#3975b8', soft: '#eaf2fb', page: '#f3f7fb', gradientStart: '#23466f', gradientEnd: '#5b95d3' },
  { name: 'rose', label: '樱花粉', primary: '#c45f83', soft: '#f9eaf0', page: '#fcf5f8', gradientStart: '#783d57', gradientEnd: '#df8eaa' },
  { name: 'peach', label: '蜜桃橘', primary: '#c66c48', soft: '#faeee8', page: '#fcf6f2', gradientStart: '#79412f', gradientEnd: '#df956d' },
  { name: 'amber', label: '暖阳黄', primary: '#a87621', soft: '#f7f0dc', page: '#fbf8ee', gradientStart: '#644916', gradientEnd: '#d0a642' },
  { name: 'forest', label: '森林绿', primary: '#4c7652', soft: '#eaf1e8', page: '#f4f7f1', gradientStart: '#29432e', gradientEnd: '#79a276' },
  { name: 'teal', label: '湖水青', primary: '#27818a', soft: '#e5f2f2', page: '#f1f8f7', gradientStart: '#174d55', gradientEnd: '#54a9aa' },
  { name: 'slate', label: '雾霾蓝', primary: '#596f91', soft: '#e9eef5', page: '#f4f6f9', gradientStart: '#35465e', gradientEnd: '#8298b8' },
  { name: 'graphite', label: '石墨灰', primary: '#62636b', soft: '#ededf0', page: '#f5f5f6', gradientStart: '#3a3b42', gradientEnd: '#888990' },
]

const THEME_KEY = 'ifs-ledger-theme'

export function getTheme(): ThemeName {
  try {
    const value = uni.getStorageSync(THEME_KEY) as ThemeName
    return themes.some(theme => theme.name === value) ? value : 'lavender'
  } catch { return 'lavender' }
}

export function setTheme(name: ThemeName) {
  uni.setStorageSync(THEME_KEY, name)
  applyTheme(name)
}

export function applyTheme(name: ThemeName = getTheme()) {
  const theme = themes.find(item => item.name === name) || themes[0]
  // #ifdef H5
  document.querySelectorAll('html, body, page, uni-page-body, .uni-page-body, .uni-page-wrapper').forEach(node => {
    const element = node as HTMLElement
    element.style.setProperty('--theme-primary', theme.primary)
    element.style.setProperty('--theme-soft', theme.soft)
    element.style.setProperty('--theme-page', theme.page)
    element.style.setProperty('--theme-gradient-start', theme.gradientStart)
    element.style.setProperty('--theme-gradient-end', theme.gradientEnd)
    element.style.setProperty('--wot-blue-6', theme.primary)
    element.style.setProperty('--wot-primary-6', theme.primary)
    element.style.setProperty('--wot-button-primary-bg', theme.primary)
    element.style.setProperty('--wot-button-primary-bg-active', theme.primary)
    element.style.setProperty('--wot-button-primary-color', theme.primary)
    element.style.setProperty('--wot-button-primary-color-active', theme.primary)
    element.style.setProperty('--wot-button-primary-plain-border', theme.primary)
    element.style.setProperty('--wot-button-primary-plain-border-active', theme.primary)
    element.style.setProperty('--wot-switch-color-active-bg', theme.primary)
    element.style.setProperty('--wot-switch-color-action-icon', theme.primary)
    element.style.backgroundColor = theme.page
  })
  document.querySelectorAll('.uni-page-head').forEach(header => {
    (header as HTMLElement).style.backgroundColor = theme.page
  })
  // #endif
}
