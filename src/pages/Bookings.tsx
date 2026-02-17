import { useState } from 'react'
import { Plane, Hotel, Car, Map, Calendar, CreditCard } from 'lucide-react'
import type { Booking } from '../types'
import { BookingType } from '../types'

export default function Bookings() {
  const [activeTab, setActiveTab] = useState<'upcoming' | 'past'>('upcoming')

  const bookings: Booking[] = [
    {
      id: '1',
      type: BookingType.FLIGHT,
      title: 'Paris - Charles de Gaulle',
      details: 'Air France AF1234 • Economy Class',
      date: new Date('2026-03-15T14:30:00'),
      confirmationNumber: 'AF1234XYZ',
    },
    {
      id: '2',
      type: BookingType.HOTEL,
      title: 'Hotel Le Marais',
      details: 'Deluxe Room • 7 nights',
      date: new Date('2026-03-15T15:00:00'),
      confirmationNumber: 'HLM789456',
    },
    {
      id: '3',
      type: BookingType.CAR,
      title: 'Compact Car Rental',
      details: 'Toyota Corolla or similar',
      date: new Date('2026-03-16T09:00:00'),
      confirmationNumber: 'CAR123ABC',
    },
    {
      id: '4',
      type: BookingType.ACTIVITY,
      title: 'Eiffel Tower Tour',
      details: 'Guided tour with priority access',
      date: new Date('2026-03-17T10:00:00'),
      confirmationNumber: 'EIF123ABC',
    },
  ]

  const pastBookings: Booking[] = [
    {
      id: '5',
      type: BookingType.FLIGHT,
      title: 'Barcelona - El Prat',
      details: 'Vueling VY8765',
      date: new Date('2026-01-05T10:00:00'),
      confirmationNumber: 'VY8765ZYX',
    },
    {
      id: '6',
      type: BookingType.HOTEL,
      title: 'Hotel Barcelona Princess',
      details: 'Standard Room • 5 nights',
      date: new Date('2026-01-05T14:00:00'),
      confirmationNumber: 'HBP456789',
    },
  ]

  const getBookingIcon = (type: BookingType) => {
    switch (type) {
      case BookingType.FLIGHT:
        return <Plane className="w-6 h-6" />
      case BookingType.HOTEL:
        return <Hotel className="w-6 h-6" />
      case BookingType.CAR:
        return <Car className="w-6 h-6" />
      case BookingType.ACTIVITY:
        return <Map className="w-6 h-6" />
    }
  }

  const getBookingColor = (type: BookingType) => {
    switch (type) {
      case BookingType.FLIGHT:
        return 'bg-blue-100 text-blue-600'
      case BookingType.HOTEL:
        return 'bg-purple-100 text-purple-600'
      case BookingType.CAR:
        return 'bg-orange-100 text-orange-600'
      case BookingType.ACTIVITY:
        return 'bg-green-100 text-green-600'
    }
  }

  const formatDateTime = (date: Date) => {
    return date.toLocaleDateString('en-US', {
      month: 'short',
      day: 'numeric',
      year: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    })
  }

  const currentBookings = activeTab === 'upcoming' ? bookings : pastBookings

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-3xl font-bold text-gray-900">My Bookings</h1>
        <p className="text-gray-600 mt-1">Track all your travel bookings in one place</p>
      </div>

      {/* Tabs */}
      <div className="border-b border-gray-200">
        <div className="flex space-x-8">
          <button
            onClick={() => setActiveTab('upcoming')}
            className={`pb-4 px-1 border-b-2 font-medium text-sm transition-colors ${
              activeTab === 'upcoming'
                ? 'border-blue-500 text-blue-600'
                : 'border-transparent text-gray-500 hover:text-gray-700'
            }`}
          >
            Upcoming ({bookings.length})
          </button>
          <button
            onClick={() => setActiveTab('past')}
            className={`pb-4 px-1 border-b-2 font-medium text-sm transition-colors ${
              activeTab === 'past'
                ? 'border-blue-500 text-blue-600'
                : 'border-transparent text-gray-500 hover:text-gray-700'
            }`}
          >
            Past ({pastBookings.length})
          </button>
        </div>
      </div>

      {/* Bookings List */}
      <div className="grid gap-4">
        {currentBookings.map((booking) => (
          <div
            key={booking.id}
            className="bg-white rounded-lg shadow-md p-6 hover:shadow-lg transition-shadow"
          >
            <div className="flex items-start">
              <div
                className={`w-12 h-12 rounded-lg flex items-center justify-center mr-4 ${getBookingColor(
                  booking.type
                )}`}
              >
                {getBookingIcon(booking.type)}
              </div>

              <div className="flex-1">
                <div className="flex items-start justify-between mb-2">
                  <div>
                    <div className="flex items-center">
                      <span
                        className={`text-xs font-medium px-2 py-1 rounded ${getBookingColor(
                          booking.type
                        )} mr-2`}
                      >
                        {booking.type}
                      </span>
                      <h3 className="text-lg font-semibold text-gray-900">{booking.title}</h3>
                    </div>
                    <p className="text-gray-600 text-sm mt-1">{booking.details}</p>
                  </div>
                </div>

                <div className="flex items-center text-sm text-gray-500 mt-3 space-x-4">
                  <div className="flex items-center">
                    <Calendar className="w-4 h-4 mr-1" />
                    {formatDateTime(booking.date)}
                  </div>
                  {booking.confirmationNumber && (
                    <div className="flex items-center">
                      <CreditCard className="w-4 h-4 mr-1" />
                      Confirmation: {booking.confirmationNumber}
                    </div>
                  )}
                </div>

                <div className="mt-4 flex space-x-3">
                  <button className="text-blue-600 hover:text-blue-700 text-sm font-medium px-4 py-2 rounded-md hover:bg-blue-50 transition-colors">
                    View Details
                  </button>
                  <button className="text-gray-600 hover:text-gray-900 text-sm px-4 py-2 rounded-md hover:bg-gray-100 transition-colors">
                    Download Voucher
                  </button>
                </div>
              </div>
            </div>
          </div>
        ))}
      </div>

      {currentBookings.length === 0 && (
        <div className="text-center py-12">
          <div className="text-gray-400 mb-4">
            <Calendar size={48} className="mx-auto" />
          </div>
          <p className="text-gray-600">No {activeTab} bookings</p>
        </div>
      )}
    </div>
  )
}
