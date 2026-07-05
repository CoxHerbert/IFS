export interface ApiResponse<T = unknown> {
  code: number
  msg: string
  data?: T
}

export interface ShipmentStatusStep {
  value: string
  label: string
  active: boolean
}

export interface ShipmentPlan {
  shipmentId: string
  shipmentNo: string
  orderNo: string
  customerName: string
  salesUserId?: string
  salesUserName?: string
  pol: string
  pod: string
  plannedEtd: string
  plannedEta: string
  actualEtd: string
  actualEta: string
  status: string
  totalWeight: number
  totalVolume: number
  totalCartons: number
  remark: string
  createTime?: string
  updateTime?: string
}

export interface ShipmentCargo {
  cargoName: string
  sku: string
  cartons: number
  weightKg: number
  volumeCbm: number
}

export interface ShipmentContainer {
  containerType: string
  quantity: number
  loadRate: number
  remark: string
}

export interface ShipmentOrder {
  orderNo: string
  status: string
}

export interface ShipmentDetail {
  plan: ShipmentPlan
  cargoList: ShipmentCargo[]
  containers: ShipmentContainer[]
  order?: ShipmentOrder
  statusFlow: ShipmentStatusStep[]
}

export async function getPortalShipmentShare(token: string): Promise<ApiResponse<ShipmentDetail>> {
  const response = await fetch(`/portal/shipment/share/${token}`)

  if (!response.ok) {
    throw new Error('网络请求失败')
  }

  return response.json()
}

export async function listWorkspaceShipments(query: {
  pageNum: number
  pageSize: number
  shipmentNo?: string
  status?: string
}, token: string): Promise<ApiResponse<{ rows: ShipmentPlan[]; total: number }>> {
  const params = new URLSearchParams()
  Object.entries(query).forEach(([key, value]) => {
    if (value !== undefined && value !== '') {
      params.set(key, String(value))
    }
  })
  const response = await fetch(`/portal/customer/shipments?${params.toString()}`, {
    headers: { Authorization: `Bearer ${token}` },
  })
  if (!response.ok) {
    throw new Error('网络请求失败')
  }
  return response.json()
}

export async function getWorkspaceShipmentDetail(shipmentId: string, token: string): Promise<ApiResponse<ShipmentDetail>> {
  const response = await fetch(`/portal/customer/shipment/${shipmentId}`, {
    headers: { Authorization: `Bearer ${token}` },
  })
  if (!response.ok) {
    throw new Error('网络请求失败')
  }
  return response.json()
}
