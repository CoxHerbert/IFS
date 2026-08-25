import type { RouteRecordRaw } from 'vue-router'
import PortalSiteLayout from '@/layouts/portal/SiteLayout.vue'
import Home from '@/views/portal/Home/index.vue'
import News from '@/views/portal/News/index.vue'
import Service from '@/views/portal/Service/index.vue'
import About from '@/views/portal/About/index.vue'
import Contact from '@/views/portal/Contact/index.vue'
import ShipmentShare from '@/views/portal/ShipmentShare/index.vue'
import ChatAgent from '@/views/portal/ChatAgent/index.vue'
import Login from '@/views/workspace/Login/index.vue'
import NotFound from '@/views/error/NotFound.vue'

export const portalRoutes: RouteRecordRaw[] = [
  {
    path: '/',
    component: PortalSiteLayout,
    children: [
      { path: '', name: 'portal-home', component: Home },
      { path: 'news', name: 'portal-news', component: News },
      { path: 'news/:slug', name: 'portal-news-detail', component: News },
      { path: 'service', name: 'portal-service', component: Service },
      { path: 'agent', name: 'portal-agent', component: ChatAgent },
      { path: 'about', name: 'portal-about', component: About },
      { path: 'contact', name: 'portal-contact', component: Contact },
      { path: 'shipment/share/:token', name: 'portal-shipment-share', component: ShipmentShare },
    ],
  },
  {
    path: '/customer-login',
    name: 'workspace-login',
    component: Login,
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'not-found',
    component: NotFound,
  },
]
