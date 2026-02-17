# Travel Easy - Complete Full Stack Application

A modern travel booking platform with React frontend and Node.js backend.

## 🚀 Quick Start Guide

### Prerequisites
- Node.js 18+
- MongoDB (local or Atlas)

### 1. Start Backend API

```bash
# Navigate to backend
cd backend

# Install dependencies (if not done)
npm install

# Start MongoDB (if using local)
mongod

# Start backend server
npm run dev
```

Backend runs on: **http://localhost:5000**

### 2. Start Frontend

```bash
# Open new terminal
# Navigate to frontend root
cd ../

# Install dependencies (if not done)
npm install

# Start frontend
npm run dev
```

Frontend runs on: **http://localhost:5173**

### 3. Seed Database (Optional)

```bash
# In backend directory
npx ts-node src/utils/seed.ts
```

## 📁 Project Structure

```
travel-easy/
├── frontend (React + TypeScript + Tailwind)
│   ├── src/
│   │   ├── components/
│   │   ├── pages/
│   │   ├── services/
│   │   └── types/
│   └── package.json
│
└── backend (Node.js + Express + MongoDB)
    ├── src/
    │   ├── models/
    │   ├── controllers/
    │   ├── routes/
    │   ├── middleware/
    │   └── config/
    └── package.json
```

## 🔧 Configuration

### Backend (.env)
```env
PORT=5000
MONGODB_URI=mongodb://localhost:27017/travel-easy
JWT_SECRET=your-secret-key
CORS_ORIGIN=http://localhost:5173
```

### Frontend (Update src/services/api.ts)
```typescript
const API_BASE_URL = 'http://localhost:5000/api/v1'
```

## 📚 API Endpoints

### Authentication
- `POST /api/v1/auth/register` - Register
- `POST /api/v1/auth/login` - Login
- `GET /api/v1/auth/profile` - Get profile

### Trips (Protected)
- `GET /api/v1/trips` - Get all trips
- `POST /api/v1/trips` - Create trip
- `PUT /api/v1/trips/:id` - Update trip
- `DELETE /api/v1/trips/:id` - Delete trip

### Bookings (Protected)
- `GET /api/v1/bookings` - Get all bookings
- `POST /api/v1/bookings` - Create booking
- `PUT /api/v1/bookings/:id/cancel` - Cancel booking

### Destinations & Activities (Public)
- `GET /api/v1/destinations` - Get destinations
- `GET /api/v1/activities` - Get activities

## 🔐 Authentication Flow

1. Register/Login → Get JWT token
2. Store token in localStorage
3. Include token in API requests:
   ```
   Authorization: Bearer <token>
   ```

## 🧪 Testing the API

### Register a User
```bash
curl -X POST http://localhost:5000/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"name":"John Doe","email":"john@example.com","password":"password123"}'
```

### Create a Trip
```bash
curl -X POST http://localhost:5000/api/v1/trips \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{"destination":"Paris","startDate":"2026-06-15","endDate":"2026-06-22","travelers":2}'
```

## 📱 Features

### Frontend
- ✅ Home with search & featured destinations
- ✅ Trip management
- ✅ Booking tracking
- ✅ User profile
- ✅ Responsive design
- ✅ React Router navigation

### Backend
- ✅ RESTful API
- ✅ JWT authentication
- ✅ MongoDB database
- ✅ User management
- ✅ CRUD operations
- ✅ Security middleware
- ✅ Input validation

## 🛠️ Development

### Backend Hot Reload
```bash
cd backend
npm run dev  # Uses nodemon
```

### Frontend Hot Reload
```bash
npm run dev  # Uses Vite HMR
```

## 🚀 Production Build

### Backend
```bash
cd backend
npm run build
npm start
```

### Frontend
```bash
npm run build
npm run preview
```

## 📖 Documentation

- Frontend README: See main README.md
- Backend README: See backend/README.md
- API Docs: See backend/README.md#API-Endpoints

## 🐛 Troubleshooting

**MongoDB Connection Error:**
- Ensure MongoDB is running
- Check MONGODB_URI in .env

**CORS Error:**
- Verify CORS_ORIGIN matches frontend URL
- Check backend is running

**Authentication Error:**
- Register/login first to get token
- Include token in Authorization header

## 📄 License

Educational purposes only

---

**Stack:** React + TypeScript + Tailwind + Node.js + Express + MongoDB
