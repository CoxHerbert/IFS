// @ts-nocheck
import request from '@/utils/request'
export const listReceipt = query => request({ url:'/freight/receipt/list', method:'get', params:query })
export const getReceipt = id => request({ url:'/freight/receipt/'+id, method:'get' })
export const addReceipt = data => request({ url:'/freight/receipt', method:'post', data, headers:{'Content-Type':'multipart/form-data'} })
export const delReceipt = id => request({ url:'/freight/receipt/'+id, method:'delete' })
