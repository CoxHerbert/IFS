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
        path: 'dashboard',
        name: 'workspace-dashboard',
        component: () => import('@/views/workspace/Dashboard/index.vue'),
        meta: {
          requiresWorkspaceAuth: true,
          title: '工作台',
          icon: 'AppstoreOutlined',
          menuId: '20001',
        },
      },
      {
        path: 'account',
        redirect: '/customer/account-profile',
      },
    ],
  },
  {
    path: '/customer-center',
    redirect: '/customer',
  },
]
