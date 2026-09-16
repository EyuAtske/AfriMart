import type { User, LoginDTO, RegisterDTO, AuthSession } from '~/types/auth'

export interface IAuthRepository {
  login(dto: LoginDTO): Promise<AuthSession>
  register(dto: RegisterDTO): Promise<AuthSession>
  logout(): Promise<void>
  getCurrentSession(): Promise<User | null>
  getProfile(): Promise<{ email: string; username: string }>
  updateUsername(username: string): Promise<User>
  updatePassword(password: string): Promise<void>
}
