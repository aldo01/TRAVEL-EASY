# Locker Storage API - Go Backend

A RESTful API for a luggage locker storage system built with Go, Gin, GORM, and PostgreSQL.

## Features

- 🔐 JWT Authentication
- 📍 Geospatial location search (Haversine formula)
- 🗄️ Locker management system
- 📅 Booking system with QR codes
- 💳 Payment tracking
- 🔓 QR code-based locker unlock
- 📊 Access logging

## Tech Stack

- **Language**: Go 1.21+
- **Framework**: Gin Web Framework
- **Database**: PostgreSQL with GORM
- **Authentication**: JWT (golang-jwt)
- **QR Codes**: go-qrcode
- **CORS**: gin-contrib/cors

## Project Structure

```
backend/
├── cmd/
│   ├── server/main.go       # Main application entry point
│   ├── migrate/main.go      # Database migrations
│   └── seed/main.go         # Seed data script
├── internal/
│   ├── handlers/            # HTTP request handlers
│   │   ├── auth.go
│   │   ├── location.go
│   │   └── booking.go
│   ├── middleware/          # Middleware functions
│   │   └── auth.go
│   ├── models/              # Database models
│   │   └── models.go
│   ├── services/            # Business logic
│   │   ├── location_service.go
│   │   └── booking_service.go
│   └── utils/               # Utility functions
│       ├── jwt.go
│       ├── password.go
│       └── response.go
├── config/
│   └── database.go          # Database configuration
├── .env                     # Environment variables
├── go.mod                   # Go dependencies
└── go.sum                   # Dependency checksums
```

## Prerequisites

- Go 1.21 or higher
- PostgreSQL 14 or higher
- Git

## Installation

### 1. Install Go

