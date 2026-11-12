declare module 'xlsx' {
  export const SSF: { parse_date_code(value: number): { y: number; m: number; d: number; H: number; M: number; S: number } | null }
  export function read(data: ArrayBuffer | Uint8Array, options: { type: 'array'; cellDates: boolean }): { SheetNames: string[]; Sheets: Record<string, unknown> }
  export const utils: {
    sheet_to_json<T>(sheet: unknown, options: { header: number; defval: null; blankrows: boolean }): T[]
  }
}
