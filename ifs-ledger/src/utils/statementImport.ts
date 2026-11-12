import type { LedgerRecord } from './ledger'

type ImportSource = 'wechat' | 'alipay'
type ParseResult = { records: LedgerRecord[]; skippedRows: number }
type SheetRow = unknown[]
type XLSXModule = typeof import('xlsx')
let xlsxModule: XLSXModule | null = null
async function getXLSX(): Promise<XLSXModule> { if (!xlsxModule) xlsxModule = await import('xlsx'); return xlsxModule }

const aliases = {
  date: ['交易时间', '交易创建时间', '创建时间', '日期', '付款时间'],
  flow: ['收/支', '收支', '交易方向', '资金流向'],
  amount: ['金额(元)', '金额', '交易金额'],
  merchant: ['交易对方', '商家名称', '商户名称', '商品说明', '商品'],
  category: ['交易类型', '交易分类', '类别'],
  status: ['当前状态', '交易状态', '状态'],
  transaction: ['交易单号', '交易订单号', '订单号', '交易号'],
}

const text = (value: unknown) => String(value ?? '').trim()
const normalized = (value: unknown) => text(value).replace(/\s+/g, '').toLowerCase()

function findColumn(header: SheetRow, names: string[]) {
  const wanted = names.map(normalized)
  return header.findIndex(cell => wanted.includes(normalized(cell)))
}

function dateText(value: unknown, xlsx: XLSXModule) {
  if (value instanceof Date && !Number.isNaN(value.getTime())) {
    const pad = (n: number) => String(n).padStart(2, '0')
    return `${value.getFullYear()}-${pad(value.getMonth() + 1)}-${pad(value.getDate())} ${pad(value.getHours())}:${pad(value.getMinutes())}:${pad(value.getSeconds())}`
  }
  if (typeof value === 'number') {
    const date = XLSX.SSF.parse_date_code(value)
    if (date) return `${date.y}-${String(date.m).padStart(2, '0')}-${String(date.d).padStart(2, '0')} ${String(date.H).padStart(2, '0')}:${String(date.M).padStart(2, '0')}:${String(Math.floor(date.S)).padStart(2, '0')}`
  }
  return text(value).replace(/年|\//g, '-').replace(/月/g, '-').replace(/日/g, '').trim()
}

function stableHash(value: string) {
  let hash = 2166136261
  for (let i = 0; i < value.length; i += 1) hash = Math.imul(hash ^ value.charCodeAt(i), 16777619)
  return (hash >>> 0).toString(36).toUpperCase()
}

export async function parseLedgerStatement(data: ArrayBuffer | Uint8Array, source: ImportSource): Promise<ParseResult> {
  const XLSX = await getXLSX()
  const workbook = XLSX.read(data, { type: 'array', cellDates: true })
  const sheet = workbook.Sheets[workbook.SheetNames[0]]
  if (!sheet) throw new Error('账单文件中没有可读取的工作表')
  const rows = XLSX.utils.sheet_to_json<SheetRow>(sheet, { header: 1, defval: null, blankrows: false })
  const headerIndex = rows.findIndex(row => findColumn(row, aliases.date) >= 0 && findColumn(row, aliases.amount) >= 0 && findColumn(row, aliases.flow) >= 0)
  if (headerIndex < 0) throw new Error('未识别到账单表头，请选择微信或支付宝导出的原始账单')
  const header = rows[headerIndex]
  const columns = Object.fromEntries(Object.entries(aliases).map(([key, names]) => [key, findColumn(header, names)])) as Record<keyof typeof aliases, number>
  const records: LedgerRecord[] = []
  let skippedRows = 0
  for (const [index, row] of rows.slice(headerIndex + 1).entries()) {
    const amountText = text(row[columns.amount]).replace(/[￥¥,\s]/g, '')
    const amount = Math.abs(Number(amountText))
    const date = dateText(row[columns.date], XLSX)
    if (!amountText || !Number.isFinite(amount) || amount <= 0 || !date) { skippedRows += 1; continue }
    const status = columns.status >= 0 ? text(row[columns.status]) : ''
    if (/失败|已取消|已撤销|已关闭/.test(status)) { skippedRows += 1; continue }
    const flow = text(row[columns.flow])
    const type = /收入|收/.test(flow) ? 'income' : /支出|付/.test(flow) ? 'expense' : Number(amountText) < 0 ? 'expense' : ''
    if (!type) { skippedRows += 1; continue }
    const transactionId = columns.transaction >= 0 ? text(row[columns.transaction]) : ''
    const merchant = columns.merchant >= 0 ? text(row[columns.merchant]) : ''
    const category = columns.category >= 0 ? text(row[columns.category]) : ''
    const rowKey = [date, flow, merchant, category, amountText, transactionId, index].join('|')
    records.push({ id: `${source === 'wechat' ? 'wx' : 'ali'}-${transactionId || stableHash(rowKey)}`, date, merchant: merchant || category || '未命名交易', category: category || flow, amount, type, source })
  }
  if (!records.length) throw new Error('文件中没有可导入的有效账单记录')
  return { records, skippedRows }
}

export function createImportBatchId(source: ImportSource) {
  const now = new Date()
  const pad = (n: number) => String(n).padStart(2, '0')
  const timestamp = `${now.getFullYear()}${pad(now.getMonth() + 1)}${pad(now.getDate())}-${pad(now.getHours())}${pad(now.getMinutes())}${pad(now.getSeconds())}`
  return `${source === 'wechat' ? 'WX' : 'ALI'}-${timestamp}-${Math.random().toString(36).slice(2, 6).toUpperCase()}`
}
