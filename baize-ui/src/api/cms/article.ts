// @ts-nocheck
import request from '@/utils/request'

export function listArticle(query) {
  return request({ url: '/cms/article/list', method: 'get', params: query })
}

export function getArticle(articleId) {
  return request({ url: '/cms/article/' + articleId, method: 'get' })
}

export function addArticle(data) {
  return request({ url: '/cms/article', method: 'post', data })
}

export function updateArticle(data) {
  return request({ url: '/cms/article', method: 'put', data })
}

export function delArticle(articleId) {
  return request({ url: '/cms/article/' + articleId, method: 'delete' })
}
