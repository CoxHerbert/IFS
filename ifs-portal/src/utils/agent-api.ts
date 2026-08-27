const AGENT_API_PREFIX = (import.meta.env.VITE_AGENT_API_PREFIX || '/agent-api').replace(/\/$/, '')

export function agentApiUrl(path: string) {
  return `${AGENT_API_PREFIX}${path.startsWith('/') ? path : `/${path}`}`
}

export function resolveAgentApiUrl(path?: string) {
  if (!path) {
    return agentApiUrl('/agent/form/submit')
  }
  if (/^https?:\/\//i.test(path)) {
    return path
  }
  return agentApiUrl(path.replace(/^\/agent-api/, ''))
}
