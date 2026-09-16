import type { IAuthRepository } from '../interfaces/IAuthRepository'
import type { User, LoginDTO, RegisterDTO, AuthSession } from '~/types/auth'
import { useMockDataStore } from './MockDataStore'

export class MockAuthRepository implements IAuthRepository {
  async login(dto: LoginDTO): Promise<AuthSession> {
    const { user, isLoggedIn } = useMockDataStore()
    const email = dto.email.trim()

    const prefix = email.split('@')[0] || 'user'
    const fallbackName = prefix
      .replace(/[._-]+/g, ' ')
      .replace(/\b\w/g, letter => letter.toUpperCase()) || 'Test User'

    user.value = {
      username: fallbackName,
      name: fallbackName,
      email,
      role: 'buyer'
    }

    isLoggedIn.value = true

    return {
      user: { ...user.value },
      token: 'mock-jwt-token-' + Date.now()
    }
  }

  async register(dto: RegisterDTO): Promise<AuthSession> {
    const { user, isLoggedIn, users } = useMockDataStore()

    const newUser: User = {
      id: `usr-${Date.now()}`,
      username: dto.username,
      name: `${dto.firstName} ${dto.lastName}`,
      email: dto.email,
      role: 'buyer',
      created_at: new Date().toISOString().split('T')[0]
    }

    user.value = newUser
    if (!users.value.some(u => u.email === dto.email)) {
      users.value.push(newUser)
    }

    isLoggedIn.value = true

    return {
      user: { ...user.value },
      token: 'mock-jwt-token-' + Date.now()
    }
  }

  async logout(): Promise<void> {
    const { user, isLoggedIn } = useMockDataStore()

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
    return isLoggedIn.value ? { ...user.value } : null
  }

  async getProfile(): Promise<{ email: string; username: string }> {
    const { user, isLoggedIn } = useMockDataStore()
    if (!isLoggedIn.value) {
      throw new Error('User not authenticated')
    }
    return {
      email: user.value.email,
      username: user.value.username || user.value.name || 'User'
    }
  }

  async updateUsername(newUsername: string): Promise<User> {
    const { user, isLoggedIn } = useMockDataStore()
    if (!isLoggedIn.value) {
      throw new Error('User not authenticated')
    }
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

    user.value = {
      ...user.value,
      username: trimmed,
      name: trimmed
    }
    return { ...user.value }
  }

  async updatePassword(newPassword: string): Promise<void> {
    const { isLoggedIn } = useMockDataStore()
    if (!isLoggedIn.value) {
      throw new Error('User not authenticated')
    }
    if (!newPassword || newPassword.length < 8) {
      throw new Error('password must be at least 8 characters long')
    }
  }
}
