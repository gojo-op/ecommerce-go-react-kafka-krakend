# Build stage
FROM golang:1.23-alpine AS builder
RUN apk add --no-cache git
ENV GOPROXY=https://goproxy.io,direct
ENV GOSUMDB=off
WORKDIR /app
COPY microservices/auth-service/ ./
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o auth-service ./cmd/server

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

# Create non-root user
RUN addgroup -g 1001 -S appuser && \
    adduser -u 1001 -S appuser -G appuser

# Set working directory
WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/auth-service .

# Change ownership to non-root user
RUN chown appuser:appuser auth-service

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 8081

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8081/health || exit 1

# Run the application
CMD ["./auth-service"]