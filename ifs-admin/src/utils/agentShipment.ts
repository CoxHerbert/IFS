// @ts-nocheck
import * as XLSX from 'xlsx'

const aliases = {
  sku: ['sku', '货号', '产品编码', '型号'],
  productName: ['品名', '产品名称', '货物名称', 'name', 'productname'],
  qty: ['箱数', '数量', 'qty', 'cartons', 'ctns'],
  dimension: ['尺寸', '规格', '长宽高', 'lwh', 'dimension', 'dimensions'],
  length: ['长', '长度', 'length', 'l'],
  width: ['宽', '宽度', 'width', 'w'],
  height: ['高', '高度', 'height', 'h'],
  weight: ['重量', '毛重', 'grossweight', 'gw', 'weight']
}

export async function parseExcelFile(file) {
  const buffer = await file.arrayBuffer()
  const workbook = XLSX.read(buffer, { type: 'array' })
  const sheetName = workbook.SheetNames[0]
  if (!sheetName) return []
  return XLSX.utils.sheet_to_json(workbook.Sheets[sheetName], { defval: '' })
}

export function normalizeShipmentRows(rows) {
  const cargoList = rows.map((row, index) => normalizeRow(row, index + 1)).filter(Boolean)
  const totalQty = cargoList.reduce((sum, item) => sum + item.qty, 0)
  const totalCBM = round(cargoList.reduce((sum, item) => sum + item.cbm, 0), 3)

  return {
    totalQty,
    totalCBM,
    containerSuggestion: planContainer(totalCBM),
    cargoList
  }
}

function normalizeRow(row, rowIndex) {
  const qty = toNumber(readField(row, 'qty'))
  const weight = toNumber(readField(row, 'weight'))
  const lengthField = toNumber(readField(row, 'length'))
  const widthField = toNumber(readField(row, 'width'))
  const heightField = toNumber(readField(row, 'height'))
  const dimension =
    lengthField > 0 && widthField > 0 && heightField > 0
      ? { length: lengthField, width: widthField, height: heightField }
      : parseDimensions(readField(row, 'dimension'))

  if (!dimension || qty <= 0) return null

  return {
    rowIndex,
    sku: stringValue(readField(row, 'sku')),
    productName: stringValue(readField(row, 'productName')),
    qty,
    length: dimension.length,
    width: dimension.width,
    height: dimension.height,
    weight,
    cbm: calculateCBM(dimension.length, dimension.width, dimension.height, qty),
    raw: row
  }
}

function parseDimensions(value) {
  const text = String(value || '')
    .trim()
    .replace(/厘米|公分/gi, 'cm')
    .replace(/[×X]/g, 'x')
    .replace(/\s+/g, ' ')

  const match = text.match(
    /(\d+(?:\.\d+)?)\s*(?:cm)?\s*(?:[*x]|\s)\s*(\d+(?:\.\d+)?)\s*(?:cm)?\s*(?:[*x]|\s)\s*(\d+(?:\.\d+)?)\s*(?:cm)?/i
  )

  if (!match) return null

  const length = Number(match[1])
  const width = Number(match[2])
  const height = Number(match[3])

  if (![length, width, height].every((item) => Number.isFinite(item) && item > 0)) {
    return null
  }

  return { length, width, height }
}

function readField(row, field) {
  const aliasSet = aliases[field].map(normalizeKey)
  const hit = Object.entries(row).find(([key]) => aliasSet.includes(normalizeKey(key)))
  return hit?.[1]
}

function normalizeKey(value) {
  return String(value || '')
    .trim()
    .toLowerCase()
    .replace(/\s+/g, '')
}

function stringValue(value) {
  return String(value || '').trim()
}

function toNumber(value) {
  const num = Number(String(value || '').replace(/,/g, '').trim())
  return Number.isFinite(num) ? num : 0
}

function calculateCBM(length, width, height, qty) {
  return round((length * width * height * qty) / 1000000, 3)
}

function planContainer(totalCBM) {
  if (totalCBM <= 28) return '1 x 20GP'
  if (totalCBM <= 58) return '1 x 40GP'
  if (totalCBM <= 68) return '1 x 40HQ'
  return `${Math.ceil(totalCBM / 68)} x 40HQ`
}

function round(value, precision = 2) {
  const factor = 10 ** precision
  return Math.round(value * factor) / factor
}
