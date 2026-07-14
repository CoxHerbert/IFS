import { shallowRef } from 'vue'
import { portalApiUrl } from '@/utils/portal-api'

export interface ApiResponse<T = unknown> {
  code: number
  msg: string
  data?: T
}

export interface WorkspaceCaptcha {
  uuid: string
  img: string
}

export interface WorkspaceAccount {
  accountId: string
  customerId: string
  customerNo: string
  customerName: string
  companyName: string
  username: string
  realName: string
  phone: string
  email: string
  isMain: string
  status: string
  lastLoginTime?: string
}

export interface WorkspaceProfile {
  user: WorkspaceAccount
  roles: string[]
  permissions: string[]
}

export type WorkspaceProfilePayload = WorkspaceProfile | WorkspaceAccount

export interface WorkspaceLoginResult {
  token: string
  user: WorkspaceAccount
}

export interface WorkspacePasswordUpdatePayload {
  oldPassword: string
  newPassword: string
  confirmPassword: string
}

export interface WorkspaceProfileUpdatePayload {
  realName: string
  phone: string
  email: string
}

export interface WorkspaceRouteMeta {
  title: string
  icon?: string
  menuId: string
  noCache?: boolean
}

export interface WorkspaceRouteItem {
  name: string
  path: string
  component?: string
  hidden?: boolean
  meta: WorkspaceRouteMeta
  children?: WorkspaceRouteItem[]
}

export const WORKSPACE_TOKEN_KEY = 'portal_customer_token'
export const WORKSPACE_PROFILE_KEY = 'portal_customer_profile'
export const WORKSPACE_ROUTES_KEY = 'portal_customer_routes'

function readLocalJson<T>(key: string): T | null {
  const raw = localStorage.getItem(key)
  if (!raw) {
    return null
  }
  try {
    return JSON.parse(raw) as T
  } catch (_error) {
    localStorage.removeItem(key)
    return null
  }
}

function writeLocalJson(key: string, value: unknown) {
  localStorage.setItem(key, JSON.stringify(value))
}

export function getWorkspaceToken() {
  return localStorage.getItem(WORKSPACE_TOKEN_KEY)
}

export function setWorkspaceToken(token: string) {
  localStorage.setItem(WORKSPACE_TOKEN_KEY, token)
}

const workspaceProfileState = shallowRef<WorkspaceProfile | null>(readLocalJson<WorkspaceProfile>(WORKSPACE_PROFILE_KEY))
const workspaceRoutesState = shallowRef<WorkspaceRouteItem[] | null>(
  readLocalJson<WorkspaceRouteItem[]>(WORKSPACE_ROUTES_KEY),
)

export function removeWorkspaceToken() {
  localStorage.removeItem(WORKSPACE_TOKEN_KEY)
  localStorage.removeItem(WORKSPACE_PROFILE_KEY)
  localStorage.removeItem(WORKSPACE_ROUTES_KEY)
  workspaceProfileState.value = null
  workspaceRoutesState.value = null
}

export async function getWorkspaceCaptcha(): Promise<ApiResponse<WorkspaceCaptcha>> {
  const response = await fetch(portalApiUrl('/captchaImage'))

  if (!response.ok) {
    throw new Error('网络请求失败')
  }

  return response.json()
}

export async function workspaceLogin(
  username: string,
  password: string,
  code: string,
  uuid: string,
): Promise<ApiResponse<WorkspaceLoginResult>> {
  const response = await fetch(portalApiUrl('/customer/login'), {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ username, password, code, uuid }),
  })

  if (!response.ok) {
    throw new Error('网络请求失败')
  }

  return response.json()
}

export async function getWorkspaceProfile(): Promise<ApiResponse<WorkspaceProfilePayload>> {
  const response = await fetch(portalApiUrl('/customer/profile'), {
    headers: {
      Authorization: `Bearer ${getWorkspaceToken() || ''}`,
    },
  })

  if (!response.ok) {
    throw new Error('网络请求失败')
  }

  return response.json()
}

export async function updateWorkspaceProfile(payload: WorkspaceProfileUpdatePayload): Promise<ApiResponse<WorkspaceAccount>> {
  const response = await fetch(portalApiUrl('/customer/profile'), {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Bearer ${getWorkspaceToken() || ''}`,
    },
    body: JSON.stringify(payload),
  })

  if (!response.ok) {
    throw new Error('网络请求失败')
  }

  return response.json()
}

export async function getWorkspaceRouters(): Promise<ApiResponse<WorkspaceRouteItem[]>> {
  const response = await fetch(portalApiUrl('/customer/routers'), {
    headers: {
      Authorization: `Bearer ${getWorkspaceToken() || ''}`,
    },
  })

  if (!response.ok) {
    throw new Error('网络请求失败')
  }

  return response.json()
}

export async function updateWorkspacePassword(payload: WorkspacePasswordUpdatePayload): Promise<ApiResponse> {
  const response = await fetch(portalApiUrl('/customer/password'), {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Bearer ${getWorkspaceToken() || ''}`,
    },
    body: JSON.stringify(payload),
  })

  if (!response.ok) {
    throw new Error('网络请求失败')
  }

  return response.json()
}

export function normalizeWorkspaceProfile(payload?: WorkspaceProfilePayload): WorkspaceProfile | null {
  if (!payload) {
    return null
  }
  if ('user' in payload) {
    return {
      user: payload.user,
      roles: payload.roles || [],
      permissions: payload.permissions || [],
    }
  }
  return {
    user: payload,
    roles: [],
    permissions: [],
  }
}

export function setWorkspaceProfileCache(profile: WorkspaceProfile | null) {
  workspaceProfileState.value = profile
  if (profile) {
    writeLocalJson(WORKSPACE_PROFILE_KEY, profile)
    return
  }
  localStorage.removeItem(WORKSPACE_PROFILE_KEY)
}

export function getWorkspaceProfileCache() {
  return workspaceProfileState.value
}

export function useWorkspaceProfileState() {
  return workspaceProfileState
}

export function setWorkspaceRoutesCache(routes: WorkspaceRouteItem[] | null) {
  workspaceRoutesState.value = routes
  if (routes) {
    writeLocalJson(WORKSPACE_ROUTES_KEY, routes)
    return
  }
  localStorage.removeItem(WORKSPACE_ROUTES_KEY)
}

export function getWorkspaceRoutesCache() {
  return workspaceRoutesState.value
}

export function useWorkspaceRoutesState() {
  return workspaceRoutesState
}
