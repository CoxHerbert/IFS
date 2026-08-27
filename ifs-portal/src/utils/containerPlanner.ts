export function planContainer(totalCBM: number) {
  if (totalCBM <= 28) return '1×20GP'
  if (totalCBM <= 58) return '1×40GP'
  if (totalCBM <= 68) return '1×40HQ'
  return `${Math.ceil(totalCBM / 68)}×40HQ`
}
