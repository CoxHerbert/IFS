import compression from 'vite-plugin-compression'

export default function createCompression(env: Record<string, string>) {
  const compressList = (env.VITE_BUILD_COMPRESS || '').split(',').filter(Boolean)
  const plugins = []

  if (compressList.includes('gzip')) {
    plugins.push(
      compression({
        ext: '.gz',
        deleteOriginFile: false
      })
    )
  }

  if (compressList.includes('brotli')) {
    plugins.push(
      compression({
        ext: '.br',
        algorithm: 'brotliCompress',
        deleteOriginFile: false
      })
    )
  }

  return plugins
}
