// @ts-nocheck
import request from '@/utils/request'

export function simpleCustomerOptions(query) {
  return request({
    url: '/simple/customer/options',
    method: 'get',
    params: query
  })
}
