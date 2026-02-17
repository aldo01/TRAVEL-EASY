import { useState } from 'react'
import { Plus, Calendar, MapPin, Users } from 'lucide-react'
import type { Trip } from '../types'
import { TripStatus } from '../types'

export default function Trips() {
  const [trips] = useState<Trip[]>([
    {
      id: '1',
      destination: 'Paris, France',
      startDate: new Date('2026-03-15'),
      endDate: new Date('2026-03-22'),
      status: TripStatus.UPCOMING,
      travelers: 2,
    },
    {
      id: '2',
      destination: 'Tokyo, Japan',
      startDate: new Date('2026-05-10'),
      endDate: new Date('2026-05-20'),
      status: TripStatus.UPCOMING,
      travelers: 1,
    },
    {
      id: '3',
      destination: 'Barcelona, Spain',
      startDate: new Date('2026-01-05'),
      endDate: new Date('2026-01-10'),
      status: TripStatus.COMPLETED,
      travelers: 3,
    },
  ])

  const getStatusColor = (status: TripStatus) => {
    switch (status) {
      case TripStatus.UPCOMING:
        return 'bg-blue-100 text-blue-800'
      case TripStatus.ONGOING:
        return 'bg-green-100 text-green-800'
      case TripStatus.COMPLETED:
        return 'bg-gray-100 text-gray-800'
      case TripStatus.CANCELLED:
        return 'bg-red-100 text-red-800'
      default:
        return 'bg-gray-100 text-gray-800'
    }
  }

  const formatDate = (date: Date) => {
    return date.toLocaleDateString('en-US', {
      month: 'short',
      day: 'numeric',
      year: 'numeric',
    })
  }

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <div>
          <h1 className="text-3xl font-bold text-gray-900">My Trips</h1>
          <p className="text-gray-600 mt-1">Manage your upcoming and past trips</p>
        </div>
        <button className="bg-blue-600 text-white px-6 py-3 rounded-lg hover:bg-blue-700 transition-colors flex items-center font-medium shadow-md">
          <Plus className="mr-2" size={20} />
          New Trip
        </button>
      </div>

      <div className="grid gap-6">
        {trips.map((trip) => (
          <div
            key={trip.id}
            className="bg-white rounded-lg shadow-md p-6 hover:shadow-lg transition-shadow"
          >
            <div className="flex justify-between items-start">
              <div className="flex-1">
                <div className="flex items-center mb-3">
                  <MapPin className="w-5 h-5 text-blue-600 mr-2" />
                  <h3 className="text-xl font-semibold text-gray-900">{trip.destination}</h3>
                </div>

                <div className="grid grid-cols-1 md:grid-cols-3 gap-4 text-sm">
                  <div className="flex items-center text-gray-600">
                    <Calendar className="w-4 h-4 mr-2" />
                    <span>
                      {formatDate(trip.startDate)} - {formatDate(trip.endDate)}
                    </span>
                  </div>
                  <div className="flex items-center text-gray-600">
                    <Users className="w-4 h-4 mr-2" />
                    <span>
                      {trip.travelers} {trip.travelers === 1 ? 'Traveler' : 'Travelers'}
                    </span>
                  </div>
                  <div className="flex items-center text-gray-600">
                    <span className="text-gray-500">
                      {Math.ceil(
                        (trip.endDate.getTime() - trip.startDate.getTime()) / (1000 * 60 * 60 * 24)
                      )}{' '}
                      nights
                    </span>
                  </div>
                </div>
              </div>

              <span
                className={`px-3 py-1 rounded-full text-sm font-medium ${getStatusColor(
                  trip.status
                )}`}
              >
                {trip.status}
              </span>
            </div>

            <div className="mt-4 pt-4 border-t border-gray-200 flex justify-end space-x-3">
              <button className="text-gray-600 hover:text-gray-900 px-4 py-2 rounded-md hover:bg-gray-100 transition-colors">
                View Details
              </button>
              <button className="text-blue-600 hover:text-blue-700 px-4 py-2 rounded-md hover:bg-blue-50 transition-colors font-medium">
                Manage Trip
              </button>
            </div>
          </div>
        ))}
      </div>
    </div>
  )
}
