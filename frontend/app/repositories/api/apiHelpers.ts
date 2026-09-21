// Shared helpers for API repositories — token management and authenticated fetch

const ACCESS_TOKEN_KEY = 'afrimart_access_token'
const REFRESH_TOKEN_KEY = 'afrimart_refresh_token'

export function getApiBase(): string {
  const config = useRuntimeConfig()
  return ((config.public?.apiBase as string) || '').replace(/\/$/, '')
}

export function getAccessToken(): string | null {
  try {
    const cookie = useCookie(ACCESS_TOKEN_KEY)
    return cookie.value || null
  } catch {
    return null
  }
}

export function setAccessToken(token: string | null): void {
  try {
    const cookie = useCookie(ACCESS_TOKEN_KEY, { maxAge: 3600, path: '/', sameSite: 'lax' })
    cookie.value = token
  } catch {
    // SSR context where useCookie may not be available
  }
}

export function getRefreshTokenValue(): string | null {
  try {
    const cookie = useCookie(REFRESH_TOKEN_KEY)
    return cookie.value || null
  } catch {
    return null
  }
}

export function setRefreshToken(token: string | null): void {
  try {
    const cookie = useCookie(REFRESH_TOKEN_KEY, { maxAge: 60 * 60 * 24 * 60, path: '/', sameSite: 'lax' })
    cookie.value = token
  } catch {
    // SSR context where useCookie may not be available
  }
}

export function clearTokens(): void {
  setAccessToken(null)
  setRefreshToken(null)
}

export function extractError(err: any, fallback: string): string {
  return err?.data?.error || err?.data?.message || err?.message || fallback
}

/**
 * Attempts to use the refresh token to get a new access token.
 * Returns true if successful, false otherwise.
 */
async function tryRefreshToken(): Promise<boolean> {
  const refreshToken = getRefreshTokenValue()
  if (!refreshToken) return false

  const apiBase = getApiBase()
  try {
    const res = await $fetch<{ token: string }>(`${apiBase}/api/refresh`, {
      method: 'POST',
      headers: {
        Authorization: `Bearer ${refreshToken}`
      }
    })

    if (res?.token) {
      setAccessToken(res.token)
      return true
    }
    return false
  } catch {
    // Refresh token is invalid or expired — clear everything
    clearTokens()
    return false
  }
}

/**
 * Performs an authenticated fetch. If the request fails with 401,
 * attempts to refresh the access token and retries the request once.
 */
export async function authenticatedFetch<T>(endpoint: string, options: Record<string, any> = {}): Promise<T> {
  const apiBase = getApiBase()
  const url = `${apiBase}/${endpoint.replace(/^\//, '')}`

  const token = getAccessToken()
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
    ...(options.headers || {})
  }
  if (token) {
    headers.Authorization = `Bearer ${token}`
  }

  try {
    return await $fetch<T>(url, {
      ...options,
      headers
    })
  } catch (err: any) {
    // If 401 and we have a refresh token, try to refresh
    const status = err?.response?.status || err?.statusCode || err?.status
    if (status === 401 && getRefreshTokenValue()) {
      const refreshed = await tryRefreshToken()
      if (refreshed) {
        // Retry with new token
        const newToken = getAccessToken()
        if (newToken) {
          headers.Authorization = `Bearer ${newToken}`
        }
        return await $fetch<T>(url, {
          ...options,
          headers
        })
      }
    }
    throw err
  }
}
