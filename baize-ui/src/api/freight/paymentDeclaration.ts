// @ts-nocheck
import request from '@/utils/request'
export const listPaymentDeclaration=query=>request({url:'/freight/payment-declaration/list',method:'get',params:query})
export const getPaymentDeclaration=id=>request({url:'/freight/payment-declaration/'+id,method:'get'})
export const approvePaymentDeclaration=(id,data)=>request({url:'/freight/payment-declaration/'+id+'/approve',method:'post',data})
export const rejectPaymentDeclaration=(id,data)=>request({url:'/freight/payment-declaration/'+id+'/reject',method:'post',data})
