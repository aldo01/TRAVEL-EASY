import { useEffect, useMemo, useState, type FormEvent } from 'react'
import { Link, useNavigate, useSearchParams } from 'react-router-dom'
import { LogOut, User as UserIcon } from 'lucide-react'
import type { User } from '../types'
import { authService, authStore } from '../services/auth'

type Mode = 'login' | 'register'

export default function Profile() {
  const navigate = useNavigate()
  const [searchParams] = useSearchParams()
  const [mode, setMode] = useState<Mode>(() => {
    const m = searchParams.get('mode')
    return m === 'register' ? 'register' : 'login'
  })
  const [user, setUser] = useState<User | null>(null)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const redirectTo = useMemo(() => {
    const value = searchParams.get('redirect')
    if (!value) return null
    // Keep it same-origin relative paths only.
    if (!value.startsWith('/')) return null
    return value
  }, [searchParams])

  const [name, setName] = useState('')
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [phoneNumber, setPhoneNumber] = useState('')

  useEffect(() => {
    const token = authStore.getToken()
    if (!token) return

    let cancelled = false
    ;(async () => {
      setLoading(true)
      setError(null)
      try {
        const me = await authService.getProfile()
        if (!cancelled) setUser(me)
      } catch (e) {
        authStore.clear()
        if (!cancelled) setError(e instanceof Error ? e.message : 'Failed to load profile')
      } finally {
        if (!cancelled) setLoading(false)
      }
    })()

    return () => {
      cancelled = true
    }
  }, [])

  const onLogout = () => {
    authStore.clear()
    setUser(null)
    setError(null)
  }

  const onSubmit = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault()

    setLoading(true)
    setError(null)
    try {
      if (mode === 'login') {
        const me = await authService.login(email.trim(), password)
        setUser(me)
      } else {
        const me = await authService.register(name.trim(), email.trim(), password, phoneNumber.trim() || undefined)
        setUser(me)
      }

      if (redirectTo) {
        navigate(redirectTo, { replace: true })
      }
    } catch (e) {
      setError(e instanceof Error ? e.message : 'Authentication failed')
    } finally {
      setLoading(false)
    }
  }

  if (user) {
    return (
      <div className="max-w-3xl space-y-6">
        <div>
          <h1 className="text-3xl font-bold text-gray-900">Profile</h1>
          <p className="text-gray-600 mt-1">Manage your account and bookings.</p>
        </div>

        {error && <div className="text-sm text-red-600">{error}</div>}

        <div className="bg-white rounded-lg shadow-md p-8">
          <div className="flex items-center">
            <div className="w-16 h-16 bg-blue-100 rounded-full flex items-center justify-center mr-6">
              <UserIcon size={32} className="text-blue-600" />
            </div>
            <div className="flex-1">
              <div className="text-xl font-semibold text-gray-900">{user.name}</div>
              <div className="text-gray-600">{user.email}</div>
              {user.phoneNumber && <div className="text-gray-600">{user.phoneNumber}</div>}
            </div>
            <button
              onClick={onLogout}
              className="inline-flex items-center gap-2 text-red-600 hover:text-red-700 px-4 py-2 rounded-md hover:bg-red-50 transition-colors font-medium"
            >
              <LogOut size={18} />
              Log out
            </button>
          </div>
        </div>

        <div className="bg-white rounded-lg shadow-md p-6 text-sm text-gray-600">
          You are logged in. Go to <Link className="text-blue-600 hover:text-blue-700" to="/bookings">Bookings</Link> to view reservations.
        </div>
      </div>
    )
  }

  return (
    <div className="max-w-lg space-y-6">
      <div>
        <h1 className="text-3xl font-bold text-gray-900">Account</h1>
        <p className="text-gray-600 mt-1">Log in to book lockers and manage bookings.</p>
      </div>

      <div className="bg-white rounded-lg shadow-md p-6">
        <div className="flex gap-2 mb-6">
          <button
            onClick={() => setMode('login')}
            className={`flex-1 py-2 rounded-md text-sm font-medium ${
              mode === 'login' ? 'bg-blue-600 text-white' : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
            }`}
          >
            Log in
          </button>
          <button
            onClick={() => setMode('register')}
            className={`flex-1 py-2 rounded-md text-sm font-medium ${
              mode === 'register' ? 'bg-blue-600 text-white' : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
            }`}
          >
            Register
          </button>
        </div>

        {error && <div className="text-sm text-red-600 mb-4">{error}</div>}

        <form onSubmit={onSubmit} className="space-y-4">
          {mode === 'register' && (
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">Name</label>
              <input
                value={name}
                onChange={(e) => setName(e.target.value)}
                className="w-full px-4 py-3 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                placeholder="Your name"
              />
            </div>
          )}

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">Email</label>
            <input
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              type="email"
              className="w-full px-4 py-3 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              placeholder="you@example.com"
            />
          </div>

          {mode === 'register' && (
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">Phone (optional)</label>
              <input
                value={phoneNumber}
                onChange={(e) => setPhoneNumber(e.target.value)}
                className="w-full px-4 py-3 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                placeholder="+1 234 567 8900"
              />
            </div>
          )}

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">Password</label>
            <input
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              type="password"
              className="w-full px-4 py-3 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              placeholder="••••••••"
            />
          </div>

          <button
            type="submit"
            disabled={loading}
            className="w-full bg-blue-600 text-white py-3 px-6 rounded-md hover:bg-blue-700 disabled:opacity-50"
          >
            {loading ? 'Please wait…' : mode === 'login' ? 'Log in' : 'Create account'}
          </button>
        </form>

        <div className="mt-6 text-xs text-gray-500">
          Seeded test user: test@example.com / password123
        </div>
      </div>
    </div>
  )
}
