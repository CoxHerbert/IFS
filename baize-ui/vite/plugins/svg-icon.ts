import svgIconsPlugin from 'vite-plugin-svg-icons'
import path from 'node:path'

export default function createSvgIcon(isBuild: boolean) {
  const pluginFactory =
    (svgIconsPlugin as unknown as { ViteSvgIconsPlugin?: typeof svgIconsPlugin; default?: typeof svgIconsPlugin })
      .ViteSvgIconsPlugin ||
    (svgIconsPlugin as unknown as { default?: typeof svgIconsPlugin }).default ||
    svgIconsPlugin

  return pluginFactory({
    iconDirs: [path.resolve(process.cwd(), 'src/assets/icons/svg')],
    symbolId: 'icon-[dir]-[name]',
    svgoOptions: isBuild
  })
}
