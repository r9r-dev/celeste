# Stage 1: Build frontend
FROM oven/bun:1.3.5-alpine AS frontend-builder

WORKDIR /app/frontend

# Copy package files and lockfile
COPY frontend/package.json frontend/bun.lockb ./

# Install dependencies with cache
RUN --mount=type=cache,target=/root/.bun/install/cache \
    bun install --frozen-lockfile

# Copy frontend source
COPY frontend/ ./

# Build frontend
RUN bun run build

# Stage 2: Build backend
FROM golang:1.25-alpine AS backend-builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY backend/go.mod backend/go.sum ./

# Download dependencies with cache
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

# Copy backend source
COPY backend/ ./

# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /celeste ./cmd/server

# Stage 3: Production image
FROM alpine:3.21

# Install ca-certificates for HTTPS and docker CLI for compose operations
RUN apk add --no-cache ca-certificates docker-cli docker-cli-compose

WORKDIR /app

# Copy binary from builder
COPY --from=backend-builder /celeste /app/celeste

# Copy frontend build
COPY --from=frontend-builder /app/frontend/build /app/static

# Create non-root user
RUN addgroup -g 1000 celeste && \
    adduser -u 1000 -G celeste -s /bin/sh -D celeste

# Environment variables
ENV GIN_MODE=release
ENV PORT=8080
ENV STACKS_PATH=/stacks
ENV STATIC_PATH=/app/static
ENV HOST_PROC=/host/proc
ENV HOST_SYS=/host/sys

# Expose port
EXPOSE 8080

# Run as non-root (note: may need root for docker socket access)
# USER celeste

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget -qO- http://localhost:8080/health || exit 1

# Run the application
CMD ["/app/celeste"]
