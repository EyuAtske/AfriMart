import { ref, type Ref } from 'vue'
import type { IAuthRepository } from '../interfaces/IAuthRepository'
import type { User, LoginDTO, RegisterDTO, AuthSession, ApiAuthResponse } from '~/types/auth'
import { useMockDataStore } from '../mock/MockDataStore'

export class ApiAuthRepository implements IAuthRepository {
  private accessToken: Ref<string | null> = ref<string | null>(null)

  getMemoryToken(): string | null {
    return this.accessToken.value
  }

  async login(dto: LoginDTO): Promise<AuthSession> {
    const { user, isLoggedIn } = useMockDataStore()
    const config = useRuntimeConfig()
    const apiBase = (config.public.apiBase as string) || 'http://localhost:8080'

    try {
      const url = `${apiBase.replace(/\/$/, '')}/api/auth/login`
      const res = await $fetch<ApiAuthResponse>(url, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: {
          email: dto.email.trim(),
          password: dto.password.trim()
        }
      })

      if (!res || !res.token) {
        throw new Error('Invalid backend login response: missing authentication token.')
      }

      this.accessToken.value = res.token

      const fallbackName = (res.username || res.name || res.email.split('@')[0] || 'User')
        .replace(/[._-]+/g, ' ')
        .replace(/\b\w/g, letter => letter.toUpperCase())

      const authenticatedUser: User = {
        id: res.id,
        username: res.username || res.email.split('@')[0] || 'User',
        name: fallbackName,
        email: res.email,
        role: 'buyer',
        created_at: res.created_at,
        updated_at: res.updated_at
      }

      user.value = authenticatedUser
      isLoggedIn.value = true

      return {
        user: { ...authenticatedUser },
        token: res.token,
        refreshToken: res.refresh_token
      }
    } catch (err: any) {
      this.accessToken.value = null
      isLoggedIn.value = false
      const errorMsg = err?.data?.message || err?.message || 'Authentication failed'
      throw new Error(errorMsg)
    }
  }

  async register(dto: RegisterDTO): Promise<AuthSession> {
    const { user, isLoggedIn } = useMockDataStore()
    const config = useRuntimeConfig()
    const apiBase = (config.public.apiBase as string) || 'http://localhost:8080'

    try {
      const url = `${apiBase.replace(/\/$/, '')}/api/auth/register`
      const res = await $fetch<ApiAuthResponse>(url, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: {
          username: dto.username.trim(),
          firstName: dto.firstName.trim(),
          lastName: dto.lastName.trim(),
          email: dto.email.trim(),
          password: dto.password.trim()
        }
      })

      const registeredUser: User = {
        id: res?.id,
        username: dto.username.trim(),
        name: `${dto.firstName.trim()} ${dto.lastName.trim()}`,
        email: dto.email.trim(),
        role: 'buyer',
        created_at: res?.created_at,
        updated_at: res?.updated_at
      }

      // Do NOT set access token or isLoggedIn = true unless response includes token
      if (res && res.token) {
        this.accessToken.value = res.token
        user.value = registeredUser
        isLoggedIn.value = true
        return {
          user: { ...registeredUser },
          token: res.token,
          refreshToken: res.refresh_token
        }
      }

      // Tokens omitted: Return created user without setting isLoggedIn
      return {
        user: { ...registeredUser },
        token: ''
      }
    } catch (err: any) {
      this.accessToken.value = null
      isLoggedIn.value = false
      const errorMsg = err?.data?.message || err?.message || 'Registration failed'
      throw new Error(errorMsg)
    }
  }

  async logout(): Promise<void> {
    const { user, isLoggedIn } = useMockDataStore()

    // Clear local in-memory state only (live refresh-token cookie logout pending backend HttpOnly transport)
    this.accessToken.value = null
    user.value = {
      username: '',
      name: '',
      email: '',
      role: 'buyer'
    }
    isLoggedIn.value = false
  }

  async getCurrentSession(): Promise<User | null> {
    const { user, isLoggedIn } = useMockDataStore()
    if (isLoggedIn.value && this.accessToken.value) {
      return { ...user.value }
    }
    return null
  }
}
