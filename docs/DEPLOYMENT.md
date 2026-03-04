# Deployment Guide

## Docker (Recommended)

```bash
# Build and run
docker compose up -d

# Access at http://localhost:3880
```

## Docker without Compose

```bash
docker build -t pineapple-music .
docker run -d \
  -p 3880:3880 \
  -v pineapple-data:/app/data \
  --name pineapple-music \
  pineapple-music
```

## Manual Build

### Prerequisites
- Go 1.24+
- Node.js 20+ with pnpm

### Build

```bash
# Frontend
cd frontend
pnpm install
pnpm build
cd ..

# Backend
cd backend
go build -o pineapple-music .
cd ..
```

### Run

```bash
cd backend
./pineapple-music serve
```

The server starts on port 3880 by default.

## Configuration

Environment variables (prefix `PM_`):

| Variable | Default | Description |
|---|---|---|
| `PM_PORT` | `3880` | Server port |
| `PM_DATA_DIR` | `./data` | Data directory (DB + music files) |
| `PM_LOG_LEVEL` | `info` | Log level |
| `PM_RATE_LIMIT_RPS` | `10` | Rate limit requests per second |
| `PM_RATE_LIMIT_BURST` | `20` | Rate limit burst |
| `PM_SESSION_MAX_AGE` | `86400` | Session TTL in seconds |
| `PM_MAX_UPLOAD_SIZE` | `52428800` | Max upload size in bytes (50MB) |
| `PM_TURNSTILE_SITE_KEY` | | Cloudflare Turnstile site key (optional) |
| `PM_TURNSTILE_SECRET` | | Cloudflare Turnstile secret (optional) |

Or create `data/config.yaml`:

```yaml
port: 3880
data_dir: ./data
rate_limit_rps: 10
```

## Reverse Proxy

### Nginx

```nginx
server {
    listen 80;
    server_name music.example.com;

    client_max_body_size 50M;

    location / {
        proxy_pass http://127.0.0.1:3880;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

### Caddy

```
music.example.com {
    reverse_proxy localhost:3880
}
```

## systemd Service

```ini
[Unit]
Description=Pineapple Music
After=network.target

[Service]
Type=simple
WorkingDirectory=/opt/pineapple-music
ExecStart=/opt/pineapple-music/pineapple-music serve
Environment=PM_DATA_DIR=/opt/pineapple-music/data
Environment=GIN_MODE=release
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

## First Setup

1. Start the server
2. Open `http://localhost:3880` in your browser
3. Click "Initialize" to create the admin key
4. **Save the admin key** - it's only shown once!
5. Log in with the admin key
6. Upload music or put files in `data/music/` and run a scan

## CLI Commands

```bash
# Start server
pineapple-music serve

# Scan music directory
pineapple-music scan

# Reset admin key (local access only)
pineapple-music reset-key admin

# Reset guest key
pineapple-music reset-key guest
```

## Subsonic Client Setup

Compatible with Subsonic API v1.16.1. Use these settings in your client:

- **Server URL**: `http://your-server:3880`
- **Username**: any (not validated)
- **Password**: your admin or guest key
- **Plain password mode** (not token auth)
