// @ts-nocheck
import request from '@/utils/agent-request'

export function analyzeShipmentInChat(sessionId, file, modelName = '') {
  const data = new FormData()
  data.append('file', file)
  if (modelName) {
    data.append('modelName', modelName)
  }
  return request({
    url: '/chat/session/' + sessionId + '/shipment-analyze',
    method: 'post',
    data
  })
}


