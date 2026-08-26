// @ts-nocheck
import request from '@/utils/request'

export const listPaymentLedger = query => request({ url: '/freight/payment-ledger/list', method: 'get', params: query })
export const updatePayableAmount = (shipmentId, payableAmount) => request({
  url: `/freight/payment-ledger/${shipmentId}/payable`, method: 'put', data: { payableAmount }
})
export const getShipmentCharges = shipmentId => request({ url: `/freight/payment-ledger/${shipmentId}/charges`, method: 'get' })
export const saveShipmentCharges = (shipmentId, items) => request({ url: `/freight/payment-ledger/${shipmentId}/charges`, method: 'put', data: { items } })
