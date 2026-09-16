#!/bin/bash
# Run this ON the EC2 instance after SSH-ing in.
# Installs Docker, Docker Compose, and Git.

set -e

echo "=== Updating system ==="
sudo apt update && sudo apt upgrade -y

echo "=== Installing Docker ==="
sudo apt install -y docker.io docker-compose

echo "=== Starting Docker ==="
sudo systemctl enable docker
sudo systemctl start docker

echo "=== Adding ubuntu user to docker group ==="
sudo usermod -aG docker ubuntu

echo "=== Installing Git ==="
sudo apt install -y git

echo ""
echo "============================================"
echo "  Server setup complete!"
echo "  LOG OUT and SSH back in, then run:"
echo ""
echo "    git clone <your-repo-url> portfolio"
echo "    cd portfolio"
echo "    docker-compose up -d --build"
echo "============================================"
