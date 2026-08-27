export interface AgentResult {
  version: string
  type: string
  title: string
  summary: string
  blocks: AgentBlock[]
}

export interface AgentBlock {
  type: 'metrics' | 'table' | 'markdown' | 'file' | 'error' | 'form' | 'action' | string
  title?: string
  content?: string
  items?: MetricItem[]
  columns?: TableColumn[]
  data?: Record<string, unknown>[]
  name?: string
  url?: string
  formCode?: string
  submitApi?: string
  fields?: FormField[]
  initialValues?: Record<string, unknown>
  actionCode?: string
  label?: string
  payload?: Record<string, unknown>
}

export interface FormField {
  field: string
  label: string
  component: 'input' | 'textarea' | 'number' | 'select' | 'date' | 'upload' | string
  required?: boolean
  placeholder?: string
  options?: FormOption[]
}

export interface FormOption {
  label: string
  value: unknown
}

export interface MetricItem {
  label: string
  value: unknown
}

export interface TableColumn {
  label: string
  field: string
}

export interface ChatSession {
  id: number
  title: string
  modelName: string
  summary?: string
  updatedAt?: string
}

export interface ChatMessage {
  id: number
  sessionId: number
  role: 'user' | 'assistant' | 'system' | 'tool'
  content: string
  blockResult?: AgentResult
  createdAt: string
}

export interface SendMessageResponse {
  messageId: number
  sessionId: number
  result: AgentResult
}
