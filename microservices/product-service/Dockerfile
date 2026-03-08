FROM golang:1.23-alpine AS builder
WORKDIR /app
RUN apk add --no-cache git
ENV GOPROXY=https://goproxy.io,direct
ENV GOSUMDB=off
COPY microservices/product-service/ ./
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o product-service ./cmd/server

FROM alpine:latest
RUN apk --no-cache add ca-certificates wget
WORKDIR /root/
COPY --from=builder /app/product-service .
EXPOSE 8082
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 CMD wget --no-verbose --tries=1 --spider http://localhost:8082/health || exit 1
CMD ["./product-service"]