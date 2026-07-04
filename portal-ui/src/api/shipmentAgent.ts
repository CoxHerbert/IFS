import type { AgentResult } from '@/types/agent'

export async function analyzeShipment(file: File): Promise<AgentResult> {
  const formData = new FormData()
  formData.append('file', file)
  formData.append('modelName', 'qwen2.5:7b')
  const response = await fetch('/api/shipment/analyze', {
    method: 'POST',
    body: formData,
  })
  const data = await response.json()
  if (!response.ok) {
    throw new Error(data?.summary || data?.message || '分析失败')
  }
  return data as AgentResult
}