Download and install Go from [golang.org](https://golang.org/dl/)

### 2. Install PostgreSQL

**Windows:**
```bash
# Using Chocolatey
choco install postgresql

# Or download from https://www.postgresql.org/download/windows/
```

**macOS:**
```bash
brew install postgresql@14
brew services start postgresql@14
```

**Linux:**
```bash
sudo apt update
sudo apt install postgresql postgresql-contrib
sudo systemctl start postgresql
```

### 3. Create Database

```bash
# Access PostgreSQL
psql -U postgres

# Create database
CREATE DATABASE locker_storage;

# Create user (optional)
CREATE USER lockeruser WITH PASSWORD 'yourpassword';
GRANT ALL PRIVILEGES ON DATABASE locker_storage TO lockeruser;

# Exit
\q
```

### 4. Clone and Setup

```bash
# Navigate to backend directory
cd backend

# Install dependencies
go mod download

# Copy environment file
cp .env.example .env

# Edit .env file with your database credentials
```

### 5. Configure Environment

Edit `.env` file:

```env
# Database
DATABASE_URL=postgresql://postgres:yourpassword@localhost:5432/locker_storage?sslmode=disable

# JWT
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production
JWT_EXPIRY=168h

# Server
PORT=8080
GIN_MODE=debug
```

### 6. Run Migrations

```bash
go run cmd/migrate/main.go
```

Expected output:
```
Running database migrations...
✓ Migrations completed successfully
```

### 7. Seed Database (Optional)

```bash
go run cmd/seed/main.go
```

This creates:
- Test user (test@example.com / password123)
- Sample locations across multiple cities (e.g. Copenhagen, Bengaluru, New Delhi, London, New York)
- 6 lockers per location
- 1 sample booking

### 8. Start Server

```bash
go run cmd/server/main.go
```

Server will start on `http://localhost:8080`.

If you run the full stack via the repo-root `docker compose.yml`, the backend is exposed on `http://localhost:8081`.

## API Endpoints

### Authentication

```
POST   /api/v1/auth/register    # Register new user
POST   /api/v1/auth/login       # Login user
GET    /api/v1/profile          # Get user profile (protected)
```

### Locations

```
GET    /api/v1/locations                # List locations (city, type filters)
GET    /api/v1/locations/nearby         # Find nearby locations (lat, lon, radius)
GET    /api/v1/locations/:id            # Get location details
GET    /api/v1/locations/:id/lockers    # Get available lockers at location
```

### Lockers

```
GET    /api/v1/lockers/available        # Search available lockers (size, city filters)
```

### Bookings (Protected)

```
POST   /api/v1/bookings                 # Create new booking
GET    /api/v1/bookings                 # Get user bookings (status filter)
GET    /api/v1/bookings/:id             # Get booking details
POST   /api/v1/bookings/:id/unlock      # Unlock locker with QR code
POST   /api/v1/bookings/:id/checkin     # Check in to locker
POST   /api/v1/bookings/:id/checkout    # Check out from locker
POST   /api/v1/bookings/:id/cancel      # Cancel booking
```

## API Usage Examples

### Register User

```bash
curl -X POST http://localhost:8081/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123",
    "name": "John Doe",
    "phoneNumber": "+4512345678"
  }'
```

### Login

```bash
curl -X POST http://localhost:8081/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'
```

Response:
```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": "uuid",
      "email": "user@example.com",
      "name": "John Doe",
      "phoneNumber": "+4512345678"
    }
  }
}
```

### Find Nearby Locations

```bash
curl "http://localhost:8081/api/v1/locations/nearby?lat=55.6761&lon=12.5683&radius=5"
```

### Create Booking

```bash
curl -X POST http://localhost:8081/api/v1/bookings \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "locationId": "location-uuid",
    "lockerId": "locker-uuid",
    "startTime": "2024-01-15T10:00:00Z",
    "duration": 24,
    "rateType": "DAILY"
  }'
```

### Unlock Locker

```bash
curl -X POST http://localhost:8081/api/v1/bookings/:id/unlock \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "qrCode": "QR-CODE-FROM-BOOKING"
  }'
```

## Database Schema

### Users
- id, email, password, name, phone_number
- JWT authentication

### Locations
- id, name, address, city, country, type
- latitude, longitude (for geospatial queries)
- hourly_rate, daily_rate, rating
- opening_time, closing_time
- amenities stored as text (JSON array string)

### Lockers
- id, location_id, locker_number, size, status
- Sizes: SMALL, MEDIUM, LARGE, XLARGE
- Status: AVAILABLE, OCCUPIED, RESERVED, MAINTENANCE

### Bookings
- id, user_id, location_id, locker_id
- booking_number, start_time, end_time
- rate_type, base_price, total_price
- status, qr_code, payment_status
- check_in_time, check_out_time

### Access Logs
- id, booking_id, action, timestamp, method
- Tracks all locker unlock/lock events

## Development

### Run with Auto-reload

```bash
# Install air for hot reload
go install github.com/cosmtrek/air@latest

# Run with air
air
```

### Run Tests

```bash
go test ./...
```

### Build for Production

```bash
# Build binary
go build -o locker-api cmd/server/main.go

# Run binary
./locker-api
```

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| DATABASE_URL | PostgreSQL connection string | Required |
| JWT_SECRET | Secret key for JWT tokens | Required |
| JWT_EXPIRY | JWT token expiry duration | 168h (7 days) |
| PORT | Server port | 8080 |
| GIN_MODE | Gin mode (debug/release) | debug |

## Troubleshooting

### Database Connection Error

```
Failed to connect to database
```

**Solution:**
- Check PostgreSQL is running: `pg_isready`
- Verify DATABASE_URL in .env
- Check database exists: `psql -U postgres -l`

### Port Already in Use

```
bind: address already in use
```

**Solution:**
- Change PORT in .env
- Or kill process on port 8080:
  ```bash
  # Windows
  netstat -ano | findstr :8080
  taskkill /PID <PID> /F
  
  # Linux/Mac
  lsof -ti:8080 | xargs kill
  ```

### Migration Failed

```
Failed to run migrations
```

**Solution:**
- Drop and recreate database
- Check GORM model definitions
- Run migrations again

## License

MIT

## Support

For issues and questions, please open an issue on GitHub.
