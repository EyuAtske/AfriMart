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
    const { user, isLoggedIn } = useMockDataStore()

    user.value = {
      username: dto.username,
      name: `${dto.firstName} ${dto.lastName}`,
      email: dto.email,
      role: 'buyer'
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
}
