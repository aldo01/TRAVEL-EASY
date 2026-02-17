# Travel Easy - Quick Start Guide

Complete guide to set up and run the Travel Easy locker storage system.

## System Overview

**Travel Easy** is a web-based luggage locker storage system where users can:
- Find nearby storage locations (gyms, clubs, schools, kiosks)
- Book lockers in different sizes
- Use QR codes to unlock/lock lockers
- Track bookings and payments

## Architecture

- **Frontend**: React 19 + TypeScript + Vite + Tailwind CSS (Port 5173)
- **Backend**: Go + Gin + GORM + PostgreSQL (Port 8080)
- **Database**: PostgreSQL with geospatial support

---

## Prerequisites

### Required Software

1. **Node.js** 18+ ([Download](https://nodejs.org/))
2. **Go** 1.21+ ([Download](https://golang.org/dl/))
3. **PostgreSQL** 14+ ([Download](https://www.postgresql.org/download/))

### Verify Installation

```bash
node --version    # Should be 18+
npm --version     # Should be 9+
go version        # Should be 1.21+
psql --version    # Should be 14+
```

---

## Setup Instructions

### Step 1: Database Setup

#### Windows (PowerShell)
```powershell
# Start PostgreSQL service
net start postgresql-x64-14

# Access PostgreSQL
psql -U postgres

# In psql console:
CREATE DATABASE locker_storage;
\q
```

#### macOS/Linux
```bash
# Start PostgreSQL
brew services start postgresql@14  # macOS
# or
sudo systemctl start postgresql    # Linux

# Access PostgreSQL
psql -U postgres

# In psql console:
CREATE DATABASE locker_storage;
\q
```

### Step 2: Backend Setup

```bash
# Navigate to backend directory
cd backend

# Install dependencies
go mod download

# Configure environment
# Edit .env file and update DATABASE_URL if needed
# Default: postgresql://postgres:postgres@localhost:5432/locker_storage?sslmode=disable

# Run migrations
go run cmd/migrate/main.go

# Seed database with sample data
go run cmd/seed/main.go
```

**Expected Output:**
```
✓ Created test user: test@example.com
✓ Created location: Copenhagen Central Station
  ✓ Created 6 lockers for Copenhagen Central Station
...
✓ Database seeding completed successfully!

Test credentials:
Email: test@example.com
Password: password123
```

### Step 3: Frontend Setup

```bash
# Navigate to root directory (where package.json is)
cd ..

# Install dependencies
npm install

# Start development server
npm run dev
```

### Step 4: Start Backend Server

Open a **new terminal window**:

```bash
cd backend
go run cmd/server/main.go
```

**Expected Output:**
```
Starting server on port 8080...
[GIN-debug] Listening and serving HTTP on :8080
```

---

## Running the Application

### Terminal 1 - Backend Server
```bash
cd backend
go run cmd/server/main.go
```
→ Backend API: http://localhost:8080

### Terminal 2 - Frontend Server
```bash
npm run dev
```
→ Frontend App: http://localhost:5173

---

## Test the Application

### 1. Health Check

```bash
curl http://localhost:8080/health
```

**Expected Response:**
```json
{
  "status": "healthy",
  "message": "Locker Storage API is running"
}
```

### 2. Login with Test User

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }'
```

**Expected Response:**
```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": "uuid",
      "email": "test@example.com",
      "firstName": "John",
      "lastName": "Doe"
    }
  }
}
```

### 3. Get Nearby Locations (Copenhagen)

```bash
curl "http://localhost:8080/api/v1/locations/nearby?lat=55.6761&lon=12.5683&radius=10"
```

### 4. Open Frontend

Visit: http://localhost:5173

**Test User Login:**
- Email: `test@example.com`
- Password: `password123`

---

## Project Structure

```
travel-easy/
├── frontend/                    # React frontend
│   ├── src/
│   │   ├── pages/              # Home, Trips, Bookings, Profile
│   │   ├── components/         # Layout, reusable components
│   │   ├── services/           # API client
│   │   └── types/              # TypeScript types
│   ├── package.json
│   └── vite.config.ts
│
└── backend/                     # Go backend
    ├── cmd/
    │   ├── server/main.go      # Main server
    │   ├── migrate/main.go     # Database migrations
    │   └── seed/main.go        # Seed data
    ├── internal/
    │   ├── handlers/           # API endpoints
    │   ├── middleware/         # Auth middleware
    │   ├── models/             # Database models
    │   ├── services/           # Business logic
    │   └── utils/              # Helpers
    ├── config/
    │   └── database.go         # DB connection
    ├── .env                    # Environment variables
    └── go.mod                  # Go dependencies
```

---

## API Endpoints

### Public Endpoints

```
POST   /api/v1/auth/register         # Register new user
POST   /api/v1/auth/login            # Login user
GET    /api/v1/locations             # List all locations
GET    /api/v1/locations/nearby      # Find nearby locations
GET    /api/v1/locations/:id         # Get location details
GET    /api/v1/lockers/available     # Search available lockers
```

### Protected Endpoints (Require JWT)

```
GET    /api/v1/profile               # Get user profile
POST   /api/v1/bookings              # Create booking
GET    /api/v1/bookings              # Get user bookings
POST   /api/v1/bookings/:id/unlock   # Unlock locker with QR
POST   /api/v1/bookings/:id/checkin  # Check in
POST   /api/v1/bookings/:id/checkout # Check out
POST   /api/v1/bookings/:id/cancel   # Cancel booking
```

---

## Sample Data

The seed script creates:

### Locations (Copenhagen)
1. **Copenhagen Central Station** - Main train station
2. **Fitness World Nørreport** - Modern gym
3. **Copenhagen University Club** - Student club
4. **Nyhavn Tourist Center** - Prime tourist location
5. **Copenhagen Airport Storage** - 24/7 airport facility

### Lockers
- 6 lockers per location
- Mix of SMALL, MEDIUM, and LARGE sizes
- All initially AVAILABLE status

### Test User
- Email: test@example.com
- Password: password123
- 1 active booking created

---

## Common Issues

### Issue: Database Connection Failed

**Error:**
```
Failed to connect to database
```

**Solution:**
```bash
# Check PostgreSQL is running
pg_isready

# Windows
net start postgresql-x64-14

# macOS
brew services start postgresql@14

# Linux
sudo systemctl start postgresql
```

### Issue: Port 8080 Already in Use

**Error:**
```
bind: address already in use
```

**Solution:**
```bash
# Windows
netstat -ano | findstr :8080
taskkill /PID <PID> /F

# macOS/Linux
lsof -ti:8080 | xargs kill -9
```

### Issue: Frontend Can't Connect to Backend

**Check:**
1. Backend is running on http://localhost:8080
2. Frontend API URL is correct in `src/services/api.ts`
3. No CORS errors in browser console

**Solution:**
```bash
# Check backend health
curl http://localhost:8080/health

# Verify ports
netstat -an | findstr "8080 5173"
```

### Issue: Migrations Failed

**Solution:**
```bash
# Drop and recreate database
psql -U postgres
DROP DATABASE IF EXISTS locker_storage;
CREATE DATABASE locker_storage;
\q

# Run migrations again
cd backend
go run cmd/migrate/main.go
go run cmd/seed/main.go
```

---

## Development Workflow

### Making Backend Changes

1. Edit files in `backend/internal/`
2. Restart server: `Ctrl+C` then `go run cmd/server/main.go`

**With hot reload:**
```bash
go install github.com/cosmtrek/air@latest
cd backend
air  # Auto-reloads on file changes
```

### Making Frontend Changes

Vite dev server auto-reloads on file changes. Just save and refresh browser.

### Database Schema Changes

1. Edit `backend/internal/models/models.go`
2. Run migrations: `go run cmd/migrate/main.go`
3. Re-seed if needed: `go run cmd/seed/main.go`

---

## Building for Production

### Backend

```bash
cd backend

# Build binary
go build -o locker-api cmd/server/main.go

# Run binary
./locker-api  # Unix
locker-api.exe  # Windows
```

### Frontend

```bash
# Build for production
npm run build

# Preview production build
npm run preview

# Output in dist/ directory
```

---

## Next Steps

1. **Customize Locations**: Edit seed data in `backend/cmd/seed/main.go`
2. **Update UI**: Modify frontend components in `src/components/` and `src/pages/`
3. **Add Features**: 
   - Payment integration (Stripe)
   - Email notifications
   - Push notifications
   - Admin dashboard
4. **Deploy**: 
   - Backend: Railway, Render, or AWS
   - Frontend: Vercel, Netlify, or Cloudflare Pages
   - Database: Supabase, Neon, or AWS RDS

---

## Resources

- [Go Gin Documentation](https://gin-gonic.com/docs/)
- [GORM Documentation](https://gorm.io/docs/)
- [React Documentation](https://react.dev/)
- [Vite Documentation](https://vitejs.dev/)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)

---

## Support

For issues or questions:
1. Check the troubleshooting section above
2. Review backend README: `backend/README.md`
3. Check project documentation: `IMPLEMENTATION-GUIDE.md`

---

## License

MIT License - See LICENSE file for details
