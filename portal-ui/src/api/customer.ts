export interface ApiResponse<T = unknown> {
  code: number
  msg: string
  data?: T
}

export interface CustomerAccount {
  accountId: string
  customerId: string
  customerNo: string
  customerName: string
  companyName: string
  username: string
  realName: string
  phone: string
  email: string
  isMain: string
  status: string
  lastLoginTime?: string
}

export interface LoginResult {
  token: string
  user: CustomerAccount
}

export const CUSTOMER_TOKEN_KEY = 'portal_customer_token'

export function getCustomerToken() {
  return localStorage.getItem(CUSTOMER_TOKEN_KEY)
}

export function setCustomerToken(token: string) {
  localStorage.setItem(CUSTOMER_TOKEN_KEY, token)
}

export function removeCustomerToken() {
  localStorage.removeItem(CUSTOMER_TOKEN_KEY)
}

export async function customerLogin(username: string, password: string): Promise<ApiResponse<LoginResult>> {
  const response = await fetch('/portal/customer/login', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ username, password })
  })

  if (!response.ok) {
    throw new Error('网络请求失败')
  }

  return response.json()
}

export async function getCustomerProfile(): Promise<ApiResponse<CustomerAccount>> {
  const response = await fetch('/portal/customer/profile', {
    headers: {
      Authorization: `Bearer ${getCustomerToken() || ''}`
    }
  })

  if (!response.ok) {
    throw new Error('网络请求失败')
  }

  return response.json()
}
