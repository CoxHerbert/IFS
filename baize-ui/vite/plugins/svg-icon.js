const svgIcon = require('vite-plugin-svg-icons')
const path = require('path')

module.exports = function createSvgIcon(isBuild) {
    const pluginFactory = svgIcon.default || svgIcon
    return pluginFactory({
		iconDirs: [path.resolve(process.cwd(), 'src/assets/icons/svg')],
        symbolId: 'icon-[dir]-[name]',
        svgoOptions: isBuild
    })
}
