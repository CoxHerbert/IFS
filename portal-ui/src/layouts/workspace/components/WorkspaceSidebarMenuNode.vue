<template>
  <a-sub-menu v-if="hasChildren" :key="String(currentItem.key)">
    <template #icon>
      <component :is="currentItem.icon" v-if="currentItem.icon" />
    </template>
    <template #title>{{ currentItem.label }}</template>

    <WorkspaceSidebarMenuNode
      v-for="child in currentItem.children"
      :key="String(child.key)"
      :item="child"
    />
  </a-sub-menu>

  <a-menu-item v-else :key="String(currentItem.key)">
    <template #icon>
      <component :is="currentItem.icon" v-if="currentItem.icon" />
    </template>
    {{ currentItem.label }}
  </a-menu-item>
</template>

<script setup lang="ts">
import { computed } from 'vue'

defineOptions({ name: 'WorkspaceSidebarMenuNode' })

interface SidebarMenuNode {
  key: string
  label: string
  icon?: unknown
  children?: SidebarMenuNode[]
}

const props = defineProps<{
  item: unknown
}>()

const currentItem = computed(() => props.item as SidebarMenuNode)
const hasChildren = computed(() => Boolean(currentItem.value.children?.length))
</script>
