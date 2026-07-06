<template>
  <a-layout
    class="workspace-shell"
    :class="[`theme-${theme}`, { 'content-fullscreen': contentFullscreen }]"
  >
    <WorkspaceSidebar
      v-if="!contentFullscreen"
      v-model:collapsed="collapsed"
      :menu-items="menuItems"
      :selected-keys="selectedKeys"
      :open-keys="openKeys"
      @menu-click="onMenuClick"
      @open-change="onOpenChange"
    />

    <a-layout class="workspace-main">
      <WorkspaceHeader
        v-if="!contentFullscreen"
        :collapsed="collapsed"
        :page-title="pageTitle"
        :username="profile?.username"
        :theme="theme"
        @toggle-collapse="collapsed = !collapsed"
        @refresh-content="refreshCurrentView"
        @theme-change="applyTheme"
        @go-portal="router.push('/')"
        @go-profile="router.push('/customer/account')"
        @logout="logout"
      />

      <WorkspaceTagsView
        :content-fullscreen="contentFullscreen"
        @refresh="refreshCurrentView"
        @toggle-content-fullscreen="contentFullscreen = !contentFullscreen"
      />

      <a-layout-content class="workspace-content">
        <router-view v-slot="{ Component, route: viewRoute }">
          <KeepAlive>
            <component
              :is="Component"
              v-if="Component && !viewRoute.meta.noCache"
              :key="viewKey(viewRoute.fullPath)"
            />
          </KeepAlive>
          <component
            :is="Component"
            v-if="Component && viewRoute.meta.noCache"
            :key="viewKey(viewRoute.fullPath)"
          />
        </router-view>
      </a-layout-content>
    </a-layout>
  </a-layout>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import {
  getWorkspaceProfile,
  getWorkspaceProfileCache,
  normalizeWorkspaceProfile,
  removeWorkspaceToken,
  setWorkspaceProfileCache,
  setWorkspaceRoutesCache,
  type WorkspaceAccount,
} from '@/api/workspace/auth'
import WorkspaceHeader from './components/WorkspaceHeader.vue'
import WorkspaceSidebar from './components/WorkspaceSidebar.vue'
import WorkspaceTagsView from './components/WorkspaceTagsView.vue'
import { getWorkspaceTheme, setWorkspaceTheme, type WorkspaceTheme } from './theme'
import { useWorkspaceMenu } from './useWorkspaceMenu'

defineOptions({ name: 'WorkspaceShellLayout' })

const route = useRoute()
const router = useRouter()
const collapsed = ref(false)
const contentFullscreen = ref(false)
const theme = ref<WorkspaceTheme>(getWorkspaceTheme())
const profile = ref<WorkspaceAccount>()
const refreshVersions = reactive<Record<string, number>>({})
const { menuItems, openKeys, selectedKeys } = useWorkspaceMenu()

const pageTitle = computed(() => String(route.meta.title || '客户端工作台'))

function onMenuClick(key: string) {
  router.push(key)
}

function onOpenChange(keys: string[]) {
  openKeys.value = keys
}

function applyTheme(nextTheme: WorkspaceTheme) {
  theme.value = nextTheme
  setWorkspaceTheme(nextTheme)
}

function logout() {
  removeWorkspaceToken()
  setWorkspaceProfileCache(null)
  setWorkspaceRoutesCache(null)
  router.push('/customer-login')
}

function refreshCurrentView() {
  refreshVersions[route.fullPath] = (refreshVersions[route.fullPath] || 0) + 1
}

function viewKey(fullPath: string) {
  return `${fullPath}:${refreshVersions[fullPath] || 0}`
}

onMounted(async () => {
  const cachedProfile = getWorkspaceProfileCache()
  if (cachedProfile?.user) {
    profile.value = cachedProfile.user
    return
  }

  try {
    const response = await getWorkspaceProfile()
    const normalizedProfile = normalizeWorkspaceProfile(response.data)
    if (response.code === 200 && normalizedProfile?.user) {
      setWorkspaceProfileCache(normalizedProfile)
      profile.value = normalizedProfile.user
    }
  } catch (_error) {
    // Route guard will handle the next redirect.
  }
})
</script>

