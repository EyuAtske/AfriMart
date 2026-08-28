type MockApiOptions = {
  delay?: number
}

export const useApiClient = () => {
  const config = useRuntimeConfig()
  const apiBase = (config.public.apiBase as string) || 'http://localhost:8080'
  const authMode = (config.public.authMode as string) || 'mock'

  const getAuthHeaders = (token?: string | null): Record<string, string> => {
    const headers: Record<string, string> = {
      'Content-Type': 'application/json'
    }
    if (token) {
      headers.Authorization = `Bearer ${token}`
    }
    return headers
  }

  const apiFetch = async <T>(endpoint: string, options: Record<string, any> = {}): Promise<T> => {
    const url = `${apiBase.replace(/\/$/, '')}/${endpoint.replace(/^\//, '')}`
    return await $fetch<T>(url, {
      ...options,
      headers: {
        ...getAuthHeaders(options.token),
        ...(options.headers || {})
      }
    })
  }

  const mockRequest = async <T>(payload: T, options: MockApiOptions = {}) => {
    const delay = options.delay ?? 250

    await new Promise(resolve => setTimeout(resolve, delay))

    return {
      data: payload,
      error: null,
      meta: {
        apiBase,
        mocked: true
      }
    }
  }

  return {
    apiBase,
    authMode,
    getAuthHeaders,
    apiFetch,
    mockRequest
  }
}
