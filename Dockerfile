# ----------------------------------------------------------------------
# 1. BUILD STAGE: Compila la aplicación Go
# ----------------------------------------------------------------------
FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Copia .env si lo usas (no es estrictamente necesario con Compose)
COPY .env .

# Compila binario estático
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-s -w' -o /app/sensor-app ./cmd/server/main.go

# ----------------------------------------------------------------------
# 2. RUN STAGE: Imagen final
# ----------------------------------------------------------------------
FROM alpine:latest

RUN apk add --no-cache ca-certificates

EXPOSE 8080

COPY --from=builder /app/sensor-app /sensor-app
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENTRYPOINT ["/sensor-app"]
