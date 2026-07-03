import { getWorkspaceToken, type ApiResponse } from '@/api/workspace/auth'

export interface ShipmentAssistantRow {
  sku: string
  cargoName: string
  packageType: string
  quantity: number
  cartons: number
  weightKg: number
  volumeCbm: number
  lengthCm: number
  widthCm: number
  heightCm: number
}

export interface ShipmentAssistantSummary {
  lineCount: number
  totalQuantity: number
  totalCartons: number
  totalWeight: number
  totalVolume: number
}

export interface ShipmentAssistantContainer {
  containerType: string
  quantity: number
  maxVolume: number
  maxWeight: number
  usedVolume: number
  usedWeight: number
  loadRate: number
  remark: string
}

export interface ShipmentAssistantLcl {
  recommended: boolean
  totalVolume: number
  remark: string
}

export interface ShipmentAssistantResult {
  summary: ShipmentAssistantSummary
  normalizedCargoList: ShipmentAssistantRow[]
  containers: ShipmentAssistantContainer[]
  lcl: ShipmentAssistantLcl
}

export interface ShipmentAssistantEstimateRequest {
  preferredType: string
  cargoList: ShipmentAssistantRow[]
}

export async function estimateWorkspaceShipment(
  payload: ShipmentAssistantEstimateRequest,
): Promise<ApiResponse<ShipmentAssistantResult>> {
  const response = await fetch('/portal/customer/shipment-assistant/estimate', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Bearer ${getWorkspaceToken() || ''}`,
    },
    body: JSON.stringify(payload),
  })

  if (!response.ok) {
    throw new Error('网络请求失败')
  }

  return response.json()
}
