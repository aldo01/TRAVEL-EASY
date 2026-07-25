import { useEffect, useMemo, useState } from 'react'
import { Link } from 'react-router-dom'
import { Calendar, MapPin, Lock } from 'lucide-react'
import type { StorageBooking } from '../types'
import { storageService } from '../services/storage'
import { authStore } from '../services/auth'

export default function Bookings() {
  const [activeTab, setActiveTab] = useState<'active' | 'history'>('active')
  const [bookings, setBookings] = useState<StorageBooking[]>([])
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const token = authStore.getToken()

  useEffect(() => {
    if (!token) return

    let cancelled = false
    ;(async () => {
      setLoading(true)
      setError(null)
      try {
        const data = await storageService.getMyBookings()
        if (!cancelled) setBookings(data)
      } catch (e) {
        if (!cancelled) setError(e instanceof Error ? e.message : 'Failed to load bookings')
      } finally {
        if (!cancelled) setLoading(false)
      }
    })()

    return () => {
      cancelled = true
    }
  }, [token])

  const { active, history } = useMemo(() => {
    const now = Date.now()
    const isHistory = (b: StorageBooking) => {
      if (b.status === 'COMPLETED' || b.status === 'CANCELLED' || b.status === 'EXPIRED') return true
      if (b.endTime && new Date(b.endTime).getTime() < now) return true
      return false
    }

    const activeBookings = bookings.filter((b) => !isHistory(b))
    const historyBookings = bookings.filter((b) => isHistory(b))
    return { active: activeBookings, history: historyBookings }
  }, [bookings])

  const current = activeTab === 'active' ? active : history

  const formatDateTime = (iso: string) => {
    const date = new Date(iso)
    return date.toLocaleString(undefined, {
      month: 'short',
      day: 'numeric',
      year: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    })
  }

  if (!token) {
    return (
      <div className="space-y-6">
        <div>
          <h1 className="text-3xl font-bold text-gray-900">My bookings</h1>
          <p className="text-gray-600 mt-1">Log in to see your locker reservations.</p>
        </div>
        <div className="bg-white rounded-lg shadow-md p-6">
          <Link to="/profile" className="text-blue-600 hover:text-blue-700 font-medium">
            Go to login →
          </Link>
        </div>
      </div>
    )
  }

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-3xl font-bold text-gray-900">My bookings</h1>
        <p className="text-gray-600 mt-1">Your luggage storage reservations and history.</p>
      </div>

      <div className="border-b border-gray-200">
        <div className="flex space-x-8">
          <button
            onClick={() => setActiveTab('active')}
            className={`pb-4 px-1 border-b-2 font-medium text-sm transition-colors ${
              activeTab === 'active'
                ? 'border-blue-500 text-blue-600'
                : 'border-transparent text-gray-500 hover:text-gray-700'
            }`}
          >
            Active ({active.length})
          </button>
          <button
            onClick={() => setActiveTab('history')}
            className={`pb-4 px-1 border-b-2 font-medium text-sm transition-colors ${
              activeTab === 'history'
                ? 'border-blue-500 text-blue-600'
                : 'border-transparent text-gray-500 hover:text-gray-700'
            }`}
          >
            History ({history.length})
          </button>
        </div>
      </div>

      {error && <div className="text-sm text-red-600">{error}</div>}
      {loading && <div className="text-sm text-gray-600">Loading…</div>}

      <div className="grid gap-4">
        {current.map((b) => (
          <div key={b.id} className="bg-white rounded-lg shadow-md p-6">
            <div className="flex items-start justify-between">
              <div>
                <div className="text-sm text-gray-500">Booking</div>
                <div className="text-lg font-semibold text-gray-900">#{b.bookingNumber}</div>
              </div>
              <div className="text-right">
                <div className="text-xs px-2 py-1 rounded bg-gray-100 text-gray-700 inline-block">{b.status}</div>
                <div className="text-sm text-gray-600 mt-2">{b.paymentStatus}</div>
              </div>
            </div>

            <div className="mt-4 grid grid-cols-1 md:grid-cols-3 gap-4 text-sm text-gray-600">
              <div className="flex items-start gap-2">
                <MapPin className="w-4 h-4 mt-0.5 text-blue-600" />
                <div>
                  <div className="font-medium text-gray-900">{b.location?.name || 'Location'}</div>
                  <div>{b.location?.address || ''}</div>
                </div>
              </div>

              <div className="flex items-start gap-2">
                <Lock className="w-4 h-4 mt-0.5 text-blue-600" />
                <div>
                  <div className="font-medium text-gray-900">
                    Locker #{b.locker?.lockerNumber || '—'}
                  </div>
                  <div>{b.locker?.size || ''}</div>
                </div>
              </div>

              <div className="flex items-start gap-2">
                <Calendar className="w-4 h-4 mt-0.5 text-blue-600" />
                <div>
                  <div>
                    {formatDateTime(b.startTime)}
                  </div>
                  {b.endTime && <div>Until {formatDateTime(b.endTime)}</div>}
                </div>
              </div>
            </div>

            <div className="mt-4 flex items-center justify-between">
              <div className="text-sm text-gray-600">
                Total: <span className="font-medium text-gray-900">{b.totalPrice}</span>
              </div>
              <div className="text-xs text-gray-500">Rate: {b.rateType}</div>
            </div>
          </div>
        ))}
      </div>

      {!loading && current.length === 0 && (
        <div className="bg-white rounded-lg shadow-md p-6 text-gray-600">
          No bookings in this view yet. <Link to="/locations" className="text-blue-600 hover:text-blue-700">Book a locker</Link>.
        </div>
      )}
    </div>
  )
}
