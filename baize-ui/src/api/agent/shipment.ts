// @ts-nocheck
import request from '@/utils/agent-request'

export function analyzeShipmentInChat(sessionId, file, modelName = 'qwen2.5:7b') {
  const data = new FormData()
  data.append('file', file)
  data.append('modelName', modelName)
  return request({
    url: '/chat/session/' + sessionId + '/shipment-analyze',
    method: 'post',
    data
  })
}


