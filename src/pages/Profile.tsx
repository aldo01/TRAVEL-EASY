import { User as UserIcon, CreditCard, Settings, HelpCircle, Mail, LogOut } from 'lucide-react'
import type { User } from '../types'

export default function Profile() {
  const user: User = {
    id: '1',
    name: 'John Doe',
    email: 'john.doe@example.com',
    phoneNumber: '+1 234 567 8900',
  }

  const menuItems = [
    {
      section: 'Account',
      items: [
        { icon: UserIcon, label: 'Personal Information', href: '#' },
        { icon: CreditCard, label: 'Payment Methods', href: '#' },
        { icon: Settings, label: 'Travel Preferences', href: '#' },
      ],
    },
    {
      section: 'Support',
      items: [
        { icon: HelpCircle, label: 'Help Center', href: '#' },
        { icon: Mail, label: 'Contact Us', href: '#' },
      ],
    },
  ]

  return (
    <div className="max-w-4xl mx-auto space-y-6">
      {/* Profile Header */}
      <div className="bg-white rounded-lg shadow-md p-8">
        <div className="flex items-center">
          <div className="w-24 h-24 bg-blue-100 rounded-full flex items-center justify-center mr-6">
            <UserIcon size={48} className="text-blue-600" />
          </div>
          <div className="flex-1">
            <h1 className="text-3xl font-bold text-gray-900">{user.name}</h1>
            <p className="text-gray-600 mt-1">{user.email}</p>
            {user.phoneNumber && <p className="text-gray-600">{user.phoneNumber}</p>}
          </div>
          <button className="text-blue-600 hover:text-blue-700 px-4 py-2 rounded-md hover:bg-blue-50 transition-colors font-medium">
            Edit Profile
          </button>
        </div>
      </div>

      {/* Account Stats */}
      <div className="grid grid-cols-3 gap-4">
        <div className="bg-white rounded-lg shadow-md p-6 text-center">
          <div className="text-3xl font-bold text-blue-600">12</div>
          <div className="text-gray-600 text-sm mt-1">Total Trips</div>
        </div>
        <div className="bg-white rounded-lg shadow-md p-6 text-center">
          <div className="text-3xl font-bold text-blue-600">28</div>
          <div className="text-gray-600 text-sm mt-1">Bookings</div>
        </div>
        <div className="bg-white rounded-lg shadow-md p-6 text-center">
          <div className="text-3xl font-bold text-blue-600">15</div>
          <div className="text-gray-600 text-sm mt-1">Countries Visited</div>
        </div>
      </div>

      {/* Menu Sections */}
      {menuItems.map((section) => (
        <div key={section.section} className="bg-white rounded-lg shadow-md overflow-hidden">
          <div className="px-6 py-4 bg-gray-50 border-b border-gray-200">
            <h2 className="text-lg font-semibold text-gray-900">{section.section}</h2>
          </div>
          <div className="divide-y divide-gray-200">
            {section.items.map((item) => (
              <a
                key={item.label}
                href={item.href}
                className="flex items-center px-6 py-4 hover:bg-gray-50 transition-colors"
              >
                <item.icon className="w-5 h-5 text-gray-400 mr-4" />
                <span className="flex-1 text-gray-900">{item.label}</span>
                <span className="text-gray-400">→</span>
              </a>
            ))}
          </div>
        </div>
      ))}

      {/* Logout Button */}
      <div className="bg-white rounded-lg shadow-md overflow-hidden">
        <button className="w-full flex items-center px-6 py-4 hover:bg-red-50 transition-colors text-left">
          <LogOut className="w-5 h-5 text-red-500 mr-4" />
          <span className="flex-1 text-red-600 font-medium">Log Out</span>
        </button>
      </div>
    </div>
  )
}
