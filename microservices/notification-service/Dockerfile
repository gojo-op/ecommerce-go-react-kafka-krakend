FROM golang:1.23-alpine AS builder
WORKDIR /app
RUN apk add --no-cache git
ENV GOPROXY=https://goproxy.io,direct
ENV GOSUMDB=off
COPY microservices/notification-service/ ./
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o notification-service ./cmd/server

FROM alpine:latest
RUN apk --no-cache add ca-certificates wget
WORKDIR /root/
COPY --from=builder /app/notification-service .
EXPOSE 8087
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 CMD wget --no-verbose --tries=1 --spider http://localhost:8087/health || exit 1
CMD ["./notification-service"]