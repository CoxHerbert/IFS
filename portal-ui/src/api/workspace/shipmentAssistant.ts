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
  internalLengthCm: number
  internalWidthCm: number
  internalHeightCm: number
  safeVolume: number
  effectiveVolume: number
  riskLevel: string
  unitCost: number
  extraFees: number
  totalCost: number
  warnings: string[]
  placements: LoadingPlacement[]
}

export interface ShipmentAssistantLcl {
  recommended: boolean
  totalVolume: number
  ratePerCbm: number
  minCharge: number
  extraFees: number
  totalCost: number
  remark: string
}

export interface ShipmentRecommendation {
  mode: string
  title: string
  reason: string
  saving: number
  riskLevel: string
  confidence: string
}

export interface LoadingPlacement {
  cargoName: string
  sku: string
  color: string
  quantity: number
  x: number
  y: number
  z: number
  length: number
  width: number
  height: number
  remark: string
}

export interface LoadingPlan {
  containerType: string
  quantity: number
  internalLengthCm: number
  internalWidthCm: number
  internalHeightCm: number
  utilization: number
  placements: LoadingPlacement[]
}

export interface ShipmentAssistantResult {
  summary: ShipmentAssistantSummary
  normalizedCargoList: ShipmentAssistantRow[]
  containers: ShipmentAssistantContainer[]
  lcl: ShipmentAssistantLcl
  recommendation?: ShipmentRecommendation
  loadingPlan?: LoadingPlan
  warnings: string[]
}

export interface ShipmentAssistantEstimateRequest {
  preferredType: string
  lclRate: number
  lclMinCharge: number
  rate20GP: number
  rate40GP: number
  rate40HQ: number
  extraFees: number
  cargoList: ShipmentAssistantRow[]
}

export interface ShipmentPlanCreateRequest {
  orderNo: string
  pol: string
  pod: string
  plannedEtd: string
  plannedEta: string
  remark: string
  preferredType: string
  cargoList: ShipmentAssistantRow[]
}

export interface ShipmentPlanCreated {
  plan: {
    shipmentId: string
    shipmentNo: string
    orderNo: string
    status: string
    paymentStatus: string
    paymentAmount: number
    totalWeight: number
    totalVolume: number
    totalCartons: number
  }
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

export async function createWorkspaceShipmentPlan(
  payload: ShipmentPlanCreateRequest,
): Promise<ApiResponse<ShipmentPlanCreated>> {
  const response = await fetch('/portal/customer/shipment-assistant/plan', {
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
