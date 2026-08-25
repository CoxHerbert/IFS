import type { RouteRecordRaw, Router } from 'vue-router'
import {
  getWorkspaceRouters,
  getWorkspaceRoutesCache,
  getWorkspaceToken,
  setWorkspaceRoutesCache,
  type WorkspaceRouteItem,
} from '@/api/workspace/auth'

const workspaceViewModules = import.meta.glob([
  '../views/workspace/**/*.vue',
  '!../views/workspace/Login/**/*.vue',
])

function normalizeWorkspaceComponentPath(component: string) {
  const normalized = component
    .trim()
    .replace(/\\/g, '/')
    .replace(/^@\//, '')
    .replace(/^\/?src\//, '')
    .replace(/^views\//, '')
    .replace(/^\/+/, '')
    .replace(/\.vue$/, '')
  const componentPath = normalized
  return componentPath.startsWith('workspace/') ? normalized : `workspace/${normalized}`
}

function resolveWorkspaceComponent(component?: string): RouteRecordRaw['component'] | undefined {
  if (!component) {
    return undefined
  }
  const componentPath = normalizeWorkspaceComponentPath(component)
  if (!componentPath.startsWith('workspace/')) {
    return undefined
  }
  return workspaceViewModules[`../views/${componentPath}.vue`] as RouteRecordRaw['component'] | undefined
}

let workspaceRoutesLoaded = false
let workspaceRoutesLoadedForToken = ''

export function resetWorkspaceRouteState() {
  workspaceRoutesLoaded = false
  workspaceRoutesLoadedForToken = ''
}

function normalizeWorkspaceRouteItems(items: WorkspaceRouteItem[]): WorkspaceRouteItem[] {
  return items.map((item) => {
    const normalized: WorkspaceRouteItem = {
      ...item,
      meta: {
        ...item.meta,
      },
    }
    if (item.children?.length) {
      normalized.children = normalizeWorkspaceRouteItems(item.children)
    }
    return normalized
  })
}

function flattenWorkspaceRoutes(router: Router, items: WorkspaceRouteItem[], parentPath = ''): RouteRecordRaw[] {
  const result: RouteRecordRaw[] = []
  for (const item of items) {
    const currentPath = parentPath ? `${parentPath}/${item.path}` : item.path
    const component = resolveWorkspaceComponent(item.component)
    const isStaticDashboard = currentPath.replace(/^\/+|\/+$/g, '') === 'dashboard'
      && router.hasRoute('workspace-dashboard')
    if (component && !isStaticDashboard && !router.hasRoute(item.name)) {
      const route: RouteRecordRaw = {
        path: currentPath,
        name: item.name,
        component,
        meta: {
          requiresWorkspaceAuth: true,
          title: item.meta.title,
          icon: item.meta.icon,
          menuId: item.meta.menuId,
          noCache: item.meta.noCache,
          hiddenMenu: item.hidden,
        },
      }
      result.push(route)
    }
    if (item.children?.length) {
      result.push(...flattenWorkspaceRoutes(router, item.children, currentPath))
    }
  }
  return result
}

function registerWorkspaceRoutes(router: Router, items: WorkspaceRouteItem[]) {
  const dynamicRoutes = flattenWorkspaceRoutes(router, items)
  dynamicRoutes.forEach((route) => {
    if (!router.hasRoute(String(route.name))) {
      router.addRoute('workspace-root', route)
    }
  })
}

export function restoreWorkspaceRoutesFromCache(router: Router) {
  const cachedRoutes = getWorkspaceRoutesCache()
  if (cachedRoutes?.length) {
    registerWorkspaceRoutes(router, normalizeWorkspaceRouteItems(cachedRoutes))
  }
}

export function resolveWorkspaceEntryPath(items: WorkspaceRouteItem[], parentPath = ''): string {
  for (const item of items) {
    const currentPath = parentPath ? `${parentPath}/${item.path}` : item.path
    if (item.component && !item.hidden) {
      return `/customer/${currentPath}`
    }
    if (item.children?.length) {
      const childPath = resolveWorkspaceEntryPath(item.children, currentPath)
      if (childPath) {
        return childPath
      }
    }
  }
  return '/customer/workspace'
}

export async function ensureWorkspaceRoutes(router: Router): Promise<WorkspaceRouteItem[]> {
  const token = getWorkspaceToken() || ''
  if (workspaceRoutesLoaded && workspaceRoutesLoadedForToken === token) {
    return getWorkspaceRoutesCache() || []
  }

  const cachedRoutes = getWorkspaceRoutesCache()
  if (cachedRoutes?.length) {
    registerWorkspaceRoutes(router, normalizeWorkspaceRouteItems(cachedRoutes))
  }

  let routeItems = normalizeWorkspaceRouteItems(cachedRoutes || [])
  try {
    const response = await getWorkspaceRouters()
    if (response.code === 200 && Array.isArray(response.data)) {
      routeItems = normalizeWorkspaceRouteItems(response.data)
      setWorkspaceRoutesCache(routeItems)
      registerWorkspaceRoutes(router, routeItems)
    }
  } catch (_error) {
    routeItems = normalizeWorkspaceRouteItems(cachedRoutes || [])
  }

  workspaceRoutesLoaded = true
  workspaceRoutesLoadedForToken = token
  return routeItems
}
