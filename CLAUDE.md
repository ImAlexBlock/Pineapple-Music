# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Pineapple-Music is a self-deployable, browser-based music player with Go+Gin+GORM+SQLite backend and Vue3+TypeScript+Vuetify3 frontend.

## Build & Dev Commands

```bash
# Backend
cd backend && go build ./...       # Build
cd backend && go test ./...        # Test
cd backend && go run . serve       # Run dev server (port 3880)
cd backend && go run . scan        # Scan music directory
cd backend && go run . reset-key admin  # Reset admin key

# Frontend
cd frontend && pnpm install        # Install deps
cd frontend && pnpm dev            # Dev server (proxies to :3880)
cd frontend && pnpm build          # Production build

# Docker
docker compose up -d
```

## Architecture

### Backend (`backend/`)
- `main.go` + `cmd/` - Cobra CLI (serve, scan, reset-key)
- `internal/config/` - Viper config with PM_ env prefix
- `internal/model/` - GORM models + SQLite init (12 tables)
- `internal/handler/api/` - REST API handlers
- `internal/handler/rest/` - Subsonic protocol handlers
- `internal/middleware/` - Auth, CSRF, rate limiting, Turnstile
- `internal/service/` - Business logic layer
- `internal/scanner/` - Music file scanning + metadata extraction
- `internal/util/` - Crypto, response helpers, pagination

### Frontend (`frontend/`)
- `src/api/` - Axios client + API modules
- `src/stores/` - Pinia stores (auth, player, theme)
- `src/plugins/` - Vuetify, i18n, router
- `src/views/` - Page components
- `src/components/player/` - PlayerBar, LyricsPanel

### Key Patterns
- Handler closure: `func ListTracks(db *gorm.DB) gin.HandlerFunc`
- Service layer: business logic separate from HTTP
- Dependency injection via router setup, no globals
- Pure-Go SQLite driver (github.com/glebarez/sqlite)
- Frontend SPA served from `../frontend/dist` by backend

## Dependencies
- Go 1.24+, Node 20+, pnpm
- SQLite via `github.com/glebarez/sqlite` (pure Go, no CGO)
- Audio metadata: `github.com/dhowden/tag`
