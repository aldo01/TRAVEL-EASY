@echo off
REM Exit on error is not available in batch, but we can check errorlevel after each command

echo 🚀 Starting deployment process...

REM Load environment variables from .env.production if it exists
if exist .env.production (
    for /f "tokens=*" %%a in (.env.production) do (
        set %%a
    )
)

echo 📦 Building Docker images...
docker-compose build
if %errorlevel% neq 0 exit /b %errorlevel%

echo 🗄️  Starting services...
docker-compose up -d
if %errorlevel% neq 0 exit /b %errorlevel%

echo ⏳ Waiting for database to be ready...
timeout /t 10 /nobreak

echo 🔄 Running database migrations...
docker-compose exec backend go run cmd/migrate/main.go
if %errorlevel% neq 0 exit /b %errorlevel%

echo 🌱 Seeding database (optional - comment out if not needed)...
docker-compose exec backend go run cmd/seed/main.go

echo ✅ Deployment complete!
echo.
echo Services:
echo   - Backend API: http://localhost:8080
echo   - Health Check: http://localhost:8080/health
echo   - PostgreSQL: localhost:5432
echo.
echo To view logs: docker-compose logs -f
echo To stop: docker-compose down
echo To stop and remove volumes: docker-compose down -v
