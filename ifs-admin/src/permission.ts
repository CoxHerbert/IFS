// @ts-nocheck
import { message } from 'ant-design-vue'
import NProgress from 'nprogress'
import router from './router'
import store from './store'
import { getToken } from '@/utils/auth'
import { isHttp } from '@/utils/validate'
import 'nprogress/nprogress.css'

NProgress.configure({ showSpinner: false })

const whiteList = ['/login', '/auth-redirect', '/bind', '/register']

router.beforeEach((to, from, next) => {
  NProgress.start()

  if (getToken()) {
    if (to.meta.title) {
      store.dispatch('settings/setTitle', to.meta.title)
    }

    if (to.path === '/login') {
      next({ path: '/' })
      NProgress.done()
      return
    }

    if (store.getters.roles.length === 0) {
      store
        .dispatch('GetInfo')
        .then(() => store.dispatch('GenerateRoutes'))
        .then((accessRoutes) => {
          accessRoutes.forEach((route) => {
            if (!isHttp(route.path)) {
              router.addRoute(route)
            }
          })
          next({ ...to, replace: true })
        })
        .catch((err) => {
          store.dispatch('LogOut').then(() => {
            message.error(String(err))
            next({ path: '/' })
          })
        })

      return
    }

    next()
    return
  }

  if (whiteList.includes(to.path)) {
    next()
  } else {
    next(`/login?redirect=${to.fullPath}`)
    NProgress.done()
  }
})

router.afterEach(() => {
  NProgress.done()
})


