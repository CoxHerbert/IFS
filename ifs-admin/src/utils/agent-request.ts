// @ts-nocheck
import request from '@/utils/request'

const baseURL = import.meta.env.VITE_AGENT_API_PREFIX || '/agent-api'

export default function agentRequest(config) {
  return request({
    baseURL,
    ...config
  })
}
