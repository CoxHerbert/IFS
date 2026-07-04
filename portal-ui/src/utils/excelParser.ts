import * as XLSX from 'xlsx'
import type { RawCargoRow } from '@/types/shipment'

export async function parseExcelFile(file: File): Promise<RawCargoRow[]> {
  const buffer = await file.arrayBuffer()
  const workbook = XLSX.read(buffer, { type: 'array' })
  const sheetName = workbook.SheetNames[0]
  if (!sheetName) {
    return []
  }
  const sheet = workbook.Sheets[sheetName]
  return XLSX.utils.sheet_to_json<RawCargoRow>(sheet, { defval: '' })
}
