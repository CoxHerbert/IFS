// @ts-nocheck
import request from '@/utils/agent-request'

export function submitAgentForm(data) {
  return request({
    url: '/agent/form/submit',
    method: 'post',
    data
  })
}

export function executeAgentAction(data) {
  return request({
    url: '/agent/action/execute',
    method: 'post',
    data
  })
}


