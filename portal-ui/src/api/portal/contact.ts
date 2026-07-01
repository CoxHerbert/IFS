export interface ContactPayload {
  contactName: string
  companyName?: string
  phone?: string
  email?: string
  route?: string
  cargoInfo?: string
  message: string
  source?: string
}

export interface ApiResponse<T = unknown> {
  code: number
  msg: string
  data?: T
}

export interface ContactSubmitResult {
  leadNo: string
}

export async function submitContact(data: ContactPayload): Promise<ApiResponse<ContactSubmitResult>> {
  const response = await fetch('/portal/contact', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify(data)
  })

  if (!response.ok) {
    throw new Error('网络请求失败')
  }

  return response.json()
}
