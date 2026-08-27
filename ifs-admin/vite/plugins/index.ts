import vue from '@vitejs/plugin-vue'
import createAutoImport from './auto-import'
import createSvgIcon from './svg-icon'
import createCompression from './compression'
import createSetupExtend from './setup-extend'

export default function createVitePlugins(viteEnv: Record<string, string>, isBuild = false) {
  const vitePlugins = [vue(), createAutoImport(), createSetupExtend(), createSvgIcon(isBuild)]

  if (isBuild) {
    vitePlugins.push(...createCompression(viteEnv))
  }

  return vitePlugins
}
