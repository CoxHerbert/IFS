import { computed, h, ref, watch, type ComputedRef, type Ref } from 'vue'
import { useRoute } from 'vue-router'
import { AppstoreOutlined, CalculatorOutlined, MessageOutlined, ProfileOutlined, RadarChartOutlined } from '@ant-design/icons-vue'
import { type WorkspaceRouteItem, useWorkspaceRoutesState } from '@/api/workspace/auth'

const iconMap: Record<string, () => ReturnType<typeof h>> = {
  AppstoreOutlined: () => h(AppstoreOutlined),
  CalculatorOutlined: () => h(CalculatorOutlined),
  MessageOutlined: () => h(MessageOutlined),
  ProfileOutlined: () => h(ProfileOutlined),
  RadarChartOutlined: () => h(RadarChartOutlined),
}

function buildMenuItems(items: WorkspaceRouteItem[], parentPath = '/customer'): Array<Record<string, unknown>> {
  return items
    .filter((item) => !item.hidden)
    .map((item) => {
      const fullPath = `${parentPath}/${item.path}`.replace(/\/+/g, '/')
      const icon = item.meta.icon && iconMap[item.meta.icon] ? iconMap[item.meta.icon] : undefined
      if (item.children?.length) {
        return {
          key: fullPath,
          icon,
          label: item.meta.title,
          children: buildMenuItems(item.children, fullPath),
        }
      }
      return {
        key: fullPath,
        icon,
        label: item.meta.title,
      }
    })
}

function resolveOpenKeys(path: string) {
  const parts = path.split('/').filter(Boolean)
  const keys: string[] = []
  if (parts.length <= 2) {
    return keys
  }
  let current = ''
  for (let index = 0; index < parts.length - 1; index += 1) {
    current += `/${parts[index]}`
    if (index >= 1) {
      keys.push(current)
    }
  }
  return keys
}

export function useWorkspaceMenu(): {
  menuItems: ComputedRef<Array<Record<string, unknown>>>
  openKeys: Ref<string[]>
  selectedKeys: ComputedRef<string[]>
} {
  const route = useRoute()
  const workspaceRoutesState = useWorkspaceRoutesState()
  const openKeys = ref<string[]>([])
  const menuItems = computed(() => buildMenuItems(workspaceRoutesState.value || []))
  const selectedKeys = computed(() => [route.path])

  watch(
    () => route.path,
    (path) => {
      openKeys.value = resolveOpenKeys(path)
    },
    { immediate: true },
  )

  return {
    menuItems,
    openKeys,
    selectedKeys,
  }
}
