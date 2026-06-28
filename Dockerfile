# Build stage
FROM golang:alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git make

# Copy go mod files
COPY go.mod go.sum* ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/.

# Final stage - minimal runtime image
FROM alpine:latest

WORKDIR /app

# Install runtime dependencies if needed
RUN apk add --no-cache ca-certificates

# Copy binary from builder
COPY --from=builder /app/main .

# Copy templates and assets
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/assets ./assets

# Expose port
EXPOSE 8090

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8090/ || exit 1

# Run the application
CMD ["./main"]