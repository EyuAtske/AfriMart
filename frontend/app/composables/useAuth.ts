import type { User, LoginDTO, RegisterDTO } from '~/types/auth'
import { useMockDataStore } from '~/repositories/mock/MockDataStore'
import { useRepositories } from '~/composables/useRepositories'

export type AuthUser = User

export const useAuth = () => {
  const { user, isLoggedIn } = useMockDataStore()
  const { authRepo } = useRepositories()

  const login = async (credentials: LoginDTO) => {
    if (!credentials.email.trim() || !credentials.password.trim()) {
      throw new Error('Please enter your email and password.')
    }

    await authRepo.login(credentials)
    await navigateTo('/profile')
  }

  const register = async (details: RegisterDTO) => {
    const session = await authRepo.register(details)
    if (isLoggedIn.value) {
      await navigateTo('/profile')
    }
    return session
  }

  const logout = async () => {
    await authRepo.logout()
    await navigateTo('/account')
  }

  return {
    isLoggedIn,
    user,
    login,
    register,
    logout
  }
}
