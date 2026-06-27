import { createRouter, createWebHistory } from 'vue-router'
import PortalLayout from '@/layouts/PortalLayout.vue'
import HomeView from '@/views/HomeView.vue'
import NewsView from '@/views/NewsView.vue'
import ServiceView from '@/views/ServiceView.vue'
import AboutView from '@/views/AboutView.vue'

const routes = [
  {
    path: '/',
    component: PortalLayout,
    children: [
      { path: '', name: 'home', component: HomeView },
      { path: 'news', name: 'news', component: NewsView },
      { path: 'service', name: 'service', component: ServiceView },
      { path: 'about', name: 'about', component: AboutView }
    ]
  }
]

export default createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior() {
    return { top: 0 }
  }
})
