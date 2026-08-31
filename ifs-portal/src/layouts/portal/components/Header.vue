<template>
  <a-layout-header class="topbar">
    <div class="brand">
      <img :src="logoUrl" alt="IFS" />
      <div>
        <h1>{{ t('brandName') }}</h1>
        <p>{{ t('brandTagline') }}</p>
      </div>
    </div>

    <a-menu :selectedKeys="selectedKeys" mode="horizontal" class="menu">
      <a-menu-item key="home"><router-link :to="localePath('/')">{{ t('home') }}</router-link></a-menu-item>
      <a-menu-item key="news"><router-link :to="localePath('/news')">{{ t('news') }}</router-link></a-menu-item>
      <a-menu-item key="service"><router-link :to="localePath('/service')">{{ t('services') }}</router-link></a-menu-item>
      <a-menu-item key="routes"><router-link :to="localePath('/routes/china-to-usa')">{{ t('route') }}</router-link></a-menu-item>
      <a-menu-item key="about"><router-link :to="localePath('/about')">{{ t('about') }}</router-link></a-menu-item>
      <a-menu-item key="contact"><router-link :to="localePath('/contact')">{{ t('contact') }}</router-link></a-menu-item>
    </a-menu>

    <a-space :size="12" class="topbar-actions">
      <router-link to="/customer">
        <a-button type="text">{{ t('customerCenter') }}</a-button>
      </router-link>
      <a-select :value="locale" class="locale-select" aria-label="Language" @change="changeLocale">
        <a-select-option value="en">English</a-select-option>
        <a-select-option value="zh-cn">简体中文</a-select-option>
      </a-select>
      <router-link :to="localePath('/contact')">
        <a-button type="primary">{{ t('getQuote') }}</a-button>
      </router-link>
    </a-space>
  </a-layout-header>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { useRouter } from 'vue-router'
import logoUrl from '@/assets/logo.svg'
import { usePortalI18n, type PortalLocale } from '@/i18n'

const route = useRoute()
const router = useRouter()
const { locale, t, localePath, setLocale } = usePortalI18n()

function changeLocale(value: PortalLocale) {
  setLocale(value)
  router.push(localePath(route.fullPath, value))
}

const selectedKeys = computed<string[]>(() => {
  const name = String(route.name || 'portal-home')
  if (name.startsWith('workspace-')) return ['workspace']
  if (name === 'portal-home') return ['home']
  if (name === 'portal-news') return ['news']
  if (name === 'portal-service') return ['service']
  if (name.startsWith('portal-route-')) return ['routes']
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
.locale-select { width: 116px; }

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
