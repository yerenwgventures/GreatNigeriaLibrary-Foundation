# Great Nigeria Library Foundation - Dockerfile
FROM golang:1.21-alpine AS builder

# Install git and ca-certificates
RUN apk add --no-cache git ca-certificates

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy foundation source code
COPY backend/ ./backend/
COPY main.go ./

# Build the foundation application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o foundation-app ./main.go

# Final stage
FROM alpine:latest

# Install ca-certificates
RUN apk --no-cache add ca-certificates

# Create app directory
WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/foundation-app .

# Copy demo content
COPY demo-content/ ./demo-content/

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Run the application
CMD ["./foundation-app"]
