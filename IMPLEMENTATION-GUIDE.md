# Locker Storage System - Complete Implementation Guide

## 🎯 System Overview

**Business Model**: Luggage locker storage system similar to LuggageHero
- Users search for nearby storage locations (gyms, clubs, schools, kiosks)
- Book lockers by size (Small, Medium, Large, XLarge)
- Pay hourly or daily rates
- Access lockers using QR codes
- Real-time availability tracking

## 🏗️ Architecture

### Frontend (React + TypeScript)
- **Map view** with location markers (like LuggageHero UI)
- **Search** by location/address with autocomplete
- **Location list** showing:
  - Distance from user
  - Hourly/daily rates
  - Ratings & reviews
  - Operating hours
  - Available lockers
- **Booking flow**:
  1. Select location
  2. Choose locker size
  3. Select date/time & duration
  4. Payment
  5. Receive QR code

### Backend (Golang + PostgreSQL)
- **REST API** with Gin framework
- **PostgreSQL** database with GORM
- **JWT** authentication
- **QR code** generation for locker access
- **Geospatial** queries for nearby locations

## 📊 Database Schema (PostgreSQL)

### Tables Created:
1. **users** - User accounts with authentication
2. **locations** - Storage points (gyms, clubs, kiosks)
   - Geospatial data (lat/long)
   - Operating hours
   - Pricing (hourly/daily rates)
   - Amenities & ratings
3. **lockers** - Individual storage units
   - Size (Small/Medium/Large/XLarge)
   - Status (Available/Occupied/Reserved)
   - Physical location within facility
4. **bookings** - Locker reservations
   - User, location, locker references
   - Time duration & pricing
   - QR code for access
   - Payment status
5. **access_logs** - QR code scan tracking
   - Lock/unlock events
   - Timestamp & device info

## 🔧 Backend Implementation (Go)

### File Structure:
```
backend/
├── cmd/
│   └── server/
│       └── main.go              # Application entry point
├── internal/
│   ├── models/
│   │   └── models.go            # ✅ GORM models (DONE)
│   ├── handlers/
│   │   ├── auth.go              # Auth endpoints
│   │   ├── location.go          # Location search & details
│   │   ├── locker.go            # Locker availability
│   │   └── booking.go           # Booking CRUD + QR
│   ├── middleware/
│   │   └── auth.go              # JWT validation
│   ├── services/
│   │   ├── location_service.go  # Nearby locations (Haversine)
│   │   ├── booking_service.go   # Booking logic
│   │   └── qr_service.go        # QR code generation
│   └── utils/
│       ├── password.go          # bcrypt hashing
│       ├── jwt.go               # JWT tokens
│       └── response.go          # API responses
├── config/
│   └── database.go              # PostgreSQL connection
├── go.mod                       # ✅ Dependencies (DONE)
└── .env                         # ✅ Configuration (DONE)
```

### Key Features to Implement:

#### 1. **Nearby Locations API**
```go
GET /api/v1/locations/nearby?lat=55.6761&lon=12.5683&radius=5
```
- Use Haversine formula to calculate distance
- Filter by radius (km)
- Return sorted by distance
- Include available locker count

#### 2. **Locker Availability**
```go
GET /api/v1/locations/:id/lockers?size=MEDIUM&date=2026-02-03
```
- Check real-time availability
- Filter by size
- Show pricing

#### 3. **Booking Flow**
```go
POST /api/v1/bookings
{
  "locationId": "uuid",
  "lockerId": "uuid",
  "startTime": "2026-02-03T10:00:00Z",
  "duration": 4, // hours
  "rateType": "HOURLY"
}
```
- Generate unique booking number
- Create QR code (UUID-based)
- Calculate pricing
- Return QR code image

#### 4. **QR Code Access**
```go
POST /api/v1/bookings/:id/unlock
{
  "qrCode": "generated-uuid",
  "action": "UNLOCK"
}
```
- Validate QR code
- Check booking status & time
- Log access event
- Return success/error

## 🗺️ Frontend Implementation (React)

### Components to Create:

#### 1. **Map View** (like LuggageHero)
```tsx
// Uses react-map-gl or @react-google-maps/api
<Map
  initialView={{center: [lat, lon], zoom: 13}}
  markers={locations.map(loc => ({
    position: [loc.latitude, loc.longitude],
    popup: <LocationCard location={loc} />
  }))}
/>
```

