import { trace } from '@opentelemetry/api'
let inMemoryAccessToken: string | null = null
let inMemoryRefreshToken: string | null = null

const TOKEN_KEY = 'afrimart_access_token'
const REFRESH_KEY = 'afrimart_refresh_token'

export function getApiBase(): string {
  const config = useRuntimeConfig()
  return ((config.public?.apiBase as string) || '').replace(/\/$/, '')
}

export function getAccessToken(): string | null {
  if (!import.meta.client) return null
  if (inMemoryAccessToken) return inMemoryAccessToken
  try {
    const cookie = useCookie<string | null>(TOKEN_KEY)
    if (cookie.value) {
      inMemoryAccessToken = cookie.value.replace(/[\r\n]/g, '').trim()
      return inMemoryAccessToken
    }
  } catch {
    // ignore cookie error
  }
  try {
    const stored = localStorage.getItem(TOKEN_KEY)
    if (stored) {
      inMemoryAccessToken = stored.replace(/[\r\n]/g, '').trim()
      return inMemoryAccessToken
    }
  } catch {
    // ignore storage error
  }
  return null
}

export function setAccessToken(token: string | null): void {
  if (!import.meta.client) return
  const clean = token ? token.replace(/[\r\n]/g, '').trim() : null
  inMemoryAccessToken = clean
  try {
    const cookie = useCookie<string | null>(TOKEN_KEY, { sameSite: 'lax', maxAge: 60 * 60 * 24 * 7 })
    cookie.value = clean
  } catch { }
  try {
    if (clean) localStorage.setItem(TOKEN_KEY, clean)
    else localStorage.removeItem(TOKEN_KEY)
  } catch { }
}

export function getRefreshTokenValue(): string | null {
  if (!import.meta.client) return null
  if (inMemoryRefreshToken) return inMemoryRefreshToken
  try {
    const cookie = useCookie<string | null>(REFRESH_KEY)
    if (cookie.value) {
      inMemoryRefreshToken = cookie.value.replace(/[\r\n]/g, '').trim()
      return inMemoryRefreshToken
    }
  } catch { }
  try {
    const stored = localStorage.getItem(REFRESH_KEY)
    if (stored) {
      inMemoryRefreshToken = stored.replace(/[\r\n]/g, '').trim()
      return inMemoryRefreshToken
    }
  } catch { }
  return null
}

export function setRefreshToken(token: string | null): void {
  if (!import.meta.client) return
  const clean = token ? token.replace(/[\r\n]/g, '').trim() : null
  inMemoryRefreshToken = clean
  try {
    const cookie = useCookie<string | null>(REFRESH_KEY, { sameSite: 'lax', maxAge: 60 * 60 * 24 * 30 })
    cookie.value = clean
  } catch { }
  try {
    if (clean) localStorage.setItem(REFRESH_KEY, clean)
    else localStorage.removeItem(REFRESH_KEY)
  } catch { }
}

export function clearTokens(): void {
  if (!import.meta.client) return
  inMemoryAccessToken = null
  inMemoryRefreshToken = null
  try {
    const cookieToken = useCookie<string | null>(TOKEN_KEY)
    cookieToken.value = null
    const cookieRefresh = useCookie<string | null>(REFRESH_KEY)
    cookieRefresh.value = null
  } catch { }
  try {
    localStorage.removeItem(TOKEN_KEY)
    localStorage.removeItem(REFRESH_KEY)
  } catch { }
}

/**
 * Determines whether an error message is a safe, user-facing validation string
 * rather than an internal stack trace, database error, or unknown payload.
 */
function isSafeValidationMessage(msg: unknown): boolean {
  if (typeof msg !== 'string') return false
  const trimmed = msg.trim()
  if (!trimmed || trimmed.length > 200) return false

  // Reject HTML error pages (e.g. from proxies or web servers)
  if (/<[^>]+>/.test(trimmed)) return false

  // Reject database, SQL, panic, and internal system error leaks
  const unsafePatterns = [
    /\brelation\s+["']?[\w.-]+["']?\s+does\s+not\s+exist/i,
    /\b(syntax error|sqlstate|postgresql|pg_|pq:|pq\b|sql:|sqlite|mysql|mongodb)\b/i,
    /\b(foreign key|unique constraint|null constraint|check constraint)\b/i,
    /\b(select|insert|update|delete|drop|alter|truncate)\s+[a-z0-9_*]/i,
    /\b(panic:|runtime error:|segmentation fault|null pointer|goroutine \d+)\b/i,
    /(?:at\s+[\w/\\.-]+:\d+|\.go:\d+|\.ts:\d+|\.js:\d+)/i
  ]

  return !unsafePatterns.some(pattern => pattern.test(trimmed))
}

export function extractError(err: any, fallback: string): string {
  const status = err?.response?.status || err?.statusCode || err?.status
  const raw = err?.data?.error || err?.data?.message || err?.message

  // 5xx: return only the generic user message (no console.error of raw error)
  if (status && status >= 500) {
    return 'Something went wrong. Please try again later.'
  }

  // 4xx: show only known safe validation messages; use a generic message for unknown payloads
  if (status && status >= 400 && status < 500) {
    if (isSafeValidationMessage(raw)) {
      return (raw as string).trim()
    }
    return isSafeValidationMessage(fallback) ? fallback : 'Invalid request. Please check your inputs and try again.'
  }

  // Non-HTTP or untyped errors (e.g. network disconnection)
  if (isSafeValidationMessage(raw)) {
    return (raw as string).trim()
  }

  return isSafeValidationMessage(fallback) ? fallback : 'Something went wrong. Please try again later.'
}

/**
 * Typed error for 403 Forbidden responses.
 * Callers can check `instanceof AuthForbiddenError` to show
 * a permission message without clearing the session.
 */
export class AuthForbiddenError extends Error {
  constructor(message: string) {
    super(message)
    this.name = 'AuthForbiddenError'
  }
}

/**
 * Clears auth state in MockDataStore and navigates to /account.
 * Called on unrecoverable 401.
 */
function clearSessionAndRedirect(): void {
  clearTokens()
  try {
    const { user, isLoggedIn, cart, orders } = useMockDataStore()
    user.value = { username: '', name: '', email: '', role: 'buyer' }
    isLoggedIn.value = false
    cart.value = []
    orders.value = []
  } catch {
    // Guard against contexts where the composable is unavailable
  }
  if (import.meta.client && typeof navigateTo === 'function') {
    try {
      navigateTo('/account')
    } catch {
      // Guard against environments without Nuxt app context
    }
  }
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
    // Refresh token is invalid or expired
    return false
  }
}

