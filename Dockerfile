# Great Nigeria Library Foundation - Multi-Stage Dockerfile
# Stage 1: Build stage
FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache \
    git \
    ca-certificates \
    tzdata \
    && update-ca-certificates

# Create non-root user for build
RUN adduser -D -g '' appuser

# Set working directory
WORKDIR /build

# Copy go mod files first for better caching
COPY go.mod go.sum ./

# Download dependencies (cached layer if go.mod/go.sum unchanged)
RUN go mod download && go mod verify

# Copy source code
COPY backend/ ./backend/
COPY main.go ./

# Build the application with optimizations
RUN CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o foundation-app \
    ./main.go

# Verify the binary
RUN ./foundation-app --version || echo "Binary built successfully"

# Stage 2: Final runtime stage
FROM alpine:3.18

# Install runtime dependencies
RUN apk --no-cache add \
    ca-certificates \
    tzdata \
    wget \
    curl \
    && update-ca-certificates

# Create non-root user
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# Create necessary directories
RUN mkdir -p /app/uploads /app/logs /app/demo-content && \
    chown -R appuser:appgroup /app

# Set working directory
WORKDIR /app

# Copy timezone data from builder
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Copy CA certificates from builder
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy the binary from builder stage
COPY --from=builder /build/foundation-app ./foundation-app

# Copy demo content with proper ownership
COPY --chown=appuser:appgroup demo-content/ ./demo-content/

# Make binary executable
RUN chmod +x ./foundation-app

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 8080

# Add labels for better container management
LABEL maintainer="Great Nigeria Library Foundation" \
      version="1.0.0" \
      description="Great Nigeria Library Foundation - Open Source Educational Platform" \
      org.opencontainers.image.source="https://github.com/yerenwgventures/GreatNigeriaLibrary-Foundation"

# Health check with improved reliability
HEALTHCHECK --interval=30s --timeout=10s --start-period=40s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider --timeout=5 http://localhost:8080/health || exit 1

# Set environment variables
ENV GIN_MODE=release \
    PORT=8080 \
    TZ=UTC

# Run the application
CMD ["./foundation-app"]
