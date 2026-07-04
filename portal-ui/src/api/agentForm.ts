import type { AgentResult } from '@/types/agent'
import { getWorkspaceToken } from '@/api/workspace/auth'

export async function submitAgentForm(payload: {
  sessionId: number
  formCode: string
  values: Record<string, unknown>
  submitApi?: string
}): Promise<AgentResult> {
  const response = await fetch(payload.submitApi || '/api/agent/form/submit', {
    method: 'POST',
    headers: buildHeaders(),
    body: JSON.stringify({
      sessionId: payload.sessionId,
      formCode: payload.formCode,
      values: payload.values,
    }),
  })
  return parseAgentResult(response)
}

export async function executeAgentAction(payload: {
  sessionId: number
  actionCode: string
  payload?: Record<string, unknown>
}): Promise<AgentResult> {
  const response = await fetch('/api/agent/action/execute', {
    method: 'POST',
    headers: buildHeaders(),
    body: JSON.stringify(payload),
  })
  return parseAgentResult(response)
}

function buildHeaders(): HeadersInit {
  const headers: Record<string, string> = { 'Content-Type': 'application/json' }
  const token = getWorkspaceToken()
  if (token) {
    headers.Authorization = `Bearer ${token}`
  }
  return headers
}

async function parseAgentResult(response: Response): Promise<AgentResult> {
  const data = await response.json()
  if (!response.ok) {
    throw new Error(data?.summary || data?.message || '请求失败')
  }
  return data as AgentResult
}
