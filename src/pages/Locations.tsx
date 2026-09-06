import { useEffect, useMemo, useRef, useState } from 'react'
import { useNavigate, useSearchParams } from 'react-router-dom'
import { MapPin, Navigation, Lock, Clock, Search, Calendar } from 'lucide-react'
import { MapContainer, Marker, Popup, TileLayer, useMap, ZoomControl } from 'react-leaflet'
import { DayPicker } from 'react-day-picker'
import L from 'leaflet'
import type { Location, LocationWithDistance } from '../types'
import { storageService } from '../services/storage'
import { geocodeService, type PlaceSuggestion } from '../services/geocode'
import BookingDialog from '../components/BookingDialog'

type AnyLocation = Location | LocationWithDistance

function isNearby(loc: AnyLocation): loc is LocationWithDistance {
  return typeof (loc as LocationWithDistance).distance === 'number'
}

// Fix Leaflet default marker icons under bundlers
const SelectedIcon = L.icon({
  iconUrl: 'https://unpkg.com/leaflet@1.9.4/dist/images/marker-icon.png',
  iconRetinaUrl: 'https://unpkg.com/leaflet@1.9.4/dist/images/marker-icon-2x.png',
  shadowUrl: 'https://unpkg.com/leaflet@1.9.4/dist/images/marker-shadow.png',
  iconSize: [30, 49],
  iconAnchor: [15, 49],
})

// Keeps the map in sync with the current results. MapContainer's `center` prop
// is only applied on mount, so we imperatively recenter/fit when data changes.
function RecenterMap({ points, center }: { points: Array<[number, number]>; center: [number, number] }) {
  const map = useMap()
  useEffect(() => {
    if (points.length === 1) {
      map.setView(points[0], 14)
    } else if (points.length > 1) {
      map.fitBounds(points, { padding: [40, 40], maxZoom: 15 })
    } else {
      map.setView(center, 12)
    }
  }, [points, center, map])
  return null
}

// Bounce-style price pin marker.
function priceIcon(label: string) {
  return L.divIcon({
    className: 'te-price-pin',
    html:
      `<div style="background:#2563eb;color:#fff;font-weight:700;font-size:12px;line-height:1;` +
      `padding:6px 10px;border-radius:9999px;box-shadow:0 2px 6px rgba(0,0,0,.35);` +
      `border:2px solid #fff;white-space:nowrap;">${label}</div>`,
    iconSize: [54, 26],
    iconAnchor: [27, 26],
  })
}

function dayLabel(d: Date) {
  const today = new Date()
  today.setHours(0, 0, 0, 0)
  const day = new Date(d)
  day.setHours(0, 0, 0, 0)
  const diff = (day.getTime() - today.getTime()) / 86400000
  if (diff === 0) return 'Today'
  if (diff === 1) return 'Tomorrow'
  return day.toLocaleDateString(undefined, { weekday: 'short', month: 'short', day: 'numeric' })
}

