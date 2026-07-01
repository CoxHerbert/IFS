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
  shipmentNo: string
  orderNo: string
  customerName: string
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
