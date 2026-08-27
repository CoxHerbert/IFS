// @ts-nocheck
import request from '@/utils/request'

// 查询缓存详情
export function getCache() {
  return request({
    url: '/monitor/cache',
    method: 'get'
  })
}

