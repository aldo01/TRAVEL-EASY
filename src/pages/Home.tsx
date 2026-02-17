import { useState } from 'react'
import { Search, MapPin, Star, TrendingUp } from 'lucide-react'
import type { Destination, Activity } from '../types'

export default function Home() {
  const [searchLocation, setSearchLocation] = useState('')
  const [selectedDate, setSelectedDate] = useState('')
  const [travelers, setTravelers] = useState(1)

  const featuredDestinations: Destination[] = [
    { id: '1', name: 'Paris', country: 'France', description: 'The City of Light awaits you', rating: 4.8 },
    { id: '2', name: 'Tokyo', country: 'Japan', description: 'Where tradition meets future', rating: 4.9 },
    { id: '3', name: 'New York', country: 'USA', description: 'The city that never sleeps', rating: 4.7 },
    { id: '4', name: 'Barcelona', country: 'Spain', description: 'Architecture and beaches', rating: 4.6 },
  ]

  const popularActivities: Activity[] = [
    { id: '1', name: 'City Tours', description: 'Explore cities with local guides', icon: 'map' },
    { id: '2', name: 'Food & Wine', description: 'Authentic culinary experiences', icon: 'utensils' },
    { id: '3', name: 'Adventure Sports', description: 'Thrilling outdoor activities', icon: 'mountain' },
    { id: '4', name: 'Cultural Sites', description: 'Museums and historical places', icon: 'landmark' },
  ]

  const handleSearch = () => {
    console.log('Searching:', { searchLocation, selectedDate, travelers })
  }

  return (
    <div className="space-y-12">
      {/* Hero Section with Search */}
      <div className="bg-gradient-to-r from-blue-600 to-blue-800 -mx-4 sm:-mx-6 lg:-mx-8 px-4 sm:px-6 lg:px-8 py-16 rounded-lg">
        <div className="max-w-4xl mx-auto">
          <h1 className="text-4xl md:text-5xl font-bold text-white text-center mb-4">
            Travel Easy, Travel Smart
          </h1>
          <p className="text-xl text-blue-100 text-center mb-8">
            Book flights, hotels, and activities at the best prices
          </p>

          {/* Search Box */}
          <div className="bg-white rounded-lg shadow-xl p-6">
            <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
              <div className="md:col-span-2">
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  <MapPin className="inline w-4 h-4 mr-1" />
                  Location
                </label>
                <input
                  type="text"
                  placeholder="Where do you want to go?"
                  value={searchLocation}
                  onChange={(e) => setSearchLocation(e.target.value)}
                  className="w-full px-4 py-3 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  Date
                </label>
                <input
                  type="date"
                  value={selectedDate}
                  onChange={(e) => setSelectedDate(e.target.value)}
                  className="w-full px-4 py-3 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  Travelers
                </label>
                <select
                  value={travelers}
                  onChange={(e) => setTravelers(Number(e.target.value))}
                  className="w-full px-4 py-3 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                >
                  {[1, 2, 3, 4, 5, 6].map((num) => (
                    <option key={num} value={num}>
                      {num} {num === 1 ? 'Traveler' : 'Travelers'}
                    </option>
                  ))}
                </select>
              </div>
            </div>
            <button
              onClick={handleSearch}
              className="w-full mt-6 bg-blue-600 text-white py-3 px-6 rounded-md hover:bg-blue-700 transition-colors font-medium flex items-center justify-center"
            >
              <Search className="mr-2" size={20} />
              Search
            </button>
          </div>
        </div>
      </div>

      {/* Stats Section */}
      <div className="grid grid-cols-2 md:grid-cols-4 gap-6 text-center">
        <div>
          <div className="text-3xl font-bold text-blue-600">2M+</div>
          <div className="text-gray-600">Bookings</div>
        </div>
        <div>
          <div className="text-3xl font-bold text-blue-600">150+</div>
          <div className="text-gray-600">Destinations</div>
        </div>
        <div>
          <div className="text-3xl font-bold text-blue-600">4.7</div>
          <div className="text-gray-600 flex items-center justify-center">
            <Star className="w-4 h-4 fill-yellow-400 text-yellow-400 mr-1" />
            Rating
          </div>
        </div>
        <div>
          <div className="text-3xl font-bold text-blue-600">24/7</div>
          <div className="text-gray-600">Support</div>
        </div>
      </div>

      {/* Featured Destinations */}
      <section>
        <div className="flex items-center justify-between mb-6">
          <h2 className="text-3xl font-bold text-gray-900">Featured Destinations</h2>
          <button className="text-blue-600 hover:text-blue-700 font-medium">View All →</button>
        </div>
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
          {featuredDestinations.map((destination) => (
            <div
              key={destination.id}
              className="bg-white rounded-lg shadow-md overflow-hidden hover:shadow-xl transition-shadow cursor-pointer"
            >
              <div className="h-48 bg-gradient-to-br from-blue-400 to-blue-600 flex items-center justify-center">
                <MapPin size={48} className="text-white" />
              </div>
              <div className="p-4">
                <h3 className="text-xl font-semibold text-gray-900">{destination.name}</h3>
                <p className="text-gray-600 text-sm mb-2">{destination.country}</p>
                <p className="text-gray-500 text-sm mb-3">{destination.description}</p>
                <div className="flex items-center">
                  <Star className="w-4 h-4 fill-yellow-400 text-yellow-400 mr-1" />
                  <span className="text-sm font-medium">{destination.rating}</span>
                </div>
              </div>
            </div>
          ))}
        </div>
      </section>

      {/* Popular Activities */}
      <section>
        <div className="flex items-center mb-6">
          <TrendingUp className="w-8 h-8 text-blue-600 mr-3" />
          <h2 className="text-3xl font-bold text-gray-900">Popular Activities</h2>
        </div>
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          {popularActivities.map((activity) => (
            <div
              key={activity.id}
              className="bg-white rounded-lg shadow-md p-6 hover:shadow-lg transition-shadow cursor-pointer flex items-center"
            >
              <div className="w-16 h-16 bg-green-100 rounded-lg flex items-center justify-center mr-4">
                <TrendingUp className="w-8 h-8 text-green-600" />
              </div>
              <div className="flex-1">
                <h3 className="text-lg font-semibold text-gray-900">{activity.name}</h3>
                <p className="text-gray-600 text-sm">{activity.description}</p>
              </div>
              <div className="text-gray-400">→</div>
            </div>
          ))}
        </div>
      </section>

      {/* Why Choose Us */}
      <section className="bg-gray-50 -mx-4 sm:-mx-6 lg:-mx-8 px-4 sm:px-6 lg:px-8 py-12 rounded-lg">
        <h2 className="text-3xl font-bold text-gray-900 text-center mb-8">Why Travel Easy?</h2>
        <div className="grid grid-cols-1 md:grid-cols-4 gap-6">
          <div className="text-center">
            <div className="w-16 h-16 bg-blue-100 rounded-full flex items-center justify-center mx-auto mb-4">
              <Star className="w-8 h-8 text-blue-600" />
            </div>
            <h3 className="font-semibold text-gray-900 mb-2">Best Prices</h3>
            <p className="text-sm text-gray-600">Guaranteed lowest prices on all bookings</p>
          </div>
          <div className="text-center">
            <div className="w-16 h-16 bg-blue-100 rounded-full flex items-center justify-center mx-auto mb-4">
              <MapPin className="w-8 h-8 text-blue-600" />
            </div>
            <h3 className="font-semibold text-gray-900 mb-2">6000+ Locations</h3>
            <p className="text-sm text-gray-600">Worldwide coverage in major cities</p>
          </div>
          <div className="text-center">
            <div className="w-16 h-16 bg-blue-100 rounded-full flex items-center justify-center mx-auto mb-4">
              <Search className="w-8 h-8 text-blue-600" />
            </div>
            <h3 className="font-semibold text-gray-900 mb-2">Easy Booking</h3>
            <p className="text-sm text-gray-600">Simple, fast online booking process</p>
          </div>
          <div className="text-center">
            <div className="w-16 h-16 bg-blue-100 rounded-full flex items-center justify-center mx-auto mb-4">
              <TrendingUp className="w-8 h-8 text-blue-600" />
            </div>
            <h3 className="font-semibold text-gray-900 mb-2">24/7 Support</h3>
            <p className="text-sm text-gray-600">Customer service always available</p>
          </div>
        </div>
      </section>
    </div>
  )
}
