import type { RouteRecordRaw, Router } from 'vue-router'
import {
  getWorkspaceRouters,
  getWorkspaceRoutesCache,
  getWorkspaceToken,
  setWorkspaceRoutesCache,
  type WorkspaceRouteItem,
} from '@/api/workspace/auth'

const workspaceComponentMap: Record<string, RouteRecordRaw['component']> = {
  'workspace/dashboard': () => import('@/views/workspace/WorkspaceDashboardView.vue'),
  'workspace/account-profile': () => import('@/views/workspace/WorkspaceAccountProfileView.vue'),
  'workspace/shipment-tracking': () => import('@/views/workspace/WorkspaceShipmentTrackingView.vue'),
  'workspace/shipment-assistant': () => import('@/views/workspace/WorkspaceShipmentAssistantView.vue'),
}

export const defaultWorkspaceRouteItems: WorkspaceRouteItem[] = [
  {
    name: 'workspace-dashboard',
    path: 'workspace',
    component: 'workspace/dashboard',
    meta: { title: '工作台', icon: 'AppstoreOutlined', menuId: '0' },
  },
  {
    name: 'workspace-account-profile',
    path: 'account',
    component: 'workspace/account-profile',
    meta: { title: '账号资料', icon: 'ProfileOutlined', menuId: '0' },
  },
  {
    name: 'workspace-shipment-tracking',
    path: 'shipment',
    component: 'workspace/shipment-tracking',
    meta: { title: '出货查询', icon: 'RadarChartOutlined', menuId: '0' },
  },
  {
    name: 'workspace-shipment-assistant',
    path: 'shipment-assistant',
    component: 'workspace/shipment-assistant',
    meta: { title: '智能出货助手', icon: 'CalculatorOutlined', menuId: '0' },
  },
]

let workspaceRoutesLoaded = false
let workspaceRoutesLoadedForToken = ''

export function resetWorkspaceRouteState() {
  workspaceRoutesLoaded = false
  workspaceRoutesLoadedForToken = ''
}

function flattenWorkspaceRoutes(router: Router, items: WorkspaceRouteItem[], parentPath = ''): RouteRecordRaw[] {
  const result: RouteRecordRaw[] = []
  for (const item of items) {
    const currentPath = parentPath ? `${parentPath}/${item.path}` : item.path
    const component = item.component ? workspaceComponentMap[item.component] : undefined
    if (component && !router.hasRoute(item.name)) {
      const route: RouteRecordRaw = {
        path: currentPath,
        name: item.name,
        component,
        meta: {
          requiresWorkspaceAuth: true,
          title: item.meta.title,
          icon: item.meta.icon,
          menuId: item.meta.menuId,
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

export function resolveWorkspaceEntryPath(items: WorkspaceRouteItem[], parentPath = ''): string {
  for (const item of items) {
    const currentPath = parentPath ? `${parentPath}/${item.path}` : item.path
    if (item.component) {
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
    return getWorkspaceRoutesCache() || defaultWorkspaceRouteItems
  }

  const cachedRoutes = getWorkspaceRoutesCache()
  if (cachedRoutes?.length) {
    registerWorkspaceRoutes(router, cachedRoutes)
  } else {
    setWorkspaceRoutesCache(defaultWorkspaceRouteItems)
    registerWorkspaceRoutes(router, defaultWorkspaceRouteItems)
  }

  let routeItems = cachedRoutes?.length ? cachedRoutes : defaultWorkspaceRouteItems
  try {
    const response = await getWorkspaceRouters()
    if (response.code === 200 && response.data?.length) {
      routeItems = response.data
      setWorkspaceRoutesCache(routeItems)
      registerWorkspaceRoutes(router, routeItems)
    }
  } catch (_error) {
    routeItems = cachedRoutes?.length ? cachedRoutes : defaultWorkspaceRouteItems
  }

  workspaceRoutesLoaded = true
  workspaceRoutesLoadedForToken = token
  return routeItems
}
