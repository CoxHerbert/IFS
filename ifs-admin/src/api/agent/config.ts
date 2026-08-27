// @ts-nocheck
import request from '@/utils/request'

export function getAgentOllamaConfig() {
  return request({
    url: '/agent/config/ollama',
    method: 'get'
  })
}

export function updateAgentOllamaConfig(data) {
  return request({
    url: '/agent/config/ollama',
    method: 'put',
    data
  })
}

export function testAgentOllamaConfig(data) {
  return request({
    url: '/agent/config/ollama/test',
    method: 'post',
    data
  })
}
