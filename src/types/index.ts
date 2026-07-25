// TypeScript interfaces for Travel Easy (Locker Storage)

export type LocationType =
  | 'KIOSK'
  | 'GYM'
  | 'CLUB'
  | 'SCHOOL'
  | 'HOTEL'
  | 'CAFE'
  | 'STORE'
  | 'OTHER'

export interface Location {
  id: string
  name: string
  type: LocationType
  address: string
  city: string
  state?: string | null
  zipCode?: string | null
  country: string
  latitude: number
  longitude: number
  openingTime?: string
  closingTime?: string
  isOpen24Hours: boolean
  phoneNumber?: string | null
  email?: string | null
  description?: string | null
  imageUrl?: string | null
  amenities: string
  hourlyRate: number
  dailyRate: number
  isActive: boolean
  rating: number
  totalReviews: number
  createdAt: string
  updatedAt: string
}

export interface LocationWithDistance extends Location {
  distance: number
  availableLockers: number
}

export type LockerSize = 'SMALL' | 'MEDIUM' | 'LARGE' | 'XLARGE'
export type LockerStatus = 'AVAILABLE' | 'OCCUPIED' | 'RESERVED' | 'MAINTENANCE' | 'OUT_OF_SERVICE'

export interface Locker {
  id: string
  locationId: string
  lockerNumber: string
  size: LockerSize
  status: LockerStatus
  isOperational: boolean
}

export type RateType = 'HOURLY' | 'DAILY' | 'MULTI_DAY'
export type BookingStatus = 'PENDING' | 'CONFIRMED' | 'ACTIVE' | 'COMPLETED' | 'CANCELLED' | 'EXPIRED'
export type PaymentStatus = 'PENDING' | 'PAID' | 'PARTIALLY_PAID' | 'REFUNDED' | 'FAILED'

export interface StorageBooking {
  id: string
  userId: string
  locationId: string
  lockerId: string
  bookingNumber: string
  startTime: string
  endTime?: string | null
  rateType: RateType
  bags?: number
  basePrice: number
  totalPrice: number
  status: BookingStatus
  qrCode: string
  paymentStatus: PaymentStatus
  location?: Location
  locker?: Locker
}

export interface User {
  id: string
  name: string
  email: string
  phoneNumber?: string | null
  profileImageUrl?: string | null
}
