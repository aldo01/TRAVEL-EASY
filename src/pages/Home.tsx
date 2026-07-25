import { useEffect, useMemo, useRef, useState } from 'react'
import { Link, useNavigate, useSearchParams } from 'react-router-dom'
import { Navigation, CheckCircle, X } from 'lucide-react'
import { DayPicker } from 'react-day-picker'
import { geocodeService, type PlaceSuggestion } from '../services/geocode'

function toRFC3339Local(date: Date) {
  const pad = (n: number) => String(n).padStart(2, '0')
  const yyyy = date.getFullYear()
  const mm = pad(date.getMonth() + 1)
  const dd = pad(date.getDate())
  const hh = pad(date.getHours())
  const min = pad(date.getMinutes())
  const tz = -date.getTimezoneOffset()
  const sign = tz >= 0 ? '+' : '-'
  const abs = Math.abs(tz)
  const tzh = pad(Math.floor(abs / 60))
  const tzm = pad(abs % 60)
  return `${yyyy}-${mm}-${dd}T${hh}:${min}:00${sign}${tzh}:${tzm}`
}

export default function Home() {
  const navigate = useNavigate()
  const [searchParams, setSearchParams] = useSearchParams()

  const bookingConfirmed = searchParams.get('bookingConfirmed') === '1'
  const bookingNumber = searchParams.get('bookingNumber')
  const bookingLocationName = searchParams.get('locationName')
  const bookingTotal = searchParams.get('total')
  const bookingBags = searchParams.get('bags')
  const bookingDropOff = searchParams.get('dropOff')
  const bookingPickUp = searchParams.get('pickUp')

  const dismissConfirmation = () => {
    setSearchParams({}, { replace: true })
  }

  const [city, setCity] = useState('')
  const [radiusKm, setRadiusKm] = useState(10)

  const [suggestions, setSuggestions] = useState<PlaceSuggestion[]>([])
  const [showSuggestions, setShowSuggestions] = useState(false)
  const [selectedSuggestion, setSelectedSuggestion] = useState<PlaceSuggestion | null>(null)
  const [suggestionError, setSuggestionError] = useState<string | null>(null)
  const debounceRef = useRef<number | null>(null)

  const today = useMemo(() => {
    const d = new Date()
    d.setHours(0, 0, 0, 0)
    return d
  }, [])

  const [startDay, setStartDay] = useState<Date>(() => new Date())
  const [dateOpen, setDateOpen] = useState(false)
  const datePopoverRef = useRef<HTMLDivElement | null>(null)
  const [startTime, setStartTime] = useState<string>(() => {
    const d = new Date()
    d.setSeconds(0, 0)
    const hh = String(d.getHours()).padStart(2, '0')
    const mm = String(d.getMinutes()).padStart(2, '0')
    return `${hh}:${mm}`
  })

  const buildStartParam = () => {
    const [hh, mm] = startTime.split(':').map((n) => Number(n))
    const startDateTime = new Date(startDay)
    startDateTime.setHours(Number.isFinite(hh) ? hh : 0, Number.isFinite(mm) ? mm : 0, 0, 0)
    return toRFC3339Local(startDateTime)
  }

  const [searchError, setSearchError] = useState<string | null>(null)
  const [searching, setSearching] = useState(false)

  useEffect(() => {
    if (!dateOpen) return

    const onDown = (e: MouseEvent) => {
      const target = e.target as Node | null
      if (!target) return
      if (datePopoverRef.current && !datePopoverRef.current.contains(target)) {
        setDateOpen(false)
      }
    }

    document.addEventListener('mousedown', onDown)
    return () => document.removeEventListener('mousedown', onDown)
  }, [dateOpen])

  useEffect(() => {
    const text = city.trim()
    setSuggestionError(null)

    // If user edits the field after selecting, unlock selection.
    if (selectedSuggestion && text !== selectedSuggestion.title && text !== selectedSuggestion.displayName) {
      setSelectedSuggestion(null)
    }

    if (!showSuggestions) return
    if (text.length < 2) {
      setSuggestions([])
      return
    }

    if (debounceRef.current) {
      window.clearTimeout(debounceRef.current)
    }

    debounceRef.current = window.setTimeout(() => {
      void (async () => {
        try {
          const items = await geocodeService.suggest(text, 6)
          setSuggestions(items)
        } catch (e) {
          setSuggestions([])
          setSuggestionError(e instanceof Error ? e.message : 'Failed to load suggestions')
        }
      })()
    }, 250)

    return () => {
      if (debounceRef.current) window.clearTimeout(debounceRef.current)
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [city, showSuggestions])

  const searchCity = async () => {
    const trimmed = city.trim()
    if (!trimmed) return
    setSearching(true)
    setSearchError(null)
    try {
      const start = buildStartParam()
      const geo = selectedSuggestion
        ? { lat: selectedSuggestion.lat, lon: selectedSuggestion.lon, displayName: selectedSuggestion.displayName }
        : (await geocodeService.suggest(trimmed, 1))[0]

      if (geo && Number.isFinite(geo.lat) && Number.isFinite(geo.lon)) {
        navigate(
          `/locations?nearby=1&lat=${encodeURIComponent(String(geo.lat))}` +
            `&lon=${encodeURIComponent(String(geo.lon))}` +
            `&radiusKm=${encodeURIComponent(String(radiusKm))}` +
            `&q=${encodeURIComponent(geo.displayName ?? trimmed)}` +
            `&start=${encodeURIComponent(start)}`
        )
        return
      }

      // Fallback: treat it as a city filter
      navigate(`/locations?city=${encodeURIComponent(trimmed)}&start=${encodeURIComponent(start)}`)
    } catch (e) {
      setSearchError(e instanceof Error ? e.message : 'Failed to find that place')
    } finally {
      setSearching(false)
    }
  }

  const useMyLocation = () => {
    const start = buildStartParam()
    navigate(`/locations?nearby=1&radiusKm=${encodeURIComponent(String(radiusKm))}&start=${encodeURIComponent(start)}`)
  }

  return (
    <div className="space-y-8">
      <div className="bg-white rounded-lg shadow-md p-8">
        <h1 className="text-3xl font-bold text-gray-900">Luggage storage, nearby</h1>
        <p className="text-gray-600 mt-1">
          Find secure lockers in hotels, cafes, shops and kiosks — then book in minutes.
        </p>

        <div className="mt-6 grid grid-cols-1 md:grid-cols-6 gap-4">
          <div className="md:col-span-2 relative">
            <label className="block text-sm font-medium text-gray-700 mb-2">
              City / area
            </label>
            <input
              value={city}
              onChange={(e) => {
                setCity(e.target.value)
                setShowSuggestions(true)
              }}
              onFocus={() => setShowSuggestions(true)}
              onBlur={() => {
                // Allow click selection before hiding.
                window.setTimeout(() => setShowSuggestions(false), 150)
              }}
              onKeyDown={(e) => {
                if (e.key === 'Enter') {
                  e.preventDefault()
                  void searchCity()
                  setShowSuggestions(false)
                }
              }}
              placeholder="e.g. London Luton Airport (LTN)"
              className="w-full px-4 py-3 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            />

            {showSuggestions && (suggestions.length > 0 || suggestionError) && (
              <div className="absolute z-50 mt-2 w-full bg-white border border-gray-200 rounded-md shadow-lg overflow-hidden">
                {suggestionError && (
                  <div className="px-4 py-3 text-sm text-red-600">{suggestionError}</div>
                )}

                {suggestions.map((s) => (
                  <button
                    type="button"
                    key={`${s.lat},${s.lon},${s.title}`}
                    onPointerDown={(e) => {
                      // Prevent input blur before we can handle selection.
                      e.preventDefault()
                    }}
                    onMouseDown={(e) => {
                      // Prevent input blur before we can handle selection.
                      e.preventDefault()
                    }}
                    onClick={() => {
                      setSelectedSuggestion(s)
                      setCity(s.displayName || s.title)
                      setSuggestions([])
                      setShowSuggestions(false)
                    }}
                    className="w-full text-left px-4 py-3 hover:bg-gray-50"
                  >
                    <div className="text-sm font-medium text-gray-900">{s.title}</div>
                    {s.subtitle && <div className="text-xs text-gray-600 mt-0.5">{s.subtitle}</div>}
                  </button>
                ))}
              </div>
            )}
          </div>

          <div className="md:col-span-2">
            <label className="block text-sm font-medium text-gray-700 mb-2">Date</label>
            <div className="relative" ref={datePopoverRef}>
              <button
                type="button"
                onClick={() => setDateOpen((v) => !v)}
                className="w-full px-4 py-3 border border-gray-300 rounded-md text-left hover:bg-gray-50 focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              >
                {startDay.toLocaleDateString(undefined, { weekday: 'short', day: '2-digit', month: 'short', year: 'numeric' })}
              </button>

              {dateOpen && (
                <div className="absolute z-50 mt-2 bg-white rounded-md border border-gray-200 shadow-lg p-3">
                  <DayPicker
                    mode="single"
                    selected={startDay}
                    onSelect={(d: Date | undefined) => {
                      if (d) {
                        setStartDay(d)
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

            <div className="mt-3">
              <label className="block text-sm font-medium text-gray-700 mb-2">Time</label>
              <input
                type="time"
                value={startTime}
                onChange={(e) => setStartTime(e.target.value)}
                className="w-full px-4 py-3 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              />
            </div>
          </div>

          <div className="md:col-span-1">
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

          <div className="md:col-span-1 flex items-end gap-2">
            <button
              onClick={() => {
                void searchCity()
              }}
              disabled={!city.trim() || searching}
              className="flex-1 bg-blue-600 text-white py-3 px-4 rounded-md hover:bg-blue-700 disabled:opacity-50"
            >
              Find storage
            </button>
            <button
              onClick={useMyLocation}
              disabled={searching}
              className="px-4 py-3 border border-gray-300 rounded-md hover:bg-gray-50 disabled:opacity-50"
              title="Use my location"
            >
              <Navigation size={18} />
            </button>
          </div>
        </div>

        {searchError && <div className="mt-3 text-sm text-red-600">{searchError}</div>}

        <div className="mt-6 flex items-center gap-4 text-sm">
          <button
            onClick={() => navigate('/locations')}
            className="text-blue-600 hover:text-blue-700 font-medium"
          >
            Browse all locations →
          </button>
          <Link to="/host" className="text-gray-600 hover:text-gray-800">
            List your space
          </Link>
        </div>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        <div className="bg-white rounded-lg shadow-md p-6">
          <div className="text-sm font-medium text-gray-900">Secure storage</div>
          <div className="text-sm text-gray-600 mt-1">Verified partner locations with staffed handover options.</div>
        </div>
        <div className="bg-white rounded-lg shadow-md p-6">
          <div className="text-sm font-medium text-gray-900">Fast booking</div>
          <div className="text-sm text-gray-600 mt-1">Book a locker instantly and get your booking code.</div>
        </div>
        <div className="bg-white rounded-lg shadow-md p-6">
          <div className="text-sm font-medium text-gray-900">Flexible time</div>
          <div className="text-sm text-gray-600 mt-1">Hourly and daily pricing depending on location.</div>
        </div>
      </div>

      <div className="bg-white rounded-lg shadow-md p-6 text-sm text-gray-600">
        Tip: Click the location arrow to use your current location and see nearby luggage drop points.
      </div>

      {bookingConfirmed && bookingNumber && (
        <div className="fixed inset-0 z-[3000]">
          <div className="absolute inset-0 bg-black/40" onClick={dismissConfirmation} />
          <div className="absolute inset-0 flex items-center justify-center p-4">
            <div className="relative z-[3001] w-full max-w-md bg-white rounded-lg shadow-lg overflow-hidden">
              <div className="flex items-center justify-between p-5 border-b">
                <div className="text-lg font-semibold text-gray-900">Booking Confirmed</div>
                <button onClick={dismissConfirmation} className="p-2 rounded-md hover:bg-gray-100">
                  <X size={18} />
                </button>
              </div>
              <div className="p-6 text-center">
                <div className="mx-auto flex h-16 w-16 items-center justify-center rounded-full bg-green-100 mb-4">
                  <CheckCircle size={36} className="text-green-600" />
                </div>
                <div className="text-xl font-bold text-gray-900 mb-1">Your booking is confirmed!</div>
                <div className="text-sm text-gray-600 mb-4">
                  A confirmation email has been sent with your booking details.
                </div>
                <div className="bg-gray-50 rounded-lg p-4 text-left space-y-2 text-sm">
                  <div className="flex justify-between">
                    <span className="text-gray-600">Booking ID</span>
                    <span className="font-semibold text-gray-900">{bookingNumber}</span>
                  </div>
                  {bookingLocationName && (
                    <div className="flex justify-between">
                      <span className="text-gray-600">Location</span>
                      <span className="font-medium text-gray-900">{bookingLocationName}</span>
                    </div>
                  )}
                  {bookingDropOff && (
                    <div className="flex justify-between">
                      <span className="text-gray-600">Drop-off</span>
                      <span className="text-gray-900">{new Date(bookingDropOff).toLocaleString(undefined, { month: 'short', day: 'numeric', hour: '2-digit', minute: '2-digit' })}</span>
                    </div>
                  )}
                  {bookingPickUp && (
                    <div className="flex justify-between">
                      <span className="text-gray-600">Pick-up</span>
                      <span className="text-gray-900">{new Date(bookingPickUp).toLocaleString(undefined, { month: 'short', day: 'numeric', hour: '2-digit', minute: '2-digit' })}</span>
                    </div>
                  )}
                  {bookingBags && (
                    <div className="flex justify-between">
                      <span className="text-gray-600">Bags</span>
                      <span className="text-gray-900">{bookingBags}</span>
                    </div>
                  )}
                  {bookingTotal && (
                    <div className="flex justify-between border-t pt-2 mt-2">
                      <span className="font-medium text-gray-900">Total</span>
                      <span className="font-bold text-gray-900">₹{bookingTotal}</span>
                    </div>
                  )}
                </div>
              </div>
              <div className="p-5 border-t space-y-2">
                <button
                  onClick={() => { dismissConfirmation(); navigate('/bookings'); }}
                  className="w-full bg-blue-600 text-white py-3 rounded-md hover:bg-blue-700 font-medium"
                >
                  View my bookings
                </button>
                <button
                  onClick={dismissConfirmation}
                  className="w-full border border-gray-300 text-gray-900 py-3 rounded-md hover:bg-gray-50"
                >
                  Continue browsing
                </button>
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}
