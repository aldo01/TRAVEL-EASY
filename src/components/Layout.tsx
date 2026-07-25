import type { ReactNode } from 'react'
import { Link, useLocation } from 'react-router-dom'
import { Home, MapPin, Calendar, User, Store, LogIn } from 'lucide-react'
import { authStore } from '../services/auth'

interface LayoutProps {
  children: ReactNode
}

export default function Layout({ children }: LayoutProps) {
  const location = useLocation()
  const isAuthed = Boolean(authStore.getToken())

  const navItems = [
    { path: '/', label: 'Home', icon: Home },
    { path: '/locations', label: 'Locations', icon: MapPin },
    { path: '/bookings', label: 'Bookings', icon: Calendar },
    { path: '/host', label: 'List your space', icon: Store },
  ]

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Header */}
      <header className="bg-white shadow-sm">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4 flex items-center justify-between">
          <div>
            <h1 className="text-2xl font-bold text-blue-600">Travel Easy</h1>
            <p className="text-sm text-gray-500">Luggage storage & locker booking</p>
          </div>
          <div>
            {isAuthed ? (
              <Link
                to="/profile"
                className="inline-flex items-center gap-2 px-4 py-2 rounded-full bg-blue-50 text-blue-700 hover:bg-blue-100 transition-colors font-medium text-sm"
              >
                <div className="w-8 h-8 bg-blue-600 text-white rounded-full flex items-center justify-center">
                  <User size={16} />
                </div>
                My Account
              </Link>
            ) : (
              <div className="flex items-center gap-3">
                <Link
                  to="/profile"
                  className="inline-flex items-center gap-2 text-gray-600 hover:text-gray-900 text-sm font-medium"
                >
                  <LogIn size={18} />
                  Log in
                </Link>
                <Link
                  to="/profile?mode=register"
                  className="inline-flex items-center gap-2 bg-blue-600 text-white px-4 py-2 rounded-md hover:bg-blue-700 transition-colors text-sm font-medium"
                >
                  Sign up
                </Link>
              </div>
            )}
          </div>
        </div>
      </header>

      {/* Navigation */}
      <nav className="bg-white border-b">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex space-x-8">
            {navItems.map(({ path, label, icon: Icon }) => (
              <Link
                key={path}
                to={path}
                className={`flex items-center space-x-2 px-3 py-4 text-sm font-medium border-b-2 transition-colors ${
                  location.pathname === path
                    ? 'border-blue-500 text-blue-600'
                    : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
                }`}
              >
                <Icon size={20} />
                <span>{label}</span>
              </Link>
            ))}
          </div>
        </div>
      </nav>

      {/* Main Content */}
      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {children}
      </main>
    </div>
  )
}
