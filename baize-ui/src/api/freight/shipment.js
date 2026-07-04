import request from '@/utils/request'

export function listShipment(query) {
  return request({
    url: '/freight/shipment/list',
    method: 'get',
    params: query
  })
}

export function importShipment(data) {
  return request({
    url: '/freight/shipment/import',
    method: 'post',
    data
  })
}

export function getShipment(shipmentId) {
  return request({
    url: '/freight/shipment/' + shipmentId,
    method: 'get'
  })
}

export function updateShipmentStatus(shipmentId, data) {
  return request({
    url: '/freight/shipment/' + shipmentId + '/status',
    method: 'put',
    data
  })
}

export function bindShipmentCustomer(shipmentId, data) {
  return request({
    url: '/freight/shipment/' + shipmentId + '/customer',
    method: 'put',
    data
  })
}

export function confirmShipment(shipmentId) {
  return request({
    url: '/freight/shipment/' + shipmentId + '/confirm',
    method: 'post'
  })
}

export function getShipmentShare(shipmentId) {
  return request({
    url: '/freight/shipment/' + shipmentId + '/share',
    method: 'get'
  })
}

export function delShipment(shipmentId) {
  return request({
    url: '/freight/shipment/' + shipmentId,
    method: 'delete'
  })
}
