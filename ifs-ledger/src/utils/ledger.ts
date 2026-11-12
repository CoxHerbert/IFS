export type LedgerRecord = {
  id: string
  date: string
  merchant: string
  category: string
  amount: number
  type: 'expense' | 'income'
  source: 'manual' | 'wechat' | 'alipay'
  importBatchId?: string
  importFileName?: string
}

const STORAGE_KEY = 'ifs-ledger-records'
const ACCOUNT_KEY = 'ifs-ledger-account'
const CLOUD_SYNC_KEY = 'ifs-ledger-cloud-sync'
const PROFILE_KEY = 'ifs-ledger-profile'
const IMPORT_BATCH_KEY = 'ifs-ledger-import-batches'

export type LocalAccount = { accountId: number; displayName: string; loginId: string; token: string; uid: string }
export type LocalProfile = { uid: string; nickname: string; phone: string; email: string }
export type LedgerImportBatch = { batchId: string; source: 'wechat' | 'alipay'; fileName: string; importedAt: string; importedCount: number; duplicateCount: number }

type ApiEnvelope<T> = { code: number; msg: string; data: T }
type AuthResult = { accountId: number; uid: string; loginId: string; token: string; cloudSyncEnabled: boolean }
const API_BASE = import.meta.env.VITE_LEDGER_API_BASE || '/ledger-api'

function request<T>(path: string, method: 'GET' | 'POST' | 'PUT' | 'DELETE', data?: unknown, token?: string): Promise<T> {
  return new Promise((resolve, reject) => {
    uni.request({
      url: `${API_BASE}${path}`,
      method,
      data,
      header: token ? { Authorization: `Bearer ${token}` } : {},
      success: response => {
        const body = response.data as ApiEnvelope<T>
        if (body?.code !== 200) reject(new Error(body?.msg || '请求失败'))
        else resolve(body.data)
      },
      fail: error => reject(new Error(error.errMsg || '网络连接失败')),
    })
  })
}

export async function authenticate(loginId: string, password: string, register: boolean) {
  const result = await request<AuthResult>(`/ledger/auth/${register ? 'register' : 'login'}`, 'POST', { loginId, password })
  const account = { accountId: result.accountId, loginId: result.loginId, displayName: result.loginId.split('@')[0] || result.loginId, token: result.token, uid: result.uid || accountUID(result.accountId) }
  setAccount(account)
  setCloudSyncEnabled(result.cloudSyncEnabled)
  return { account, cloudSyncEnabled: result.cloudSyncEnabled }
}

export function accountUID(accountId: number | string) {
  return `QJ-${Number(accountId).toString(36).toUpperCase()}`
}

function createLocalUID() {
  const randomPart = () => Math.random().toString(36).slice(2, 8).toUpperCase()
  return `QJ-L-${Date.now().toString(36).toUpperCase()}-${randomPart()}${randomPart()}`
}

export function getLocalProfile(): LocalProfile {
  try {
    const saved = uni.getStorageSync(PROFILE_KEY) as Partial<LocalProfile> | null
    if (saved?.uid) return { uid: saved.uid, nickname: saved.nickname || '', phone: saved.phone || '', email: saved.email || '' }
  } catch { /* Fall through and initialize a device profile. */ }
  const profile: LocalProfile = { uid: createLocalUID(), nickname: '', phone: '', email: '' }
  setLocalProfile(profile)
  return profile
}

export function setLocalProfile(profile: LocalProfile) {
  uni.setStorageSync(PROFILE_KEY, profile)
}

export async function updateCloudConsent(enabled: boolean, account = getAccount()) {
  if (!account?.token) throw new Error('请先登录')
  const result = await request<{ cloudSyncEnabled: boolean }>('/ledger/sync-consent', 'PUT', { enabled }, account.token)
  setCloudSyncEnabled(result.cloudSyncEnabled)
  return result.cloudSyncEnabled
}

export async function uploadCloudRecords(account = getAccount()) {
  if (!account?.token) throw new Error('请先登录')
  return request<{ uploaded: number }>('/ledger/records', 'PUT', { records: getRecords() }, account.token)
}

export async function fetchCloudRecords(account = getAccount()) {
  if (!account?.token) throw new Error('请先登录')
  return request<{ records: LedgerRecord[]; cloudSyncEnabled: boolean }>('/ledger/records', 'GET', undefined, account.token)
}

