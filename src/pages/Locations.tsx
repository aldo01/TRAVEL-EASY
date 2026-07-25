import { useEffect, useMemo, useState } from 'react'
import { useNavigate, useSearchParams } from 'react-router-dom'
import { MapPin, Navigation, Lock, Clock } from 'lucide-react'
import { MapContainer, Marker, Popup, TileLayer } from 'react-leaflet'
import L from 'leaflet'
import type { Location, LocationWithDistance } from '../types'
import { storageService } from '../services/storage'
import BookingDialog from '../components/BookingDialog'

type AnyLocation = Location | LocationWithDistance

function isNearby(loc: AnyLocation): loc is LocationWithDistance {
  return typeof (loc as LocationWithDistance).distance === 'number'
}

// Fix Leaflet default marker icons under bundlers
const DefaultIcon = L.icon({
  iconUrl: 'https://unpkg.com/leaflet@1.9.4/dist/images/marker-icon.png',
  iconRetinaUrl: 'https://unpkg.com/leaflet@1.9.4/dist/images/marker-icon-2x.png',
  shadowUrl: 'https://unpkg.com/leaflet@1.9.4/dist/images/marker-shadow.png',
  iconSize: [25, 41],
  iconAnchor: [12, 41],
})

const SelectedIcon = L.icon({
  iconUrl: 'https://unpkg.com/leaflet@1.9.4/dist/images/marker-icon.png',
  iconRetinaUrl: 'https://unpkg.com/leaflet@1.9.4/dist/images/marker-icon-2x.png',
  shadowUrl: 'https://unpkg.com/leaflet@1.9.4/dist/images/marker-shadow.png',
  iconSize: [30, 49],
  iconAnchor: [15, 49],
})

