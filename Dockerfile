# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git make

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/agent-sandbox ./cmd/server/

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS and basic shell utilities
RUN apk --no-cache add ca-certificates bash curl

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/bin/agent-sandbox .

# Copy config file
COPY configs/config.yaml ./configs/

# Create workspace directory
RUN mkdir -p /tmp/agent-sandbox

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD curl -f http://localhost:8080/health || exit 1

# Run the application
CMD ["./agent-sandbox"]
