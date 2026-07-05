import type { ChatMessage, ChatSession, SendMessageResponse } from '@/types/agent'
import { getWorkspaceToken } from '@/api/workspace/auth'

export interface AgentModelOption {
  label: string
  value: string
  description: string
  default: boolean
}

export async function listAgentModels(): Promise<AgentModelOption[]> {
  const response = await fetch('/api/chat/models', {
    headers: authHeaders(),
  })
  return parseResponse(response)
}

export async function createChatSession(payload: {
  title?: string
  modelName?: string
}): Promise<ChatSession> {
  const response = await fetch('/api/chat/session', {
    method: 'POST',
    headers: jsonHeaders(),
    body: JSON.stringify(payload),
  })
  return parseResponse(response)
}

export async function listChatSessions(): Promise<ChatSession[]> {
  const response = await fetch('/api/chat/sessions', {
    headers: authHeaders(),
  })
  return parseResponse(response)
}

export async function listChatMessages(sessionId: number): Promise<ChatMessage[]> {
  const response = await fetch(`/api/chat/session/${sessionId}/messages`, {
    headers: authHeaders(),
  })
  return parseResponse(response)
}

export async function updateChatSessionTitle(sessionId: number, title: string): Promise<void> {
  const response = await fetch(`/api/chat/session/${sessionId}/title`, {
    method: 'PUT',
    headers: jsonHeaders(),
    body: JSON.stringify({ title }),
  })
  await parseResponse(response)
}

export async function deleteChatSession(sessionId: number): Promise<void> {
  const response = await fetch(`/api/chat/session/${sessionId}`, {
    method: 'DELETE',
    headers: authHeaders(),
  })
  await parseResponse(response)
}

export async function sendChatMessage(payload: {
  sessionId: number
  message: string
  modelName?: string
}): Promise<SendMessageResponse> {
  const response = await fetch('/api/chat/send', {
    method: 'POST',
    headers: jsonHeaders(),
    body: JSON.stringify(payload),
  })
  return parseResponse(response)
}

export async function analyzeShipmentInChat(payload: {
  sessionId: number
  file: File
  modelName?: string
}): Promise<SendMessageResponse> {
  const formData = new FormData()
  formData.append('file', payload.file)
  formData.append('modelName', payload.modelName || 'qwen2.5:7b')
  const response = await fetch(`/api/chat/session/${payload.sessionId}/shipment-analyze`, {
    method: 'POST',
    headers: authHeaders(),
    body: formData,
  })
  return parseResponse(response)
}

function authHeaders(): HeadersInit {
  const token = getWorkspaceToken()
  return token ? { Authorization: `Bearer ${token}` } : {}
}

function jsonHeaders(): HeadersInit {
  return {
    'Content-Type': 'application/json',
    ...authHeaders(),
  }
}

async function parseResponse<T>(response: Response): Promise<T> {
  const data = await response.json()
  if (!response.ok) {
    throw new Error(data?.message || 'Request failed')
  }
  return data as T
}
