# Travel Easy - Production Ready Deployment 🚀

## What's Been Set Up

✅ **Production Dockerfile** - Multi-stage build for optimized backend
✅ **Docker Compose** - Full stack with PostgreSQL + Backend
✅ **Environment Config** - Production environment template
✅ **Deployment Scripts** - Automated deployment (Windows & Linux)
✅ **Nginx Config** - Frontend with SSL-ready configuration
✅ **Complete Documentation** - Detailed deployment guide

---

## Quick Start (Local Testing)

### 1. Start Services

**Windows:**
```bash
cd backend
docker-compose up -d --build
```

**The command is building now...**

### 2. Run Migrations (after build completes)

```bash
# Wait for containers to be healthy (30 seconds)
timeout /t 30

# Run migrations
docker-compose exec backend go run cmd/migrate/main.go

# Seed database
docker-compose exec backend go run cmd/seed/main.go
```

### 3. Test API

```bash
# Health check
curl http://localhost:8080/health

# Or visit in browser:
# http://localhost:8080/health
```

### 4. View Logs

```bash
docker-compose logs -f backend
```

---

## Deployment to EC2

### Option 1: Quick Deploy (Automated)

```bash
# On your EC2 instance:
cd backend
chmod +x deploy.sh
./deploy.sh
```

### Option 2: Manual Deploy

```bash
# 1. Upload files to EC2
scp -i your-key.pem -r backend/ ec2-user@your-ec2-ip:/home/ec2-user/

# 2. SSH to EC2
ssh -i your-key.pem ec2-user@your-ec2-ip

# 3. Configure environment
cd backend
cp .env.production .env
nano .env  # Update passwords and secrets

# 4. Deploy
docker-compose up -d --build

# 5. Run migrations
docker-compose exec backend go run cmd/migrate/main.go
docker-compose exec backend go run cmd/seed/main.go

# 6. Test
curl http://localhost:8080/health
```

---

## Production Files Created

### Backend
- ✅ `Dockerfile` - Multi-stage Go build
- ✅ `docker-compose.yml` - PostgreSQL + Backend orchestration
- ✅ `.dockerignore` - Optimize build context
- ✅ `.env.production` - Production environment template
- ✅ `deploy.sh` - Linux deployment script
- ✅ `deploy.bat` - Windows deployment script
- ✅ `DEPLOYMENT.md` - Complete deployment guide

### Frontend
- ✅ `Dockerfile` - React build + Nginx
- ✅ `nginx.conf` - Production-ready Nginx config

---

## Architecture

```
┌─────────────────┐
│   Frontend      │ (Port 80/443)
│   React + Nginx │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│   Backend API   │ (Port 8080)
│   Go + Gin      │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│   PostgreSQL    │ (Port 5432)
│   Database      │
└─────────────────┘
```

---

## Environment Variables

### Required for Production

```bash
# Strong password (20+ characters)
DB_PASSWORD=your-super-strong-password-change-this

# Random 64+ character string
JWT_SECRET=your-random-64-char-jwt-secret-for-production

# Your domain or EC2 IP
CORS_ORIGIN=https://yourdomain.com
```

### All Variables

See `.env.production` for complete list.

---

## Security Checklist

Before deploying to production:

- [ ] Change `DB_PASSWORD` in `.env`
- [ ] Generate strong `JWT_SECRET` (64+ characters)
- [ ] Update `CORS_ORIGIN` with your domain
- [ ] Set `GIN_MODE=release`
- [ ] Configure EC2 Security Group (ports 80, 443, 8080, 22)
- [ ] Set up SSL/HTTPS (Let's Encrypt)
- [ ] Enable automated backups
- [ ] Set up monitoring (CloudWatch)
- [ ] Configure log rotation
- [ ] Review and restrict database access

---

## Useful Commands

### Docker Operations

```bash
# Start services
docker-compose up -d

# Stop services
docker-compose down

# Rebuild
docker-compose up -d --build

# View logs
docker-compose logs -f

# Check status
docker-compose ps

# Restart service
docker-compose restart backend

# Execute command in container
docker-compose exec backend sh
```

### Database Operations

```bash
# Access PostgreSQL
docker-compose exec postgres psql -U postgres -d locker_storage

# Backup database
docker-compose exec postgres pg_dump -U postgres locker_storage > backup.sql

# Restore database
docker-compose exec -T postgres psql -U postgres locker_storage < backup.sql

# Run migrations
docker-compose exec backend go run cmd/migrate/main.go
```

### Monitoring

```bash
# Resource usage
docker stats

# Disk usage
docker system df

# Container details
docker inspect locker-storage-api

# Network info
docker network inspect backend_locker-network
```

---

## Current Build Status

The Docker build is currently in progress. It's:
1. ✅ Pulling PostgreSQL 15 Alpine image
2. 🔄 Building Go backend image (multi-stage build)
3. ⏳ Creating network and volumes

**Estimated time:** 2-5 minutes (depending on internet speed)

---

## After Build Completes

### Verify Services

```bash
# Check running containers
docker ps

# Should see:
# - locker-storage-db (postgres)
# - locker-storage-api (backend)
```

### Run Migrations

```bash
docker-compose exec backend go run cmd/migrate/main.go
```

### Test API

```bash
curl http://localhost:8080/health

# Expected response:
# {"status":"healthy","message":"Locker Storage API is running"}
```

### Test Login

```bash
# First seed the database
docker-compose exec backend go run cmd/seed/main.go

# Then test login
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}'
```

---

## Troubleshooting

### Build Failed

```bash
# Check logs
docker-compose logs backend

# Rebuild from scratch
docker-compose down -v
docker-compose up -d --build
```

### Can't Connect to Database

```bash
# Check if PostgreSQL is healthy
docker-compose ps

# Check database logs
docker-compose logs postgres

# Verify connection string
docker-compose exec backend env | grep DATABASE_URL
```

### Port Already in Use

```bash
# Find process using port
netstat -ano | findstr :8080

# Kill process
taskkill /PID <PID> /F

# Or change port in .env
echo PORT=8081 >> .env
docker-compose up -d
```

---

## Next Steps

1. ⏳ Wait for Docker build to complete
2. ✅ Run migrations
3. ✅ Seed database (optional)
4. ✅ Test API endpoints
5. ✅ Update frontend API URL
6. 🚀 Deploy to EC2 using `DEPLOYMENT.md`

---

## Documentation

- **[DEPLOYMENT.md](DEPLOYMENT.md)** - Complete EC2 deployment guide
- **[README.md](README.md)** - Backend API documentation
- **[SETUP-HELP.md](SETUP-HELP.md)** - Setup troubleshooting

---

## Production Ready ✅

Your application is now containerized and ready for production deployment!

- Multi-stage Docker builds for optimization
- Health checks configured
- Automated restarts on failure
- PostgreSQL with persistent volumes
- Environment-based configuration
- Security best practices
- Complete deployment documentation

**Deploy to EC2:** See [DEPLOYMENT.md](DEPLOYMENT.md) for step-by-step guide.
