# ============================================================================
# MULTI-STAGE DOCKERFILE FOR GO APPLICATION
# ============================================================================
# Stage 1: Build the Go binary
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

# Install required packages for CGO (if needed) and git
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy go mod files first for better caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
# CGO_ENABLED=0 creates a static binary
# -ldflags="-w -s" strips debug info for smaller binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -o /app/api \
    ./cmd/api

# -----------------------------------------------------------------------------
# STAGE 2: RUNTIME
# -----------------------------------------------------------------------------
FROM alpine:3.21

# Install ca-certificates for HTTPS and tzdata for timezone support
RUN apk --no-cache add ca-certificates tzdata

# Create non-root user for security
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# Set working directory
WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/api .

# Copy timezone data
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Set timezone (can be overridden by TZ env var)
ENV TZ=Asia/Jakarta

# Change ownership to non-root user
RUN chown -R appuser:appgroup /app

# Switch to non-root user
USER appuser

# Expose application port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Run the application
CMD ["./api"]
