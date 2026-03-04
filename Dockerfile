# Build stage - Frontend
FROM node:20-alpine AS frontend-builder
WORKDIR /app/frontend
RUN corepack enable && corepack prepare pnpm@latest --activate
COPY frontend/package.json frontend/pnpm-lock.yaml ./
RUN pnpm install --frozen-lockfile
COPY frontend/ ./
RUN pnpm build

# Build stage - Backend
FROM golang:1.24-alpine AS backend-builder
WORKDIR /app/backend
COPY backend/go.mod backend/go.sum ./
RUN go mod download
COPY backend/ ./
RUN CGO_ENABLED=0 go build -o /pineapple-music .

# Runtime
FROM alpine:3.19
RUN apk add --no-cache ca-certificates
WORKDIR /app

COPY --from=backend-builder /pineapple-music /app/pineapple-music
COPY --from=frontend-builder /app/frontend/dist /app/frontend/dist

EXPOSE 3880
VOLUME ["/app/data"]

ENV PM_DATA_DIR=/app/data
ENV PM_PORT=3880
ENV GIN_MODE=release

CMD ["/app/pineapple-music", "serve"]
