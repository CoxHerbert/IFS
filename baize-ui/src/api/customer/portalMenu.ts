// @ts-nocheck
import request from '@/utils/request'

export function listPortalMenu(query) {
  return request({
    url: '/customer/portal/menu/list',
    method: 'get',
    params: query
  })
}

export function getPortalMenu(menuId) {
  return request({
    url: '/customer/portal/menu/' + menuId,
    method: 'get'
  })
}

export function addPortalMenu(data) {
  return request({
    url: '/customer/portal/menu',
    method: 'post',
    data
  })
}

export function updatePortalMenu(data) {
  return request({
    url: '/customer/portal/menu',
    method: 'put',
    data
  })
}

export function delPortalMenu(menuId) {
  return request({
    url: '/customer/portal/menu/' + menuId,
    method: 'delete'
  })
}


