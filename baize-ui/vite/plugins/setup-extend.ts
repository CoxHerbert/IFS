import setupExtend from 'vite-plugin-vue-setup-extend'

export default function createSetupExtend() {
  const pluginFactory = (setupExtend as { default?: typeof setupExtend }).default || setupExtend
  return pluginFactory()
}
