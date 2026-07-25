import { useEffect, useMemo, useState } from 'react'
import { Link, useNavigate, useParams, useSearchParams } from 'react-router-dom'
import { MapContainer, Marker, TileLayer } from 'react-leaflet'
import L from 'leaflet'
import { ArrowLeft, ExternalLink, Star } from 'lucide-react'
import type { Location } from '../types'
import { storageService } from '../services/storage'
import BookingDialog from '../components/BookingDialog'

const MarkerIcon = L.icon({
  iconUrl: 'https://unpkg.com/leaflet@1.9.4/dist/images/marker-icon.png',
  iconRetinaUrl: 'https://unpkg.com/leaflet@1.9.4/dist/images/marker-icon-2x.png',
  shadowUrl: 'https://unpkg.com/leaflet@1.9.4/dist/images/marker-shadow.png',
  iconSize: [25, 41],
  iconAnchor: [12, 41],
})

function parseAmenities(raw: string): string[] {
  if (!raw) return []
  try {
    const parsed = JSON.parse(raw)
    if (Array.isArray(parsed)) return parsed.filter((x) => typeof x === 'string')
  } catch {
    // ignore
  }
  return []
}

function parseStartParam(start: string | null): Date | undefined {
  if (!start) return undefined
  const d = new Date(start)
  if (Number.isNaN(d.getTime())) return undefined
  return d
}

function dayLabel(d: Date) {
  const today = new Date()
  today.setHours(0, 0, 0, 0)
  const day = new Date(d)
  day.setHours(0, 0, 0, 0)
  const diff = (day.getTime() - today.getTime()) / (1000 * 60 * 60 * 24)
  if (diff === 0) return 'Today'
  if (diff === 1) return 'Tomorrow'
  return day.toLocaleDateString(undefined, { weekday: 'short', month: 'short', day: 'numeric' })
}

