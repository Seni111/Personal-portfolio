# Personal Portfolio Website — Project Report

## 1. Project Overview

**Project:** A personal portfolio website for Adeseni Adefila, built from scratch as a full-stack web application and deployed to AWS cloud infrastructure.

**Objective:** To design, develop, containerize, and deploy a production-grade portfolio website that demonstrates cloud engineering, DevOps, and software development skills. The project serves dual purposes: (1) a live portfolio to showcase projects and skills to recruiters and collaborators, and (2) a hands-on learning exercise covering the complete lifecycle of a cloud-deployed application — from writing application code to configuring DNS and SSL certificates.

**Live URL:** https://adeseni.duckdns.org

---

## 2. Architecture

The application follows a three-tier architecture running entirely in Docker containers on a single AWS EC2 instance:

```
Internet
   │
   ▼
┌─────────────────────────────────────────────────────┐
│  AWS EC2 (t3.micro, Ubuntu 26.04)                   │
│  Elastic IP: 100.24.188.210                         │
│                                                     │
│  ┌───────────────────────────────────────────────┐  │
│  │  Docker Compose                               │  │
│  │                                               │  │
│  │  ┌─────────┐   ┌──────────┐   ┌───────────┐  │  │
│  │  │  Nginx  │──▶│  Go/Gin  │──▶│PostgreSQL │  │  │
│  │  │ :80/:443│   │  :8080   │   │  :5432    │  │  │
│  │  └─────────┘   └──────────┘   └───────────┘  │  │
│  │   (public)      (internal)     (internal)     │  │
│  └───────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────┘
   ▲
   │
DNS: adeseni.duckdns.org → 100.24.188.210
SSL: Let's Encrypt (Certbot)
```

**Request flow:**

1. A user visits `https://adeseni.duckdns.org`
2. DNS (DuckDNS) resolves the domain to the EC2 Elastic IP `100.24.188.210`
3. Nginx receives the request on port 443 (HTTPS), terminates SSL, and forwards the request to the Go application on port 8080 via Docker's internal network
4. The Go/Gin application processes the request, queries PostgreSQL if needed, renders an HTML template, and returns the response
5. Nginx sends the response back to the user over the encrypted HTTPS connection

HTTP requests on port 80 are automatically redirected to HTTPS via a 301 redirect in the Nginx configuration.

---

## 3. Technology Stack

| Layer | Technology | Role |
|-------|-----------|------|
| **Language** | Go 1.25+ | Application language |
| **Web Framework** | Gin | HTTP routing, middleware, template rendering |
| **Database** | PostgreSQL 16 | Persistent storage for blog posts, projects, and contact messages |
| **Database Driver** | lib/pq | Go PostgreSQL driver for `database/sql` |
| **Templating** | Go `html/template` | Server-side HTML rendering with layout inheritance |
| **Styling** | CSS3 | Dark-themed responsive design with Inter + JetBrains Mono fonts |
| **Containerization** | Docker | Multi-stage build for the Go binary; Alpine-based production image |
| **Orchestration** | Docker Compose | Multi-container management (Nginx + App + PostgreSQL) |
| **Reverse Proxy** | Nginx (Alpine) | SSL termination, HTTP→HTTPS redirect, request forwarding |
| **Cloud Provider** | AWS EC2 (t3.micro) | Compute instance hosting all containers |
| **Operating System** | Ubuntu 26.04 LTS | EC2 instance OS |
| **Static IP** | AWS Elastic IP | Persistent public IP that survives instance stop/start |
| **Storage** | AWS EBS (16 GB) | Block storage, expanded from default 8 GB |
| **DNS** | DuckDNS | Free dynamic DNS mapping domain to Elastic IP |
| **SSL/TLS** | Let's Encrypt (Certbot) | Free SSL certificate for HTTPS |
| **Version Control** | Git + GitHub | Source code management and deployment pipeline |

---

## 4. Project Structure

```
portfolio/
├── main.go                    # Entry point: custom template renderer, route definitions, server startup
├── go.mod                     # Go module definition and dependency versions
├── go.sum                     # Dependency checksum lockfile
├── Dockerfile                 # Multi-stage Docker build (Go compile → Alpine runtime)
├── docker-compose.yml         # Three-service orchestration: Nginx, App, PostgreSQL
├── .gitignore                 # Excludes binary, certs, env files, pem keys
├── .dockerignore              # Excludes git history, markdown, compose file from Docker build
├── CLAUDE.md                  # Project documentation for AI-assisted development
│
├── database/
│   └── db.go                  # PostgreSQL connection, env-based config, auto table creation
│
├── models/
│   └── models.go              # Data structs: BlogPost, Project, ContactMessage
│
├── handlers/
│   ├── home.go                # Home page handler (static content)
│   ├── projects.go            # Projects list handler (DB query with COALESCE for NULLs)
│   ├── blog.go                # Blog list + single post handlers
│   └── contact.go             # Contact form display + submission handler
│
├── templates/
│   ├── layout.html            # Shared HTML shell: nav, footer, fonts, scripts, scroll-reveal
│   ├── home.html              # Hero section, about, stats, skills grid, experience
│   ├── projects.html          # Project cards from database
│   ├── blog.html              # Blog post list from database
│   ├── post.html              # Single blog post view
│   └── contact.html           # Contact form with GitHub/LinkedIn/Email links
│
├── static/
│   └── style.css              # Full dark-theme stylesheet: CSS variables, responsive, animations
│
├── nginx/
│   └── nginx.conf             # Nginx reverse proxy config (HTTP→HTTPS redirect + SSL + proxy_pass)
│
└── deploy/
    ├── setup-server.sh        # EC2 bootstrap: installs Docker, Docker Compose, Git
    └── setup-ssl.sh           # SSL setup: runs Certbot, updates Nginx config, restarts Nginx
```