export default function Locations() {
  const [searchParams] = useSearchParams()
  const navigate = useNavigate()
  const [query, setQuery] = useState('')
  const [radiusKm, setRadiusKm] = useState(25)
  const [locations, setLocations] = useState<AnyLocation[]>([])
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const [bookingOpen, setBookingOpen] = useState(false)
  const [bookingLocation, setBookingLocation] = useState<Location | null>(null)

  const [searchCenter, setSearchCenter] = useState<{ lat: number; lon: number } | null>(null)

  const [suggestions, setSuggestions] = useState<PlaceSuggestion[]>([])
  const [showSuggestions, setShowSuggestions] = useState(false)
  const [selectedSuggestion, setSelectedSuggestion] = useState<PlaceSuggestion | null>(null)
  const debounceRef = useRef<number | null>(null)

  const today = useMemo(() => {
    const d = new Date()
    d.setHours(0, 0, 0, 0)
    return d
  }, [])
  const [dropOffDay, setDropOffDay] = useState<Date>(() => new Date())
  const [dropTime, setDropTime] = useState<string>(() => {
    const d = new Date()
    d.setSeconds(0, 0)
    return `${String(d.getHours()).padStart(2, '0')}:${String(d.getMinutes()).padStart(2, '0')}`
  })
  const [dateOpen, setDateOpen] = useState(false)
  const datePopoverRef = useRef<HTMLDivElement | null>(null)

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

    const startStr = searchParams.get('start')
    if (startStr) {
      const sd = new Date(startStr)
      if (!Number.isNaN(sd.getTime())) {
        setDropOffDay(sd)
        setDropTime(`${String(sd.getHours()).padStart(2, '0')}:${String(sd.getMinutes()).padStart(2, '0')}`)
      }
    }

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

  // Close the date popover on outside click.
  useEffect(() => {
    if (!dateOpen) return
    const onDown = (e: MouseEvent) => {
      const t = e.target as Node | null
      if (t && datePopoverRef.current && !datePopoverRef.current.contains(t)) {
        setDateOpen(false)
      }
    }
    document.addEventListener('mousedown', onDown)
    return () => document.removeEventListener('mousedown', onDown)
  }, [dateOpen])

  // Debounced place autosuggest (mirrors the Home page search).
  useEffect(() => {
    const text = query.trim()
    if (
      selectedSuggestion &&
      text !== selectedSuggestion.title &&
      text !== selectedSuggestion.displayName
    ) {
      setSelectedSuggestion(null)
    }

    if (!showSuggestions) return
    if (text.length < 2) {
      setSuggestions([])
      return
    }

    if (debounceRef.current) window.clearTimeout(debounceRef.current)
    debounceRef.current = window.setTimeout(() => {
      void (async () => {
        try {
          setSuggestions(await geocodeService.suggest(text, 6))
        } catch {
          setSuggestions([])
        }
      })()
    }, 250)

    return () => {
      if (debounceRef.current) window.clearTimeout(debounceRef.current)
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [query, showSuggestions])

  // Search: geocode the query (or use a picked suggestion) and search nearby.
  const runSearch = async () => {
    const text = query.trim()
    if (!text) return
    setShowSuggestions(false)
    try {
      const geo = selectedSuggestion
        ? { lat: selectedSuggestion.lat, lon: selectedSuggestion.lon }
        : (await geocodeService.suggest(text, 1))[0]
      if (geo && Number.isFinite(geo.lat) && Number.isFinite(geo.lon)) {
        await loadNearbyAt(geo.lat, geo.lon)
        return
      }
      await loadByCity(text)
    } catch {
      await loadByCity(text)
    }
  }

  const openDetails = (locationId: string) => {
    const qs = searchParams.toString()
    navigate(`/locations/${locationId}${qs ? `?${qs}` : ''}`)
  }

  const initialStart = useMemo(() => {
    const [hh, mm] = dropTime.split(':').map(Number)
    const d = new Date(dropOffDay)
    d.setHours(Number.isFinite(hh) ? hh : 0, Number.isFinite(mm) ? mm : 0, 0, 0)
    return d
  }, [dropOffDay, dropTime])

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

  const mapPoints = useMemo<Array<[number, number]>>(() => {
    const pts = mapMarkers.map((m) => [m.lat, m.lon] as [number, number])
    if (searchCenter) pts.push([searchCenter.lat, searchCenter.lon])
    return pts
  }, [mapMarkers, searchCenter])

  return (
    <div className="absolute inset-0 bg-gray-100">
      {/* Full-screen map */}
      <div className="absolute inset-0">
        <MapContainer
          center={mapCenter}
          zoom={13}
          zoomControl={false}
          style={{ height: '100%', width: '100%' }}
          scrollWheelZoom
        >
          <TileLayer
            attribution='&copy; OpenStreetMap contributors'
            url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
          />
          <ZoomControl position="topright" />
          <RecenterMap points={mapPoints} center={mapCenter} />

          {searchCenter && (
            <Marker position={[searchCenter.lat, searchCenter.lon]} icon={SelectedIcon}>
              <Popup>
                <div className="text-sm font-medium">Search area</div>
              </Popup>
            </Marker>
          )}

          {mapMarkers.map((m) => (
            <Marker
              key={m.id}
              position={[m.lat, m.lon]}
              icon={priceIcon(`₹${m.dailyRate}`)}
              eventHandlers={{ click: () => openDetails(m.id) }}
            >
              <Popup>
                <div className="text-sm">
                  <div className="font-medium">{m.name}</div>
                  <div className="text-gray-600">{m.city}, {m.country}</div>
                  <div className="mt-1">₹{m.hourlyRate}/hr • ₹{m.dailyRate}/day</div>
                  {m.distance !== null && <div className="mt-1 text-gray-600">{m.distance} km away</div>}
                  {m.available !== null && <div className="text-gray-600">{m.available} available</div>}
                  <button
                    type="button"
                    onClick={() => openDetails(m.id)}
                    className="mt-2 rounded bg-blue-600 px-3 py-1.5 text-xs font-medium text-white hover:bg-blue-700"
                  >
                    View details
                  </button>
                </div>
              </Popup>
            </Marker>
          ))}
        </MapContainer>
      </div>

      {/* Floating search + results panel */}
      <div className="absolute left-4 top-4 z-[1000] flex max-h-[calc(100%-2rem)] w-[calc(100%-2rem)] max-w-sm flex-col overflow-hidden rounded-2xl bg-white shadow-2xl ring-1 ring-black/5">
        <div className="overflow-y-auto">
          <div className="space-y-6 p-5">
            {/* Where */}
            <div className="relative">
              <div className="text-lg font-bold text-gray-900">Where?</div>
              <div className="relative mt-2">
                <Search className="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400" size={18} />
                <input
                  value={query}
                  onChange={(e) => {
                    setQuery(e.target.value)
                    setShowSuggestions(true)
                  }}
                  onFocus={() => setShowSuggestions(true)}
                  onBlur={() => window.setTimeout(() => setShowSuggestions(false), 150)}
                  onKeyDown={(e) => {
                    if (e.key === 'Enter') {
                      e.preventDefault()
                      void runSearch()
                    }
                  }}
                  placeholder="Search a city, area or place"
                  className="w-full rounded-xl border border-gray-300 bg-gray-50 py-3 pl-10 pr-4 focus:border-transparent focus:bg-white focus:ring-2 focus:ring-blue-500"
                />
              </div>

              {showSuggestions && suggestions.length > 0 && (
                <div className="absolute z-[1000] mt-2 w-full overflow-hidden rounded-xl border border-gray-200 bg-white shadow-lg">
                  {suggestions.map((s) => (
                    <button
                      type="button"
                      key={`${s.lat},${s.lon},${s.title}`}
                      onMouseDown={(e) => e.preventDefault()}
                      onClick={() => {
                        setSelectedSuggestion(s)
                        setQuery(s.displayName || s.title)
                        setSuggestions([])
                        setShowSuggestions(false)
                        void loadNearbyAt(s.lat, s.lon)
                      }}
                      className="flex w-full items-start gap-3 px-4 py-3 text-left hover:bg-gray-50"
                    >
                      <MapPin size={16} className="mt-0.5 shrink-0 text-gray-400" />
                      <span className="min-w-0">
                        <span className="block truncate text-sm font-medium text-gray-900">{s.title}</span>
                        {s.subtitle && <span className="block truncate text-xs text-gray-500">{s.subtitle}</span>}
                      </span>
                    </button>
                  ))}
                </div>
              )}
            </div>

            {/* When */}
            <div>
              <div className="text-lg font-bold text-gray-900">When?</div>
              <div className="mt-2 grid grid-cols-2 gap-2">
                <div className="relative" ref={datePopoverRef}>
                  <button
                    type="button"
                    onClick={() => setDateOpen((v) => !v)}
                    className="flex w-full items-center gap-2 rounded-xl border border-gray-300 bg-gray-50 px-3 py-3 text-left text-sm hover:bg-white"
                  >
                    <Calendar size={16} className="text-gray-400" />
                    {dayLabel(dropOffDay)}
                  </button>
                  {dateOpen && (
                    <div className="absolute z-[1000] mt-2 rounded-xl border border-gray-200 bg-white p-3 shadow-lg">
                      <DayPicker
                        mode="single"
                        selected={dropOffDay}
                        onSelect={(d?: Date) => {
                          if (d) {
                            setDropOffDay(d)
                            setDateOpen(false)
                          }
                        }}
                        disabled={{ before: today }}
                        showOutsideDays
                        classNames={{
                          months: 'flex flex-col',
                          month: 'space-y-3',
                          caption: 'flex items-center justify-between',
                          caption_label: 'text-sm font-medium text-gray-900',
                          nav: 'flex items-center gap-2',
                          nav_button: 'h-8 w-8 rounded-md border border-gray-200 hover:bg-gray-50',
                          table: 'w-full border-collapse',
                          head_row: 'flex',
                          head_cell: 'w-9 text-[11px] font-medium text-gray-500',
                          row: 'flex w-full mt-1',
                          cell: 'w-9 h-9 text-sm',
                          day: 'w-9 h-9 rounded-md hover:bg-gray-100',
                          day_selected: 'bg-blue-600 text-white hover:bg-blue-600',
                          day_today: 'border border-blue-600',
                          day_outside: 'text-gray-400',
                          day_disabled: 'text-gray-300 line-through',
                        }}
                      />
                    </div>
                  )}
                </div>
                <input
                  type="time"
                  value={dropTime}
                  onChange={(e) => setDropTime(e.target.value)}
                  className="w-full rounded-xl border border-gray-300 bg-gray-50 px-3 py-3 text-sm focus:border-transparent focus:bg-white focus:ring-2 focus:ring-blue-500"
                />
              </div>
            </div>

            {/* Action */}
            <button
              onClick={() => void runSearch()}
              disabled={loading || !query.trim()}
              className="w-full rounded-xl bg-blue-600 py-3 font-medium text-white hover:bg-blue-700 disabled:opacity-50"
            >
              {loading ? 'Searching…' : 'Find storage'}
            </button>

            {error && <div className="text-sm text-red-600">{error}</div>}

            {/* Results */}
            <div>
              <div className="text-sm text-gray-500">
                {loading ? 'Searching…' : `${locations.length} storage ${locations.length === 1 ? 'spot' : 'spots'}`}
              </div>

              {locations.length === 0 && !loading && (
                <div className="mt-3 rounded-xl border border-dashed border-gray-200 p-6 text-sm text-gray-600">
                  Search a city or use “Near me” to see storage spots.
                </div>
              )}

              <div className="mt-3 space-y-3">
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
                    className="cursor-pointer rounded-xl border border-gray-200 p-4 transition-shadow hover:shadow-md"
                  >
                    <div className="flex items-start justify-between gap-3">
                      <div className="min-w-0">
                        <h3 className="truncate text-base font-semibold text-gray-900">{loc.name}</h3>
                        <p className="mt-0.5 truncate text-sm text-gray-500">{loc.city}, {loc.country}</p>
                        <div className="mt-2 flex items-center gap-2 text-sm text-gray-600">
                          <Lock className="h-4 w-4" />
                          ₹{loc.hourlyRate}/hr • ₹{loc.dailyRate}/day
                        </div>
                        <div className="mt-1 flex items-center gap-2 text-xs text-gray-500">
                          <Clock className="h-4 w-4" />
                          {loc.isOpen24Hours ? 'Open 24/7' : `${loc.openingTime || '—'} - ${loc.closingTime || '—'}`}
                        </div>
                      </div>
                      <div className="shrink-0 text-right">
                        {isNearby(loc) && <div className="text-xs text-gray-500">{loc.distance} km</div>}
                        {isNearby(loc) && <div className="text-xs font-medium text-gray-900">{loc.availableLockers} free</div>}
                        <button
                          type="button"
                          onClick={(e) => {
                            e.stopPropagation()
                            openBooking(loc)
                          }}
                          className="mt-2 rounded-lg bg-blue-600 px-3 py-2 text-xs font-medium text-white hover:bg-blue-700"
                        >
                          Book
                        </button>
                      </div>
                    </div>
                  </div>
                ))}
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* Floating "Near me" */}
      <button
        onClick={useMyLocation}
        disabled={loading}
        title="Use my location"
        className="absolute bottom-6 right-6 z-[1000] inline-flex items-center gap-2 rounded-full bg-white px-4 py-2.5 text-sm font-medium text-gray-800 shadow-lg ring-1 ring-black/5 hover:bg-gray-50 disabled:opacity-50"
      >
        <Navigation size={16} className="text-blue-600" />
        Near me
      </button>

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
