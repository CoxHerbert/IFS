// @ts-nocheck
import request from '@/utils/request'

export const listLedgerAccounts = query => request({ url: '/ledger/admin/accounts', method: 'get', params: query })
export const getLedgerAccountRecords = (accountId, query) => request({ url: `/ledger/admin/accounts/${accountId}/records`, method: 'get', params: query })
export const deleteLedgerAccountRecord = (accountId, recordId) => request({ url: `/ledger/admin/accounts/${accountId}/records/${encodeURIComponent(recordId)}`, method: 'delete' })
export const getLedgerAccountImportBatches = (accountId, query) => request({ url: `/ledger/admin/accounts/${accountId}/import-batches`, method: 'get', params: query })
export const deleteLedgerAccountImportBatch = (accountId, batchId) => request({ url: `/ledger/admin/accounts/${accountId}/import-batches/${encodeURIComponent(batchId)}`, method: 'delete' })
