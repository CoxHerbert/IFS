import type { RouteRecordRaw } from 'vue-router'
import WorkspaceShellLayout from '@/layouts/workspace/WorkspaceShellLayout.vue'
import WorkspaceDashboardView from '@/views/workspace/WorkspaceDashboardView.vue'
import WorkspaceAccountProfileView from '@/views/workspace/WorkspaceAccountProfileView.vue'
import WorkspaceShipmentTrackingView from '@/views/workspace/WorkspaceShipmentTrackingView.vue'

export const workspaceBaseRoutes: RouteRecordRaw[] = [
  {
    path: '/customer',
    name: 'workspace-root',
    component: WorkspaceShellLayout,
    meta: { requiresWorkspaceAuth: true },
    children: [
      {
        path: 'workspace',
        name: 'workspace-dashboard',
        component: WorkspaceDashboardView,
        meta: { requiresWorkspaceAuth: true, title: '工作台', icon: 'AppstoreOutlined', menuId: '0' },
      },
      {
        path: 'account',
        name: 'workspace-account-profile',
        component: WorkspaceAccountProfileView,
        meta: { requiresWorkspaceAuth: true, title: '账号资料', icon: 'ProfileOutlined', menuId: '0' },
      },
      {
        path: 'shipment',
        name: 'workspace-shipment-tracking',
        component: WorkspaceShipmentTrackingView,
        meta: { requiresWorkspaceAuth: true, title: '出货查询', icon: 'RadarChartOutlined', menuId: '0' },
      },
    ],
  },
  {
    path: '/customer-center',
    redirect: '/customer',
  },
]
