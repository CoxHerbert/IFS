import type { RawCargoRow, ShipmentSummary, StandardCargoItem } from '@/types/shipment'
import { calculateCBM, round } from './cbmCalculator'
import { planContainer } from './containerPlanner'
import { parseDimensions } from './dimensionParser'

type FieldName = 'sku' | 'productName' | 'qty' | 'dimension' | 'length' | 'width' | 'height' | 'weight'

const aliases: Record<FieldName, string[]> = {
  sku: ['sku', '货号', '产品编码', '型号'],
  productName: ['品名', '产品名称', '货物名称', 'name', 'productname'],
  qty: ['箱数', '数量', 'qty', 'cartons', 'ctns'],
  dimension: ['尺寸', '规格', '长宽高', 'lwh', 'dimension', 'dimensions'],
  length: ['长', '长度', 'length', 'l'],
  width: ['宽', '宽度', 'width', 'w'],
  height: ['高', '高度', 'height', 'h'],
  weight: ['重量', '毛重', 'grossweight', 'gw', 'weight'],
}

export function normalizeShipmentRows(rows: RawCargoRow[]): ShipmentSummary {
  const cargoList = rows
    .map((row, index) => normalizeRow(row, index + 1))
    .filter((item): item is StandardCargoItem => Boolean(item))

  const totalQty = cargoList.reduce((sum, item) => sum + item.qty, 0)
  const totalCBM = round(cargoList.reduce((sum, item) => sum + item.cbm, 0), 3)
  return {
    totalQty,
    totalCBM,
    containerSuggestion: planContainer(totalCBM),
    cargoList,
  }
}

function normalizeRow(row: RawCargoRow, rowIndex: number): StandardCargoItem | null {
  const sku = stringValue(readField(row, 'sku'))
  const productName = stringValue(readField(row, 'productName'))
  const qty = toNumber(readField(row, 'qty'))
  const weight = toNumber(readField(row, 'weight'))

  const lengthField = toNumber(readField(row, 'length'))
  const widthField = toNumber(readField(row, 'width'))
  const heightField = toNumber(readField(row, 'height'))
  const dimension = lengthField > 0 && widthField > 0 && heightField > 0
    ? { length: lengthField, width: widthField, height: heightField, unit: 'cm' as const }
    : parseDimensions(readField(row, 'dimension'))

  if (!dimension || qty <= 0) {
    return null
  }

  return {
    rowIndex,
    sku,
    productName,
    qty,
    length: dimension.length,
    width: dimension.width,
    height: dimension.height,
    weight,
    cbm: calculateCBM(dimension.length, dimension.width, dimension.height, qty),
    raw: row,
  }
}

function readField(row: RawCargoRow, field: FieldName) {
  const entries = Object.entries(row)
  const aliasSet = aliases[field].map(normalizeKey)
  const hit = entries.find(([key]) => aliasSet.includes(normalizeKey(key)))
  return hit?.[1]
}

function normalizeKey(value: string) {
  return value.trim().toLowerCase().replace(/\s+/g, '')
}

function stringValue(value: unknown) {
  return String(value ?? '').trim()
}

function toNumber(value: unknown) {
  const num = Number(String(value ?? '').replace(/,/g, '').trim())
  return Number.isFinite(num) ? num : 0
}