export default function LocationDetails() {
  const params = useParams()
  const [searchParams] = useSearchParams()
  const navigate = useNavigate()

  const locationId = params.id

  const [location, setLocation] = useState<Location | null>(null)
  const [availableLockers, setAvailableLockers] = useState<number | null>(null)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const [bookingOpen, setBookingOpen] = useState(false)

  const initialStart = useMemo(() => parseStartParam(searchParams.get('start')), [searchParams])
  const bags = Number(searchParams.get('bags') ?? '2')
  const initialBags = Number.isFinite(bags) && bags > 0 ? Math.floor(bags) : 2
  const bagsLabel = `${initialBags} bag${initialBags === 1 ? '' : 's'}`

  useEffect(() => {
    if (!locationId) return

    let cancelled = false
    ;(async () => {
      setLoading(true)
      setError(null)
      try {
        const data = await storageService.getLocation(locationId)
        if (cancelled) return
        setLocation(data.location)
        setAvailableLockers(data.availableLockers)
      } catch (e) {
        if (cancelled) return
        setError(e instanceof Error ? e.message : 'Failed to load location')
      } finally {
        if (!cancelled) setLoading(false)
      }
    })()

    return () => {
      cancelled = true
    }
  }, [locationId])

  const openingHours = useMemo(() => {
    if (!location) return null
    if (location.isOpen24Hours) return 'Open 24 hours'
    const o = location.openingTime || '—'
    const c = location.closingTime || '—'
    return `Open ${o}-${c}`
  }, [location])

  const amenities = useMemo(() => (location ? parseAmenities(location.amenities) : []), [location])

  const mapsUrl = useMemo(() => {
    if (!location) return '#'
    const q = `${location.latitude},${location.longitude}`
    return `https://www.google.com/maps/search/?api=1&query=${encodeURIComponent(q)}`
  }, [location])

  if (!locationId) {
    return (
      <div className="bg-white rounded-lg shadow-md p-6">
        <div className="text-gray-900 font-medium">Missing location id</div>
      </div>
    )
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-3">
          <button
            onClick={() => navigate(-1)}
            className="p-2 rounded-md border border-gray-300 hover:bg-gray-50"
            aria-label="Back"
          >
            <ArrowLeft size={18} />
          </button>
          <div>
            <div className="text-sm text-gray-600">Location details</div>
            <div className="text-xl font-semibold text-gray-900">{location?.name ?? (loading ? 'Loading…' : '—')}</div>
          </div>
        </div>

        <div className="flex items-center gap-3">
          <div className="hidden md:flex items-center gap-2 text-sm text-gray-700">
            <span className="inline-flex items-center gap-1">
              {initialStart ? dayLabel(initialStart) : 'Select dates'}
            </span>
            <span className="text-gray-400">•</span>
            <span>{bagsLabel}</span>
          </div>
          <button
            onClick={() => setBookingOpen(true)}
            disabled={!location}
            className="bg-orange-600 text-white px-5 py-2 rounded-md hover:bg-orange-700 disabled:opacity-50"
          >
            Book
          </button>
        </div>
      </div>

      {error && (
        <div className="bg-white rounded-lg shadow-md p-4 text-sm text-red-600">{error}</div>
      )}

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <div className="lg:col-span-2 space-y-6">
          <div className="bg-white rounded-lg shadow-md p-6">
            <div className="flex items-start justify-between gap-4">
              <div>
                <div className="text-2xl font-semibold text-gray-900">{location?.name}</div>
                <div className="mt-1 text-sm text-gray-600">{location?.city}, {location?.country}</div>
                <div className="mt-2 text-sm text-green-700">{openingHours}</div>

                <div className="mt-3 flex items-center gap-2 text-sm text-gray-700">
                  <span className="inline-flex items-center gap-1">
                    <Star size={16} className="text-yellow-500" />
                    <span className="font-medium">{location?.rating?.toFixed?.(1) ?? '—'}</span>
                  </span>
                  <span className="text-gray-400">(</span>
                  <span className="text-gray-600">{location?.totalReviews ?? 0} ratings</span>
                  <span className="text-gray-400">)</span>
                  {availableLockers !== null && (
                    <>
                      <span className="text-gray-300">•</span>
                      <span className="text-gray-600">{availableLockers} available</span>
                    </>
                  )}
                </div>
              </div>

              <div className="text-right text-sm text-gray-700">
                <div className="font-medium">Pricing</div>
                <div className="mt-1">{location ? `${location.hourlyRate}/hour per bag` : '—'}</div>
                <div>{location ? `${location.dailyRate}/day per bag` : '—'}</div>
              </div>
            </div>
          </div>

          <div className="bg-white rounded-lg shadow-md p-6">
            <div className="text-lg font-semibold text-gray-900">Opening hours</div>
            <div className="mt-4 grid grid-cols-2 gap-y-2 text-sm">
              {['Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday', 'Sunday'].map((day) => (
                <div key={day} className="contents">
                  <div className="text-gray-700">{day}</div>
                  <div className="text-gray-900 text-right">
                    {location?.isOpen24Hours
                      ? '00-24'
                      : `${location?.openingTime || '—'}-${location?.closingTime || '—'}`}
                  </div>
                </div>
              ))}
            </div>
          </div>

          <div className="bg-white rounded-lg shadow-md p-6">
            <div className="text-lg font-semibold text-gray-900">About</div>
            {location?.description && (
              <div className="mt-2 text-sm text-gray-700">{location.description}</div>
            )}

            <div className="mt-4 grid grid-cols-1 md:grid-cols-2 gap-3 text-sm text-gray-700">
              {amenities.length > 0 ? (
                amenities.map((a) => (
                  <div key={a} className="border rounded-md px-3 py-2">{a}</div>
                ))
              ) : (
                <div className="text-gray-600">No amenities listed.</div>
              )}
            </div>

            <div className="mt-5 text-sm text-gray-600">
              <div className="font-medium text-gray-900">Address</div>
              <div className="mt-1">{location?.address}</div>
              <div className="mt-2">
                <a
                  href={mapsUrl}
                  target="_blank"
                  rel="noreferrer"
                  className="inline-flex items-center gap-2 text-blue-600 hover:underline"
                >
                  Open in Google Maps <ExternalLink size={16} />
                </a>
              </div>
            </div>
          </div>

          <div className="text-sm text-gray-600">
            <Link to="/locations" className="text-blue-600 hover:underline">See more locations</Link>
          </div>
        </div>

        <div className="lg:col-span-1">
          <div className="bg-white rounded-lg shadow-md overflow-hidden">
            <div className="p-4 border-b">
              <div className="text-sm font-medium text-gray-900">Map</div>
              <div className="text-sm text-gray-600 mt-1">Tap to open Google Maps.</div>
            </div>
            <div className="h-[260px]">
              {location && (
                <MapContainer
                  center={[location.latitude, location.longitude]}
                  zoom={15}
                  style={{ height: '100%', width: '100%' }}
                  scrollWheelZoom
                >
                  <TileLayer
                    attribution='&copy; OpenStreetMap contributors'
                    url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
                  />
                  <Marker
                    position={[location.latitude, location.longitude]}
                    icon={MarkerIcon}
                    eventHandlers={{
                      click: () => {
                        window.open(mapsUrl, '_blank', 'noreferrer')
                      },
                    }}
                  />
                </MapContainer>
              )}
            </div>
            <div className="p-4">
              <a
                href={mapsUrl}
                target="_blank"
                rel="noreferrer"
                className="w-full inline-flex items-center justify-center gap-2 border border-gray-300 rounded-md py-2 hover:bg-gray-50"
              >
                Share storage location
              </a>
            </div>
          </div>
        </div>
      </div>

      {location && (
        <BookingDialog
          open={bookingOpen}
          onClose={() => setBookingOpen(false)}
          location={location}
          initialStart={initialStart}
          initialBags={initialBags}
        />
      )}
    </div>
  )
}
