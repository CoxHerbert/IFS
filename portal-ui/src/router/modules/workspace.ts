import type { RouteRecordRaw } from 'vue-router'
import ShellLayout from '@/layouts/workspace/ShellLayout.vue'

export const workspaceBaseRoutes: RouteRecordRaw[] = [
  {
    path: '/customer',
    name: 'workspace-root',
    component: ShellLayout,
    meta: { requiresWorkspaceAuth: true },
    children: [
      {
        path: 'workspace',
        name: 'workspace-dashboard',
        component: () => import('@/views/workspace/Dashboard/index.vue'),
        meta: { requiresWorkspaceAuth: true, title: '工作台', icon: 'AppstoreOutlined', menuId: '0' },
      },
      {
        path: 'account',
        name: 'workspace-account-profile',
        component: () => import('@/views/workspace/AccountProfile/index.vue'),
        meta: { requiresWorkspaceAuth: true, title: '账号资料', icon: 'ProfileOutlined', menuId: '0', hiddenMenu: true },
      },
      {
        path: 'shipment',
        name: 'workspace-shipment-tracking',
        component: () => import('@/views/workspace/ShipmentTracking/index.vue'),
        meta: { requiresWorkspaceAuth: true, title: '出货查询', icon: 'RadarChartOutlined', menuId: '0' },
      },
      {
        path: 'shipment/:shipmentId',
        name: 'workspace-shipment-detail',
        component: () => import('@/views/workspace/ShipmentDetail/index.vue'),
        meta: { requiresWorkspaceAuth: true, title: '出货详情', icon: 'RadarChartOutlined', menuId: '0', hiddenMenu: true },
      },
      {
        path: 'shipment-assistant',
        name: 'workspace-shipment-assistant',
        component: () => import('@/views/workspace/ShipmentAssistant/index.vue'),
        meta: { requiresWorkspaceAuth: true, title: '智能出货助手', icon: 'CalculatorOutlined', menuId: '0' },
      },
      {
        path: 'agent-chat',
        name: 'workspace-agent-chat',
        component: () => import('@/views/workspace/AgentChat/index.vue'),
        meta: { requiresWorkspaceAuth: true, title: 'Agent 对话', icon: 'MessageOutlined', menuId: '0' },
      },
    ],
  },
  {
    path: '/customer-center',
    redirect: '/customer',
  },
]
