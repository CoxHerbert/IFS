import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'
import PortalLayout from '@/layouts/PortalLayout.vue'
import CustomerLayout from '@/layouts/CustomerLayout.vue'
import HomeView from '@/views/HomeView.vue'
import NewsView from '@/views/NewsView.vue'
import ServiceView from '@/views/ServiceView.vue'
import AboutView from '@/views/AboutView.vue'
import ContactView from '@/views/ContactView.vue'
import CustomerLoginView from '@/views/CustomerLoginView.vue'
import CustomerCenterView from '@/views/CustomerCenterView.vue'
import ShipmentShareView from '@/views/ShipmentShareView.vue'
import { getCustomerToken } from '@/api/customer'

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    component: PortalLayout,
    children: [
      { path: '', name: 'home', component: HomeView },
      { path: 'news', name: 'news', component: NewsView },
      { path: 'service', name: 'service', component: ServiceView },
      { path: 'about', name: 'about', component: AboutView },
      { path: 'contact', name: 'contact', component: ContactView },
      { path: 'customer-login', name: 'customer-login', component: CustomerLoginView },
      { path: 'shipment/share/:token', name: 'shipment-share', component: ShipmentShareView }
    ]
  },
  {
    path: '/',
    component: CustomerLayout,
    children: [
      { path: 'customer-center', name: 'customer-center', component: CustomerCenterView, meta: { requiresCustomerAuth: true } }
    ]
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior() {
    return { top: 0 }
  }
})

router.beforeEach((to, _from, next) => {
  if (to.meta.requiresCustomerAuth && !getCustomerToken()) {
    next('/customer-login')
    return
  }
  next()
})

export default router
