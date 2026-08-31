import { createRouter, createWebHistory } from 'vue-router'
import { portalRoutes } from './modules/portal'
import { workspaceBaseRoutes } from './modules/workspace'
import {
  ensureWorkspaceRoutes,
  resetWorkspaceRouteState,
  resolveWorkspaceEntryPath,
  restoreWorkspaceRoutesFromCache,
} from './workspace-runtime'
import { getWorkspaceToken, setWorkspaceProfileCache, setWorkspaceRoutesCache } from '@/api/workspace/auth'
import { updateRouteSeo } from '@/utils/seo'
import { initialPortalLocale, normalizeLocale, setPortalLocale } from '@/i18n'

const router = createRouter({
  history: createWebHistory(),
  routes: [...portalRoutes, ...workspaceBaseRoutes],
  scrollBehavior() {
    return { top: 0 }
  },
})

// Restore persisted routes before app.use(router) triggers initial URL resolution.
restoreWorkspaceRoutesFromCache(router)

router.beforeEach(async (to, _from, next) => {
  const isWorkspaceRoute = to.path === '/customer' || to.path.startsWith('/customer/')
  if (!isWorkspaceRoute && to.path !== '/customer-login') {
    const pathLocale = normalizeLocale(to.params.locale)
    if (pathLocale) {
      setPortalLocale(pathLocale)
    } else if (!to.path.startsWith('/shipment/share/')) {
      next({ path: `/${initialPortalLocale()}${to.fullPath === '/' ? '' : to.fullPath}`, replace: true })
      return
    }
  }
  if (!isWorkspaceRoute) {
    next()
    return
  }

  if (!getWorkspaceToken()) {
    next('/customer-login')
    return
  }

  try {
    const routeItems = await ensureWorkspaceRoutes(router)
    if (to.path === '/customer' || to.path === '/customer/') {
      next(resolveWorkspaceEntryPath(routeItems))
      return
    }

    if (to.name === 'not-found') {
      const resolvedRoute = router.resolve(to.fullPath)
      if (resolvedRoute.name !== 'not-found') {
        next({ path: to.fullPath, replace: true })
        return
      }
    }

    next()
  } catch (_error) {
    setWorkspaceRoutesCache(null)
    setWorkspaceProfileCache(null)
    resetWorkspaceRouteState()
    next('/customer-login')
  }
})

router.afterEach((to) => updateRouteSeo(to))

export default router
