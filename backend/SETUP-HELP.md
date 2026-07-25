# Quick Setup - Backend Server

## Issue: Server Not Starting

The backend requires a PostgreSQL database. You have two options:

---

## Option 1: Use Docker (Recommended - Easiest)

### Step 1: Start Docker Desktop
1. Open **Docker Desktop** application on Windows
2. Wait for it to fully start (whale icon in system tray should be stable)

### Step 2: Start the full stack from the repo root

```bash
cd ..
docker compose up -d --build
```

This starts PostgreSQL, runs migrations + seed automatically (the one-shot `backend_init` service), then starts the backend, notification-service, and frontend.

Expected output includes:
```
Container travel-easy-postgres-1        Healthy
Container travel-easy-backend_init-1    Exited (0)
Container travel-easy-backend-1         Started
```

Server will start on http://localhost:8081 (default)

---

## Option 2: Install PostgreSQL Directly

### Step 1: Download PostgreSQL
- Visit: https://www.postgresql.org/download/windows/
- Download installer for Windows
- Run installer with default settings
- Remember the password you set for postgres user

### Step 2: Create Database
```bash
# Open Command Prompt
"C:\Program Files\PostgreSQL\15\bin\psql.exe" -U postgres

# In psql console:
CREATE DATABASE locker_storage;
\q
```

### Step 3: Update .env File
Edit `backend/.env`:
```
DATABASE_URL=postgres://postgres:YOUR_PASSWORD@localhost:5432/locker_storage?sslmode=disable
```

### Step 4: Run Migrations & Start Server
```bash
cd backend
go run cmd/migrate/main.go
go run cmd/seed/main.go
go run cmd/server/main.go
```

---

## Option 3: Use SQLite (No Install Required)

If you want to skip PostgreSQL setup completely, I can convert the backend to use SQLite which requires no installation.

**Pros:**
- No database installation needed
- Works immediately
- Good for development

**Cons:**
- No geospatial queries (Haversine formula still works in code)
- Not recommended for production

Let me know if you want to switch to SQLite!

---

## Verification

Once database is running, test:

```bash
# Check database connection
cd backend
go run cmd/migrate/main.go
```

Expected output:
```
✅ Database connected successfully
Running database migrations...
✓ Migrations completed successfully
```

---

## Current Status

Based on the error, **Docker Desktop is not running**. 

**Quick Fix:**
1. Open Docker Desktop from Start Menu
2. Wait 30 seconds for it to start
3. Run: `docker compose up -d --build` from the repo root
4. Migrations and seeding run automatically before the server starts

---

## Need Help?

If you're still having issues, let me know which option you prefer:
1. ✅ Use Docker (just need to start Docker Desktop)
2. Install PostgreSQL directly
3. Switch to SQLite (easiest, no install)
