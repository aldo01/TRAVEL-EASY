import { useEffect, useMemo, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { X, Minus, Plus, CalendarDays } from 'lucide-react'
import type { Location, RateType } from '../types'
import { storageService } from '../services/storage'
import { authStore } from '../services/auth'

function pad2(n: number) {
  return String(n).padStart(2, '0')
}

function toRFC3339Local(date: Date) {
  const yyyy = date.getFullYear()
  const mm = pad2(date.getMonth() + 1)
  const dd = pad2(date.getDate())
  const hh = pad2(date.getHours())
  const min = pad2(date.getMinutes())
  const tz = -date.getTimezoneOffset()
  const sign = tz >= 0 ? '+' : '-'
  const abs = Math.abs(tz)
  const tzh = pad2(Math.floor(abs / 60))
  const tzm = pad2(abs % 60)
  return `${yyyy}-${mm}-${dd}T${hh}:${min}:00${sign}${tzh}:${tzm}`
}

function roundCurrency(value: number) {
  return Math.round((value + Number.EPSILON) * 100) / 100
}

export default function BookingDialog(props: {
  open: boolean
  onClose: () => void
  location: Location
  initialStart?: Date
  initialBags?: number
}) {
  const { open, onClose, location, initialStart, initialBags } = props
  const navigate = useNavigate()

  const [bags, setBags] = useState(2)
  const [rateType, setRateType] = useState<RateType>('DAILY')

  const [dropOff, setDropOff] = useState<Date>(() => initialStart ?? new Date())
  const [pickUp, setPickUp] = useState<Date>(() => {
    const base = initialStart ?? new Date()
    const d = new Date(base)
    d.setHours(d.getHours() + 24)
    return d
  })

  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [success, setSuccess] = useState<string | null>(null)
  const [loginPromptOpen, setLoginPromptOpen] = useState(false)

  useEffect(() => {
    if (!open) return
    const base = initialStart ?? new Date()
    base.setSeconds(0, 0)
    const d = new Date(base)
    const p = new Date(base)
    p.setHours(p.getHours() + 24)
    setDropOff(d)
    setPickUp(p)
    setBags(Number.isFinite(initialBags) && (initialBags ?? 0) > 0 ? Math.floor(initialBags!) : 2)
    setRateType('DAILY')
    setError(null)
    setSuccess(null)
    setLoginPromptOpen(false)
  }, [open, initialBags, initialStart])

  const durationHours = useMemo(() => {
    const ms = pickUp.getTime() - dropOff.getTime()
    const hours = ms / (1000 * 60 * 60)
    if (!Number.isFinite(hours) || hours <= 0) return 1
    return Math.max(1, Math.ceil(hours))
  }, [dropOff, pickUp])

  const storagePrice = useMemo(() => {
    const base = rateType === 'HOURLY'
      ? Math.ceil(durationHours) * location.hourlyRate
      : Math.ceil(durationHours / 24) * location.dailyRate

    return roundCurrency(base * Math.max(1, bags))
  }, [bags, durationHours, location.dailyRate, location.hourlyRate, rateType])

  const serviceFee: number = 0
  const insuranceFee: number = 0
  const total = roundCurrency(storagePrice + serviceFee + insuranceFee)

  const validationErrors = useMemo(() => {
    const errs: string[] = []
    const now = new Date()
    if (dropOff.getTime() <= now.getTime()) {
      errs.push('Drop-off date & time must be in the future.')
    }
    const diffMs = pickUp.getTime() - dropOff.getTime()
    const diffHours = diffMs / (1000 * 60 * 60)
    if (diffMs <= 0) {
      errs.push('Pick-up must be after drop-off.')
    } else if (diffHours < 3) {
      errs.push('Minimum booking duration is 3 hours.')
    }
    if (bags < 1) {
      errs.push('Please add at least 1 bag.')
    }
    return errs
  }, [dropOff, pickUp, bags])

  const canConfirm = open && !loading && validationErrors.length === 0

  const confirmBooking = async () => {
    setError(null)
    setSuccess(null)
    setLoginPromptOpen(false)
    setLoading(true)
    try {
      const token = authStore.getToken()
      if (!token) {
        setLoginPromptOpen(true)
        setError('Please log in to confirm booking.')
        return
      }

      const available = await storageService.getAvailableLockers(location.id)
      const lockerId = available[0]?.id
      if (!lockerId) {
        setError('No available lockers right now')
        return
      }

      const booking = await storageService.createBooking({
        locationId: location.id,
        lockerId,
        startTime: toRFC3339Local(dropOff),
        durationHours,
        rateType,
        bags,
      })

      onClose()
      navigate(`/?bookingConfirmed=1&bookingNumber=${encodeURIComponent(booking.bookingNumber)}&locationName=${encodeURIComponent(location.name)}&total=${encodeURIComponent(total.toFixed(2))}&bags=${bags}&dropOff=${encodeURIComponent(dropOff.toISOString())}&pickUp=${encodeURIComponent(pickUp.toISOString())}`)
    } catch (e) {
      const message = e instanceof Error ? e.message : 'Booking failed'
      if (typeof message === 'string' && message.includes('status: 401')) {
        setLoginPromptOpen(true)
        setError('Please log in to confirm booking.')
      } else {
        setError(message)
      }
    } finally {
      setLoading(false)
    }
  }

  if (!open) return null

  return (
    <div className="fixed inset-0 z-[2000]">
      <div className="absolute inset-0 bg-black/40" onClick={onClose} />
      <div className="absolute inset-0 flex items-center justify-center p-4">
        <div className="relative z-[2001] w-full max-w-3xl bg-white rounded-lg shadow-lg overflow-hidden">
          <div className="flex items-start justify-between p-5 border-b">
            <div>
              <div className="text-sm text-gray-600">Luggage storage</div>
              <div className="text-xl font-semibold text-gray-900">{location.name}</div>
              <div className="text-sm text-gray-600">{location.city}, {location.country}</div>
            </div>
            <button
              onClick={onClose}
              className="p-2 rounded-md hover:bg-gray-100"
              aria-label="Close"
            >
              <X size={18} />
            </button>
          </div>

          <div className="p-5">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div className="rounded-lg border border-gray-200 bg-white p-4 shadow-sm">
                <div className="flex items-center justify-between">
                  <div className="text-sm font-medium text-gray-900">Drop off</div>
                  <CalendarDays size={16} className="text-gray-400" />
                </div>
                <div className="mt-3 grid grid-cols-2 gap-3">
                  <input
                    type="date"
                    value={`${dropOff.getFullYear()}-${pad2(dropOff.getMonth() + 1)}-${pad2(dropOff.getDate())}`}
                    onChange={(e) => {
                      const [y, m, d] = e.target.value.split('-').map(Number)
                      if (!y || !m || !d) return
                      const next = new Date(dropOff)
                      next.setFullYear(y, m - 1, d)
                      setDropOff(next)
                    }}
                    className="w-full rounded-md border border-gray-200 bg-gray-50 px-3 py-2 text-sm"
                  />
                  <input
                    type="time"
                    value={`${pad2(dropOff.getHours())}:${pad2(dropOff.getMinutes())}`}
                    onChange={(e) => {
                      const [hh, mm] = e.target.value.split(':').map(Number)
                      const next = new Date(dropOff)
                      next.setHours(Number.isFinite(hh) ? hh : 0, Number.isFinite(mm) ? mm : 0, 0, 0)
                      setDropOff(next)
                    }}
                    className="w-full rounded-md border border-gray-200 bg-gray-50 px-3 py-2 text-sm"
                  />
                </div>
                <div className="mt-2 inline-flex items-center gap-2 text-xs text-gray-600">
                  <span className="h-2 w-2 rounded-full bg-green-600" />
                  Open {location.isOpen24Hours ? '24/7' : `${location.openingTime || '—'}-${location.closingTime || '—'}`}
                </div>
              </div>

              <div className="rounded-lg border border-gray-200 bg-white p-4 shadow-sm">
                <div className="flex items-center justify-between">
                  <div className="text-sm font-medium text-gray-900">Pick up</div>
                  <CalendarDays size={16} className="text-gray-400" />
                </div>
                <div className="mt-3 grid grid-cols-2 gap-3">
                  <input
                    type="date"
                    value={`${pickUp.getFullYear()}-${pad2(pickUp.getMonth() + 1)}-${pad2(pickUp.getDate())}`}
                    onChange={(e) => {
                      const [y, m, d] = e.target.value.split('-').map(Number)
                      if (!y || !m || !d) return
                      const next = new Date(pickUp)
                      next.setFullYear(y, m - 1, d)
                      setPickUp(next)
                    }}
                    className="w-full rounded-md border border-gray-200 bg-gray-50 px-3 py-2 text-sm"
                  />
                  <input
                    type="time"
                    value={`${pad2(pickUp.getHours())}:${pad2(pickUp.getMinutes())}`}
                    onChange={(e) => {
                      const [hh, mm] = e.target.value.split(':').map(Number)
                      const next = new Date(pickUp)
                      next.setHours(Number.isFinite(hh) ? hh : 0, Number.isFinite(mm) ? mm : 0, 0, 0)
                      setPickUp(next)
                    }}
                    className="w-full rounded-md border border-gray-200 bg-gray-50 px-3 py-2 text-sm"
                  />
                </div>
                <div className="mt-2 inline-flex items-center gap-2 text-xs text-gray-600">
                  <span className="h-2 w-2 rounded-full bg-green-600" />
                  Open {location.isOpen24Hours ? '24/7' : `${location.openingTime || '—'}-${location.closingTime || '—'}`}
                </div>
              </div>
            </div>

            <div className="mt-4 text-sm text-gray-600">
              Choose your dates. You can cancel your booking until midnight the day before arrival.
            </div>

            <div className="mt-6">
              <div className="flex items-center justify-between gap-3">
                <div>
                  <div className="text-sm font-medium text-gray-900">How many bags?</div>
                  <div className="text-xs text-gray-600 mt-1">Suitcases, backpacks, shopping bags — we’re not picky.</div>
                </div>
                <button
                  onClick={() => setBags((b) => Math.max(0, b - 1))}
                  className="h-9 w-9 rounded-full border border-gray-300 hover:bg-gray-50"
                  aria-label="Decrease bags"
                >
                  <Minus size={16} className="mx-auto" />
                </button>
                <div className="text-lg font-semibold text-gray-900 w-8 text-center">{bags}</div>
                <button
                  onClick={() => setBags((b) => Math.min(20, b + 1))}
                  className="h-9 w-9 rounded-full border border-gray-300 hover:bg-gray-50"
                  aria-label="Increase bags"
                >
                  <Plus size={16} className="mx-auto" />
                </button>
              </div>
            </div>

            <div className="mt-5 grid grid-cols-2 gap-3">
              <button
                onClick={() => setRateType('HOURLY')}
                className={`py-2 rounded-md border ${rateType === 'HOURLY' ? 'bg-blue-600 text-white border-blue-600' : 'bg-white text-gray-900 border-gray-300 hover:bg-gray-50'}`}
              >
                Hourly pricing
              </button>
              <button
                onClick={() => setRateType('DAILY')}
                className={`py-2 rounded-md border ${rateType === 'DAILY' ? 'bg-blue-600 text-white border-blue-600' : 'bg-white text-gray-900 border-gray-300 hover:bg-gray-50'}`}
              >
                Daily pricing
              </button>
            </div>

            <div className="mt-6 border-t pt-4">
              <div className="text-lg font-semibold text-gray-900">Summary</div>
              <div className="mt-3 space-y-2 text-sm">
                <div className="flex items-center justify-between">
                  <div className="text-gray-700">Storage ({bags} bag{bags === 1 ? '' : 's'})</div>
                  <div className="text-gray-900">₹{storagePrice.toFixed(2)}</div>
                </div>
                <div className="flex items-center justify-between">
                  <div className="text-gray-700">Service fee</div>
                  <div className="text-gray-900">{serviceFee === 0 ? 'Included' : `₹${serviceFee.toFixed(2)}`}</div>
                </div>
                <div className="flex items-center justify-between">
                  <div className="text-gray-700">Insurance</div>
                  <div className="text-gray-900">{insuranceFee === 0 ? 'Included' : `₹${insuranceFee.toFixed(2)}`}</div>
                </div>
                <div className="flex items-center justify-between border-t pt-2 mt-2">
                  <div className="font-medium text-gray-900">Total</div>
                  <div className="font-semibold text-gray-900">₹{total.toFixed(2)}</div>
                </div>
              </div>

              {validationErrors.length > 0 && (
                <div className="mt-3 rounded-md bg-red-50 border border-red-200 p-3">
                  <div className="text-sm font-medium text-red-800">Please fix the following:</div>
                  <ul className="mt-1 list-disc list-inside text-sm text-red-700">
                    {validationErrors.map((msg) => (
                      <li key={msg}>{msg}</li>
                    ))}
                  </ul>
                </div>
              )}
              {error && <div className="mt-3 text-sm text-red-600">{error}</div>}
              {success && <div className="mt-3 text-sm text-green-700">{success}</div>}

              <button
                onClick={() => {
                  void confirmBooking()
                }}
                disabled={!canConfirm}
                className="mt-4 w-full bg-orange-600 text-white py-3 rounded-md hover:bg-orange-700 disabled:opacity-50 disabled:cursor-not-allowed"
              >
                {loading ? 'Confirming…' : 'Confirm booking'}
              </button>

              <button
                type="button"
                onClick={onClose}
                className="mt-3 w-full border border-gray-300 text-gray-900 py-3 rounded-md hover:bg-gray-50"
              >
                Cancel
              </button>

              <div className="mt-3 text-xs text-gray-600">
                Booking requires login. Use seeded test user: test@example.com / password123
              </div>
            </div>
          </div>
        </div>
      </div>

      {loginPromptOpen && (
        <div className="absolute inset-0 z-[2100] flex items-center justify-center p-4">
          <div
            className="absolute inset-0 bg-black/20"
            onClick={() => setLoginPromptOpen(false)}
          />
          <div className="relative z-[2101] w-full max-w-md rounded-lg bg-white shadow-lg border border-gray-200 p-5">
            <div className="text-lg font-semibold text-gray-900">Log in to confirm booking</div>
            <div className="mt-1 text-sm text-gray-600">
              You need an account to book storage. Log in to continue.
            </div>

            <div className="mt-4 flex gap-2">
              <button
                type="button"
                onClick={() => {
                  setLoginPromptOpen(false)
                  onClose()
                  const redirect = `${window.location.pathname}${window.location.search}`
                  navigate(`/profile?redirect=${encodeURIComponent(redirect)}`)
                }}
                className="flex-1 bg-blue-600 text-white py-2.5 rounded-md hover:bg-blue-700"
              >
                Log in
              </button>
              <button
                type="button"
                onClick={() => setLoginPromptOpen(false)}
                className="flex-1 border border-gray-300 text-gray-900 py-2.5 rounded-md hover:bg-gray-50"
              >
                Not now
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}
