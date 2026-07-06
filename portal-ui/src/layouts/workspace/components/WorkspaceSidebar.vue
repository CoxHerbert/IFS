<template>
  <a-layout-sider
    v-model:collapsed="innerCollapsed"
    collapsible
    :trigger="null"
    :width="248"
    class="workspace-sider"
    breakpoint="lg"
  >
    <div class="brand-panel" :class="{ compact: innerCollapsed }">
      <div class="brand-mark">客</div>
      <div v-if="!innerCollapsed" class="brand-copy">
        <strong>客户端工作台</strong>
      </div>
    </div>

    <a-menu
      mode="inline"
      theme="dark"
      :selected-keys="selectedKeys"
      :open-keys="openKeys"
      @click="onMenuClick"
      @openChange="onOpenChange"
    >
      <WorkspaceSidebarMenuNode
        v-for="item in menuItems"
        :key="String(item.key)"
        :item="item"
      />
    </a-menu>
  </a-layout-sider>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import WorkspaceSidebarMenuNode from './WorkspaceSidebarMenuNode.vue'

defineOptions({ name: 'WorkspaceSidebar' })

const props = defineProps<{
  collapsed: boolean
  menuItems: Array<Record<string, unknown>>
  openKeys: string[]
  selectedKeys: string[]
}>()

const emit = defineEmits<{
  'update:collapsed': [value: boolean]
  'menu-click': [key: string]
  'open-change': [keys: string[]]
}>()

const innerCollapsed = computed({
  get: () => props.collapsed,
  set: (value: boolean) => emit('update:collapsed', value),
})

function onMenuClick({ key }: { key: string }) {
  emit('menu-click', key)
}

function onOpenChange(keys: string[]) {
  emit('open-change', keys)
}
</script>

<style scoped>
.workspace-sider {
  height: 100vh;
  overflow: hidden;
  background: var(--ws-sidebar-bg);
  flex: 0 0 auto;
}

.workspace-sider :deep(.ant-layout-sider-children) {
  min-height: 0;
  display: flex;
  flex-direction: column;
}

.brand-panel {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 22px 18px 14px;
}

.brand-panel.compact {
  justify-content: center;
  padding-inline: 10px;
}

.brand-mark {
  width: 46px;
  height: 46px;
  border-radius: 16px;
  display: grid;
  place-items: center;
  background: linear-gradient(135deg, #fef3c7, #7dd3fc);
  color: #082f49;
  font-size: 18px;
  font-weight: 800;
  box-shadow: 0 12px 28px rgba(125, 211, 252, 0.22);
}

.brand-copy strong {
  display: block;
  color: #f8fafc;
  font-size: 16px;
  letter-spacing: 0.02em;
}

.workspace-sider :deep(.ant-layout-sider-trigger),
.workspace-sider :deep(.ant-menu) {
  background: transparent;
}

.workspace-sider :deep(.ant-menu) {
  min-height: 0;
  flex: 1;
  overflow-y: auto;
  overflow-x: hidden;
}

.workspace-sider :deep(.ant-menu::-webkit-scrollbar) {
  width: 6px;
}

.workspace-sider :deep(.ant-menu::-webkit-scrollbar-thumb) {
  border-radius: 999px;
  background: rgba(226, 232, 240, 0.28);
}

.workspace-sider :deep(.ant-menu-item),
.workspace-sider :deep(.ant-menu-submenu-title) {
  height: 46px;
  line-height: 46px;
  margin-inline: 10px;
  width: auto;
  border-radius: 12px;
  color: var(--ws-sidebar-menu-color);
}

.workspace-sider :deep(.ant-menu-item:hover),
.workspace-sider :deep(.ant-menu-submenu-title:hover) {
  color: #fff;
  background: var(--ws-sidebar-hover);
}

.workspace-sider :deep(.ant-menu-item-selected) {
  color: #fff;
  background: var(--ws-sidebar-active);
  box-shadow: inset 0 0 0 1px rgba(125, 211, 252, 0.16);
}

.workspace-sider :deep(.ant-menu-item .ant-menu-title-content),
.workspace-sider :deep(.ant-menu-submenu-title .ant-menu-title-content) {
  font-size: 14px;
  font-weight: 500;
}
</style>
