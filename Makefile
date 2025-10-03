# Nombre de la aplicación y la imagen Docker
APP_NAME := sensor-app
DOCKER_IMAGE := $(APP_NAME):latest

# Comandos de Go
GO_CMD := go
GO_TEST_CMD := $(GO_CMD) test -v ./...
GO_BUILD_CMD := $(GO_CMD) build -o $(APP_NAME) ./cmd/app/main.go

# Variables de Cobertura de Tests
COVERAGE_FILE := coverage.out
COVERAGE_HTML := coverage.html

# ----------------------------------------------------------------------
# COMANDOS PRINCIPALES
# ----------------------------------------------------------------------

.PHONY: build run test clean setup infra

# Construye el binario localmente (para desarrollo rápido)
build:
	@echo "🛠️ Compilando la aplicación Go..."
	$(GO_BUILD_CMD)

# Ejecuta los tests de la aplicación
test:
	@echo "🧪 Ejecutando tests y calculando cobertura..."
	$(GO_TEST_CMD) -race -coverprofile=$(COVERAGE_FILE)
	@echo "✅ Tests completados."
	@echo "Generando reporte de cobertura HTML en $(COVERAGE_HTML)..."
	$(GO_CMD) tool cover -html=$(COVERAGE_FILE) -o $(COVERAGE_HTML)
	@echo "Abre $(COVERAGE_HTML) en tu navegador para ver el reporte."

# Limpia los archivos generados
clean:
	@echo "🗑️ Limpiando binarios y archivos de cobertura..."
	@rm -f $(APP_NAME) $(COVERAGE_FILE) $(COVERAGE_HTML)
	@echo "Limpiando contenedores de infraestructura..."
	@docker-compose -f docker-compose.yml down --remove-orphans
	@echo "Limpieza completada."

# ----------------------------------------------------------------------
# COMANDOS DE DOCKER Y INFRAESTRUCTURA
# ----------------------------------------------------------------------

# Levanta el servidor NATS y la base de datos (PostgreSQL/Redis)
infra:
	@echo "🚀 Levantando infraestructura (NATS, DB) con Docker Compose..."
	@docker-compose -f docker-compose.yml up -d

# Detiene la infraestructura
infra-down:
	@echo "🛑 Deteniendo infraestructura con Docker Compose..."
	@docker-compose -f docker-compose.yml down

# Construye la imagen Docker de la aplicación
docker-build:
	@echo "🏗️ Construyendo imagen Docker: $(DOCKER_IMAGE)"
	@docker build -t $(DOCKER_IMAGE) .

# Ejecuta la aplicación en un contenedor (requiere 'infra' corriendo)
docker-run: docker-build infra
	@echo "🏃 Ejecutando contenedor $(APP_NAME)..."
	@docker run --rm -d --network="host" --name $(APP_NAME)-container $(DOCKER_IMAGE)

# ----------------------------------------------------------------------
# COMANDO PRINCIPAL PARA INICIAR EL PROYECTO
# ----------------------------------------------------------------------

# Inicia todo el stack (infraestructura y aplicación)
start: infra docker-run
	@echo "✨ Proyecto iniciado. NATS y DB listos. Aplicación corriendo."
	@echo "Para ver logs: docker logs -f $(APP_NAME)-container"
	@echo "Para detener todo: make clean"