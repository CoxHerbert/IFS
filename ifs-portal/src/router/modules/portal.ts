import type { RouteRecordRaw } from 'vue-router'
import PortalSiteLayout from '@/layouts/portal/SiteLayout.vue'
import Home from '@/views/portal/Home/index.vue'
import News from '@/views/portal/News/index.vue'
import Service from '@/views/portal/Service/index.vue'
import About from '@/views/portal/About/index.vue'
import Contact from '@/views/portal/Contact/index.vue'
import ShipmentShare from '@/views/portal/ShipmentShare/index.vue'
import Login from '@/views/workspace/Login/index.vue'
import NotFound from '@/views/error/NotFound.vue'
import ChinaToUsa from '@/views/portal/Routes/ChinaToUsa.vue'

export const portalRoutes: RouteRecordRaw[] = [
  {
    path: '/:locale(en|zh-cn)?',
    component: PortalSiteLayout,
    children: [
      { path: '', name: 'portal-home', component: Home, meta: { seoTitle: '中美货运转运平台 | 海运、FBA头程、清关与派送', seoDescription: '提供中美线海运、FBA头程、清关与派送服务，在线获取运输方案、物流咨询和报价。', seoEn: { title: 'China Freight Forwarder | Ocean, Air and Door-to-Door Shipping', description: 'International freight forwarding from China by ocean, air, FCL, LCL and door-to-door services. Request a tailored shipping plan.' } } },
      { path: 'news', name: 'portal-news', component: News, meta: { seoTitle: '国际物流新闻与航线动态 | 中美货运转运平台', seoDescription: '了解国际物流市场、航线运价、港口舱位与清关政策变化，及时规划出货节奏。', seoEn: { title: 'Freight Insights, Shipping Routes and Market Updates', description: 'Practical international freight insights covering routes, transit times, ports, capacity and customs updates.' } } },
      { path: 'news/:slug', name: 'portal-news-detail', component: News, meta: { seoTitle: '物流资讯 | 中美货运转运平台', seoDescription: '国际物流航线、运价、港口与政策资讯。', seoType: 'article' } },
      { path: 'service', name: 'portal-service', component: Service, meta: { seoTitle: '国际货运服务 | 海运、空运、FBA头程与清关派送', seoDescription: '提供海运整柜、海运拼箱、空运快件、FBA头程、清关及门到门派送服务。', seoEn: { title: 'International Freight Services | Ocean, Air, FCL and LCL', description: 'Explore ocean freight, air freight, FCL, LCL, warehousing, customs coordination and door-to-door shipping services.' } } },
      { path: 'about', name: 'portal-about', component: About, meta: { seoTitle: '关于我们 | 中美货运转运平台', seoDescription: '了解中美货运转运平台的国际物流服务能力、服务理念与业务范围。', seoEn: { title: 'About IFS International Logistics', description: 'Learn about our international freight forwarding capabilities, service approach and covered shipping routes.' } } },
      { path: 'contact', name: 'portal-contact', component: Contact, meta: { seoTitle: '联系询价 | 获取国际物流运输方案', seoDescription: '提交起运港、目的港、货物信息与时效要求，获取可执行的国际货运方案与报价。', seoEn: { title: 'Get a Freight Quote | Contact IFS Logistics', description: 'Share your origin, destination, cargo details and timeline to receive a tailored international shipping plan.' } } },
      { path: 'routes/china-to-usa', name: 'portal-route-china-usa', component: ChinaToUsa, meta: { seoTitle: 'Freight Forwarder China to USA | Ocean, Air & DDP Shipping', seoDescription: 'Ship from China to the USA by ocean, air, FCL, LCL, DDP or door to door. Compare transit times, routes and shipping options, then request a tailored freight quote.' } },
      { path: 'shipment/share/:token', name: 'portal-shipment-share', component: ShipmentShare, meta: { noIndex: true } },
    ],
  },
  {
    path: '/customer-login',
    name: 'workspace-login',
    component: Login,
    meta: { noIndex: true },
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'not-found',
    component: NotFound,
    meta: { noIndex: true },
  },
]
