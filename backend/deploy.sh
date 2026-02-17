#!/bin/bash

# Exit on error
set -e

echo "🚀 Starting deployment process..."

# Load environment variables
if [ -f .env.production ]; then
    export $(cat .env.production | grep -v '^#' | xargs)
fi

echo "📦 Building Docker images..."
docker-compose build

echo "🗄️  Starting services..."
docker-compose up -d

echo "⏳ Waiting for database to be ready..."
sleep 10

echo "🔄 Running database migrations..."
docker-compose exec backend go run cmd/migrate/main.go

echo "🌱 Seeding database (optional - comment out if not needed)..."
docker-compose exec backend go run cmd/seed/main.go

echo "✅ Deployment complete!"
echo ""
echo "Services:"
echo "  - Backend API: http://localhost:8080"
echo "  - Health Check: http://localhost:8080/health"
echo "  - PostgreSQL: localhost:5432"
echo ""
echo "To view logs: docker-compose logs -f"
echo "To stop: docker-compose down"
echo "To stop and remove volumes: docker-compose down -v"
