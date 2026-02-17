import { apiService } from './api'
import type { Destination, Trip, Booking, Activity } from '../types'

export class TravelService {
  // Destinations
  async getDestinations(): Promise<Destination[]> {
    return apiService.get<Destination[]>('/destinations')
  }

  async getDestinationById(id: string): Promise<Destination> {
    return apiService.get<Destination>(`/destinations/${id}`)
  }

  async searchDestinations(query: string): Promise<Destination[]> {
    return apiService.get<Destination[]>(`/destinations/search?q=${encodeURIComponent(query)}`)
  }

  // Trips
  async getTrips(): Promise<Trip[]> {
    return apiService.get<Trip[]>('/trips')
  }

  async getTripById(id: string): Promise<Trip> {
    return apiService.get<Trip>(`/trips/${id}`)
  }

  async createTrip(trip: Omit<Trip, 'id'>): Promise<Trip> {
    return apiService.post<Trip>('/trips', trip)
  }

  async updateTrip(id: string, trip: Partial<Trip>): Promise<Trip> {
    return apiService.put<Trip>(`/trips/${id}`, trip)
  }

  async deleteTrip(id: string): Promise<void> {
    return apiService.delete<void>(`/trips/${id}`)
  }

  // Bookings
  async getBookings(): Promise<Booking[]> {
    return apiService.get<Booking[]>('/bookings')
  }

  async getBookingById(id: string): Promise<Booking> {
    return apiService.get<Booking>(`/bookings/${id}`)
  }

  async createBooking(booking: Omit<Booking, 'id'>): Promise<Booking> {
    return apiService.post<Booking>('/bookings', booking)
  }

  async cancelBooking(id: string): Promise<void> {
    return apiService.delete<void>(`/bookings/${id}`)
  }

  // Activities
  async getActivities(): Promise<Activity[]> {
    return apiService.get<Activity[]>('/activities')
  }

  async getActivitiesByDestination(destinationId: string): Promise<Activity[]> {
    return apiService.get<Activity[]>(`/activities?destination=${destinationId}`)
  }
}

export const travelService = new TravelService()
