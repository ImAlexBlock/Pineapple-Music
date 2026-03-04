# Pineapple Music

Self-deployable, browser-based music player with Subsonic protocol support.

## Features

- Web-based music player with responsive UI
- Admin/Guest dual-mode access control
- Music library scanning with metadata extraction
- File upload support (MP3, FLAC, OGG, M4A, WAV, AAC)
- Synced lyrics display (LRC format)
- Playlist management
- Subsonic API v1.16.1 compatible
- Dark/Light theme with Material Design 3
- Internationalization (English / Chinese)
- Docker deployment ready

## Quick Start

```bash
# Docker
docker compose up -d

# Or manual
cd frontend && pnpm install && pnpm build && cd ..
cd backend && go build -o pineapple-music . && ./pineapple-music serve
```

Open http://localhost:3880, initialize, and start listening!

## Tech Stack

- **Backend**: Go + Gin + GORM + SQLite
- **Frontend**: Vue 3 + TypeScript + Vuetify 3
- **API**: REST + Subsonic protocol

## Documentation

- [Deployment Guide](docs/DEPLOYMENT.md)

## License

Apache 2.0
