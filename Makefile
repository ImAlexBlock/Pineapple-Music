.PHONY: dev build test clean

# Backend
dev:
	cd backend && go run . serve

build-backend:
	cd backend && go build -o ../pineapple-music .

test-backend:
	cd backend && go test ./...

# Frontend
install-frontend:
	cd frontend && pnpm install

dev-frontend:
	cd frontend && pnpm dev

build-frontend:
	cd frontend && pnpm build

test-frontend:
	cd frontend && pnpm test

# All
build: build-frontend build-backend

test: test-backend test-frontend

clean:
	rm -f pineapple-music
	rm -rf frontend/dist
	rm -rf data/
