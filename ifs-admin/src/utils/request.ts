// @ts-nocheck
import axios from 'axios'
import { message, Modal, notification } from 'ant-design-vue'
import store from '@/store'
import { getToken } from '@/utils/auth'
import errorCode from '@/utils/errorCode'
import { tansParams, blobValidate } from '@/utils/ruoyi'
import { saveAs } from 'file-saver'

let closeDownloadLoading

axios.defaults.headers['Content-Type'] = 'application/json;charset=utf-8'

const service = axios.create({
  baseURL: import.meta.env.VITE_APP_BASE_API,
  timeout: 10000
})

service.interceptors.request.use(
  (config) => {
    const isToken = config.headers?.isToken === false

    if (getToken() && !isToken) {
      config.headers = config.headers || {}
      config.headers.Authorization = `Bearer ${getToken()}`
    }

    if (config.method === 'get' && config.params) {
      let url = `${config.url}?${tansParams(config.params)}`
      url = url.slice(0, -1)
      config.params = {}
      config.url = url
    }

    return config
  },
  (error) => Promise.reject(error)
)

service.interceptors.response.use(
  (res) => {
    const code = res.data.code || 200
    const msg = errorCode[code] || res.data.msg || errorCode.default

    if (res.request.responseType === 'blob' || res.request.responseType === 'arraybuffer') {
      return res.data
    }

    if (code === 401) {
      Modal.confirm({
        title: '系统提示',
        content: '登录状态已过期，您可以继续留在该页面，或者重新登录。',
        okText: '重新登录',
        cancelText: '取消',
        onOk: () =>
          store.dispatch('LogOut').then(() => {
            location.href = '/index'
          })
      })
      return Promise.reject(new Error(msg))
    }

    if (code === 500) {
      message.error(msg)
      return Promise.reject(new Error(msg))
    }

    if (code !== 200) {
      notification.error({
        message: '请求失败',
        description: msg
      })
      return Promise.reject(new Error(msg))
    }

    return res.data
  },
  (error) => {
    let errorMessage = error.message

    if (errorMessage === 'Network Error') {
      errorMessage = '后端接口连接异常'
    } else if (errorMessage.includes('timeout')) {
      errorMessage = '系统接口请求超时'
    } else if (errorMessage.includes('Request failed with status code')) {
      errorMessage = `系统接口 ${errorMessage.slice(-3)} 异常`
    }

    message.error(errorMessage)
    return Promise.reject(error)
  }
)

export function download(url, params, filename) {
  closeDownloadLoading = message.loading({
    content: '正在下载数据，请稍候',
    duration: 0
  })

  return service
    .post(url, params, {
      transformRequest: [(requestParams) => tansParams(requestParams)],
      headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
      responseType: 'blob'
    })
    .then(async (data) => {
      const isBlob = await blobValidate(data)

      if (isBlob) {
        saveAs(new Blob([data]), filename)
      } else {
        const resText = await data.text()
        const rspObj = JSON.parse(resText)
        const errMsg = errorCode[rspObj.code] || rspObj.msg || errorCode.default
        message.error(errMsg)
      }

      closeDownloadLoading?.()
      closeDownloadLoading = undefined
    })
    .catch((error) => {
      console.error(error)
      message.error('下载文件出现错误，请联系管理员')
      closeDownloadLoading?.()
      closeDownloadLoading = undefined
      return Promise.reject(error)
    })
}

export default service
