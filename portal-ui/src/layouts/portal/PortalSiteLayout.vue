<template>
  <a-layout class="portal-shell">
    <PortalHeader />

    <a-layout-content class="portal-content">
      <transition name="portal-loader-fade">
        <div v-if="pageLoading" class="portal-loader">
          <div class="portal-loader-card">
            <span class="portal-loader-spinner" />
            <strong>加载中</strong>
            <p>请稍候，正在准备页面内容</p>
          </div>
        </div>
      </transition>
      <router-view />
    </a-layout-content>

    <PortalFooter />
    <PortalFloatingAgent />
  </a-layout>
</template>

<script setup lang="ts">
import { onBeforeUnmount, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import PortalFooter from './components/PortalFooter.vue'
import PortalFloatingAgent from './components/PortalFloatingAgent.vue'
import PortalHeader from './components/PortalHeader.vue'

const route = useRoute()
const pageLoading = ref(true)
let loadingTimer: ReturnType<typeof setTimeout> | undefined

watch(
  () => route.fullPath,
  () => {
    pageLoading.value = true
    if (loadingTimer) {
      clearTimeout(loadingTimer)
    }
    loadingTimer = setTimeout(() => {
      pageLoading.value = false
    }, 320)
  },
  { immediate: true },
)

onBeforeUnmount(() => {
  if (loadingTimer) {
    clearTimeout(loadingTimer)
  }
})
</script>

<style scoped>
.portal-content {
  position: relative;
}

.portal-loader {
  position: fixed;
  inset: 88px 0 0;
  z-index: 30;
  display: grid;
  place-items: center;
  padding: 24px;
  background: rgba(248, 250, 252, 0.74);
  backdrop-filter: blur(8px);
}

.portal-loader-card {
  min-width: min(320px, calc(100vw - 48px));
  padding: 26px 28px;
  border-radius: 24px;
  background: rgba(255, 255, 255, 0.96);
  border: 1px solid rgba(15, 23, 42, 0.08);
  box-shadow: 0 22px 48px rgba(15, 23, 42, 0.12);
  text-align: center;
}

.portal-loader-card strong,
.portal-loader-card p {
  display: block;
}

.portal-loader-card strong {
  margin-top: 16px;
  color: #0f172a;
  font-size: 18px;
}

.portal-loader-card p {
  margin: 8px 0 0;
  color: #64748b;
}

.portal-loader-spinner {
  width: 42px;
  height: 42px;
  display: inline-block;
  border: 3px solid rgba(22, 119, 255, 0.16);
  border-top-color: #1677ff;
  border-radius: 50%;
  animation: portal-loader-spin 0.8s linear infinite;
}

.portal-loader-fade-enter-active,
.portal-loader-fade-leave-active {
  transition: opacity 0.18s ease;
}

.portal-loader-fade-enter-from,
.portal-loader-fade-leave-to {
  opacity: 0;
}

@keyframes portal-loader-spin {
  to {
    transform: rotate(360deg);
  }
}

@media (max-width: 960px) {
  .portal-loader {
    inset: 132px 0 0;
  }
}
</style>
