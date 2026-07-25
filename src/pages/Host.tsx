import { useMemo, useState, type FormEvent } from 'react'
import { storageService } from '../services/storage'
import type { LocationType } from '../types'

const LOCATION_TYPES: LocationType[] = ['HOTEL', 'CAFE', 'STORE', 'KIOSK', 'GYM', 'CLUB', 'SCHOOL', 'OTHER']

export default function Host() {
  const [name, setName] = useState('')
  const [type, setType] = useState<LocationType>('HOTEL')
  const [address, setAddress] = useState('')
  const [city, setCity] = useState('')
  const [country, setCountry] = useState('')
  const [latitude, setLatitude] = useState<string>('')
  const [longitude, setLongitude] = useState<string>('')
  const [openingTime, setOpeningTime] = useState('')
  const [closingTime, setClosingTime] = useState('')
  const [isOpen24Hours, setIsOpen24Hours] = useState(false)
  const [hourlyRate, setHourlyRate] = useState<number>(3)
  const [dailyRate, setDailyRate] = useState<number>(12)
  const [amenities, setAmenities] = useState<string>('CCTV, Staffed, Insurance')

  const [smallCount, setSmallCount] = useState<number>(5)
  const [mediumCount, setMediumCount] = useState<number>(5)
  const [largeCount, setLargeCount] = useState<number>(2)
  const [xlargeCount, setXlargeCount] = useState<number>(0)

  const [submitting, setSubmitting] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [success, setSuccess] = useState<string | null>(null)

  const amenitiesArray = useMemo(() => {
    return amenities
      .split(',')
      .map((x) => x.trim())
      .filter(Boolean)
  }, [amenities])

  const canSubmit = useMemo(() => {
    if (!name.trim() || !address.trim() || !city.trim() || !country.trim()) return false
    if (!latitude.trim() || !longitude.trim()) return false
    const lat = Number(latitude)
    const lon = Number(longitude)
    if (!Number.isFinite(lat) || !Number.isFinite(lon)) return false
    return true
  }, [name, address, city, country, latitude, longitude])

  const onUseMyLocation = async () => {
    setError(null)
    setSuccess(null)
    try {
      const pos = await new Promise<GeolocationPosition>((resolve, reject) => {
        if (!navigator.geolocation) reject(new Error('Geolocation not supported'))
        navigator.geolocation.getCurrentPosition(resolve, reject, { enableHighAccuracy: true, timeout: 10000 })
      })
      setLatitude(String(pos.coords.latitude))
      setLongitude(String(pos.coords.longitude))
    } catch (e) {
      setError(e instanceof Error ? e.message : 'Failed to get your location')
    }
  }

  const onSubmit = async (e: FormEvent) => {
    e.preventDefault()
    if (!canSubmit) return

    setSubmitting(true)
    setError(null)
    setSuccess(null)

    try {
      const created = await storageService.createLocation({
        name: name.trim(),
        type,
        address: address.trim(),
        city: city.trim(),
        country: country.trim(),
        latitude: Number(latitude),
        longitude: Number(longitude),
        isOpen24Hours,
        openingTime: isOpen24Hours ? undefined : openingTime || undefined,
        closingTime: isOpen24Hours ? undefined : closingTime || undefined,
        hourlyRate,
        dailyRate,
        amenities: amenitiesArray,
        lockers: {
          small: smallCount,
          medium: mediumCount,
          large: largeCount,
          xlarge: xlargeCount,
        },
      })

      setSuccess(`Location created: ${created.location.name}`)
    } catch (e) {
      setError(e instanceof Error ? e.message : 'Failed to create location (are you logged in?)')
    } finally {
      setSubmitting(false)
    }
  }

  return (
    <div className="max-w-3xl space-y-6">
      <div>
        <h1 className="text-3xl font-bold text-gray-900">List your space</h1>
        <p className="text-gray-600 mt-1">Hotels, cafes, shops and kiosks can add luggage storage locations.</p>
      </div>

      <div className="bg-white rounded-lg shadow-md p-6">
        <div className="text-sm text-gray-600">
          This action requires a logged-in account (JWT). Use seeded test user: test@example.com / password123
        </div>
      </div>

      <form onSubmit={onSubmit} className="bg-white rounded-lg shadow-md p-6 space-y-5">
        {error && <div className="text-sm text-red-600">{error}</div>}
        {success && <div className="text-sm text-green-700">{success}</div>}

        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">Business name</label>
            <input
              value={name}
              onChange={(e) => setName(e.target.value)}
              className="w-full px-4 py-3 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              placeholder="e.g. Central Station Cafe"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">Type</label>
            <select
              value={type}
              onChange={(e) => setType(e.target.value as LocationType)}
              className="w-full px-4 py-3 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            >
              {LOCATION_TYPES.map((t) => (
                <option key={t} value={t}>
                  {t}
                </option>
              ))}
            </select>
          </div>
        </div>

        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">Address</label>
          <input
            value={address}
            onChange={(e) => setAddress(e.target.value)}
            className="w-full px-4 py-3 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            placeholder="Street and number"
          />
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">City</label>
            <input
              value={city}
              onChange={(e) => setCity(e.target.value)}
              className="w-full px-4 py-3 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              placeholder="e.g. Copenhagen"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">Country</label>
            <input
              value={country}
              onChange={(e) => setCountry(e.target.value)}
              className="w-full px-4 py-3 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              placeholder="e.g. Denmark"
            />
          </div>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <div className="flex items-center justify-between">
              <label className="block text-sm font-medium text-gray-700 mb-2">Latitude</label>
              <button
                type="button"
                onClick={onUseMyLocation}
                className="text-sm text-blue-600 hover:text-blue-700"
              >
                Use my location
              </button>
            </div>
            <input
              value={latitude}
              onChange={(e) => setLatitude(e.target.value)}
              className="w-full px-4 py-3 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              placeholder="e.g. 55.6761"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">Longitude</label>
            <input
              value={longitude}
              onChange={(e) => setLongitude(e.target.value)}
              className="w-full px-4 py-3 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              placeholder="e.g. 12.5683"
            />
          </div>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          <div className="md:col-span-1">
            <label className="block text-sm font-medium text-gray-700 mb-2">Open 24/7</label>
            <label className="inline-flex items-center gap-2 text-sm text-gray-700">
              <input
                type="checkbox"
                checked={isOpen24Hours}
                onChange={(e) => setIsOpen24Hours(e.target.checked)}
              />
              Yes
            </label>
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">Opening time</label>
            <input
              value={openingTime}
              onChange={(e) => setOpeningTime(e.target.value)}
              disabled={isOpen24Hours}
              placeholder="09:00"
              className="w-full px-4 py-3 border border-gray-300 rounded-md disabled:bg-gray-50 focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">Closing time</label>
            <input
              value={closingTime}
              onChange={(e) => setClosingTime(e.target.value)}
              disabled={isOpen24Hours}
              placeholder="18:00"
              className="w-full px-4 py-3 border border-gray-300 rounded-md disabled:bg-gray-50 focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            />
          </div>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">Hourly rate</label>
            <input
              type="number"
              min={0}
              step={0.5}
              value={hourlyRate}
              onChange={(e) => setHourlyRate(Number(e.target.value))}
              className="w-full px-4 py-3 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">Daily rate</label>
            <input
              type="number"
              min={0}
              step={1}
              value={dailyRate}
              onChange={(e) => setDailyRate(Number(e.target.value))}
              className="w-full px-4 py-3 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            />
          </div>
        </div>

        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">Amenities (comma separated)</label>
          <input
            value={amenities}
            onChange={(e) => setAmenities(e.target.value)}
            className="w-full px-4 py-3 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-transparent"
          />
        </div>

        <div>
          <h2 className="text-lg font-semibold text-gray-900">Locker inventory</h2>
          <p className="text-sm text-gray-600 mt-1">How many lockers of each size can you host?</p>

          <div className="grid grid-cols-2 md:grid-cols-4 gap-4 mt-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">Small</label>
              <input
                type="number"
                min={0}
                value={smallCount}
                onChange={(e) => setSmallCount(Number(e.target.value))}
                className="w-full px-4 py-3 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">Medium</label>
              <input
                type="number"
                min={0}
                value={mediumCount}
                onChange={(e) => setMediumCount(Number(e.target.value))}
                className="w-full px-4 py-3 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">Large</label>
              <input
                type="number"
                min={0}
                value={largeCount}
                onChange={(e) => setLargeCount(Number(e.target.value))}
                className="w-full px-4 py-3 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">XLarge</label>
              <input
                type="number"
                min={0}
                value={xlargeCount}
                onChange={(e) => setXlargeCount(Number(e.target.value))}
                className="w-full px-4 py-3 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              />
            </div>
          </div>
        </div>

        <button
          type="submit"
          disabled={submitting || !canSubmit}
          className="w-full bg-blue-600 text-white py-3 px-6 rounded-md hover:bg-blue-700 disabled:opacity-50"
        >
          {submitting ? 'Creating…' : 'Create location'}
        </button>
      </form>
    </div>
  )
}
