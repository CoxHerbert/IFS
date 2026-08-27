// @ts-nocheck
import store from '@/store'

export default {
  mounted(el, binding) {
    const { value } = binding
    const allPermission = '*:*:*'
    const permissions = store.getters?.permissions || []

    if (Array.isArray(value) && value.length > 0) {
      const hasPermissions = permissions.some((permission) => {
        return allPermission === permission || value.includes(permission)
      })

      if (!hasPermissions) {
        el.parentNode && el.parentNode.removeChild(el)
      }
      return
    }

    throw new Error('请设置操作权限标签值')
  }
}