---

## 5. Infrastructure Setup

### EC2 Instance
- **Instance type:** t3.micro (1 vCPU, 1 GB RAM) — AWS free-tier eligible
- **AMI:** Ubuntu 26.04 LTS
- **Elastic IP:** 100.24.188.210 — allocated and associated to ensure a static public IP that persists across instance stop/start cycles

### Security Groups (Inbound Rules)
| Port | Protocol | Source | Purpose |
|------|----------|--------|---------|
| 22 | TCP | Restricted to owner IP | SSH access |
| 80 | TCP | 0.0.0.0/0 | HTTP (redirects to HTTPS) |
| 443 | TCP | 0.0.0.0/0 | HTTPS (serves the site) |

Port 8080 is **not** exposed publicly — the Go application is only reachable via Nginx's internal Docker network.

### DNS (DuckDNS)
- **Domain:** adeseni.duckdns.org
- **Provider:** DuckDNS (free dynamic DNS service)
- **Configuration:** A-record pointing to 100.24.188.210, configured via the DuckDNS web dashboard

### Storage (EBS Volume)
- **Original size:** 8 GB (default)
- **Expanded to:** 16 GB via AWS Console (Modify Volume), then extended on the OS with:
  ```bash
  sudo growpart /dev/nvme0n1 1
  sudo resize2fs /dev/nvme0n1p1
  ```

### Swap Space
- **Size:** 1 GB (`/swapfile`)
- **Purpose:** Provides overflow memory for Go compilation inside Docker, which exceeds the 1 GB physical RAM of t3.micro instances
- **Persistence:** Added to `/etc/fstab` so it activates automatically on reboot

---

## 6. Deployment Pipeline

The deployment follows a manual pull-based workflow:

```
Local Development          GitHub             EC2 Production
─────────────────          ──────             ──────────────
  Edit code
       │
  git add + commit
       │
  git push origin main ──▶ Repository
                              │
                     SSH into EC2, then:
                              │
                    git pull origin main ◀──┘
                              │
                    docker compose up -d --build
                              │
                    ┌─────────┴─────────┐
                    │  Docker builds:    │
                    │  1. Go compile     │
                    │  2. Alpine image   │
                    │  3. Pull Postgres  │
                    │  4. Pull Nginx     │
                    └─────────┬─────────┘
                              │
                    Containers start:
                    db → app → nginx
                              │
                    Site live at
                    https://adeseni.duckdns.org
```

**Initial server setup** (one-time):
```bash
bash deploy/setup-server.sh   # Installs Docker, Docker Compose, Git
```

**Deploying updates:**
```bash
cd portfolio
git pull origin main
docker compose up -d --build
```

Docker caching ensures that subsequent builds are fast — only changed layers (typically just the application code) are rebuilt. The dependency download layer (`go mod download`) is cached.

---

## 7. Key Technical Challenges & Solutions

### Challenge 1: Go Template Collision Bug
**Problem:** `r.LoadHTMLGlob("templates/*")` loads all template files into a single Go template set. Since every page template defined a `{{ define "content" }}` block, Go's template engine only kept the last one parsed. All pages rendered the same content (whichever file was parsed last — `projects.html`).

**Solution:** Replaced `LoadHTMLGlob` with a custom `templateRenderer` that parses `layout.html` + each page file separately using `template.ParseFiles()`. Each page gets its own isolated template set, so `"content"` definitions don't collide. The renderer implements Gin's `render.HTMLRender` interface, executing the `"layout"` named template for each page.

### Challenge 2: NULL Scan in PostgreSQL Query
**Problem:** The `projects` table has nullable `url` and `image_url` columns. When `rows.Scan()` encounters a SQL NULL and tries to put it into a Go `string`, it returns an error. The handler's error-handling code (`continue`) silently skipped every row, resulting in an empty projects page despite data being present in the database.

**Solution:** Used `COALESCE(url, '')` and `COALESCE(image_url, '')` in the SQL query to convert NULL values to empty strings before Go sees them. Go templates handle empty strings correctly (`{{ if .URL }}` treats `""` as falsy).

