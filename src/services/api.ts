// API Service Layer for Travel Easy (Locker Storage)

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8081/api/v1'

type ApiEnvelope<T> = {
  success: boolean
  message?: string
  data?: T
  error?: string
}

const getAuthToken = () => {
  try {
    return localStorage.getItem('travel_easy_token')
  } catch {
    return null
  }
}

class ApiService {
  private async request<T>(
    endpoint: string,
    options?: RequestInit
  ): Promise<T> {
    const url = `${API_BASE_URL}${endpoint}`

    const token = getAuthToken()
    
    const config: RequestInit = {
      ...options,
      headers: {
        'Content-Type': 'application/json',
        ...(token ? { Authorization: `Bearer ${token}` } : {}),
        ...options?.headers,
      },
    }

    try {
      const response = await fetch(url, config)

      if (response.status === 401) {
        localStorage.removeItem('travel_easy_token')
        throw new Error('Session expired. Please log in again.')
      }

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`)
      }

      const envelope = (await response.json()) as ApiEnvelope<T>

      if (!envelope.success) {
        throw new Error(envelope.error || 'Request failed')
      }

      return envelope.data as T
    } catch (error) {
      console.error('API request failed:', error)
      throw error
    }
  }

  async get<T>(endpoint: string): Promise<T> {
    return this.request<T>(endpoint, { method: 'GET' })
  }

  async post<T>(endpoint: string, body: unknown): Promise<T> {
    return this.request<T>(endpoint, {
      method: 'POST',
      body: JSON.stringify(body),
    })
  }

  async put<T>(endpoint: string, body: unknown): Promise<T> {
    return this.request<T>(endpoint, {
      method: 'PUT',
      body: JSON.stringify(body),
    })
  }

  async delete<T>(endpoint: string): Promise<T> {
    return this.request<T>(endpoint, { method: 'DELETE' })
  }
}

export const apiService = new ApiService()
