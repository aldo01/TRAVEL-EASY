import { useState } from 'react'
import { Link } from 'react-router-dom'
import {
  Search,
  MapPin,
  QrCode,
  ShieldCheck,
  CalendarCheck,
  Headset,
  Wallet,
  Star,
  ChevronDown,
} from 'lucide-react'

const steps = [
  {
    icon: Search,
    title: 'Book storage online',
    body: 'Choose a convenient locker or partner shop and save your spot in a couple of minutes. Prices and hours vary by location.',
  },
  {
    icon: MapPin,
    title: 'Go to the storage spot',
    body: 'Follow the directions to the locker or local shop and drop off your bags whenever suits your day.',
  },
  {
    icon: QrCode,
    title: 'Show your QR code',
    body: 'Use your booking confirmation to check in. Your bags are stored securely until you come back for them.',
  },
]

const benefits = [
  {
    icon: CalendarCheck,
    title: 'Free cancellation',
    body: 'Plans change. Cancel any time before drop-off for a full refund — no fees, ever.',
  },
  {
    icon: ShieldCheck,
    title: 'Bag protection',
    body: 'Every booking is covered by protection of up to ₹8,00,000 for total peace of mind.',
  },
  {
    icon: Headset,
    title: '24/7 support',
    body: 'Our team is on hand around the clock in case you need a hand with your booking.',
  },
  {
    icon: Wallet,
    title: 'Great value',
    body: 'Transparent daily and hourly rates with no surprise charges for larger bags.',
  },
]

const testimonials = [
  {
    name: 'Alexia C.',
    city: 'Berlin',
    text: 'Very uncomplicated and friendly. Drop off and pick up was smooth and quick.',
  },
  {
    name: 'Kiran P.',
    city: 'Tokyo',
    text: 'Very easy to find and the whole process was seamless. Bags felt completely safe.',
  },
  {
    name: 'Susan H.',
    city: 'Florence',
    text: 'Nearby our accommodation, felt secure, and let us keep exploring before our flight home.',
  },
]

const cities = [
  'New York',
  'London',
  'Paris',
  'Rome',
  'Tokyo',
  'Barcelona',
  'Amsterdam',
  'Dubai',
  'Singapore',
  'Sydney',
  'Bangkok',
  'Mumbai',
]

const faqs = [
  {
    q: 'How does Travel Easy luggage storage work?',
    a: 'Tell us where and when you need storage and how many bags you have. We show nearby verified locations. Book online, head to the spot, and show your QR code to drop off your bags. Pick them up the same way.',
  },
  {
    q: 'Can I cancel and get a refund?',
    a: 'Yes. Cancel any time before your drop-off time for a full refund to your original payment method. We never charge cancellation fees.',
  },
  {
    q: 'Is there a size or weight limit?',
    a: 'No strict limit — from backpacks to suitcases and bulky items, most locations can accommodate what you need. Smaller bags are charged less.',
  },
  {
    q: 'How long can I store my luggage?',
    a: 'From a couple of hours to several days. Just choose your drop-off and pick-up times when booking, and extend later from your bookings page if needed.',
  },
  {
    q: 'Are my bags protected?',
    a: 'Every online booking includes bag protection of up to ₹8,00,000, and locations are verified partners with secure storage areas.',
  },
]

function FaqItem({ q, a }: { q: string; a: string }) {
  const [open, setOpen] = useState(false)
  return (
    <div className="border-b border-gray-200">
      <button
        type="button"
        onClick={() => setOpen((v) => !v)}
        className="w-full flex items-center justify-between gap-4 py-4 text-left"
      >
        <span className="text-base font-medium text-gray-900">{q}</span>
        <ChevronDown
          size={20}
          className={`shrink-0 text-gray-500 transition-transform ${open ? 'rotate-180' : ''}`}
        />
      </button>
      {open && <p className="pb-4 text-sm leading-relaxed text-gray-600">{a}</p>}
    </div>
  )
}

