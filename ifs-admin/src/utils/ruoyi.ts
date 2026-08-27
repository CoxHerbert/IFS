// @ts-nocheck
export function parseTime(time, pattern) {
  if (arguments.length === 0 || !time) {
    return null
  }

  const format = pattern || '{y}-{m}-{d} {h}:{i}:{s}'
  let date

  if (typeof time === 'object') {
    date = time
  } else {
    if (typeof time === 'string' && /^[0-9]+$/.test(time)) {
      time = parseInt(time)
    } else if (typeof time === 'string') {
      time = time.replace(/-/g, '/').replace('T', ' ').replace(/\.[\d]{3}/g, '')
    }
    if (typeof time === 'number' && time.toString().length === 10) {
      time = time * 1000
    }
    date = new Date(time)
  }

  const formatObj = {
    y: date.getFullYear(),
    m: date.getMonth() + 1,
    d: date.getDate(),
    h: date.getHours(),
    i: date.getMinutes(),
    s: date.getSeconds(),
    a: date.getDay()
  }

  return format.replace(/{(y|m|d|h|i|s|a)+}/g, (result, key) => {
    let value = formatObj[key]
    if (key === 'a') {
      return ['日', '一', '二', '三', '四', '五', '六'][value]
    }
    if (result.length > 0 && value < 10) {
      value = '0' + value
    }
    return value || 0
  })
}

export function resetForm(refName) {
  if (this.$refs[refName]) {
    this.$refs[refName].resetFields()
  }
}

export function addDateRange(params, dateRange, propName) {
  const search = params
  search.params = typeof search.params === 'object' && search.params !== null && !Array.isArray(search.params) ? search.params : {}
  dateRange = Array.isArray(dateRange) ? dateRange : []

  if (typeof propName === 'undefined') {
    search.params.beginTime = dateRange[0]
    search.params.endTime = dateRange[1]
  } else {
    search.params['begin' + propName] = dateRange[0]
    search.params['end' + propName] = dateRange[1]
  }

  return search
}

export function selectDictLabel(datas, value) {
  const actions = []
  Object.keys(datas).some((key) => {
    if (datas[key].value == `${value}`) {
      actions.push(datas[key].label)
      return true
    }
    return false
  })
  return actions.join('')
}

export function selectDictLabels(datas, value, separator = ',') {
  const actions = []
  const temp = String(value).split(separator)
  Object.keys(temp).some((val) => {
    Object.keys(datas).some((key) => {
      if (datas[key].value == `${temp[val]}`) {
        actions.push(datas[key].label + separator)
      }
      return false
    })
    return false
  })
  return actions.join('').slice(0, -1)
}

export function sprintf(str) {
  const args = arguments
  let flag = true
  let i = 1
  str = str.replace(/%s/g, () => {
    const arg = args[i++]
    if (typeof arg === 'undefined') {
      flag = false
      return ''
    }
    return arg
  })
  return flag ? str : ''
}

export function praseStrEmpty(str) {
  if (!str || str === 'undefined' || str === 'null') {
    return ''
  }
  return str
}

export function mergeRecursive(source, target) {
  for (const p in target) {
    try {
      if (target[p].constructor == Object) {
        source[p] = mergeRecursive(source[p], target[p])
      } else {
        source[p] = target[p]
      }
    } catch (e) {
      source[p] = target[p]
    }
  }
  return source
}

export function handleTree(data, id = 'id', parentId = 'parentId', children = 'children') {
  const config = {
    id,
    parentId,
    childrenList: children
  }

  const childrenListMap = {}
  const nodeIds = {}
  const tree = []

  for (const d of data) {
    const parentIdValue = d[config.parentId]
    if (childrenListMap[parentIdValue] == null) {
      childrenListMap[parentIdValue] = []
    }
    nodeIds[d[config.id]] = d
    childrenListMap[parentIdValue].push(d)
  }

  for (const d of data) {
    const parentIdValue = d[config.parentId]
    if (nodeIds[parentIdValue] == null) {
      tree.push(d)
    }
  }

  for (const t of tree) {
    adaptToChildrenList(t)
  }

  function adaptToChildrenList(o) {
    if (childrenListMap[o[config.id]] != null) {
      o[config.childrenList] = childrenListMap[o[config.id]]
    }
    if (o[config.childrenList]) {
      for (const c of o[config.childrenList]) {
        adaptToChildrenList(c)
      }
    }
  }

  return tree
}

export function tansParams(params) {
  let result = ''
  for (const propName of Object.keys(params)) {
    const value = params[propName]
    const part = encodeURIComponent(propName) + '='
    if (value !== null && typeof value !== 'undefined') {
      if (typeof value === 'object') {
        for (const key of Object.keys(value)) {
          if (value[key] !== null && typeof value[key] !== 'undefined') {
            const paramsKey = `${propName}[${key}]`
            const subPart = encodeURIComponent(paramsKey) + '='
            result += subPart + encodeURIComponent(value[key]) + '&'
          }
        }
      } else {
        result += part + encodeURIComponent(value) + '&'
      }
    }
  }
  return result
}

export function getNormalPath(path) {
  if (!path || path === 'undefined' || path.length === 0) {
    return path
  }
  const res = path.replace('//', '/')
  if (res[res.length - 1] === '/') {
    return res.slice(0, -1)
  }
  return res
}

export async function blobValidate(data) {
  try {
    const text = await data.text()
    JSON.parse(text)
    return false
  } catch (error) {
    return true
  }
}
