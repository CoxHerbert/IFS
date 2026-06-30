const vue = require('@vitejs/plugin-vue')

const createAutoImport = require('./auto-import')
const createSvgIcon = require('./svg-icon')
const createCompression = require('./compression')
const createSetupExtend = require('./setup-extend')

module.exports = function createVitePlugins(viteEnv, isBuild = false) {
    const vitePlugins = [vue()]
    vitePlugins.push(createAutoImport())
	vitePlugins.push(createSetupExtend())
    vitePlugins.push(createSvgIcon(isBuild))
	isBuild && vitePlugins.push(...createCompression(viteEnv))
    return vitePlugins
}
