<template>
  <a-layout class="workspace-shell">
    <WorkspaceSidebar
      v-model:collapsed="collapsed"
      :menu-items="menuItems"
      :selected-keys="selectedKeys"
      :open-keys="openKeys"
      @menu-click="onMenuClick"
      @open-change="onOpenChange"
    />

    <a-layout>
      <WorkspaceHeader
        :collapsed="collapsed"
        :page-title="pageTitle"
        :username="profile?.username"
        @toggle-collapse="collapsed = !collapsed"
        @go-portal="router.push('/')"
        @go-profile="router.push('/customer/account')"
        @logout="logout"
      />

      <a-layout-content class="workspace-content">
        <router-view />
      </a-layout-content>
    </a-layout>
  </a-layout>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
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
import { useWorkspaceMenu } from './useWorkspaceMenu'

const route = useRoute()
const router = useRouter()
const collapsed = ref(false)
const profile = ref<WorkspaceAccount>()
const { menuItems, openKeys, selectedKeys } = useWorkspaceMenu()
const pageTitle = computed(() => String(route.meta.title || '客户端工作台'))

function onMenuClick(key: string) {
  router.push(key)
}

function onOpenChange(keys: string[]) {
  openKeys.value = keys
}

function logout() {
  removeWorkspaceToken()
  setWorkspaceProfileCache(null)
  setWorkspaceRoutesCache(null)
  router.push('/customer-login')
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
  min-height: 100vh;
  background:
    radial-gradient(circle at top left, rgba(8, 145, 178, 0.12), transparent 24%),
    radial-gradient(circle at bottom right, rgba(14, 165, 233, 0.1), transparent 28%),
    linear-gradient(180deg, #f5f7fb 0%, #eef3f8 100%);
}

.workspace-content {
  padding: 24px 28px 30px;
}

@media (max-width: 992px) {
  .workspace-content {
    padding: 18px;
  }
}
</style>
