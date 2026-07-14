import { Modal } from 'ant-design-vue'
import packageJson from '../../package.json'
import { portalApiUrl } from '@/utils/portal-api'

interface VersionInfo {
  name?: string
  version?: string
  startTime?: string
}

interface VersionResponse {
  code: number
  msg: string
  data?: VersionInfo
}

const DISMISSED_VERSION_KEY = 'ifs_portal_version_dismissed'

function normalizeVersion(version?: string) {
  return (version || '').trim().replace(/^[vV]/, '')
}

function buildRefreshUrl() {
  const url = new URL(window.location.href)
  url.searchParams.set('_v', `${Date.now()}`)
  return url.toString()
}

export async function checkAppVersion() {
  try {
    const response = await fetch(portalApiUrl('/version'))
    if (!response.ok) {
      return
    }

    const result = (await response.json()) as VersionResponse
    const serverVersion = result?.data?.version
    const localVersion = packageJson.version

    if (!serverVersion || normalizeVersion(serverVersion) === normalizeVersion(localVersion)) {
      sessionStorage.removeItem(DISMISSED_VERSION_KEY)
      return
    }

    if (sessionStorage.getItem(DISMISSED_VERSION_KEY) === serverVersion) {
      return
    }

    Modal.confirm({
      title: '发现新版本',
      content: `当前浏览器版本为 ${localVersion}，服务器版本为 ${serverVersion}。是否立即刷新到最新版本？`,
      okText: '确认更新',
      cancelText: '取消',
      onOk: () => {
        sessionStorage.removeItem(DISMISSED_VERSION_KEY)
        window.location.replace(buildRefreshUrl())
      },
      onCancel: () => {
        sessionStorage.setItem(DISMISSED_VERSION_KEY, serverVersion)
      }
    })
  } catch (_error) {
    // Ignore version-check failures to avoid blocking normal page usage.
  }
}
