import type { RouteRecordRaw } from 'vue-router'
import PortalSiteLayout from '@/layouts/portal/PortalSiteLayout.vue'
import PortalHomeView from '@/views/portal/PortalHomeView.vue'
import PortalNewsView from '@/views/portal/PortalNewsView.vue'
import PortalServiceView from '@/views/portal/PortalServiceView.vue'
import PortalAboutView from '@/views/portal/PortalAboutView.vue'
import PortalContactView from '@/views/portal/PortalContactView.vue'
import PortalShipmentShareView from '@/views/portal/PortalShipmentShareView.vue'
import ChatAgent from '@/views/portal/ChatAgent.vue'
import WorkspaceLoginView from '@/views/workspace/WorkspaceLoginView.vue'

export const portalRoutes: RouteRecordRaw[] = [
  {
    path: '/',
    component: PortalSiteLayout,
    children: [
      { path: '', name: 'portal-home', component: PortalHomeView },
      { path: 'news', name: 'portal-news', component: PortalNewsView },
      { path: 'news/:slug', name: 'portal-news-detail', component: PortalNewsView },
      { path: 'service', name: 'portal-service', component: PortalServiceView },
      { path: 'agent', name: 'portal-agent', component: ChatAgent },
      { path: 'about', name: 'portal-about', component: PortalAboutView },
      { path: 'contact', name: 'portal-contact', component: PortalContactView },
      { path: 'shipment/share/:token', name: 'portal-shipment-share', component: PortalShipmentShareView },
    ],
  },
  {
    path: '/customer-login',
    name: 'workspace-login',
    component: WorkspaceLoginView,
  },
]
