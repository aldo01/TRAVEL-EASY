// TypeScript interfaces for the application

export interface Destination {
  id: string;
  name: string;
  country: string;
  imageUrl?: string;
  description: string;
  rating: number;
}

export const TripStatus = {
  UPCOMING: 'Upcoming',
  ONGOING: 'Ongoing',
  COMPLETED: 'Completed',
  CANCELLED: 'Cancelled'
} as const

export type TripStatus = typeof TripStatus[keyof typeof TripStatus]

export interface Trip {
  id: string;
  destination: string;
  startDate: Date;
  endDate: Date;
  status: TripStatus;
  travelers: number;
}

export const BookingType = {
  FLIGHT: 'Flight',
  HOTEL: 'Hotel',
  CAR: 'Car Rental',
  ACTIVITY: 'Activity'
} as const

export type BookingType = typeof BookingType[keyof typeof BookingType]

export interface Booking {
  id: string;
  type: BookingType;
  title: string;
  details: string;
  date: Date;
  confirmationNumber?: string;
}

export interface Activity {
  id: string;
  name: string;
  description: string;
  icon: string;
  price?: number;
}

export interface User {
  id: string;
  name: string;
  email: string;
  phoneNumber?: string;
  profileImageUrl?: string;
}
