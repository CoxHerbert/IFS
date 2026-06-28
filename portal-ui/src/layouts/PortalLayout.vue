<template>
  <a-layout class="portal-shell">
    <a-layout-header class="topbar">
      <div class="brand">
        <img :src="logoUrl" alt="海运通" />
        <div>
          <h1>海运通</h1>
          <p>货代销售官网模板</p>
        </div>
      </div>

      <a-menu v-model:selectedKeys="selectedKeys" mode="horizontal" class="menu">
        <a-menu-item key="home"><router-link to="/">首页</router-link></a-menu-item>
        <a-menu-item key="news"><router-link to="/news">航线资讯</router-link></a-menu-item>
        <a-menu-item key="service"><router-link to="/service">服务能力</router-link></a-menu-item>
        <a-menu-item key="about"><router-link to="/about">关于我们</router-link></a-menu-item>
        <a-menu-item key="contact"><router-link to="/contact">联系我们</router-link></a-menu-item>
        <a-menu-item key="customer-center"><router-link to="/customer-center">客户中心</router-link></a-menu-item>
      </a-menu>

      <a-space :size="12" class="topbar-actions">
        <router-link to="/customer-login">
          <a-button type="text">客户登录</a-button>
        </router-link>
        <router-link to="/contact">
          <a-button type="primary">获取报价</a-button>
        </router-link>
      </a-space>
    </a-layout-header>

    <a-layout-content>
      <router-view />
    </a-layout-content>

    <a-layout-footer class="footer">
      <div>
        <strong>海运通货代官网</strong>
        <p>整合询价、航线、案例和联系入口的一体化销售页面。</p>
      </div>
      <div class="footer-meta">
        <span>quote@seawaypro.com</span>
        <span>400-888-2026</span>
      </div>
    </a-layout-footer>
  </a-layout>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import logoUrl from '@/assets/logo.svg'

const route = useRoute()
const selectedKeys = computed<string[]>(() => [String(route.name || 'home')])
</script>

<style scoped>
.topbar {
  position: sticky;
  top: 0;
  z-index: 20;
  display: flex;
  align-items: center;
  gap: 24px;
  height: auto;
  padding: 18px 40px;
  line-height: 1;
  background: rgba(255, 255, 255, 0.86);
  backdrop-filter: blur(16px);
  border-bottom: 1px solid rgba(16, 35, 63, 0.08);
}

.brand {
  display: flex;
  align-items: center;
  gap: 14px;
  min-width: 240px;
}

.brand img {
  width: 42px;
  height: 42px;
}

.brand h1 {
  margin: 0;
  font-size: 18px;
  font-weight: 800;
}

.brand p {
  margin: 6px 0 0;
  color: #66748b;
  font-size: 12px;
}

.menu {
  flex: 1;
  min-width: 0;
  border-bottom: 0;
  background: transparent;
}

.topbar-actions {
  flex-shrink: 0;
}

.footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 20px;
  padding: 28px 40px;
  background: rgba(255, 255, 255, 0.88);
  border-top: 1px solid rgba(16, 35, 63, 0.08);
}

.footer p {
  margin: 6px 0 0;
  color: #66748b;
}

.footer-meta {
  display: flex;
  gap: 18px;
  color: #55657d;
}

@media (max-width: 960px) {
  .topbar,
  .footer {
    flex-direction: column;
    align-items: flex-start;
  }

  .menu {
    width: 100%;
  }
}
</style>
