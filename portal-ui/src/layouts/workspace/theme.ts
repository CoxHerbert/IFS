export type WorkspaceTheme = 'light' | 'dark'

export const WORKSPACE_THEME_KEY = 'portal_workspace_theme'

export function normalizeWorkspaceTheme(theme?: string | null): WorkspaceTheme {
  return theme === 'dark' ? 'dark' : 'light'
}

export function getWorkspaceTheme(): WorkspaceTheme {
  return normalizeWorkspaceTheme(localStorage.getItem(WORKSPACE_THEME_KEY))
}

export function setWorkspaceTheme(theme: WorkspaceTheme) {
  localStorage.setItem(WORKSPACE_THEME_KEY, theme)
}
