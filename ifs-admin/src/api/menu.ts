// @ts-nocheck
import request from '@/utils/request'

// 获取路由
export const getRouters = () => {
  return request({
    url: '/getRouters',
    method: 'get',
    params: { _t: Date.now() },
    headers: { 'Cache-Control': 'no-cache' }
  })
}

