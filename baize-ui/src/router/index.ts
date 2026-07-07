// @ts-nocheck
import { createRouter, createWebHistory } from 'vue-router'
import Layout from '@/layout/index.vue'

/**
 * 路由配置说明
 *
 * hidden: true
 * 当设置为 true 时，该路由不会在侧边栏中显示，例如 login、401 等页面。
 *
 * alwaysShow: true
 * 当一个路由下面的 children 多于 1 个时，默认会显示为嵌套菜单。
 * 当只有 1 个子路由时，默认会将该子路由显示为根菜单。
 * 如果无论 children 数量多少都希望显示根路由，则设置 alwaysShow: true。
 *
 * redirect: noRedirect
 * 当设置为 noRedirect 时，该路由在面包屑导航中不可点击。
 *
 * name: 'router-name'
 * 路由名称，使用 <keep-alive> 时必须设置。
 *
 * query: '{"id": 1, "name": "ry"}'
 * 访问路由时默认携带的参数。
 *
 * meta: {
 *   noCache: true
 *   设置为 true 时不会被 <keep-alive> 缓存，默认 false。
 *
 *   title: 'title'
 *   侧边栏和面包屑中显示的名称。
 *
 *   icon: 'svg-name'
 *   路由图标，对应 src/assets/icons/svg。
 *
 *   breadcrumb: false
 *   设置为 false 时不会在面包屑中显示。
 *
 *   activeMenu: '/system/user'
 *   当前路由激活时，对应高亮的侧边栏菜单。
 * }
 */

// 公共路由
export const constantRoutes = [
  {
    path: '/portal',
    component: () => import('@/views/portal/index.vue'),
    hidden: true
  },
  {
    path: '/redirect',
    component: Layout,
    hidden: true,
    children: [
      {
        path: '/redirect/:path(.*)',
        component: () => import('@/views/redirect/index.vue')
      }
    ]
  },
  {
    path: '/login',
    component: () => import('@/views/login.vue'),
    hidden: true
  },
  {
    path: '/register',
    component: () => import('@/views/register.vue'),
    hidden: true
  },
  {
    path: '/:pathMatch(.*)*',
    component: () => import('@/views/error/404.vue'),
    hidden: true
  },
  {
    path: '/401',
    component: () => import('@/views/error/401.vue'),
    hidden: true
  },
  {
    path: '',
    component: Layout,
    redirect: 'index',
    children: [
      {
        path: '/index',
        component: () => import('@/views/index.vue'),
        name: 'Index',
        meta: { title: '首页', icon: 'dashboard', affix: true }
      }
    ]
  },
  {
    path: '/user',
    component: Layout,
    hidden: true,
    redirect: 'noredirect',
    children: [
      {
        path: 'profile',
        component: () => import('@/views/system/user/profile/index.vue'),
        name: 'Profile',
        meta: { title: '个人中心', icon: 'user' }
      }
    ]
  },
  {
    path: '/system/user-auth',
    component: Layout,
    hidden: true,
    children: [
      {
        path: 'role/:userId(\\d+)',
        component: () => import('@/views/system/user/authRole.vue'),
        name: 'AuthRole',
        meta: { title: '分配角色', activeMenu: '/system/user' }
      }
    ]
  },
  {
    path: '/system/role-auth',
    component: Layout,
    hidden: true,
    children: [
      {
        path: 'user/:roleId(\\d+)',
        component: () => import('@/views/system/role/authUser.vue'),
        name: 'AuthUser',
        meta: { title: '分配用户', activeMenu: '/system/role' }
      }
    ]
  },
  {
    path: '/system/dict-data',
    component: Layout,
    hidden: true,
    children: [
      {
        path: 'index/:dictId(\\d+)',
        component: () => import('@/views/system/dict/data.vue'),
        name: 'Data',
        meta: { title: '字典数据', activeMenu: '/system/dict' }
      }
    ]
  },
  {
    path: '/monitor/job-log',
    component: Layout,
    hidden: true,
    children: [
      {
        path: 'index',
        component: () => import('@/views/monitor/job/log.vue'),
        name: 'JobLog',
        meta: { title: '调度日志', activeMenu: '/monitor/job' }
      }
    ]
  },
  {
    path: '/tool/gen-edit',
    component: Layout,
    hidden: true,
    children: [
      {
        path: 'index',
        component: () => import('@/views/tool/gen/editTable.vue'),
        name: 'GenEdit',
        meta: { title: '修改生成配置', activeMenu: '/tool/gen' }
      }
    ]
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes: constantRoutes,
  scrollBehavior(_to, _from, savedPosition) {
    if (savedPosition) {
      return savedPosition
    }
    return { top: 0 }
  }
})

export default router
