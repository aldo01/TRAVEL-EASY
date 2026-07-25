import { apiService } from './api'
import type {
  Location,
  LocationWithDistance,
  Locker,
  StorageBooking,
  RateType,
  LocationType,
} from '../types'

export const storageService = {
  async getLocations(params?: { city?: string; type?: LocationType }) {
    const qs = new URLSearchParams()
    if (params?.city) qs.set('city', params.city)
    if (params?.type) qs.set('type', params.type)

    const data = await apiService.get<{ count: number; locations: Location[] }>(
      `/locations${qs.toString() ? `?${qs.toString()}` : ''}`
    )
    return data.locations
  },

  async getNearbyLocations(params: { lat: number; lon: number; radiusKm?: number }) {
    const qs = new URLSearchParams({
      lat: String(params.lat),
      lon: String(params.lon),
      radius: String(params.radiusKm ?? 5),
    })

    const data = await apiService.get<{ count: number; locations: LocationWithDistance[] }>(
      `/locations/nearby?${qs.toString()}`
    )
    return data.locations
  },

  async getLocation(id: string) {
    const data = await apiService.get<{ location: Location; availableLockers: number }>(`/locations/${id}`)
    return data
  },

  async getAvailableLockers(locationId: string) {
    const data = await apiService.get<{ count: number; lockers: Locker[] }>(
      `/locations/${locationId}/lockers?status=AVAILABLE`
    )
    return data.lockers
  },

  async createBooking(params: {
    locationId: string
    lockerId: string
    startTime: string
    durationHours: number
    rateType: RateType
    bags?: number
  }) {
    return apiService.post<StorageBooking>('/bookings', {
      locationId: params.locationId,
      lockerId: params.lockerId,
      startTime: params.startTime,
      duration: params.durationHours,
      rateType: params.rateType,
      bags: params.bags,
    })
  },

  async getMyBookings() {
    const data = await apiService.get<{ count: number; bookings: StorageBooking[] }>(`/bookings`)
    return data.bookings
  },

  async createLocation(params: {
    name: string
    type: LocationType
    address: string
    city: string
    country: string
    latitude: number
    longitude: number
    openingTime?: string
    closingTime?: string
    isOpen24Hours?: boolean
    phoneNumber?: string
    email?: string
    description?: string
    imageUrl?: string
    amenities?: string[]
    hourlyRate?: number
    dailyRate?: number
    lockers: { small: number; medium: number; large: number; xlarge: number }
  }) {
    return apiService.post<{ location: Location }>('/locations', {
      ...params,
      isOpen24Hours: params.isOpen24Hours ?? false,
      hourlyRate: params.hourlyRate ?? 0,
      dailyRate: params.dailyRate ?? 0,
      amenities: params.amenities ?? [],
    })
  },
}
