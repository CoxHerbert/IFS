// @ts-nocheck
import request from '@/utils/request'

export function submitAgentForm(data) {
  return request({
    url: '/agent/chat/form/submit',
    method: 'post',
    data
  })
}

export function executeAgentAction(data) {
  return request({
    url: '/api/agent/action/execute',
    method: 'post',
    data
  })
}


