<template>
  <div class="workspace-tags-view">
    <div ref="scrollRef" class="tags-scroll" @scroll="closeMenu">
      <button
        v-for="tag in visitedViews"
        :key="tag.fullPath"
        :data-current="isActive(tag) ? 'true' : 'false'"
        class="tag-item"
        :class="{ active: isActive(tag), affix: isAffix(tag) }"
        type="button"
        @click="goTag(tag)"
        @click.middle="closeTag(tag)"
        @contextmenu.prevent="openMenu(tag, $event)"
      >
        <span class="tag-dot" />
        <span class="tag-title">{{ tag.title }}</span>
        <PushpinOutlined v-if="isAffix(tag)" class="tag-pin" />
        <CloseOutlined v-else class="tag-close" @click.prevent.stop="closeTag(tag)" />
      </button>
    </div>

    <div class="tags-tools">
      <button class="tags-tools-button" type="button" @click="toggleContentFullscreen">
        <FullscreenExitOutlined v-if="contentFullscreen" />
        <FullscreenOutlined v-else />
      </button>

      <a-dropdown placement="bottomRight" trigger="click">
        <button class="tags-tools-button" type="button" @click="selectCurrentTag">
          <MoreOutlined />
        </button>
        <template #overlay>
          <a-menu class="tags-tools-menu">
            <a-menu-item key="close" :disabled="!selectedTag || isAffix(selectedTag)" @click="closeTag(selectedTag)">
              <CloseOutlined />
              <span>关闭</span>
            </a-menu-item>
            <a-menu-item key="pin" :disabled="!selectedTag" @click="toggleAffix(selectedTag)">
              <PushpinOutlined />
              <span>{{ isAffix(selectedTag) ? '取消固定' : '固定' }}</span>
            </a-menu-item>
            <a-menu-item key="fullscreen" @click="toggleContentFullscreen">
              <component :is="contentFullscreen ? FullscreenExitOutlined : FullscreenOutlined" />
              <span>{{ contentFullscreen ? '退出全屏' : '内容全屏' }}</span>
            </a-menu-item>
            <a-menu-item key="reload" :disabled="!selectedTag" @click="refreshTag(selectedTag)">
              <ReloadOutlined />
              <span>重新加载</span>
            </a-menu-item>
            <a-menu-item key="window" :disabled="!selectedTag" @click="openInNewWindow(selectedTag)">
              <ExportOutlined />
              <span>在新窗口打开</span>
            </a-menu-item>
            <a-menu-divider />
            <a-menu-item key="close-left" :disabled="isFirstTag" @click="closeLeft">
              <ArrowLeftOutlined />
              <span>关闭左侧标签页</span>
            </a-menu-item>
            <a-menu-item key="close-right" :disabled="isLastTag" @click="closeRight">
              <ArrowRightOutlined />
              <span>关闭右侧标签页</span>
            </a-menu-item>
            <a-menu-item key="close-others" :disabled="!selectedTag" @click="closeOthers">
              <SwapOutlined />
              <span>关闭其它标签页</span>
            </a-menu-item>
            <a-menu-item key="close-all" @click="closeAll">
              <ColumnWidthOutlined />
              <span>关闭全部标签页</span>
            </a-menu-item>
          </a-menu>
        </template>
      </a-dropdown>
    </div>

    <ul
      v-show="contextVisible"
      class="context-menu"
      :style="{ left: `${contextLeft}px`, top: `${contextTop}px` }"
    >
      <li :class="{ disabled: !selectedTag || isAffix(selectedTag) }" @click="closeTag(selectedTag)">
        <CloseOutlined />
        <span>关闭</span>
      </li>
      <li :class="{ disabled: !selectedTag }" @click="toggleAffix(selectedTag)">
        <PushpinOutlined />
        <span>{{ isAffix(selectedTag) ? '取消固定' : '固定' }}</span>
      </li>
      <li @click="toggleContentFullscreen">
        <component :is="contentFullscreen ? FullscreenExitOutlined : FullscreenOutlined" />
        <span>{{ contentFullscreen ? '退出全屏' : '内容全屏' }}</span>
      </li>
      <li :class="{ disabled: !selectedTag }" @click="refreshTag(selectedTag)">
        <ReloadOutlined />
        <span>重新加载</span>
      </li>
      <li :class="{ disabled: !selectedTag }" @click="openInNewWindow(selectedTag)">
        <ExportOutlined />
        <span>在新窗口打开</span>
      </li>
      <li :class="{ disabled: isFirstTag }" @click="closeLeft">
        <ArrowLeftOutlined />
        <span>关闭左侧标签页</span>
      </li>
      <li :class="{ disabled: isLastTag }" @click="closeRight">
        <ArrowRightOutlined />
        <span>关闭右侧标签页</span>
      </li>
      <li :class="{ disabled: !selectedTag }" @click="closeOthers">
        <SwapOutlined />
        <span>关闭其它标签页</span>
      </li>
      <li @click="closeAll">
        <ColumnWidthOutlined />
        <span>关闭全部标签页</span>
      </li>
    </ul>
  </div>
