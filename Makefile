.PHONY: dev dev-frontend dev-backend build clean docker-build docker-up docker-down

# Development
dev: dev-backend dev-frontend

dev-frontend:
	cd frontend && npm run dev -- --port 3001

dev-backend:
	cd backend && go run ./cmd/server

# Build
build: build-frontend build-backend

build-frontend:
	cd frontend && npm run build

build-backend:
	cd backend && go build -o bin/server ./cmd/server

# Docker
docker-build:
	container build -t aperture-science-network-backend ./backend
	container build -t aperture-science-network-frontend ./frontend

docker-up:
	container compose up -d

docker-down:
	container compose down

# Clean
clean:
	rm -rf frontend/build frontend/.svelte-kit
	rm -rf backend/bin

# Install dependencies
install:
	cd frontend && npm install
	cd backend && go mod download

# Format
fmt:
	cd backend && go fmt ./...

# Lint
lint:
	cd frontend && npm run check
	cd backend && go vet ./...
