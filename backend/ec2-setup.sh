#!/bin/bash

###############################################################################
# EC2 Instance Setup Script for Travel Easy Backend
# This script automates the complete setup of Docker and the application
# Run this on a fresh EC2 instance (Amazon Linux 2 or Ubuntu)
###############################################################################

set -e  # Exit on error

echo "============================================"
echo "🚀 Travel Easy Backend - EC2 Setup"
echo "============================================"
echo ""

# Detect OS
if [ -f /etc/os-release ]; then
    . /etc/os-release
    OS=$ID
else
    echo "❌ Cannot detect OS"
    exit 1
fi

echo "📋 Detected OS: $OS"
echo ""

# Update system
echo "📦 Updating system packages..."
if [ "$OS" = "amzn" ]; then
    sudo yum update -y
elif [ "$OS" = "ubuntu" ]; then
    sudo apt update && sudo apt upgrade -y
fi

# Install Docker
echo "🐳 Installing Docker..."
if [ "$OS" = "amzn" ]; then
    sudo yum install docker -y
    sudo systemctl start docker
    sudo systemctl enable docker
elif [ "$OS" = "ubuntu" ]; then
    sudo apt install docker.io -y
    sudo systemctl start docker
    sudo systemctl enable docker
fi

# Install Docker Compose
echo "📦 Installing Docker Compose..."
sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

# Add current user to docker group
echo "👤 Adding user to docker group..."
sudo usermod -aG docker $USER

# Install additional tools
echo "🛠️  Installing additional tools..."
if [ "$OS" = "amzn" ]; then
    sudo yum install git htop wget curl -y
elif [ "$OS" = "ubuntu" ]; then
    sudo apt install git htop wget curl -y
fi

# Create application directory
echo "📁 Creating application directory..."
mkdir -p ~/travel-easy
cd ~/travel-easy

# Verify installations
echo ""
echo "✅ Installation Summary:"
echo "   Docker: $(docker --version)"
echo "   Docker Compose: $(docker-compose --version)"
echo "   Git: $(git --version)"
echo ""

# Setup instructions
echo "============================================"
echo "✅ Setup Complete!"
echo "============================================"
echo ""
echo "📝 Next Steps:"
echo ""
echo "1. Upload your application:"
echo "   scp -i your-key.pem -r backend/ ec2-user@$(curl -s http://169.254.169.254/latest/meta-data/public-ipv4):/home/ec2-user/travel-easy/"
echo ""
echo "2. Configure environment:"
echo "   cd ~/travel-easy/backend"
echo "   cp .env.production .env"
echo "   nano .env  # Update passwords and secrets"
echo ""
echo "3. Deploy application:"
echo "   ./deploy.sh"
echo ""
echo "4. Configure security group to allow:"
echo "   - Port 22 (SSH)"
echo "   - Port 80 (HTTP)"
echo "   - Port 443 (HTTPS)"
echo "   - Port 8080 (API)"
echo ""
echo "⚠️  IMPORTANT: You need to log out and log back in for docker group changes to take effect!"
echo "   Run: exit"
echo "   Then reconnect: ssh -i your-key.pem ec2-user@your-instance-ip"
echo ""
echo "Your instance public IP: $(curl -s http://169.254.169.254/latest/meta-data/public-ipv4)"
echo ""