</template>

<script setup lang="ts">
import { computed, nextTick, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import {
  ArrowLeftOutlined,
  ArrowRightOutlined,
  CloseOutlined,
  ColumnWidthOutlined,
  ExportOutlined,
  FullscreenExitOutlined,
  FullscreenOutlined,
  MoreOutlined,
  PushpinOutlined,
  ReloadOutlined,
  SwapOutlined,
} from '@ant-design/icons-vue'
import { type WorkspaceRouteItem, useWorkspaceRoutesState } from '@/api/workspace/auth'

defineOptions({ name: 'WorkspaceTagsView' })

defineProps<{
  contentFullscreen?: boolean
}>()

const emit = defineEmits<{
  refresh: []
  toggleContentFullscreen: []
}>()

interface WorkspaceTag {
  name?: string
  path: string
  fullPath: string
  title: string
  affix?: boolean
}

const route = useRoute()
const router = useRouter()
const workspaceRoutesState = useWorkspaceRoutesState()
const scrollRef = ref<HTMLDivElement>()
const visitedViews = ref<WorkspaceTag[]>([])
const selectedTag = ref<WorkspaceTag>()
const contextVisible = ref(false)
const contextLeft = ref(0)
const contextTop = ref(0)

const selectedIndex = computed(() => {
  if (!selectedTag.value) return -1
  return visitedViews.value.findIndex((tag) => tag.fullPath === selectedTag.value?.fullPath)
})

const isFirstTag = computed(() => selectedIndex.value <= 0)
const isLastTag = computed(() => selectedIndex.value === visitedViews.value.length - 1)

function normalizePath(path: string) {
  return path.replace(/\/+/g, '/')
}

function findFirstLeaf(items: WorkspaceRouteItem[], parentPath = '/customer'): WorkspaceTag | undefined {
  for (const item of items) {
    const fullPath = normalizePath(`${parentPath}/${item.path}`)
    if (item.component && !item.hidden) {
      return {
        name: item.name,
        path: fullPath,
        fullPath,
        title: item.meta.title,
        affix: true,
      }
    }
    if (item.children?.length) {
      const child = findFirstLeaf(item.children, fullPath)
      if (child) return child
    }
  }
  return undefined
}

function ensureAffixTag() {
  const firstLeaf = findFirstLeaf(workspaceRoutesState.value || [])
  if (!firstLeaf || visitedViews.value.some((tag) => tag.fullPath === firstLeaf.fullPath)) {
    return
  }
  visitedViews.value.unshift(firstLeaf)
}

function currentRouteTag(): WorkspaceTag | undefined {
  if (!route.name || !route.path.startsWith('/customer/') || !route.meta.title) {
    return undefined
  }
  return {
    name: String(route.name),
    path: route.path,
    fullPath: route.fullPath,
    title: String(route.meta.title),
  }
}

function addCurrentTag() {
  ensureAffixTag()
  const tag = currentRouteTag()
  if (!tag) return

  const existing = visitedViews.value.find((item) => item.path === tag.path)
  if (existing) {
    existing.fullPath = tag.fullPath
    existing.title = tag.title
    existing.name = tag.name
    moveToCurrentTag()
    return
  }

  visitedViews.value.push(tag)
  moveToCurrentTag()
}

function moveToCurrentTag() {
  nextTick(() => {
    const target = scrollRef.value?.querySelector<HTMLElement>('[data-current="true"]')
    target?.scrollIntoView({ behavior: 'smooth', inline: 'center', block: 'nearest' })
  })
}

function isActive(tag?: WorkspaceTag) {
  return Boolean(tag && tag.path === route.path)
}

function isAffix(tag?: WorkspaceTag) {
  return Boolean(tag?.affix)
}

function selectCurrentTag() {
  selectedTag.value =
    visitedViews.value.find((tag) => isActive(tag)) || visitedViews.value[visitedViews.value.length - 1]
}

function goTag(tag: WorkspaceTag) {
  closeMenu()
  if (tag.fullPath !== route.fullPath) {
    router.push(tag.fullPath)
  }
}

function latestTag(except?: WorkspaceTag) {
  const tags = except
    ? visitedViews.value.filter((tag) => tag.fullPath !== except.fullPath)
    : visitedViews.value
  return tags[tags.length - 1]
}

function goLatestTag(closedTag?: WorkspaceTag) {
  const latest = latestTag(closedTag)
  if (latest) {
    router.push(latest.fullPath)
    return
  }
  router.push('/customer')
}

function closeTag(tag?: WorkspaceTag) {
  if (!tag || isAffix(tag)) {
    closeMenu()
    return
  }

  const wasActive = isActive(tag)
  visitedViews.value = visitedViews.value.filter((item) => item.fullPath !== tag.fullPath)
  closeMenu()

  if (wasActive) {
    goLatestTag(tag)
  }
}

function toggleAffix(tag?: WorkspaceTag) {
  if (!tag) {
    closeMenu()
    return
  }

  const target = visitedViews.value.find((item) => item.fullPath === tag.fullPath)
  if (target) {
    target.affix = !target.affix
  }
  closeMenu()
}

function toggleContentFullscreen() {
  emit('toggleContentFullscreen')
  closeMenu()
}

function openInNewWindow(tag?: WorkspaceTag) {
  if (!tag) {
    closeMenu()
    return
  }

  window.open(tag.fullPath, '_blank', 'noopener,noreferrer')
  closeMenu()
}

function refreshTag(tag?: WorkspaceTag) {
  if (!tag) {
    closeMenu()
    return
  }

  closeMenu()
  if (!isActive(tag)) {
    router.push(tag.fullPath).then(() => emit('refresh'))
    return
  }
  emit('refresh')
}

function closeOthers() {
  if (!selectedTag.value) {
    closeMenu()
    return
  }

  visitedViews.value = visitedViews.value.filter((tag) => tag.affix || tag.fullPath === selectedTag.value?.fullPath)
  router.push(selectedTag.value.fullPath)
  closeMenu()
}

function closeLeft() {
  if (!selectedTag.value || selectedIndex.value < 0) {
    closeMenu()
    return
  }

  const keep = new Set(
    visitedViews.value
      .filter((tag, index) => tag.affix || index >= selectedIndex.value)
      .map((tag) => tag.fullPath),
  )
  visitedViews.value = visitedViews.value.filter((tag) => keep.has(tag.fullPath))

  if (!visitedViews.value.some((tag) => tag.path === route.path)) {
    router.push(selectedTag.value.fullPath)
  }
  closeMenu()
}

function closeRight() {
  if (!selectedTag.value || selectedIndex.value < 0) {
    closeMenu()
    return
  }

  const keep = new Set(
    visitedViews.value
      .filter((tag, index) => tag.affix || index <= selectedIndex.value)
      .map((tag) => tag.fullPath),
  )
  visitedViews.value = visitedViews.value.filter((tag) => keep.has(tag.fullPath))

  if (!visitedViews.value.some((tag) => tag.path === route.path)) {
    router.push(selectedTag.value.fullPath)
  }
  closeMenu()
}

function closeAll() {
  visitedViews.value = visitedViews.value.filter((tag) => tag.affix)
  goLatestTag()
  closeMenu()
}

function openMenu(tag: WorkspaceTag, event: MouseEvent) {
  const maxLeft = window.innerWidth - 180
  contextLeft.value = Math.max(8, Math.min(event.clientX + 8, maxLeft))
  contextTop.value = event.clientY
  selectedTag.value = tag
  contextVisible.value = true
}

function closeMenu() {
  contextVisible.value = false
}

watch(
  () => route.fullPath,
  () => addCurrentTag(),
  { immediate: true },
)

watch(
  workspaceRoutesState,
  () => ensureAffixTag(),
  { immediate: true },
)

onMounted(() => {
  document.body.addEventListener('click', closeMenu)
})

onBeforeUnmount(() => {
  document.body.removeEventListener('click', closeMenu)
})
</script>

<style scoped>
.workspace-tags-view {
  position: relative;
  height: 38px;
  display: grid;
  grid-template-columns: minmax(0, 1fr) auto;
  align-items: center;
  background: var(--ws-surface);
  border-bottom: 1px solid var(--ws-border);
  box-shadow: 0 8px 22px rgba(15, 23, 42, 0.04);
}

.tags-scroll {
  height: 100%;
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 0 14px;
  overflow-x: auto;
  overflow-y: hidden;
  scrollbar-width: none;
}

.tags-scroll::-webkit-scrollbar {
  display: none;
}

.tag-item {
  height: 28px;
  min-width: 86px;
  max-width: 188px;
  padding: 0 9px;
  border: 1px solid var(--ws-border);
  border-radius: 4px;
  background: var(--ws-surface-strong);
  color: var(--ws-text-secondary);
  display: inline-flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
  line-height: 1;
  cursor: pointer;
  transition: background-color 0.2s ease, border-color 0.2s ease, color 0.2s ease;
  flex: 0 0 auto;
}

.tag-item:hover {
  background: var(--ws-surface-hover);
  color: var(--ws-text-primary);
}

.tag-item.active {
  color: #fff;
  background: #0891b2;
  border-color: #0891b2;
}

.tag-item.affix {
  padding-right: 10px;
}

.tag-dot {
  width: 7px;
  height: 7px;
  border-radius: 50%;
  background: currentColor;
  opacity: 0.35;
  flex: 0 0 auto;
}

.tag-item.active .tag-dot {
  opacity: 1;
}

.tag-title {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.tag-close,
.tag-pin {
  width: 13px;
  height: 13px;
  padding: 2px;
  border-radius: 50%;
  flex: 0 0 auto;
}

.tag-close:hover {
  background: rgba(15, 23, 42, 0.16);
}

.tags-tools {
  height: 100%;
  display: flex;
  align-items: stretch;
}

.tags-tools-button {
  width: 38px;
  height: 100%;
  border: 0;
  border-left: 1px solid var(--ws-border);
  background: var(--ws-surface-strong);
  color: var(--ws-icon-color);
  cursor: pointer;
  transition: background-color 0.2s ease, color 0.2s ease;
}

.tags-tools-button:hover {
  background: var(--ws-surface-hover);
  color: var(--ws-text-primary);
}

.context-menu {
  position: fixed;
  z-index: 3000;
  min-width: 168px;
  margin: 0;
  padding: 6px 0;
  list-style: none;
  border-radius: 6px;
  border: 1px solid var(--ws-border);
  background: var(--ws-surface-hover);
  box-shadow: 0 16px 34px rgba(15, 23, 42, 0.18);
}

.context-menu li {
  height: 32px;
  padding: 0 12px;
  display: flex;
  align-items: center;
  gap: 8px;
  color: var(--ws-text-secondary);
  font-size: 12px;
  cursor: pointer;
}

.context-menu li:hover {
  background: rgba(148, 163, 184, 0.12);
}

.context-menu li.disabled {
  color: var(--ws-text-muted);
  cursor: not-allowed;
}

.context-menu li.disabled:hover {
  background: transparent;
}

@media (max-width: 640px) {
  .workspace-tags-view {
    height: 42px;
  }

  .tag-item {
    min-width: 78px;
    max-width: 146px;
  }
}
</style>
