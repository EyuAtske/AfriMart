import type { User, LoginDTO, RegisterDTO } from '~/types/auth'
import { useMockDataStore } from '~/repositories/mock/MockDataStore'
import { useRepositories } from '~/composables/useRepositories'

export type AuthUser = User

export const useAuth = () => {
  const { user, isLoggedIn } = useMockDataStore()
  const { authRepo } = useRepositories()
  const { gtag } = useGtag()

  // Restore session from persisted tokens on first load (client-side only)
  const sessionRestored = useState<boolean>('auth-session-restored', () => false)
  if (import.meta.client && !sessionRestored.value && !isLoggedIn.value) {
    sessionRestored.value = true
    authRepo.getCurrentSession().catch(() => {
      // Silently fail — user stays logged out
    })
  }

const login = async (credentials: LoginDTO) => {
  if (!credentials.email.trim() || !credentials.password.trim()) {
    throw new Error('Please enter your email and password.')
  }

  await authRepo.login(credentials)

  if (import.meta.client) {
    gtag('event', 'login', {
      method: 'email'
    })
  }

  await navigateTo('/profile')
}
 const register = async (details: RegisterDTO) => {
  const session = await authRepo.register(details)

  if (import.meta.client) {
    gtag('event', 'sign_up', {
      method: 'email'
    })
  }

  if (isLoggedIn.value) {
    await navigateTo('/profile')
  }

  return session
}

  const logout = async () => {
    await authRepo.logout()
    await navigateTo('/account')
  }

  const fetchProfile = async () => {
    return await authRepo.getProfile()
  }

  const updateUsername = async (username: string) => {
    return await authRepo.updateUsername(username)
  }

  const updatePassword = async (password: string) => {
    return await authRepo.updatePassword(password)
  }

  return {
    isLoggedIn,
    user,
    login,
    register,
    logout,
    fetchProfile,
    updateUsername,
    updatePassword
  }
}
