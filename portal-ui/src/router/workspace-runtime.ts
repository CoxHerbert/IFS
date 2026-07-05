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
  'workspace/shipment-detail': () => import('@/views/workspace/WorkspaceShipmentDetailView.vue'),
  'workspace/shipment-assistant': () => import('@/views/workspace/WorkspaceShipmentAssistantView.vue'),
  'workspace/agent-chat': () => import('@/views/workspace/WorkspaceAgentChatView.vue'),
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
  {
    name: 'workspace-agent-chat',
    path: 'agent-chat',
    component: 'workspace/agent-chat',
    meta: { title: 'Agent 对话', icon: 'MessageOutlined', menuId: '0' },
  },
]

const requiredWorkspaceRouteItems: WorkspaceRouteItem[] = [
  {
    name: 'workspace-shipment-detail',
    path: 'shipment/:shipmentId',
    component: 'workspace/shipment-detail',
    meta: { title: '出货详情', icon: 'RadarChartOutlined', menuId: '0' },
  },
  {
    name: 'workspace-agent-chat',
    path: 'agent-chat',
    component: 'workspace/agent-chat',
    meta: { title: 'Agent 对话', icon: 'MessageOutlined', menuId: '0' },
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
          noCache: item.meta.noCache,
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

function sameWorkspaceRoute(item: WorkspaceRouteItem, route: WorkspaceRouteItem) {
  return item.name === route.name || item.path === route.path || Boolean(item.component && item.component === route.component)
}

function hasWorkspaceRoute(items: WorkspaceRouteItem[], route: WorkspaceRouteItem): boolean {
  return items.some((item) => sameWorkspaceRoute(item, route) || Boolean(item.children?.length && hasWorkspaceRoute(item.children, route)))
}

function withRequiredWorkspaceRoutes(items: WorkspaceRouteItem[]): WorkspaceRouteItem[] {
  const result = [...items]
  for (const route of requiredWorkspaceRouteItems) {
    if (!hasWorkspaceRoute(result, route)) {
      result.push(route)
    }
  }
  return result
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
    return withRequiredWorkspaceRoutes(getWorkspaceRoutesCache() || defaultWorkspaceRouteItems)
  }

  const cachedRoutes = getWorkspaceRoutesCache()
  if (cachedRoutes?.length) {
    registerWorkspaceRoutes(router, withRequiredWorkspaceRoutes(cachedRoutes))
  } else {
    const routes = withRequiredWorkspaceRoutes(defaultWorkspaceRouteItems)
    setWorkspaceRoutesCache(routes)
    registerWorkspaceRoutes(router, routes)
  }

  let routeItems = withRequiredWorkspaceRoutes(cachedRoutes?.length ? cachedRoutes : defaultWorkspaceRouteItems)
  try {
    const response = await getWorkspaceRouters()
    if (response.code === 200 && response.data?.length) {
      routeItems = withRequiredWorkspaceRoutes(response.data)
      setWorkspaceRoutesCache(routeItems)
      registerWorkspaceRoutes(router, routeItems)
    }
  } catch (_error) {
    routeItems = withRequiredWorkspaceRoutes(cachedRoutes?.length ? cachedRoutes : defaultWorkspaceRouteItems)
  }

  workspaceRoutesLoaded = true
  workspaceRoutesLoadedForToken = token
  return routeItems
}
