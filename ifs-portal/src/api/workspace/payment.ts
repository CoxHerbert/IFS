import { getWorkspaceToken, type ApiResponse } from './auth'
import { portalApiUrl } from '@/utils/portal-api'

export interface WorkspacePaymentItem {
  shipmentId: string
  shipmentNo: string
  orderNo: string
  customerName: string
  payableAmount: number
  paidAmount: number
  unpaidAmount: number
  paymentStatus: 'UNPAID' | 'PARTIAL' | 'PAID'
}

export async function listWorkspacePayments(query: { pageNum: number; pageSize: number; paymentStatus?: string }): Promise<ApiResponse<{ rows: WorkspacePaymentItem[]; total: number }>> {
  const params = new URLSearchParams()
  Object.entries(query).forEach(([key,value]) => value !== undefined && value !== '' && params.set(key,String(value)))
  const response = await fetch(portalApiUrl(`/customer/payments?${params}`), { headers: { Authorization: `Bearer ${getWorkspaceToken() || ''}` } })
  if (!response.ok) throw new Error('网络请求失败')
  return response.json()
}