<style scoped>
.workspace-shell {
  --ws-text-primary: #0f172a;
  --ws-text-secondary: #334155;
  --ws-text-muted: #64748b;
  --ws-border: rgba(148, 163, 184, 0.18);
  --ws-border-strong: rgba(15, 23, 42, 0.16);
  --ws-header-border: rgba(148, 163, 184, 0.14);
  --ws-surface: rgba(255, 255, 255, 0.9);
  --ws-surface-strong: rgba(255, 255, 255, 0.82);
  --ws-surface-hover: #ffffff;
  --ws-icon-color: #334155;
  --ws-content-bg:
    radial-gradient(circle at top left, rgba(8, 145, 178, 0.12), transparent 24%),
    radial-gradient(circle at bottom right, rgba(14, 165, 233, 0.1), transparent 28%),
    linear-gradient(180deg, #f5f7fb 0%, #eef3f8 100%);
  --ws-sidebar-bg:
    linear-gradient(180deg, rgba(6, 23, 44, 0.98) 0%, rgba(15, 52, 96, 0.98) 100%),
    linear-gradient(135deg, rgba(56, 189, 248, 0.18), transparent 45%);
  --ws-sidebar-menu-color: rgba(241, 245, 249, 0.82);
  --ws-sidebar-hover: rgba(255, 255, 255, 0.08);
  --ws-sidebar-active:
    linear-gradient(90deg, rgba(56, 189, 248, 0.34), rgba(37, 99, 235, 0.26));
  height: 100vh;
  min-height: 0;
  overflow: hidden;
  background: var(--ws-content-bg);
  color: var(--ws-text-primary);
}

.workspace-shell.theme-dark {
  --ws-text-primary: #e2e8f0;
  --ws-text-secondary: #cbd5e1;
  --ws-text-muted: #94a3b8;
  --ws-border: rgba(71, 85, 105, 0.4);
  --ws-border-strong: rgba(148, 163, 184, 0.3);
  --ws-header-border: rgba(71, 85, 105, 0.35);
  --ws-surface: rgba(15, 23, 42, 0.9);
  --ws-surface-strong: rgba(15, 23, 42, 0.82);
  --ws-surface-hover: rgba(30, 41, 59, 0.98);
  --ws-icon-color: #cbd5e1;
  --ws-content-bg:
    radial-gradient(circle at top left, rgba(14, 165, 233, 0.12), transparent 24%),
    radial-gradient(circle at bottom right, rgba(59, 130, 246, 0.1), transparent 28%),
    linear-gradient(180deg, #020617 0%, #0f172a 100%);
  --ws-sidebar-bg:
    linear-gradient(180deg, rgba(2, 6, 23, 0.98) 0%, rgba(15, 23, 42, 0.98) 100%),
    linear-gradient(135deg, rgba(14, 165, 233, 0.12), transparent 45%);
  --ws-sidebar-menu-color: rgba(226, 232, 240, 0.82);
  --ws-sidebar-hover: rgba(148, 163, 184, 0.12);
  --ws-sidebar-active:
    linear-gradient(90deg, rgba(14, 165, 233, 0.3), rgba(37, 99, 235, 0.24));
}

.workspace-main {
  min-width: 0;
  min-height: 0;
  height: 100vh;
  overflow: hidden;
}

.workspace-content {
  min-height: 0;
  padding: 24px 28px 30px;
  overflow-y: auto;
  overflow-x: hidden;
}

.content-fullscreen .workspace-main,
.content-fullscreen .workspace-content {
  height: 100vh;
}

.content-fullscreen .workspace-content {
  padding-top: 18px;
}

@media (max-width: 992px) {
  .workspace-content {
    padding: 18px;
  }
}
</style>
