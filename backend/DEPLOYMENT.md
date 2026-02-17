# Production Deployment Guide - EC2 Instance

## Prerequisites

- ✅ AWS EC2 instance running Ubuntu 20.04+ or Amazon Linux 2
- ✅ Docker and Docker Compose installed on EC2
- ✅ Security group configured (ports 80, 443, 8080, 5432, 22)
- ✅ Domain name (optional, for HTTPS)

---

## Quick Deploy on EC2

### 1. Connect to EC2 Instance

```bash
ssh -i your-key.pem ec2-user@your-ec2-ip
# or for Ubuntu:
ssh -i your-key.pem ubuntu@your-ec2-ip
```

### 2. Install Docker (if not installed)

```bash
# Update system
sudo yum update -y  # Amazon Linux
# or
sudo apt update && sudo apt upgrade -y  # Ubuntu

# Install Docker
sudo yum install docker -y  # Amazon Linux
# or
sudo apt install docker.io -y  # Ubuntu

# Start Docker
sudo systemctl start docker
sudo systemctl enable docker

# Install Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

# Add user to docker group
sudo usermod -aG docker $USER
newgrp docker

# Verify installation
docker --version
docker-compose --version
```

### 3. Clone/Upload Project

**Option A: From Git**
```bash
git clone https://github.com/yourusername/travel-easy.git
cd travel-easy/backend
```

**Option B: Upload via SCP (from your local machine)**
```bash
# From your local machine (Windows)
scp -i your-key.pem -r c:\Users\91945\workspace\live-presence\frontend\src\pages\travel-easy\backend ec2-user@your-ec2-ip:/home/ec2-user/
```

### 4. Configure Environment

```bash
cd backend

# Create production environment file
cp .env.production .env

# Edit with your values
nano .env
```

**Update these values:**
```bash
DB_PASSWORD=your-super-strong-password-123
JWT_SECRET=your-64-character-random-jwt-secret-key-for-production-use
CORS_ORIGIN=http://your-ec2-ip:80,http://your-domain.com
```

### 5. Deploy

```bash
# Make deploy script executable
chmod +x deploy.sh

# Run deployment
./deploy.sh
```

### 6. Verify Deployment

```bash
# Check running containers
docker ps

# Check logs
docker-compose logs -f backend

# Test health endpoint
curl http://localhost:8080/health

# Test from outside (replace with your EC2 public IP)
curl http://your-ec2-ip:8080/health
```

---

## EC2 Security Group Configuration

Open these ports in your EC2 Security Group:

| Port | Protocol | Source | Purpose |
|------|----------|--------|---------|
| 22 | TCP | Your IP | SSH access |
| 80 | TCP | 0.0.0.0/0 | HTTP (Frontend) |
| 443 | TCP | 0.0.0.0/0 | HTTPS (SSL) |
| 8080 | TCP | 0.0.0.0/0 | Backend API |
| 5432 | TCP | Security Group | PostgreSQL (internal only) |

---

## Database Management

### Backup Database

```bash
# Create backup
docker-compose exec postgres pg_dump -U postgres locker_storage > backup_$(date +%Y%m%d).sql

# Restore backup
docker-compose exec -T postgres psql -U postgres locker_storage < backup_20260203.sql
```

### Access Database

```bash
# Connect to PostgreSQL
docker-compose exec postgres psql -U postgres -d locker_storage

# Common queries
\dt                    # List tables
\d users              # Describe users table
SELECT * FROM users;  # Query users
\q                    # Quit
```

---

## Monitoring & Logs

### View Logs

```bash
# All services
docker-compose logs -f

# Backend only
docker-compose logs -f backend

# PostgreSQL only
docker-compose logs -f postgres

# Last 100 lines
docker-compose logs --tail=100 backend
```

### Container Status

```bash
# Check status
docker-compose ps

# Resource usage
docker stats

# Disk usage
docker system df
```

---

## SSL/HTTPS Setup (Production)

### Option 1: Using Nginx + Let's Encrypt

```bash
# Install Certbot
sudo yum install certbot python3-certbot-nginx -y  # Amazon Linux
# or
sudo apt install certbot python3-certbot-nginx -y  # Ubuntu

# Get SSL certificate
sudo certbot --nginx -d yourdomain.com -d www.yourdomain.com

# Auto-renewal (already set up by certbot)
sudo certbot renew --dry-run
```

### Option 2: Using Traefik (Docker)

Create `docker-compose.prod.yml`:

```yaml
version: '3.8'

services:
  traefik:
    image: traefik:v2.10
    command:
      - "--api.insecure=false"
      - "--providers.docker=true"
      - "--entrypoints.web.address=:80"
      - "--entrypoints.websecure.address=:443"
      - "--certificatesresolvers.myresolver.acme.tlschallenge=true"
      - "--certificatesresolvers.myresolver.acme.email=your@email.com"
      - "--certificatesresolvers.myresolver.acme.storage=/letsencrypt/acme.json"
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - "/var/run/docker.sock:/var/run/docker.sock:ro"
      - "./letsencrypt:/letsencrypt"
    networks:
      - locker-network

  backend:
    labels:
      - "traefik.enable=true"
      - "traefik.http.routers.backend.rule=Host(`api.yourdomain.com`)"
      - "traefik.http.routers.backend.entrypoints=websecure"
      - "traefik.http.routers.backend.tls.certresolver=myresolver"
```