#### 2. **Location Search**
```tsx
<SearchBar
  placeholder="Copenhagen Airport (CPH)"
  onSearch={(query) => searchLocations(query)}
  onLocationSelect={(coords) => findNearby(coords)}
/>
```

#### 3. **Location Card**
```tsx
<LocationCard>
  <Badge>Best choice / Highly rated</Badge>
  <Title>LuggageHero Central Station</Title>
  <Status>Open 24 hours • 8.7 km away</Status>
  <Rating>⭐ 4.6 (4477)</Rating>
  <Pricing>
    kr 15/hour • kr 59.95/day
  </Pricing>
  <Button>Book now</Button>
</LocationCard>
```

#### 4. **Booking Modal**
```tsx
<BookingModal>
  <LockerSizeSelector />
  <DateTimePicker />
  <DurationSelector /> // Hours or days
  <PricingSummary />
  <PaymentMethod />
  <ConfirmButton />
</BookingModal>
```

#### 5. **QR Code Display**
```tsx
<BookingConfirmation>
  <QRCodeImage src={booking.qrCodeImageUrl} />
  <BookingNumber>{booking.bookingNumber}</BookingNumber>
  <Instructions>
    Show this QR code at the locker to unlock
  </Instructions>
  <LockerDetails>
    Location: {location.name}
    Locker: {locker.lockerNumber}
  </LockerDetails>
</BookingConfirmation>
```

## 🔐 Security Features

1. **JWT Authentication**
   - Register/Login with bcrypt passwords
   - Protected routes require Bearer token

2. **QR Code Validation**
   - Unique UUID per booking
   - Time-bound access (only during booking period)
   - One-time unlock prevention

3. **Payment Integration** (Future)
   - Stripe/PayPal for online payment
   - Track payment status

## 📡 API Endpoints Summary

### Auth
- `POST /api/v1/auth/register`
- `POST /api/v1/auth/login`
- `GET /api/v1/auth/profile` (Protected)

### Locations
- `GET /api/v1/locations` - All locations
- `GET /api/v1/locations/nearby` - Nearby search
- `GET /api/v1/locations/:id` - Location details
- `GET /api/v1/locations/:id/lockers` - Available lockers

### Bookings
- `POST /api/v1/bookings` - Create booking (Protected)
- `GET /api/v1/bookings` - User's bookings (Protected)
- `GET /api/v1/bookings/:id` - Booking details (Protected)
- `POST /api/v1/bookings/:id/unlock` - QR unlock (Protected)
- `POST /api/v1/bookings/:id/checkin` - Check-in
- `POST /api/v1/bookings/:id/checkout` - Check-out
- `PUT /api/v1/bookings/:id/cancel` - Cancel booking

## 💰 Pricing Logic

```go
// Calculate total price
func CalculatePrice(rateType RateType, duration float64, hourlyRate, dailyRate float64) float64 {
    switch rateType {
    case RateHourly:
        return duration * hourlyRate
    case RateDaily:
        days := math.Ceil(duration / 24)
        return days * dailyRate
    case RateMultiDay:
        days := math.Ceil(duration / 24)
        if days > 1 {
            // Discount for multi-day
            return (dailyRate * days) * 0.9
        }
        return dailyRate
    }
}
```

## 🚀 Deployment Steps

### Backend (Go)
1. Install PostgreSQL
2. Create database: `CREATE DATABASE locker_storage;`
3. Update `.env` with DATABASE_URL
4. Run migrations: `go run cmd/migrate/main.go`
5. Seed data: `go run cmd/seed/main.go`
6. Start server: `go run cmd/server/main.go`

### Frontend (React)
1. Install dependencies: `npm install`
2. Add map library: `npm install react-map-gl mapbox-gl`
3. Update API URL in services
4. Start dev server: `npm run dev`

## 🎯 Next Steps

Since the complete Go codebase is extensive (20+ files), here's what I recommend:

### Option 1: I'll create all Go files now
I'll generate all handlers, services, middleware, and the main server file.

### Option 2: Focus on specific parts first
Let me know which part you want me to build first:
- Authentication system
- Location & nearby search
- Booking system with QR codes
- Frontend UI components

### Option 3: Use existing backend frameworks
Consider using a Go boilerplate/template to speed up development.

**What would you like me to do next?**

---

**Status**: 
- ✅ Database schema designed (models.go)
- ✅ Go project structure created
- ⏳ Handlers & services pending
- ⏳ Frontend UI redesign pending

