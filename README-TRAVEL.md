# Travel Easy - Web Application

A modern web application for travel planning and booking management built with React, TypeScript, and Tailwind CSS. Inspired by leading travel booking platforms.

## Features

- 🏠 **Home Page**: Search destinations, browse featured locations, and discover popular activities
- ✈️ **Trips Management**: Plan and track your upcoming and past trips
- 📅 **Bookings**: Manage flights, hotels, car rentals, and activities in one place
- 👤 **User Profile**: Manage account settings and travel preferences
- 📱 **Responsive Design**: Works seamlessly on desktop, tablet, and mobile devices

## Tech Stack

- **Frontend**: React 19 + TypeScript
- **Build Tool**: Vite
- **Styling**: Tailwind CSS
- **Icons**: Lucide React
- **Routing**: React Router v6
- **State Management**: React Hooks

## Getting Started on Windows

### Prerequisites

- Node.js 18+ installed on your machine
- npm or yarn package manager

### Installation

1. **Install dependencies:**
```bash
npm install
```

2. **Start the development server:**
```bash
npm run dev
```

The app will open at `http://localhost:3000`

### Available Scripts

- `npm run dev` - Start development server
- `npm run build` - Build for production
- `npm run preview` - Preview production build
- `npm run lint` - Run ESLint

## Project Structure

```
travel-easy/
├── src/
│   ├── components/           # Reusable UI components
│   │   └── Layout.tsx       # Main layout with navigation
│   ├── pages/               # Page components
│   │   ├── Home.tsx         # Home page with search
│   │   ├── Trips.tsx        # Trips management
│   │   ├── Bookings.tsx     # Bookings list
│   │   └── Profile.tsx      # User profile
│   ├── services/            # API services
│   │   ├── api.ts           # Base API client
│   │   └── travel.ts        # Travel-specific services
│   ├── types/               # TypeScript interfaces
│   │   └── index.ts         # Data models
│   ├── App.tsx              # Main app component
│   ├── main.tsx             # App entry point
│   └── index.css            # Global styles
└── public/                  # Static assets
```

## Design Inspiration

This application is inspired by modern travel booking platforms like **LuggageHero**, focusing on:
- Clean, intuitive user interface
- Easy booking flow
- Clear information hierarchy
- Mobile-first responsive design
- Trust indicators (ratings, stats, reviews)

## Key Features Implemented

### Home Page
- **Search Bar**: Find destinations by location, date, and travelers
- **Stats Section**: 2M+ bookings, 4.7★ rating, 24/7 support
- **Featured Destinations**: Showcase popular locations
- **Popular Activities**: Browse top-rated experiences
- **Why Choose Us**: Platform benefits

### Trips Management
- View all trips with status badges
- Filter by upcoming/ongoing/completed
- Trip details (dates, travelers, duration)
- Create and manage trips

### Bookings System
- Separate tabs for upcoming and past bookings
- Multiple booking types (flights, hotels, cars, activities)
- Confirmation numbers and details
- Download vouchers

### User Profile
- Personal information
- Account statistics
- Payment methods
- Travel preferences
- Help center access

## API Integration

Currently using mock data. To integrate with a real backend:

1. Update `API_BASE_URL` in `src/services/api.ts`
2. Add authentication headers
3. Replace mock data with API calls

```typescript
import { travelService } from './services/travel'

// Fetch trips
const trips = await travelService.getTrips()

// Create booking
const booking = await travelService.createBooking({...})
```

## Deployment

### Build for Production
```bash
npm run build
```

### Deploy to:
- **Vercel**: `vercel deploy`
- **Netlify**: Drag & drop `dist/` folder
- **GitHub Pages**: Deploy `dist/` folder

## Future Enhancements

- [ ] User authentication
- [ ] Real-time availability
- [ ] Payment integration
- [ ] Email confirmations
- [ ] Multi-language support
- [ ] Dark mode
- [ ] PWA support
- [ ] Reviews system

## License

This project is for educational purposes.

---

Built with ❤️ using React + TypeScript + Tailwind CSS
