// @ts-nocheck
import request from '@/utils/request'

export const listLedgerAccounts = query => request({ url: '/ledger/admin/accounts', method: 'get', params: query })
export const getLedgerAccountRecords = (accountId, query) => request({ url: `/ledger/admin/accounts/${accountId}/records`, method: 'get', params: query })
export const deleteLedgerAccountRecord = (accountId, recordId) => request({ url: `/ledger/admin/accounts/${accountId}/records/${encodeURIComponent(recordId)}`, method: 'delete' })
