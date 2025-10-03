# ----------------------------------------------------------------------
# 1. BUILD STAGE: Compila la aplicación Go
# ----------------------------------------------------------------------
FROM golang:1.21-alpine AS builder

# Establecer el directorio de trabajo
WORKDIR /app

# Copiar archivos de módulos y descargar dependencias
COPY go.mod go.sum ./
RUN go mod download

# Copiar todo el código fuente
COPY . .

# Compilar la aplicación. Usamos CGO_ENABLED=0 para crear un binario estático
# que puede correr en la imagen 'scratch' (sin dependencias C).
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-s -w' -o /app/sensor-app ./cmd/app/main.go

# ----------------------------------------------------------------------
# 2. RUN STAGE: Imagen final muy ligera
# ----------------------------------------------------------------------
# Usamos 'scratch' (vacía) para la máxima seguridad y el menor tamaño.
FROM scratch

# Exponer el puerto de la aplicación (si expone API HTTP/GRPC)
# Si es puramente NATS, este puerto no es estrictamente necesario, pero es buena práctica.
EXPOSE 8080

# Copiar el binario compilado de la etapa 'builder'
COPY --from=builder /app/sensor-app /sensor-app
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Comando que se ejecuta al iniciar el contenedor
ENTRYPOINT ["/sensor-app"]