# Build stage
FROM golang:1.25.0-alpine AS builder

# Set working directory
WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Copy go mod files first (better Docker layer caching)
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy only necessary source code
COPY cmd/ ./cmd/
COPY internal/ ./internal/
COPY pkg/ ./pkg/

# Generate Swagger docs
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN swag init -g cmd/api/main.go -o ./docs --parseDependency --parseInternal

# Build the application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o api ./cmd/api

# Final stage
FROM alpine:3.19

# Install minimal runtime dependencies
RUN apk --no-cache add ca-certificates tzdata && \
    apk --no-cache add wget

# Create non-root user
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# Set working directory
WORKDIR /app

# Copy only the binary and docs from builder
COPY --from=builder --chown=appuser:appgroup /app/api .
COPY --from=builder --chown=appuser:appgroup /app/docs ./docs

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/api/v1/health || exit 1

# Run the application
CMD ["./api"]