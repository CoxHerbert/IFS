export interface RawCargoRow {
  [key: string]: unknown
}

export interface StandardCargoItem {
  rowIndex: number
  sku?: string
  productName?: string
  qty: number
  length: number
  width: number
  height: number
  weight?: number
  cbm: number
  raw: RawCargoRow
}

export interface ShipmentSummary {
  totalQty: number
  totalCBM: number
  containerSuggestion: string
  cargoList: StandardCargoItem[]
}
