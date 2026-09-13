import type { IAuthRepository } from '../interfaces/IAuthRepository'
import type { User, LoginDTO, RegisterDTO, AuthSession, ApiAuthResponse } from '~/types/auth'
import { useMockDataStore } from '../mock/MockDataStore'
import {
  getApiBase,
  getAccessToken,
  setAccessToken,
  getRefreshTokenValue,
  setRefreshToken,
  clearTokens,
  extractError,
  authenticatedFetch
} from './apiHelpers'

export class ApiAuthRepository implements IAuthRepository {

  async login(dto: LoginDTO): Promise<AuthSession> {
    const { user, isLoggedIn } = useMockDataStore()
    const apiBase = getApiBase()

    try {
      const url = `${apiBase}/api/auth/login`
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

      // Persist tokens in cookies
      setAccessToken(res.token)
      setRefreshToken(res.refresh_token || null)

      // Resolve the username — try fetching profile if not in the login response
      let username = res.username || ''
      if (!username) {
        try {
          const profile = await this.getProfile()
          if (profile?.username) {
            username = profile.username
          }
        } catch {
          // Fallback if profile endpoint is not available
        }
      }

      if (!username) {
        username = res.name || (res.email ? res.email.split('@')[0] : '') || 'User'
      }

      const displayName = (username || 'User')
        .replace(/[._-]+/g, ' ')
        .replace(/\b\w/g, letter => letter.toUpperCase())

      const authenticatedUser: User = {
        id: res.id,
        username,
        name: displayName,
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
      clearTokens()
      isLoggedIn.value = false
      throw new Error(extractError(err, 'Authentication failed'))
    }
  }

  async register(dto: RegisterDTO): Promise<AuthSession> {
    const { user, isLoggedIn } = useMockDataStore()
    const apiBase = getApiBase()

    try {
      const url = `${apiBase}/api/auth/register`
      const res = await $fetch<ApiAuthResponse>(url, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: {
          username: dto.username.trim(),
          firstname: dto.firstName.trim(),
          Lastname: dto.lastName.trim(),
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

      // Only set logged-in state if the backend returned a token
      if (res && res.token) {
        setAccessToken(res.token)
        setRefreshToken(res.refresh_token || null)
        user.value = registeredUser
        isLoggedIn.value = true
        return {
          user: { ...registeredUser },
          token: res.token,
          refreshToken: res.refresh_token
        }
      }

      // Backend did not return tokens (register-only flow) — user must log in
      return {
        user: { ...registeredUser },
        token: ''
      }
    } catch (err: any) {
      clearTokens()
      isLoggedIn.value = false
      throw new Error(extractError(err, 'Registration failed'))
    }
  }

  async logout(): Promise<void> {
    const { user, isLoggedIn } = useMockDataStore()
    const refreshToken = getRefreshTokenValue()

    try {
      if (refreshToken) {
        const apiBase = getApiBase()
        await $fetch(`${apiBase}/api/auth/logout`, {
          method: 'POST',
          headers: {
            Authorization: `Bearer ${refreshToken}`
          }
        })
      }
    } catch {
      // Proceed with local logout even if revocation fails
    } finally {
      clearTokens()
      user.value = {
        username: '',
        name: '',
        email: '',
        role: 'buyer'
      }
      isLoggedIn.value = false
    }
  }

  async getCurrentSession(): Promise<User | null> {
    const { user, isLoggedIn } = useMockDataStore()

    // If already hydrated in memory, return it
    if (isLoggedIn.value && user.value.email) {
      return { ...user.value }
    }

    // Attempt to restore from persisted cookies
    const token = getAccessToken()
    if (!token) return null

    try {
      const profile = await this.getProfile()
      if (profile) {
        const displayName = (profile.username || 'User')
          .replace(/[._-]+/g, ' ')
          .replace(/\b\w/g, letter => letter.toUpperCase())

        const restoredUser: User = {
          username: profile.username,
          name: displayName,
          email: profile.email,
          role: 'buyer'
        }
        user.value = restoredUser
        isLoggedIn.value = true
        return { ...restoredUser }
      }
    } catch {
      // Token is invalid — clear state
      clearTokens()
      isLoggedIn.value = false
    }

    return null
  }

  async getProfile(): Promise<{ email: string; username: string }> {
    const { user } = useMockDataStore()

    const res = await authenticatedFetch<{ email: string; username: string }>('api/user/profile', {
      method: 'GET'
    })

    if (res) {
      user.value.email = res.email
      if (res.username) {
        user.value.username = res.username
        user.value.name = res.username
      }
    }

    return res
  }

  async updateUsername(newUsername: string): Promise<User> {
    const { user } = useMockDataStore()

    const trimmed = newUsername.trim()
    if (!trimmed) {
      throw new Error('username is required')
    }
    if (trimmed.length < 3) {
      throw new Error('username must be at least 3 characters long')
    }
    if (trimmed.length > 50) {
      throw new Error('username must not exceed 50 characters')
    }

    await authenticatedFetch('api/auth/username', {
      method: 'PUT',
      body: {
        username: trimmed
      }
    })

    user.value = {
      ...user.value,
      username: trimmed,
      name: trimmed
    }

    return { ...user.value }
  }

  async updatePassword(newPassword: string): Promise<void> {
    if (!newPassword || newPassword.length < 8) {
      throw new Error('password must be at least 8 characters long')
    }

    await authenticatedFetch('api/auth/password', {
      method: 'PUT',
      body: {
        password: newPassword
      }
    })
  }
}
