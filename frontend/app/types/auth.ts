export type UserRole = 'buyer' | 'seller'

export interface User {
  id?: string | number
  username: string
  name: string
  email: string
  phoneNumber?: string
  role: UserRole
  created_at?: string
  updated_at?: string
}

export interface LoginDTO {
  email: string
  password: string
}

export interface RegisterDTO {
  username: string
  firstName: string
  lastName: string
  email: string
  phoneNumber: string
  password: string
}

export interface ApiAuthResponse {
  id?: string | number
  email: string
  username?: string
  name?: string
  created_at?: string
  updated_at?: string
  token?: string
  refresh_token?: string
}

export interface AuthSession {
  user: User
  token?: string
  refreshToken?: string
}
