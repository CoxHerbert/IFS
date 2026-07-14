const RESOURCE_PREFIX = '/profile'

export function resourceUrl(path = '') {
  return `${RESOURCE_PREFIX}/${String(path).replace(/^\/?(profile\/)?/, '')}`
}
