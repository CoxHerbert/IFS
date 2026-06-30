const autoImport = require('unplugin-auto-import/vite')

module.exports = function createAutoImport() {
    return autoImport({
        imports: [
            'vue',
            'vue-router',
            {
                'vuex': ['useStore']
            }
        ],
        dts: false
    })
}