### Challenge 3: OOM Killer During Docker Build
**Problem:** Go compilation is memory-intensive — it loads the entire dependency tree into memory. On a t3.micro instance with only 1 GB RAM, the Linux kernel's OOM (Out of Memory) killer terminated the Go compiler mid-build with `signal: killed`.

**Evidence:** The error message `compile: signal: killed` after 355 seconds indicated the kernel sent SIGKILL (signal 9), which only the OOM killer does in this context.

**Solution:** Created a swap file to provide overflow memory:
```bash
sudo fallocate -l 1G /swapfile
sudo chmod 600 /swapfile
sudo mkswap /swapfile
sudo swapon /swapfile
```
Made permanent via `/etc/fstab`. Build times increased (swap is disk-backed, slower than RAM) but compilation completes successfully.

### Challenge 4: Disk Space Exhaustion
**Problem:** The default EC2 EBS volume is 8 GB (6.7 GB usable). Docker build cache from failed builds + the 2 GB swap file + OS consumed all available space, causing the next build to fail with `no space left on device`.

**Solution:** Two-part fix:
1. **Immediate:** Cleaned Docker build cache (`docker system prune -af`) and reduced swap from 2 GB to 1 GB, freeing ~3-4 GB
2. **Permanent:** Expanded the EBS volume from 8 GB to 16 GB via the AWS Console, then grew the partition and filesystem on the OS

### Challenge 5: Port Conflict in Local Development
**Problem:** Port 8080 was already in use by another Docker container (`capstone_api_1`) during local testing, preventing the portfolio containers from starting.

**Solution:** Mapped the app to port 8081 on the host in `docker-compose.yml` (`8081:8080`). In production, this is irrelevant because Nginx handles external traffic on ports 80/443 and forwards internally to the app on 8080 via Docker networking.

### Challenge 6: `docker-compose` vs `docker compose`
**Problem:** The SSL setup script used the hyphenated `docker-compose` command, but the EC2 instance only had Docker's built-in `docker compose` plugin (no hyphen). The script failed with `command not found`.

**Solution:** Updated the script to use `docker compose` (the modern Docker CLI plugin syntax).

---

## 8. SSL/HTTPS Setup

### Certificate Acquisition
SSL is provided by **Let's Encrypt**, a free Certificate Authority, via the **Certbot** tool running as a Docker container.

**Process (automated in `deploy/setup-ssl.sh`):**
1. Certbot runs in webroot mode, placing a challenge file in `/var/www/certbot`
2. Let's Encrypt's servers verify domain ownership by requesting `http://adeseni.duckdns.org/.well-known/acme-challenge/<token>`
3. Nginx serves this challenge file (configured in the `location /.well-known/acme-challenge/` block)
4. Upon verification, Certbot receives and stores the certificate and private key

**Certificate files:**
- Certificate: `/etc/letsencrypt/live/adeseni.duckdns.org/fullchain.pem`
- Private key: `/etc/letsencrypt/live/adeseni.duckdns.org/privkey.pem`

### Nginx SSL Configuration
After certificate acquisition, the script rewrites `nginx.conf` to:

**Port 80 (HTTP):**
- Serves Let's Encrypt challenge files for certificate renewal
- Redirects all other traffic to HTTPS with a 301 redirect

**Port 443 (HTTPS):**
- Terminates SSL using the Let's Encrypt certificate
- Forwards decrypted traffic to the Go application at `http://app:8080`
- Passes client headers (`X-Real-IP`, `X-Forwarded-For`, `X-Forwarded-Proto`) so the application knows the real client IP and protocol

### Certificate Renewal
- **Expiry date:** December 15, 2026
- **Renewal:** Re-run `bash deploy/setup-ssl.sh` before expiry

---

## 9. Key Details

| Item | Value |
|------|-------|
| **Live URL** | https://adeseni.duckdns.org |
| **Domain** | adeseni.duckdns.org (DuckDNS, free) |
| **EC2 Elastic IP** | 100.24.188.210 |
| **Instance Type** | t3.micro (1 vCPU, 1 GB RAM) |
| **Operating System** | Ubuntu 26.04 LTS |
| **EBS Volume** | 16 GB (expanded from 8 GB) |
| **Swap** | 1 GB, persistent via /etc/fstab |
| **SSL Certificate** | Let's Encrypt, expires December 15, 2026 |
| **GitHub Repository** | github.com/Seni111/Personal-portfolio |
| **SSH Key** | portfolio.pem (stored at ~/.ssh/portfolio.pem) |
| **Database** | PostgreSQL 16 (Alpine), Docker volume `pgdata` |
| **Docker Images** | golang:1.26-alpine (build), alpine:3.20 (runtime), postgres:16-alpine, nginx:alpine |

---

*Project completed September 16, 2026 by Adeseni Adefila.*
