import { createRouter, createWebHistory } from 'vue-router'
import { portalRoutes } from './modules/portal'
import { workspaceBaseRoutes } from './modules/workspace'
import { ensureWorkspaceRoutes, resetWorkspaceRouteState, resolveWorkspaceEntryPath } from './workspace-runtime'
import { getWorkspaceToken, setWorkspaceProfileCache, setWorkspaceRoutesCache } from '@/api/workspace/auth'

const router = createRouter({
  history: createWebHistory(),
  routes: [...portalRoutes, ...workspaceBaseRoutes],
  scrollBehavior() {
    return { top: 0 }
  },
})

router.beforeEach(async (to, _from, next) => {
  const isWorkspaceRoute = to.path === '/customer' || to.path.startsWith('/customer/')
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

    if (to.matched.length === 1) {
      next({ path: to.fullPath, replace: true })
      return
    }

    next()
  } catch (_error) {
    setWorkspaceRoutesCache(null)
    setWorkspaceProfileCache(null)
    resetWorkspaceRouteState()
    next('/customer-login')
  }
})

export default router
