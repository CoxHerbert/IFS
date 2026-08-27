// @ts-nocheck
import request from '@/utils/request'

export function listNotification(query) {
  return request({
    url: '/system/notification/list',
    method: 'get',
    params: query
  })
}

export function getUnreadNotificationCount() {
  return request({
    url: '/system/notification/unread-count',
    method: 'get'
  })
}

export function readNotification(notificationId) {
  return request({
    url: '/system/notification/' + notificationId + '/read',
    method: 'put'
  })
}

export function readAllNotification() {
  return request({
    url: '/system/notification/read-all',
    method: 'put'
  })
}

export function delNotification(notificationId) {
  return request({
    url: '/system/notification/' + notificationId,
    method: 'delete'
  })
}
