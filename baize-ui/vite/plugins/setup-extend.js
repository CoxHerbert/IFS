const setupExtend = require('vite-plugin-vue-setup-extend')

module.exports = function createSetupExtend() {
    const pluginFactory = setupExtend.default || setupExtend
    return pluginFactory()
}
