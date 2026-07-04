export interface Dimensions {
  length: number
  width: number
  height: number
  unit: 'cm'
}

export function parseDimensions(value: unknown): Dimensions | null {
  const text = String(value ?? '').trim()
  if (!text) {
    return null
  }

  const normalized = text
    .replace(/厘米|公分/gi, 'cm')
    .replace(/[×X]/g, 'x')
    .replace(/\s+/g, ' ')

  const separated = normalized.match(
    /(\d+(?:\.\d+)?)\s*(?:cm)?\s*(?:[*x]|\s)\s*(\d+(?:\.\d+)?)\s*(?:cm)?\s*(?:[*x]|\s)\s*(\d+(?:\.\d+)?)\s*(?:cm)?/i,
  )
  if (!separated) {
    return null
  }

  const length = Number(separated[1])
  const width = Number(separated[2])
  const height = Number(separated[3])
  if (![length, width, height].every((item) => Number.isFinite(item) && item > 0)) {
    return null
  }
  return { length, width, height, unit: 'cm' }
}
