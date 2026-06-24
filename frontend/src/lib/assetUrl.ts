const BASE = (import.meta.env.VITE_UPLOADS_BASE_URL as string | undefined) ?? ''

export function assetUrl(path?: string | null): string | undefined {
  if (!path) return undefined
  if (path.startsWith('http://') || path.startsWith('https://') || path.startsWith('//')) return path
  return `${BASE}${path}`
}