export default function Locations() {
  const [searchParams] = useSearchParams()
  const navigate = useNavigate()
  const [query, setQuery] = useState('')
  const [radiusKm, setRadiusKm] = useState(10)
  const [locations, setLocations] = useState<AnyLocation[]>([])
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const [bookingOpen, setBookingOpen] = useState(false)
  const [bookingLocation, setBookingLocation] = useState<Location | null>(null)

  const [searchCenter, setSearchCenter] = useState<{ lat: number; lon: number } | null>(null)

  const loadByCity = async (cityOverride?: string) => {
    const city = (cityOverride ?? query).trim()
    setLoading(true)
    setError(null)
    setSearchCenter(null)
    try {
      const data = await storageService.getLocations({ city })
      setLocations(data)
    } catch (e) {
      setError(e instanceof Error ? e.message : 'Failed to load locations')
    } finally {
      setLoading(false)
    }
  }

  const loadNearby = async (radiusOverrideKm?: number) => {
    const radius = radiusOverrideKm ?? radiusKm

    setLoading(true)
    setError(null)

    try {
      const pos = await new Promise<GeolocationPosition>((resolve, reject) => {
        if (!navigator.geolocation) reject(new Error('Geolocation not supported'))
        navigator.geolocation.getCurrentPosition(resolve, reject, { enableHighAccuracy: true, timeout: 10000 })
      })

      const lat = pos.coords.latitude
      const lon = pos.coords.longitude
      setSearchCenter({ lat, lon })

      const data = await storageService.getNearbyLocations({ lat, lon, radiusKm: radius })
      setLocations(data)
    } catch (e) {
      setError(e instanceof Error ? e.message : 'Failed to find nearby locations')
    } finally {
      setLoading(false)
    }
  }

  const loadNearbyAt = async (lat: number, lon: number, radiusOverrideKm?: number) => {
    const radius = radiusOverrideKm ?? radiusKm

    setLoading(true)
    setError(null)
    setSearchCenter({ lat, lon })

    try {
      const data = await storageService.getNearbyLocations({ lat, lon, radiusKm: radius })
      setLocations(data)
    } catch (e) {
      setError(e instanceof Error ? e.message : 'Failed to find nearby locations')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    const city = searchParams.get('city')
    const nearby = searchParams.get('nearby')
    const radius = searchParams.get('radiusKm') || searchParams.get('radius')
    const latStr = searchParams.get('lat')
    const lonStr = searchParams.get('lon')
    const q = searchParams.get('q')

    let parsedRadius: number | undefined

    if (radius) {
      const parsed = Number(radius)
      if (Number.isFinite(parsed) && parsed > 0) {
        parsedRadius = parsed
        setRadiusKm(parsed)
      }
    }

    if (city && city.trim()) {
      setQuery(city)
      void loadByCity(city)
      return
    }

    if (q && q.trim()) {
      setQuery(q)
    }

    if (nearby === '1' || nearby === 'true') {
      if (latStr && lonStr) {
        const lat = Number(latStr)
        const lon = Number(lonStr)
        if (Number.isFinite(lat) && Number.isFinite(lon)) {
          void loadNearbyAt(lat, lon, parsedRadius)
        } else {
          void loadNearby(parsedRadius)
        }
      } else {
        void loadNearby(parsedRadius)
      }
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [searchParams])

  const useMyLocation = async () => loadNearby()

  const openDetails = (locationId: string) => {
    const qs = searchParams.toString()
    navigate(`/locations/${locationId}${qs ? `?${qs}` : ''}`)
  }

  const initialStart = useMemo(() => {
    const start = searchParams.get('start')
    if (!start) return undefined
    const d = new Date(start)
    return Number.isFinite(d.getTime()) ? d : undefined
  }, [searchParams])

  const openBooking = (loc: Location) => {
    setBookingLocation(loc)
    setBookingOpen(true)
  }

  const mapCenter = useMemo<[number, number]>(() => {
    if (searchCenter) return [searchCenter.lat, searchCenter.lon]
    if (locations.length > 0) return [locations[0].latitude, locations[0].longitude]
    return [55.6761, 12.5683]
  }, [locations, searchCenter])

  const mapMarkers = useMemo(() => {
    return locations.map((loc) => ({
      id: loc.id,
      lat: loc.latitude,
      lon: loc.longitude,
      name: loc.name,
      city: loc.city,
      country: loc.country,
      hourlyRate: loc.hourlyRate,
      dailyRate: loc.dailyRate,
      distance: isNearby(loc) ? loc.distance : null,
      available: isNearby(loc) ? loc.availableLockers : null,
    }))
  }, [locations])

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-3xl font-bold text-gray-900">Find luggage storage</h1>
        <p className="text-gray-600 mt-1">Search nearby lockers in hotels, cafes, shops, and kiosks.</p>
      </div>

      <div className="bg-white rounded-lg shadow-md p-6 space-y-4">
        <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
          <div className="md:col-span-2">
            <label className="block text-sm font-medium text-gray-700 mb-2">
              City / area
            </label>
            <input
              value={query}
              onChange={(e) => setQuery(e.target.value)}
              placeholder="e.g. Copenhagen"
              className="w-full px-4 py-3 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">Radius (km)</label>
            <input
              type="number"
              min={1}
              max={50}
              value={radiusKm}
              onChange={(e) => setRadiusKm(Number(e.target.value))}
              className="w-full px-4 py-3 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            />
          </div>
          <div className="flex items-end gap-2">
            <button
              onClick={() => {
                void loadByCity()
              }}
              disabled={loading || !query.trim()}
              className="flex-1 bg-blue-600 text-white py-3 px-4 rounded-md hover:bg-blue-700 disabled:opacity-50"
            >
              Search
            </button>
            <button
              onClick={useMyLocation}
              disabled={loading}
              className="px-4 py-3 border border-gray-300 rounded-md hover:bg-gray-50 disabled:opacity-50"
              title="Use my location"
            >
              <Navigation size={18} />
            </button>
          </div>
        </div>

        {error && <div className="text-sm text-red-600">{error}</div>}
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <div className="lg:col-span-1 space-y-4">
          <div className="bg-white rounded-lg shadow-md p-4">
            <div className="text-sm font-medium text-gray-900">Results</div>
            <div className="text-sm text-gray-600 mt-1">
              {loading ? 'Searching…' : `${locations.length} locations`}
            </div>
          </div>

          {locations.length === 0 && !loading && (
            <div className="bg-white rounded-lg shadow-md p-6 text-gray-600">
              No locations yet. Try “Use my location” or search a city.
            </div>
          )}

          {locations.map((loc) => (
            <div
              key={loc.id}
              role="button"
              tabIndex={0}
              onClick={() => openDetails(loc.id)}
              onKeyDown={(e) => {
                if (e.key === 'Enter' || e.key === ' ') {
                  e.preventDefault()
                  openDetails(loc.id)
                }
              }}
              className={`w-full text-left bg-white rounded-lg shadow-md p-4 hover:shadow-lg transition-shadow border cursor-pointer ${
                'border-transparent'
              }`}
            >
              <div className="flex items-start justify-between gap-3">
                <div>
                  <div className="flex items-center gap-2">
                    <MapPin className="w-5 h-5 text-blue-600" />
                    <h3 className="text-base font-semibold text-gray-900">{loc.name}</h3>
                  </div>
                  <p className="text-sm text-gray-500 mt-1">{loc.city}, {loc.country}</p>
                  <div className="mt-2 text-sm text-gray-600 flex items-center gap-2">
                    <Lock className="w-4 h-4" />
                    {loc.hourlyRate}/hr • {loc.dailyRate}/day
                  </div>
                </div>

                <div className="text-right text-xs text-gray-600">
                  {isNearby(loc) && <div>{loc.distance} km</div>}
                  {isNearby(loc) && <div className="mt-1 font-medium text-gray-900">{loc.availableLockers} available</div>}

                  <button
                    type="button"
                    onClick={(e) => {
                      e.stopPropagation()
                      openBooking(loc)
                    }}
                    className="mt-2 inline-flex items-center justify-center rounded-md bg-blue-600 px-3 py-2 text-xs font-medium text-white hover:bg-blue-700"
                  >
                    Book
                  </button>
                </div>
              </div>

              <div className="mt-3 text-xs text-gray-500 flex items-center gap-2">
                <Clock className="w-4 h-4" />
                {loc.isOpen24Hours ? 'Open 24/7' : `${loc.openingTime || '—'} - ${loc.closingTime || '—'}`}
              </div>
            </div>
          ))}
        </div>

        <div className="lg:col-span-2">
          <div className="bg-white rounded-lg shadow-md overflow-hidden lg:sticky lg:top-6">
            <div className="h-[320px] lg:h-[360px]">
              <MapContainer center={mapCenter} zoom={13} style={{ height: '100%', width: '100%' }} scrollWheelZoom>
                <TileLayer
                  attribution='&copy; OpenStreetMap contributors'
                  url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
                />

                {searchCenter && (
                  <Marker position={[searchCenter.lat, searchCenter.lon]} icon={SelectedIcon}>
                    <Popup>
                      <div className="text-sm">
                        <div className="font-medium">Search area</div>
                        <div className="text-gray-600">{searchCenter.lat.toFixed(4)}, {searchCenter.lon.toFixed(4)}</div>
                      </div>
                    </Popup>
                  </Marker>
                )}

                {mapMarkers.map((m) => (
                  <Marker
                    key={m.id}
                    position={[m.lat, m.lon]}
                    icon={DefaultIcon}
                    eventHandlers={{
                      click: () => {
                        openDetails(m.id)
                      },
                    }}
                  >
                    <Popup>
                      <div className="text-sm">
                        <div className="font-medium">{m.name}</div>
                        <div className="text-gray-600">{m.city}, {m.country}</div>
                        <div className="mt-1">{m.hourlyRate}/hr • {m.dailyRate}/day</div>
                        {m.distance !== null && <div className="text-gray-600 mt-1">{m.distance} km away</div>}
                        {m.available !== null && <div className="text-gray-600">{m.available} available</div>}
                      </div>
                    </Popup>
                  </Marker>
                ))}
              </MapContainer>
            </div>
          </div>
        </div>
      </div>

      {bookingLocation && (
        <BookingDialog
          open={bookingOpen}
          onClose={() => setBookingOpen(false)}
          location={bookingLocation}
          initialStart={initialStart}
        />
      )}
    </div>
  )
}
