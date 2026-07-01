import request from '@/utils/request'

export function listAccount(query) {
  return request({
    url: '/customer/account/list',
    method: 'get',
    params: query
  })
}

export function getAccount(accountId) {
  return request({
    url: '/customer/account/' + accountId,
    method: 'get'
  })
}

export function addAccount(data) {
  return request({
    url: '/customer/account',
    method: 'post',
    data
  })
}

export function updateAccount(data) {
  return request({
    url: '/customer/account',
    method: 'put',
    data
  })
}

export function resetAccountPwd(accountId, password) {
  return request({
    url: '/customer/account/' + accountId + '/resetPwd',
    method: 'put',
    data: { password }
  })
}

export function updateAccountRoles(accountId, roleIds) {
  return request({
    url: '/customer/account/' + accountId + '/roles',
    method: 'put',
    data: { roleIds }
  })
}

export function delAccount(accountId) {
  return request({
    url: '/customer/account/' + accountId,
    method: 'delete'
  })
}
