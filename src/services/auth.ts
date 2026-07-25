import { apiService } from './api'
import type { User } from '../types'

type LoginResponse = {
  user: User
  token: string
}

const TOKEN_KEY = 'travel_easy_token'

export const authStore = {
  getToken(): string | null {
    try {
      return localStorage.getItem(TOKEN_KEY)
    } catch {
      return null
    }
  },
  setToken(token: string) {
    localStorage.setItem(TOKEN_KEY, token)
  },
  clear() {
    localStorage.removeItem(TOKEN_KEY)
  },
}

export const authService = {
  async login(email: string, password: string) {
    const data = await apiService.post<LoginResponse>('/auth/login', { email, password })
    authStore.setToken(data.token)
    return data.user
  },

  async register(name: string, email: string, password: string, phoneNumber?: string) {
    const data = await apiService.post<LoginResponse>('/auth/register', {
      name,
      email,
      password,
      phoneNumber: phoneNumber || '',
    })
    authStore.setToken(data.token)
    return data.user
  },

  async getProfile() {
    return apiService.get<User>('/profile')
  },
}
