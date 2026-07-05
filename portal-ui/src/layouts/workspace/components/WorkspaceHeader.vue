<template>
  <a-layout-header class="workspace-header">
    <div class="header-main">
      <div class="header-top">
        <div class="header-meta">
          <a-button type="text" class="collapse-button" @click="$emit('toggle-collapse')">
            <menu-unfold-outlined v-if="collapsed" />
            <menu-fold-outlined v-else />
          </a-button>
          <a-breadcrumb class="header-breadcrumb">
            <a-breadcrumb-item>客户端</a-breadcrumb-item>
            <a-breadcrumb-item>{{ pageTitle }}</a-breadcrumb-item>
          </a-breadcrumb>
        </div>

        <div class="header-right">
          <a-tooltip :title="isFullscreen ? '退出全屏' : '进入全屏'">
            <a-button type="text" class="header-icon-button" @click="toggleFullscreen">
              <fullscreen-exit-outlined v-if="isFullscreen" />
              <fullscreen-outlined v-else />
            </a-button>
          </a-tooltip>
          <a-dropdown>
            <a class="account-pill" @click.prevent>
              <div class="account-copy">
                <strong>{{ username || '客户账号' }}</strong>
                <span>{{ username ? '' : '未登录' }}</span>
              </div>
              <down-outlined />
            </a>
            <template #overlay>
              <a-menu>
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
  DownOutlined,
  FullscreenExitOutlined,
  FullscreenOutlined,
  MenuFoldOutlined,
  MenuUnfoldOutlined,
} from '@ant-design/icons-vue'

defineOptions({ name: 'WorkspaceHeader' })

defineProps<{
  collapsed: boolean
  pageTitle: string
  username?: string
}>()

defineEmits<{
  'toggle-collapse': []
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
  height: auto;
  padding: 10px 24px;
  background: rgba(255, 255, 255, 0.82);
  backdrop-filter: blur(18px);
  border-bottom: 1px solid rgba(148, 163, 184, 0.14);
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

.collapse-button {
  width: 32px;
  height: 32px;
  border-radius: 8px;
  background: linear-gradient(180deg, #ffffff, #f8fbff);
  border: 1px solid rgba(148, 163, 184, 0.14);
  box-shadow: 0 8px 18px rgba(15, 23, 42, 0.04);
  flex-shrink: 0;
}

.header-breadcrumb {
  min-width: 0;
}

.header-breadcrumb :deep(.ant-breadcrumb-link),
.header-breadcrumb :deep(.ant-breadcrumb-separator) {
  color: #64748b;
}

.header-main h2 {
  display: block;
  margin: 0;
  color: #0f172a;
  font-size: 30px;
  line-height: 1.15;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 10px;
  flex-shrink: 0;
}

.header-icon-button {
  width: 34px;
  height: 34px;
  border-radius: 8px;
  border: 1px solid rgba(148, 163, 184, 0.16);
  background: linear-gradient(180deg, #ffffff, #f8fbff);
  box-shadow: 0 8px 18px rgba(15, 23, 42, 0.05);
  color: #334155;
  flex-shrink: 0;
}

.account-pill {
  min-width: 168px;
  padding: 6px 10px;
  border-radius: 8px;
  border: 1px solid rgba(148, 163, 184, 0.16);
  background: linear-gradient(180deg, #ffffff, #f8fbff);
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
  color: #0f172a;
  font-size: 12px;
}

.account-pill span {
  color: #64748b;
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
