export interface ApiResponse<T = unknown> {
  code: number
  msg: string
  data?: T
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

export interface WorkspaceRouteMeta {
  title: string
  icon?: string
  menuId: string
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

export function removeWorkspaceToken() {
  localStorage.removeItem(WORKSPACE_TOKEN_KEY)
  localStorage.removeItem(WORKSPACE_PROFILE_KEY)
  localStorage.removeItem(WORKSPACE_ROUTES_KEY)
}

export async function workspaceLogin(username: string, password: string): Promise<ApiResponse<WorkspaceLoginResult>> {
  const response = await fetch('/portal/customer/login', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ username, password }),
  })

  if (!response.ok) {
    throw new Error('网络请求失败')
  }

  return response.json()
}

export async function getWorkspaceProfile(): Promise<ApiResponse<WorkspaceProfilePayload>> {
  const response = await fetch('/portal/customer/profile', {
    headers: {
      Authorization: `Bearer ${getWorkspaceToken() || ''}`,
    },
  })

  if (!response.ok) {
    throw new Error('网络请求失败')
  }

  return response.json()
}

export async function getWorkspaceRouters(): Promise<ApiResponse<WorkspaceRouteItem[]>> {
  const response = await fetch('/portal/customer/routers', {
    headers: {
      Authorization: `Bearer ${getWorkspaceToken() || ''}`,
    },
  })

  if (!response.ok) {
    throw new Error('网络请求失败')
  }

  return response.json()
}

let workspaceProfileCache: WorkspaceProfile | null = readLocalJson<WorkspaceProfile>(WORKSPACE_PROFILE_KEY)
let workspaceRoutesCache: WorkspaceRouteItem[] | null = readLocalJson<WorkspaceRouteItem[]>(WORKSPACE_ROUTES_KEY)

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
  workspaceProfileCache = profile
  if (profile) {
    writeLocalJson(WORKSPACE_PROFILE_KEY, profile)
    return
  }
  localStorage.removeItem(WORKSPACE_PROFILE_KEY)
}

export function getWorkspaceProfileCache() {
  return workspaceProfileCache
}

export function setWorkspaceRoutesCache(routes: WorkspaceRouteItem[] | null) {
  workspaceRoutesCache = routes
  if (routes) {
    writeLocalJson(WORKSPACE_ROUTES_KEY, routes)
    return
  }
  localStorage.removeItem(WORKSPACE_ROUTES_KEY)
}

export function getWorkspaceRoutesCache() {
  return workspaceRoutesCache
}