/**
 * Sanitizes endpoint paths to prevent Markdown syntax, brackets, parentheses,
 * angle brackets, quotes, and path traversal in request URLs.
 */
export function sanitizeEndpoint(endpoint: string): string {
  return endpoint
    .replace(/[()[\]<>{}`'"]/g, '')
    .replace(/^\//, '')
    .replace(/\.\.\//g, '')
}

export async function authenticatedFetch<T>(
  endpoint: string,
  options: Record<string, any> = {},
): Promise<T> {
  if (!import.meta.client) {
    throw new Error(
      'Authenticated requests are not available during server-side rendering.',
    )
  }

  const apiBase = getApiBase()
  const cleanEndpoint = sanitizeEndpoint(endpoint)
  const url = `${apiBase}/${cleanEndpoint}`

  const token = getAccessToken()
  const isFormData =
    typeof FormData !== 'undefined' &&
    options.body instanceof FormData

  const headers: Record<string, string> = {
    ...(isFormData ? {} : { 'Content-Type': 'application/json' }),
    ...(options.headers || {}),
  }

  if (!headers.Authorization && token) {
    headers.Authorization = `Bearer ${token}`
  }

  // ---------------------------------------------------------------------------
  // Frontend API observability
  // ---------------------------------------------------------------------------

  const tracer = trace.getTracer('afrimart-frontend')
  const method = String(options.method || 'GET').toUpperCase()

  return tracer.startActiveSpan(
    `HTTP ${method} ${cleanEndpoint}`,
    async (span) => {
      span.setAttributes({
        'http.request.method': method,
        'url.path': `/${cleanEndpoint}`,
      })

      let emitLog: any
      try {
        if (typeof useNuxtApp === 'function') {
          const nuxtApp = useNuxtApp() as any
          emitLog = nuxtApp?.$emitOtelLog
        }
      } catch {}

      try {
        const response = await $fetch<T>(url, {
          ...options,
          headers,
        })

        span.setAttribute(
          'http.response.status_code',
          200,
        )

        emitLog?.(
          span,
          'INFO',
          'API request completed',
          {
            method,
            path: `/${cleanEndpoint}`,
            status: '200',
          },
        )

        span.end()

        return response
      } catch (err: any) {
        const status =
          err?.response?.status ||
          err?.statusCode ||
          err?.status ||
          0

        span.setAttribute(
          'http.response.status_code',
          Number(status),
        )

        emitLog?.(
          span,
          'ERROR',
          'API request failed',
          {
            method,
            path: `/${cleanEndpoint}`,
            status: String(status),
          },
        )

        // --- 403 Forbidden: keep session, surface permission error ---
        if (status === 403) {
          span.end()

          throw new AuthForbiddenError(
            extractError(
              err,
              "You don't have permission to perform this action.",
            ),
          )
        }

        // --- 401 Unauthorized: try refresh, then clear session ---
        if (status === 401) {
          if (getRefreshTokenValue()) {
            const refreshed = await tryRefreshToken()

            if (refreshed) {
              const newToken = getAccessToken()

              if (newToken) {
                headers.Authorization = `Bearer ${newToken}`
              }

              try {
                const retryResponse = await $fetch<T>(url, {
                  ...options,
                  headers,
                })

                span.setAttribute(
                  'http.response.status_code',
                  200,
                )

                emitLog?.(
                  span,
                  'INFO',
                  'API request completed after token refresh',
                  {
                    method,
                    path: `/${cleanEndpoint}`,
                    status: '200',
                  },
                )

                span.end()

                return retryResponse
              } catch (retryErr: any) {
                const retryStatus =
                  retryErr?.response?.status ||
                  retryErr?.statusCode ||
                  retryErr?.status ||
                  0

                span.setAttribute(
                  'http.response.status_code',
                  Number(retryStatus),
                )

                emitLog?.(
                  span,
                  'ERROR',
                  'API request failed after token refresh',
                  {
                    method,
                    path: `/${cleanEndpoint}`,
                    status: String(retryStatus),
                  },
                )

                span.end()

                throw retryErr
              }
            }
          }

          clearSessionAndRedirect()
        }

        span.end()

        throw err
      }
    },
  )
}


import { useMockDataStore } from '../mock/MockDataStore'
