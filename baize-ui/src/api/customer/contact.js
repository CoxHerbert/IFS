import request from '@/utils/request'

// 查询官网线索列表
export function listContact(query) {
  return request({
    url: '/portal/contact/list',
    method: 'get',
    params: query
  })
}

// 查询官网线索详细
export function getContact(contactId) {
  return request({
    url: '/portal/contact/' + contactId,
    method: 'get'
  })
}

// 修改官网线索
export function updateContact(data) {
  return request({
    url: '/portal/contact',
    method: 'put',
    data: data
  })
}

// 删除官网线索
export function delContact(contactId) {
  return request({
    url: '/portal/contact/' + contactId,
    method: 'delete'
  })
}
