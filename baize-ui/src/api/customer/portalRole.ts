// @ts-nocheck
import request from '@/utils/request'

export function listPortalRole(query) {
  return request({
    url: '/customer/portal/role/list',
    method: 'get',
    params: query
  })
}

export function getPortalRole(roleId) {
  return request({
    url: '/customer/portal/role/' + roleId,
    method: 'get'
  })
}

export function addPortalRole(data) {
  return request({
    url: '/customer/portal/role',
    method: 'post',
    data
  })
}

export function updatePortalRole(data) {
  return request({
    url: '/customer/portal/role',
    method: 'put',
    data
  })
}

export function changePortalRoleStatus(roleId, status) {
  return request({
    url: '/customer/portal/role/changeStatus',
    method: 'put',
    data: { roleId, status }
  })
}

export function delPortalRole(roleId) {
  return request({
    url: '/customer/portal/role/' + roleId,
    method: 'delete'
  })
}

export function rolePortalMenuTreeselect(roleId) {
  return request({
    url: '/customer/portal/role/roleMenuTreeselect/' + roleId,
    method: 'get'
  })
}

export function listPortalRoleOptions() {
  return request({
    url: '/customer/portal/role/options',
    method: 'get'
  })
}


