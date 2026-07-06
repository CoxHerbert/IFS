<template>
  <a-layout-header class="workspace-header">
    <div class="header-main">
      <div class="header-top">
        <div class="header-meta">
          <span class="header-icon-trigger" role="button" tabindex="0" @click="$emit('toggle-collapse')">
            <MenuUnfoldOutlined v-if="collapsed" />
            <MenuFoldOutlined v-else />
          </span>
          <span class="header-icon-trigger" role="button" tabindex="0" @click="$emit('refresh-content')">
            <ReloadOutlined />
          </span>
          <a-breadcrumb class="header-breadcrumb">
            <a-breadcrumb-item>客户端</a-breadcrumb-item>
            <a-breadcrumb-item>{{ pageTitle }}</a-breadcrumb-item>
          </a-breadcrumb>
        </div>

        <div class="header-right">
          <a-tooltip :title="isFullscreen ? '退出全屏' : '进入全屏'">
            <span class="header-icon-trigger" role="button" tabindex="0" @click="toggleFullscreen">
              <FullscreenExitOutlined v-if="isFullscreen" />
              <FullscreenOutlined v-else />
            </span>
          </a-tooltip>

          <a-dropdown>
            <a class="account-pill" @click.prevent>
              <div class="account-copy">
                <strong>{{ username || '客户账号' }}</strong>
                <span>{{ username ? '工作台设置' : '未登录' }}</span>
              </div>
              <DownOutlined />
            </a>
            <template #overlay>
              <a-menu>
                <a-menu-item key="theme-light" @click="$emit('theme-change', 'light')">
                  <CheckOutlined v-if="theme === 'light'" />
                  <span>浅色主题</span>
                </a-menu-item>
                <a-menu-item key="theme-dark" @click="$emit('theme-change', 'dark')">
                  <CheckOutlined v-if="theme === 'dark'" />
                  <span>深色主题</span>
                </a-menu-item>
                <a-menu-divider />
                <a-menu-item key="portal" @click="$emit('go-portal')">返回门户</a-menu-item>
                <a-menu-item key="profile" @click="$emit('go-profile')">账号资料</a-menu-item>
                <a-menu-item key="logout" danger @click="$emit('logout')">退出登录</a-menu-item>
              </a-menu>
            </template>
          </a-dropdown>
        </div>
      </div>
    </div>
  </a-layout-header>
</template>

<script setup lang="ts">
import { onBeforeUnmount, onMounted, ref } from 'vue'
import {
  CheckOutlined,
  DownOutlined,
  FullscreenExitOutlined,
  FullscreenOutlined,
  MenuFoldOutlined,
  MenuUnfoldOutlined,
  ReloadOutlined,
} from '@ant-design/icons-vue'
import type { WorkspaceTheme } from '../theme'

defineOptions({ name: 'WorkspaceHeader' })

defineProps<{
  collapsed: boolean
  pageTitle: string
  username?: string
  theme: WorkspaceTheme
}>()

defineEmits<{
  'toggle-collapse': []
  'refresh-content': []
  'theme-change': [theme: WorkspaceTheme]
  'go-portal': []
  'go-profile': []
  logout: []
}>()

const isFullscreen = ref(false)

function syncFullscreenState() {
  isFullscreen.value = Boolean(document.fullscreenElement)
}

async function toggleFullscreen() {
  if (document.fullscreenElement) {
    await document.exitFullscreen()
    return
  }
  await document.documentElement.requestFullscreen()
}

onMounted(() => {
  syncFullscreenState()
  document.addEventListener('fullscreenchange', syncFullscreenState)
})

onBeforeUnmount(() => {
  document.removeEventListener('fullscreenchange', syncFullscreenState)
})
</script>

<style scoped>
.workspace-header {
  --ws-header-bg: var(--ws-surface);
  height: auto;
  padding: 10px 24px;
  background: var(--ws-header-bg);
  backdrop-filter: blur(18px);
  border-bottom: 1px solid var(--ws-header-border);
  line-height: 1;
}

.header-main {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.header-top {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 14px;
}

.header-meta {
  display: flex;
  align-items: center;
  gap: 10px;
  min-width: 0;
}

.header-icon-trigger {
  width: 34px;
  height: 34px;
  border-radius: 10px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  color: var(--ws-icon-color);
  background: var(--ws-surface-strong);
  border: 1px solid var(--ws-header-border);
  box-shadow: 0 8px 18px rgba(15, 23, 42, 0.04);
  cursor: pointer;
  flex-shrink: 0;
  transition:
    transform 0.18s ease,
    color 0.18s ease,
    background-color 0.18s ease,
    border-color 0.18s ease,
    box-shadow 0.18s ease;
}

.header-icon-trigger:hover {
  color: var(--ws-text-primary);
  background: var(--ws-surface-hover);
  border-color: var(--ws-border-strong);
  box-shadow: 0 12px 24px rgba(15, 23, 42, 0.08);
  transform: translateY(-1px);
}

.header-icon-trigger:focus-visible {
  outline: 2px solid rgba(15, 23, 42, 0.2);
  outline-offset: 2px;
}

.header-breadcrumb {
  min-width: 0;
}

.header-breadcrumb :deep(.ant-breadcrumb-link),
.header-breadcrumb :deep(.ant-breadcrumb-separator) {
  color: var(--ws-text-muted);
}

.header-right {
  display: flex;
  align-items: center;
  gap: 10px;
  flex-shrink: 0;
}

.account-pill {
  min-width: 168px;
  padding: 6px 10px;
  border-radius: 8px;
  border: 1px solid var(--ws-header-border);
  background: var(--ws-surface-strong);
  box-shadow: 0 8px 18px rgba(15, 23, 42, 0.05);
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
  line-height: 1.2;
}

.account-copy {
  min-width: 0;
}

.account-pill strong {
  display: block;
  color: var(--ws-text-primary);
  font-size: 12px;
}

.account-pill span {
  color: var(--ws-text-muted);
  font-size: 11px;
}

@media (max-width: 992px) {
  .workspace-header {
    padding: 8px 14px;
  }

  .header-top {
    flex-direction: column;
    align-items: stretch;
  }

  .header-right,
  .account-pill {
    width: 100%;
  }

  .header-right {
    justify-content: space-between;
  }

  .account-pill {
    min-width: 0;
  }
}

@media (max-width: 640px) {
  .header-meta {
    align-items: flex-start;
  }
}
</style>
