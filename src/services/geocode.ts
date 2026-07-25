import { apiService } from './api'

export type PlaceSuggestion = {
  lat: number
  lon: number
  title: string
  subtitle: string
  displayName: string
}

export const geocodeService = {
  async suggest(q: string, limit = 5) {
    const qs = new URLSearchParams({ q, limit: String(limit) })
    const data = await apiService.get<{ count: number; suggestions: PlaceSuggestion[] }>(`/geocode/suggest?${qs.toString()}`)
    return data.suggestions
  },
}