export async function deleteCloudImportBatch(batchId: string, account = getAccount()) {
  if (!account?.token) throw new Error('请先登录')
  return request<{ deleted: number }>(`/ledger/records/batches/${encodeURIComponent(batchId)}`, 'DELETE', undefined, account.token)
}

export function getAccount(): LocalAccount | null {
  try {
    const saved = uni.getStorageSync(ACCOUNT_KEY) as Partial<LocalAccount> | null
    const account = saved as LocalAccount | null
    if (!account?.token) return null
    if (!account.uid || !account.accountId) {
      const accountId = Number(account.accountId) || readAccountIdFromToken(account.token)
      if (accountId > 0) {
        const migrated = { ...account, accountId, uid: accountUID(accountId) }
        setAccount(migrated)
        return migrated
      }
    }
    return account
  } catch { return null }
}

function readAccountIdFromToken(token: string) {
  try {
    const encoded = token.split('.')[1].replace(/-/g, '+').replace(/_/g, '/')
    const alphabet = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/'
    let buffer = 0
    let bits = 0
    let payload = ''
    for (const char of encoded) {
      if (char === '=') break
      const value = alphabet.indexOf(char)
      if (value < 0) continue
      buffer = (buffer << 6) | value
      bits += 6
      if (bits >= 8) {
        bits -= 8
        payload += String.fromCharCode((buffer >> bits) & 0xff)
      }
    }
    return Number(JSON.parse(payload).accountId) || 0
  } catch { return 0 }
}

export function setAccount(account: LocalAccount | null) {
  if (account) uni.setStorageSync(ACCOUNT_KEY, account)
  else uni.removeStorageSync(ACCOUNT_KEY)
}

export function getCloudSyncEnabled(): boolean {
  try { return uni.getStorageSync(CLOUD_SYNC_KEY) === true } catch { return false }
}

export function setCloudSyncEnabled(enabled: boolean) {
  uni.setStorageSync(CLOUD_SYNC_KEY, enabled)
}

export function getRecords(): LedgerRecord[] {
  try { return uni.getStorageSync(STORAGE_KEY) || [] } catch { return [] }
}

export function saveRecords(records: LedgerRecord[]) {
  uni.setStorageSync(STORAGE_KEY, records)
  if (getCloudSyncEnabled() && getAccount()?.token) {
    void uploadCloudRecords().catch(error => console.warn('Ledger cloud sync failed:', error))
  }
}

export function getImportBatches(): LedgerImportBatch[] {
  try { return uni.getStorageSync(IMPORT_BATCH_KEY) || [] } catch { return [] }
}

export function saveImportBatch(records: LedgerRecord[], source: 'wechat' | 'alipay', fileName: string, batchId: string): LedgerImportBatch {
  const current = getRecords()
  const existingIds = new Set(current.map(item => item.id))
  const seen = new Set<string>()
  const unique = records.filter(record => {
    if (existingIds.has(record.id) || seen.has(record.id)) return false
    seen.add(record.id)
    return true
  }).map(record => ({ ...record, source, importBatchId: batchId, importFileName: fileName }))
  const batch: LedgerImportBatch = {
    batchId,
    source,
    fileName,
    importedAt: new Date().toISOString(),
    importedCount: unique.length,
    duplicateCount: records.length - unique.length,
  }
  saveRecords([...unique, ...current].sort((a, b) => b.date.localeCompare(a.date)))
  uni.setStorageSync(IMPORT_BATCH_KEY, [batch, ...getImportBatches()].slice(0, 100))
  return batch
}

export function deleteImportBatch(batchId: string) {
  const nextRecords = getRecords().filter(record => record.importBatchId !== batchId)
  const nextBatches = getImportBatches().filter(batch => batch.batchId !== batchId)
  uni.setStorageSync(IMPORT_BATCH_KEY, nextBatches)
  saveRecords(nextRecords)
  return nextRecords.length
}

export async function removeImportBatch(batchId: string) {
  if (getCloudSyncEnabled() && getAccount()?.token) await deleteCloudImportBatch(batchId)
  return deleteImportBatch(batchId)
}

export function appendRecords(records: LedgerRecord[]) {
  const current = getRecords()
  const seen = new Set(current.map(item => item.id))
  const unique = records.filter(item => !seen.has(item.id))
  saveRecords([...unique, ...current].sort((a, b) => b.date.localeCompare(a.date)))
  return unique.length
}
