export function calculateCBM(length: number, width: number, height: number, qty: number) {
  const cbm = (length * width * height * qty) / 1000000
  return round(cbm, 3)
}

export function round(value: number, precision = 2) {
  const factor = 10 ** precision
  return Math.round(value * factor) / factor
}