Deploy with:
```bash
docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d
```

---

## Scaling & Performance

### Increase Container Resources

Edit `docker-compose.yml`:

```yaml
services:
  backend:
    deploy:
      resources:
        limits:
          cpus: '2'
          memory: 2G
        reservations:
          cpus: '1'
          memory: 1G
```

### Run Multiple Backend Instances

```bash
docker-compose up -d --scale backend=3
```

Add a load balancer (Nginx) in front.

---

## Maintenance Commands

### Update Application

```bash
# Pull latest code
git pull

# Rebuild and restart
docker-compose up -d --build

# Or use deploy script
./deploy.sh
```

### Clean Up

```bash
# Remove stopped containers
docker-compose down

# Remove with volumes (⚠️ deletes data)
docker-compose down -v

# Clean up old images
docker image prune -a

# Full cleanup
docker system prune -a --volumes
```

### Restart Services

```bash
# Restart all
docker-compose restart

# Restart backend only
docker-compose restart backend

# Stop and start (full restart)
docker-compose down && docker-compose up -d
```

---

## Troubleshooting

### Container Won't Start

```bash
# Check logs
docker-compose logs backend

# Check if port is already in use
sudo netstat -tulpn | grep 8080

# Inspect container
docker inspect locker-storage-api
```

### Database Connection Issues

```bash
# Check if PostgreSQL is running
docker-compose ps postgres

# Check database logs
docker-compose logs postgres

# Test connection from backend
docker-compose exec backend ping postgres

# Verify environment variables
docker-compose exec backend env | grep DATABASE
```

### Out of Disk Space

```bash
# Check disk usage
df -h

# Clean Docker
docker system prune -a --volumes

# Check Docker disk usage
docker system df
```

### Memory Issues

```bash
# Check memory
free -h

# Check container memory usage
docker stats

# Restart with limits
docker-compose down
# Edit docker-compose.yml to add memory limits
docker-compose up -d
```

---

## Automated Backups

Create backup script `/home/ec2-user/backup.sh`:

```bash
#!/bin/bash
BACKUP_DIR="/home/ec2-user/backups"
DATE=$(date +%Y%m%d_%H%M%S)

mkdir -p $BACKUP_DIR

# Backup database
docker-compose exec -T postgres pg_dump -U postgres locker_storage > $BACKUP_DIR/db_$DATE.sql

# Compress
gzip $BACKUP_DIR/db_$DATE.sql

# Keep only last 7 days
find $BACKUP_DIR -name "*.gz" -mtime +7 -delete

echo "Backup completed: db_$DATE.sql.gz"
```

Add to crontab:
```bash
# Edit crontab
crontab -e

# Add daily backup at 2 AM
0 2 * * * /home/ec2-user/backup.sh
```

---

## Cost Optimization

### EC2 Instance Recommendations

| Instance Type | vCPU | RAM | Use Case | Est. Cost/Month |
|---------------|------|-----|----------|-----------------|
| t3.micro | 2 | 1 GB | Development | ~$7.50 |
| t3.small | 2 | 2 GB | Low traffic | ~$15 |
| t3.medium | 2 | 4 GB | Production | ~$30 |
| t3.large | 2 | 8 GB | High traffic | ~$60 |

### Use Spot Instances

Save up to 90% with EC2 Spot Instances (for development):

```bash
# Launch spot instance with AWS CLI
aws ec2 run-instances \
  --instance-type t3.medium \
  --image-id ami-xxxxxxxx \
  --instance-market-options MarketType=spot
```

---

## Monitoring with CloudWatch

```bash
# Install CloudWatch agent
wget https://s3.amazonaws.com/amazoncloudwatch-agent/amazon_linux/amd64/latest/amazon-cloudwatch-agent.rpm
sudo rpm -U ./amazon-cloudwatch-agent.rpm

# Configure metrics
sudo /opt/aws/amazon-cloudwatch-agent/bin/amazon-cloudwatch-agent-ctl \
  -a fetch-config \
  -m ec2 \
  -s \
  -c ssm:CloudWatch-Config
```

---

## Production Checklist

Before going live:

- [ ] Change all default passwords
- [ ] Set strong JWT_SECRET (64+ characters)
- [ ] Configure proper CORS origins
- [ ] Set GIN_MODE=release
- [ ] Enable HTTPS/SSL
- [ ] Set up automated backups
- [ ] Configure monitoring/alerting
- [ ] Set up log rotation
- [ ] Configure firewall/security groups
- [ ] Test disaster recovery
- [ ] Document your setup
- [ ] Set up CI/CD pipeline (optional)

---

## Support & Resources

- **Docker Docs**: https://docs.docker.com/
- **AWS EC2 Docs**: https://docs.aws.amazon.com/ec2/
- **PostgreSQL Docs**: https://www.postgresql.org/docs/
- **Gin Framework**: https://gin-gonic.com/docs/

---

## Quick Reference Commands

```bash
# Deploy
./deploy.sh

# View logs
docker-compose logs -f backend

# Restart
docker-compose restart backend

# Stop all
docker-compose down

# Backup DB
docker-compose exec postgres pg_dump -U postgres locker_storage > backup.sql

# Update app
git pull && docker-compose up -d --build

# Check health
curl http://localhost:8080/health
```

---

**Your application is now production-ready! 🚀**
