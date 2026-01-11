# ============================================================================
# MULTI-STAGE DOCKERFILE FOR GO APPLICATION
# ============================================================================
# Stage 1: Build both api and cli binaries
# Stage 2: Create minimal runtime image
#
# USAGE:
#   docker build -t go-boilerplate .
#   docker run -p 8080:8080 --env-file .env go-boilerplate
# ============================================================================

# -----------------------------------------------------------------------------
# STAGE 1: BUILD
# -----------------------------------------------------------------------------
FROM golang:1.25-alpine AS builder

# Install required packages
RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /app

# Copy go mod files first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build API binary (HTTP server)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -o /app/build/api \
    ./cmd/api

# Build CLI binary (migrations & seeders)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -o /app/build/cli \
    ./cmd/cli

# -----------------------------------------------------------------------------
# STAGE 2: RUNTIME
# -----------------------------------------------------------------------------
FROM alpine:3.21

RUN apk --no-cache add ca-certificates tzdata

# Create non-root user
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

WORKDIR /app

# Copy binaries from builder
COPY --from=builder /app/build/api .
COPY --from=builder /app/build/cli .

# Copy migrations and seeders for CLI
COPY --from=builder /app/internal/modules/*/migrations ./internal/modules/
COPY --from=builder /app/internal/modules/*/seeders ./internal/modules/

# Copy timezone data
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

ENV TZ=Asia/Jakarta

RUN chown -R appuser:appgroup /app

USER appuser

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Default: run API server
CMD ["./api"]
