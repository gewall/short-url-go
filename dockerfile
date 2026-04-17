# ===== Build stage =====
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install dependency tambahan (kalau go mod butuh git)
RUN apk add --no-cache git

# Copy dependency dulu (biar cache kepakai)
COPY go.mod go.sum ./
RUN go mod download

# Copy semua source
COPY . .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/server

# ===== Run stage =====
FROM alpine:latest

WORKDIR /app

# Sertifikat HTTPS (penting kalau akses API luar / DB TLS)
RUN apk add --no-cache ca-certificates

# Copy binary dari builder
COPY --from=builder /app/main .
COPY internal/geoip/GeoLite2-Country.mmdb internal/geoip/

# Port default (docker-compose pakai 8080)
EXPOSE 8080

# Run app
CMD ["./main"]
