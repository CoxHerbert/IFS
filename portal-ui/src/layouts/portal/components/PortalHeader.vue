<template>
  <a-layout-header class="topbar">
    <div class="brand">
      <img :src="logoUrl" alt="IFS" />
      <div>
        <h1>IFS 国际物流</h1>
        <p>海运整柜与拼箱服务</p>
      </div>
    </div>

    <a-menu v-model:selectedKeys="selectedKeys" mode="horizontal" class="menu">
      <a-menu-item key="home"><router-link to="/">首页</router-link></a-menu-item>
      <a-menu-item key="news"><router-link to="/news">航线资讯</router-link></a-menu-item>
      <a-menu-item key="service"><router-link to="/service">服务能力</router-link></a-menu-item>
      <a-menu-item key="shipment-agent"><router-link to="/shipment-agent">出货分析</router-link></a-menu-item>
      <a-menu-item key="about"><router-link to="/about">关于我们</router-link></a-menu-item>
      <a-menu-item key="contact"><router-link to="/contact">联系我们</router-link></a-menu-item>
    </a-menu>

    <a-space :size="12" class="topbar-actions">
      <router-link to="/customer">
        <a-button type="text">客户中心</a-button>
      </router-link>
      <router-link to="/contact">
        <a-button type="primary">获取报价</a-button>
      </router-link>
    </a-space>
  </a-layout-header>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import logoUrl from '@/assets/logo.svg'

const route = useRoute()

const selectedKeys = computed<string[]>(() => {
  const name = String(route.name || 'portal-home')
  if (name.startsWith('workspace-')) return ['workspace']
  if (name === 'portal-home') return ['home']
  if (name === 'portal-news') return ['news']
  if (name === 'portal-service') return ['service']
  if (name === 'portal-shipment-agent') return ['shipment-agent']
  if (name === 'portal-about') return ['about']
  if (name === 'portal-contact') return ['contact']
  return ['home']
})
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

@media (max-width: 960px) {
  .topbar {
    flex-direction: column;
    align-items: flex-start;
  }

  .menu {
    width: 100%;
  }
}
</style>
