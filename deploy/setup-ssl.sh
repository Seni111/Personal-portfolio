#!/bin/bash
# Run this ON the EC2 instance AFTER docker-compose is running.
# Gets a free SSL certificate from Let's Encrypt via Certbot.

set -e

DOMAIN="adeseni.duckdns.org"
EMAIL="ezekielfila@gmail.com"

echo "=== Getting SSL certificate for $DOMAIN ==="
docker run --rm \
  -v $(pwd)/certbot/conf:/etc/letsencrypt \
  -v $(pwd)/certbot/www:/var/www/certbot \
  certbot/certbot certonly \
  --webroot \
  --webroot-path=/var/www/certbot \
  --email $EMAIL \
  --agree-tos \
  --no-eff-email \
  -d $DOMAIN

echo "=== Certificate obtained! ==="
echo "Now updating Nginx for HTTPS..."

# Replace nginx config with SSL version
cat > nginx/nginx.conf << 'NGINXCONF'
server {
    listen 80;
    server_name adeseni.duckdns.org;

    location /.well-known/acme-challenge/ {
        root /var/www/certbot;
    }

    location / {
        return 301 https://$host$request_uri;
    }
}

server {
    listen 443 ssl;
    server_name adeseni.duckdns.org;

    ssl_certificate /etc/letsencrypt/live/adeseni.duckdns.org/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/adeseni.duckdns.org/privkey.pem;

    location / {
        proxy_pass http://app:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
NGINXCONF

echo "=== Restarting Nginx ==="
docker compose restart nginx

echo ""
echo "============================================"
echo "  HTTPS is live!"
echo "  Visit: https://adeseni.duckdns.org"
echo "============================================"
