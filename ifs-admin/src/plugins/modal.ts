// @ts-nocheck
import { message, Modal, notification } from 'ant-design-vue'
import { h } from 'vue'
import { ExclamationCircleOutlined } from '@ant-design/icons-vue'

let closeLoading: (() => void) | undefined

function resolveContent(content: string) {
  return typeof content === 'string' ? content : String(content)
}

export default {
  msg(content: string) {
    message.info(resolveContent(content))
  },
  msgError(content: string) {
    message.error(resolveContent(content))
  },
  msgSuccess(content: string) {
    message.success(resolveContent(content))
  },
  msgWarning(content: string) {
    message.warning(resolveContent(content))
  },
  alert(content: string) {
    Modal.info({ title: '系统提示', content: resolveContent(content) })
  },
  alertError(content: string) {
    Modal.error({ title: '系统提示', content: resolveContent(content) })
  },
  alertSuccess(content: string) {
    Modal.success({ title: '系统提示', content: resolveContent(content) })
  },
  alertWarning(content: string) {
    Modal.warning({ title: '系统提示', content: resolveContent(content) })
  },
  notify(content: string) {
    notification.info({ message: '系统提示', description: resolveContent(content) })
  },
  notifyError(content: string) {
    notification.error({ message: '系统提示', description: resolveContent(content) })
  },
  notifySuccess(content: string) {
    notification.success({ message: '系统提示', description: resolveContent(content) })
  },
  notifyWarning(content: string) {
    notification.warning({ message: '系统提示', description: resolveContent(content) })
  },
  confirm(content: string) {
    return new Promise<void>((resolve, reject) => {
      Modal.confirm({
        title: '系统提示',
        icon: h(ExclamationCircleOutlined),
        content: resolveContent(content),
        okText: '确定',
        cancelText: '取消',
        onOk: () => resolve(),
        onCancel: () => reject(new Error('cancel'))
      })
    })
  },
  prompt(content: string) {
    return new Promise<string>((resolve, reject) => {
      const result = window.prompt('系统提示\n\n' + resolveContent(content))
      if (result === null) {
        reject(new Error('cancel'))
      } else {
        resolve(result)
      }
    })
  },
  loading(content: string) {
    closeLoading = message.loading({
      content: resolveContent(content),
      duration: 0
    })
  },
  closeLoading() {
    closeLoading?.()
    closeLoading = undefined
  }
}


