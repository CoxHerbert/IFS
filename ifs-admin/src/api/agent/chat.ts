// @ts-nocheck
import request from '@/utils/agent-request'

export function listAgentModels() {
  return request({
    url: '/chat/models',
    method: 'get'
  })
}

export function createChatSession(data) {
  return request({
    url: '/chat/session',
    method: 'post',
    data
  })
}

export function listChatSessions() {
  return request({
    url: '/chat/sessions',
    method: 'get'
  })
}

export function listChatMessages(sessionId) {
  return request({
    url: '/chat/session/' + sessionId + '/messages',
    method: 'get'
  })
}

export function updateChatSessionTitle(sessionId, title) {
  return request({
    url: '/chat/session/' + sessionId + '/title',
    method: 'put',
    data: { title }
  })
}

export function deleteChatSession(sessionId) {
  return request({
    url: '/chat/session/' + sessionId,
    method: 'delete'
  })
}

export function sendChatMessage(data) {
  return request({
    url: '/chat/send',
    method: 'post',
    data
  })
}


