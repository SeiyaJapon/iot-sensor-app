# Nombre de la aplicación y la imagen Docker
APP_NAME := sensor-app
DOCKER_IMAGE := $(APP_NAME):latest

# Comandos de Go
GO_CMD := go
GO_TEST_CMD := $(GO_CMD) test -v ./...
GO_BUILD_CMD := $(GO_CMD) build -o $(APP_NAME) ./cmd/server/main.go

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
	@docker compose -f docker-compose.yml down --remove-orphans --volumes  # Limpia DB si quieres reset total
	@echo "Limpieza completada."

# Limpieza total: contenedores, volúmenes, redes, procesos en puertos críticos
nuke:
	@echo "💣 Eliminando todo: contenedores, volúmenes, redes, procesos en puertos críticos..."
	@echo "🛑 Deteniendo contenedores..."
	@docker compose -f docker-compose.yml down --remove-orphans --volumes
	@echo "🧹 Eliminando contenedores zombie..."
	@docker rm -f $$(docker ps -aq) 2>/dev/null || true
	@echo "🧹 Eliminando imágenes dangling..."
	@docker image prune -f
	@echo "🧹 Eliminando redes no usadas..."
	@docker network prune -f
	@echo "🔪 Matando procesos en puerto 8080..."
	@kill -9 $$(lsof -ti tcp:8080) 2>/dev/null || echo "⚠️ Nada ocupando 8080"
	@echo "✅ Entorno limpio. Puedes ejecutar: make start-docker o make run-local"

# ----------------------------------------------------------------------
# COMANDOS DE DOCKER Y INFRAESTRUCTURA
# ----------------------------------------------------------------------

# Levanta SOLO infra (NATS + DB) para correr app en local
infra:
	@echo "🚀 Levantando infraestructura (NATS, DB) con Docker Compose..."
	@docker compose -f docker-compose.yml up -d nats postgres  # Solo infra, no app

# Detiene la infraestructura
infra-down:
	@echo "🛑 Deteniendo infraestructura con Docker Compose..."
	@docker compose -f docker-compose.yml down

# Construye la imagen Docker de la aplicación (opcional, si usas run suelto)
docker-build:
	@echo "🏗️ Construyendo imagen Docker: $(DOCKER_IMAGE)"
	@docker build -t $(DOCKER_IMAGE) .

# Ejecuta la app SOLA en contenedor (con host net, para debug rápido)
docker-run: docker-build infra
	@echo "🏃 Ejecutando contenedor $(APP_NAME)..."
	@docker rm -f $(APP_NAME)-container 2>/dev/null || true  # Mata zombie si existe
	@docker run --rm -d --network=host --name $(APP_NAME)-container \
		-e POSTGRES_DSN='host=127.0.0.1 user=user password=password dbname=iot_db port=5432 sslmode=disable' \
		-e NATS_URL='nats://localhost:4222' \
		$(DOCKER_IMAGE)

# ----------------------------------------------------------------------
# NUEVOS: ENTRAR EN CONTENEDORES Y RUN LOCAL
# ----------------------------------------------------------------------

# Entra en el contenedor de NATS (bash interactivo)
exec-nats:
	@echo "🔓 Entrando en NATS como un ninja..."
	@docker exec -it nats-server sh

# Entra en el contenedor de Postgres (psql interactivo)
exec-postgres:
	@echo "🔓 Entrando en Postgres para charlar con la DB..."
	@docker exec -it iot-postgres psql -U user -d iot_db

# Entra en el contenedor de la app (si está corriendo via Compose)
exec-app:
	@echo "🔓 Entrando en la app para debuggear el alma..."
	@docker exec -it sensor-app sh

# Entra en el contenedor de la app (si está corriendo via docker-run)
exec-docker-app:
	@echo "🔓 Entrando en la app docker-run para husmear..."
	@docker exec -it $(APP_NAME)-container sh

# Lanza la app en local con go run (¡usa tu .env y asume infra arriba!)
run-local:
	@echo "⚡ Lanzando app en local: go run ./cmd/server/main.go"
	@$(MAKE) infra  # ¡Fix! Invoca el target Make correctamente
	@$(GO_CMD) run ./cmd/server/main.go

# ----------------------------------------------------------------------
# COMANDO PRINCIPAL PARA INICIAR EL PROYECTO
# ----------------------------------------------------------------------

# Opción 1: Full stack en Docker (fácil y conectado)
start-docker:
	@echo "🚀 Iniciando TODO con Docker Compose (app + infra)..."
	@docker compose -f docker-compose.yml up -d --build
	@echo "✨ Stack completo arriba. Logs app: docker compose logs -f sensor-app"
	@echo "Para parar: docker compose down"

# Opción 2: Infra en Docker + app en local (tu flow preferido)
start-local: infra build
	@echo "🚀 Infra arriba, app compilada. Corre: ./$(APP_NAME) (usa tu .env local)"
	@echo "O para run directo: make run-local"
	@echo "Logs DB: docker logs -f iot-postgres | Para NATS: docker logs -f nats-server"

# Alias corto para full local flow
start: start-local

# Reinicia todo (limpia, reconstruye, relanza)
restart-all: clean nuke docker-run start-docker