export default function LandingSections() {
  return (
    <div className="space-y-16">
      {/* Trust badges */}
      <section className="grid grid-cols-1 gap-4 sm:grid-cols-3">
        {[
          { value: '32,000+', label: 'Storage spots worldwide' },
          { value: '4.9 / 5', label: 'Average customer rating' },
          { value: 'From ₹99', label: 'Per bag, per day' },
        ].map((s) => (
          <div key={s.label} className="rounded-lg bg-white p-6 text-center shadow-md">
            <div className="text-2xl font-bold text-blue-600">{s.value}</div>
            <div className="mt-1 text-sm text-gray-600">{s.label}</div>
          </div>
        ))}
      </section>

      {/* How it works */}
      <section>
        <div className="text-center">
          <h2 className="text-2xl font-bold text-gray-900">How it works</h2>
          <p className="mt-2 text-gray-600">Store your bags in three simple steps.</p>
        </div>
        <div className="mt-8 grid grid-cols-1 gap-6 md:grid-cols-3">
          {steps.map((step, i) => (
            <div key={step.title} className="relative rounded-lg bg-white p-6 shadow-md">
              <div className="flex h-12 w-12 items-center justify-center rounded-full bg-blue-50 text-blue-600">
                <step.icon size={24} />
              </div>
              <div className="mt-4 flex items-center gap-2">
                <span className="text-sm font-semibold text-blue-600">Step {i + 1}</span>
              </div>
              <h3 className="mt-1 text-lg font-semibold text-gray-900">{step.title}</h3>
              <p className="mt-2 text-sm leading-relaxed text-gray-600">{step.body}</p>
            </div>
          ))}
        </div>
      </section>

      {/* Benefits */}
      <section>
        <div className="text-center">
          <h2 className="text-2xl font-bold text-gray-900">Why travelers choose Travel Easy</h2>
          <p className="mt-2 text-gray-600">Flexibility and protection every time you store.</p>
        </div>
        <div className="mt-8 grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-4">
          {benefits.map((b) => (
            <div key={b.title} className="rounded-lg bg-white p-6 shadow-md">
              <div className="flex h-12 w-12 items-center justify-center rounded-full bg-green-50 text-green-600">
                <b.icon size={24} />
              </div>
              <h3 className="mt-4 text-base font-semibold text-gray-900">{b.title}</h3>
              <p className="mt-2 text-sm leading-relaxed text-gray-600">{b.body}</p>
            </div>
          ))}
        </div>
      </section>

      {/* Testimonials */}
      <section>
        <div className="text-center">
          <h2 className="text-2xl font-bold text-gray-900">Trusted by travelers everywhere</h2>
          <div className="mt-2 flex items-center justify-center gap-1 text-yellow-500">
            {Array.from({ length: 5 }).map((_, i) => (
              <Star key={i} size={18} fill="currentColor" />
            ))}
            <span className="ml-2 text-sm text-gray-600">4.9 out of 5 based on thousands of reviews</span>
          </div>
        </div>
        <div className="mt-8 grid grid-cols-1 gap-6 md:grid-cols-3">
          {testimonials.map((t) => (
            <div key={t.name} className="rounded-lg bg-white p-6 shadow-md">
              <div className="flex items-center gap-1 text-yellow-500">
                {Array.from({ length: 5 }).map((_, i) => (
                  <Star key={i} size={16} fill="currentColor" />
                ))}
              </div>
              <p className="mt-3 text-sm leading-relaxed text-gray-700">“{t.text}”</p>
              <div className="mt-4 flex items-center gap-3">
                <div className="flex h-9 w-9 items-center justify-center rounded-full bg-blue-600 text-sm font-semibold text-white">
                  {t.name.charAt(0)}
                </div>
                <div>
                  <div className="text-sm font-medium text-gray-900">{t.name}</div>
                  <div className="text-xs text-gray-500">{t.city}</div>
                </div>
              </div>
            </div>
          ))}
        </div>
      </section>

      {/* FAQ */}
      <section className="mx-auto max-w-3xl">
        <div className="text-center">
          <h2 className="text-2xl font-bold text-gray-900">Frequently asked questions</h2>
        </div>
        <div className="mt-6 rounded-lg bg-white px-6 shadow-md">
          {faqs.map((f) => (
            <FaqItem key={f.q} q={f.q} a={f.a} />
          ))}
        </div>
      </section>

      {/* Popular cities */}
      <section>
        <h2 className="text-2xl font-bold text-gray-900">Store your luggage all over the world</h2>
        <div className="mt-6 flex flex-wrap gap-3">
          {cities.map((c) => (
            <Link
              key={c}
              to={`/locations?city=${encodeURIComponent(c)}`}
              className="rounded-full border border-gray-200 bg-white px-4 py-2 text-sm text-gray-700 shadow-sm transition-colors hover:border-blue-300 hover:text-blue-600"
            >
              {c}
            </Link>
          ))}
        </div>
      </section>

      {/* Host CTA */}
      <section className="overflow-hidden rounded-2xl bg-blue-600">
        <div className="flex flex-col items-start gap-4 p-8 md:flex-row md:items-center md:justify-between">
          <div>
            <h2 className="text-2xl font-bold text-white">Earn with your space</h2>
            <p className="mt-2 max-w-xl text-blue-100">
              Turn spare space in your shop, hotel or cafe into income. Join as a partner and start
              accepting bag drop-offs today.
            </p>
          </div>
          <Link
            to="/host"
            className="shrink-0 rounded-md bg-white px-6 py-3 font-medium text-blue-600 transition-colors hover:bg-blue-50"
          >
            Become a partner
          </Link>
        </div>
      </section>
    </div>
  )
}
