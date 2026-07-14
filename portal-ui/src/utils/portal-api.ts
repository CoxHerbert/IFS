const PORTAL_API_PREFIX = (import.meta.env.VITE_PORTAL_API_PREFIX || '/portal-api').replace(/\/$/, '')

export function portalApiUrl(path: string) {
  return `${PORTAL_API_PREFIX}${path.startsWith('/') ? path : `/${path}`}`
}
