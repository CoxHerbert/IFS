// @ts-nocheck
import request from '@/utils/request'

export function listCustomer(query) {
  return request({
    url: '/customer/list',
    method: 'get',
    params: query
  })
}

export function customerOptions(query) {
  return request({
    url: '/customer/options',
    method: 'get',
    params: query
  })
}

export function getCustomer(customerId) {
  return request({
    url: '/customer/' + customerId,
    method: 'get'
  })
}

export function addCustomer(data) {
  return request({
    url: '/customer',
    method: 'post',
    data
  })
}

export function updateCustomer(data) {
  return request({
    url: '/customer',
    method: 'put',
    data
  })
}

export function delCustomer(customerId) {
  return request({
    url: '/customer/' + customerId,
    method: 'delete'
  })
}

export function listCustomerContact(customerId, query) {
  return request({
    url: '/customer/' + customerId + '/contact/list',
    method: 'get',
    params: query
  })
}

export function getCustomerContact(contactId) {
  return request({
    url: '/customer/contact/' + contactId,
    method: 'get'
  })
}

export function addCustomerContact(data) {
  return request({
    url: '/customer/contact',
    method: 'post',
    data
  })
}

export function updateCustomerContact(data) {
  return request({
    url: '/customer/contact',
    method: 'put',
    data
  })
}

export function delCustomerContact(contactId) {
  return request({
    url: '/customer/contact/' + contactId,
    method: 'delete'
  })
}


