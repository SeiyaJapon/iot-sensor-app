# 🌡️ IoT Sensor Management System

Una aplicación **Golang** completa para la gestión de dispositivos IoT con múltiples sensores, implementando arquitectura hexagonal y mensajería NATS.

## 📋 Descripción del Proyecto

Este proyecto implementa un sistema de gestión de sensores IoT que permite:

- **Registro y configuración** de sensores de diferentes tipos (temperatura, humedad, presión)
- **Simulación de lecturas** periódicas con parámetros configurables
- **Mensajería asíncrona** mediante NATS para eventos del sistema
- **Persistencia de datos** en PostgreSQL
- **API REST** para gestión completa del sistema
- **Métricas y monitoreo** con Prometheus

## 🏗️ Arquitectura del Sistema

El proyecto sigue los principios de **Arquitectura Hexagonal (Clean Architecture)** con una clara separación de responsabilidades:

### Estructura de Directorios

```
em3world/
├── cmd/                          # Punto de entrada de la aplicación
│   ├── app/                      # Container de dependencias
│   └── server/                   # Servidor HTTP principal
├── internal/
│   ├── iotcontext/              # Contexto de negocio IoT
│   │   ├── application/         # Casos de uso
│   │   ├── domain/              # Entidades y reglas de negocio
│   │   └── infrastructure/      # Implementaciones concretas
│   │       ├── http/            # Handlers HTTP
│   │       └── persistence/     # Repositorios y DB
│   ├── metricscontext/          # Contexto de métricas
│   └── routes.go                # Configuración de rutas
├── config/                      # Scripts de inicialización
└── docker-compose.yml           # Infraestructura local
```

## 🔧 Tecnologías Utilizadas

- **Go 1.24.5** - Lenguaje principal
- **NATS** - Sistema de mensajería
- **PostgreSQL** - Base de datos
- **GORM** - ORM para Go
- **Prometheus** - Métricas y monitoreo
- **Docker** - Containerización

## 🚀 Instalación y Configuración

### ⚡ Inicio Rápido

```bash
# 1. Clonar el repositorio
git clone <repository-url>
cd em3world

# 2. Iniciar todo el stack (infra + app)
make start-docker

# 3. ¡Listo! La app está en http://localhost:8080
```

### Prerrequisitos

- Go 1.24.5 o superior
- Docker y Docker Compose
- Git

### Configuración Local

1. **Clonar el repositorio:**
```bash
git clone <repository-url>
cd em3world
```

2. **Opción A - Infraestructura + App en Docker (Recomendado):**
```bash
make start-docker
```

3. **Opción B - Infraestructura en Docker + App local:**
```bash
make start-local
# Luego ejecutar: ./sensor-app
# O directamente: make run-local
```

### Variables de Entorno

Crea un archivo `.env` en la raíz del proyecto:

```env
POSTGRES_DSN=host=localhost user=user password=password dbname=iot_db port=55432 sslmode=disable
NATS_URL=nats://localhost:4222
```

## 📊 API Endpoints

### Sensores

| Método | Endpoint | Descripción |
|--------|----------|-------------|
| `GET` | `/sensors` | Obtener todos los sensores |
| `GET` | `/sensors?id={id}` | Obtener sensor por ID |
| `POST` | `/sensors` | Crear nuevo sensor |
| `PUT` | `/sensors?id={id}` | Actualizar configuración |

### Simulador

| Método | Endpoint | Descripción |
|--------|----------|-------------|
| `POST` | `/simulator/?sensor_id={id}&action=start` | Iniciar simulación |
| `POST` | `/simulator/?sensor_id={id}&action=stop` | Detener simulación |
| `POST` | `/simulator/?sensor_id={id}&action=inject_error` | Inyectar error |

### Otros

| Método | Endpoint | Descripción |
|--------|----------|-------------|
| `GET` | `/health` | Health check |
| `GET` | `/metrics` | Métricas Prometheus |

## 🔄 Flujo de Datos

### Diagrama de Arquitectura

```mermaid
graph TD
    subgraph "HTTP Handlers (Infra)"
        A[POST /sensors] --> B[SensorHandlers]
        C[PUT /sensors/{id}/config] --> B
        D[POST /simulator/{id}/control] --> E[SimulatorHandlers]
        F[GET /sensors] --> B
    end

    subgraph "Application (Usecases)"
        B --> G[SensorUseCase]
        E --> H[SimulatorUseCase]
    end

    subgraph "Domain"
        G --> I[Sensor Entity]
        G --> J[EventPublisher Interface]
        H --> K[SimulatorRepository Interface]
        J --> L[SensorsCreatedEvent]
    end

    subgraph "Infrastructure"
        subgraph "Persistence (Postgres)"
            I --> M[SensorRepository Impl]
            K --> N[SimulatorRepo Impl (Ticker + Rand)]
        end
        subgraph "NATS"
            J --> O[NatsPublisher Impl]
            L --> P[Publish to "sensor.created"]
            N --> Q[Publish Reading Events]
        end
    end

    subgraph "External"
        M --> R[Postgres DB]
        O --> S[NATS Server]
    end

    style A fill:#e1f5fe
    style D fill:#f3e5f5
    style R fill:#e8f5e8
    style S fill:#fff3e0
```

## 🧪 Testing

### Ejecutar Tests

```bash
# Ejecutar todos los tests con cobertura automática
make test

# El comando anterior genera automáticamente:
# - coverage.out (archivo de cobertura)
# - coverage.html (reporte visual)
```

### Cobertura de Código

El proyecto incluye tests unitarios para:
- ✅ Casos de uso de sensores
- ✅ Repositorios de persistencia
- ✅ Handlers HTTP
- ✅ Simulador de sensores

## 📈 Monitoreo y Métricas

### Prometheus Metrics

El sistema expone métricas en `/metrics`:

- `sensor_readings_total` - Total de lecturas generadas
- `sensor_errors_total` - Total de errores de sensores
- `active_sensors` - Sensores activos actualmente

### Health Checks

- **Endpoint:** `GET /health`
- **Respuesta:** `200 OK` si el sistema está funcionando

## 🔧 Comandos Make

### 🚀 Comandos Principales

```bash
# Iniciar todo el stack (infra + app en Docker)
make start-docker

# Iniciar infraestructura + app local
make start-local

# Ejecutar app local (asume infra arriba)
make run-local

# Construir binario local
make build

# Ejecutar tests con cobertura
make test
```

### 🏗️ Infraestructura

```bash
# Levantar solo infraestructura (NATS + PostgreSQL)
make infra

# Detener infraestructura
make infra-down

# Construir imagen Docker
make docker-build

# Ejecutar app en contenedor
make docker-run
```

### 🧹 Limpieza

```bash
# Limpiar archivos generados
make clean

# Limpieza total (contenedores, volúmenes, redes)
make nuke

# Reiniciar todo
make restart-all
```

### 🔍 Debug y Desarrollo

```bash
# Entrar en contenedor NATS
make exec-nats

# Entrar en contenedor PostgreSQL
make exec-postgres

# Entrar en contenedor de la app
make exec-app
```

## 📝 Ejemplos de Uso

### Crear un Sensor

```bash
curl -X POST http://localhost:8080/sensors \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Sensor Temperatura",
    "type": "temperature",
    "device_id": "device-123",
    "config": {
      "sampling_rate_ms": 1000,
      "error_rate": 0.1,
      "enabled": true,
      "thresholds": {
        "min": 0,
        "max": 50
      }
    }
  }'
```

### Iniciar Simulación

```bash
curl -X POST "http://localhost:8080/simulator/?sensor_id=sensor-123&action=start"
```

### Obtener Lecturas

```bash
curl -X GET "http://localhost:8080/readings?sensor_id=sensor-123"
```

## 🏛️ Principios de Diseño

### Arquitectura Hexagonal

- **Domain Layer:** Entidades puras sin dependencias externas
- **Application Layer:** Casos de uso y lógica de negocio
- **Infrastructure Layer:** Implementaciones concretas (HTTP, DB, NATS)

### Patrones Implementados

- **Repository Pattern:** Abstracción de persistencia
- **Use Case Pattern:** Encapsulación de lógica de negocio
- **Event-Driven Architecture:** Comunicación asíncrona
- **Dependency Injection:** Inversión de dependencias

## 🔒 Consideraciones de Seguridad

- Validación de entrada en todos los endpoints
- Sanitización de datos de configuración
- Manejo seguro de errores sin exposición de detalles internos
- Rate limiting implícito mediante configuración de sensores

## 🚀 Despliegue

### Docker

```bash
# Opción 1: Todo con Docker Compose (recomendado)
make start-docker

# Opción 2: Solo infraestructura + app en contenedor
make infra
make docker-build
make docker-run

# Opción 3: Infraestructura + app local
make start-local
```

### Variables de Producción

```env
POSTGRES_DSN=postgres://user:password@db:5432/iot_db?sslmode=disable
NATS_URL=nats://nats:4222
```

### Comandos de Despliegue

```bash
# Para desarrollo local
make start-local

# Para producción con Docker
make start-docker

# Para limpiar y reiniciar
make restart-all
```

## 🔧 Troubleshooting

### Problemas Comunes

```bash
# Si la app no arranca, limpiar todo y reiniciar
make nuke
make start-docker

# Si hay problemas con la base de datos
make exec-postgres
# Dentro del contenedor: \dt (ver tablas)

# Si hay problemas con NATS
make exec-nats
# Dentro del contenedor: nats server info

# Ver logs de la aplicación
docker compose logs -f sensor-app

# Ver logs de infraestructura
docker compose logs -f postgres
docker compose logs -f nats
```

### Comandos de Diagnóstico

```bash
# Verificar que todo esté corriendo
docker ps

# Verificar conectividad de la app
curl http://localhost:8080/health

# Ver métricas
curl http://localhost:8080/metrics

# Entrar en la app para debug
make exec-app
```

## 📚 Documentación Adicional

- [Arquitectura Hexagonal](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [NATS Documentation](https://docs.nats.io/)
- [GORM Documentation](https://gorm.io/docs/)

## 🤝 Contribución

1. Fork el proyecto
2. Crea una rama para tu feature (`git checkout -b feature/AmazingFeature`)
3. Commit tus cambios (`git commit -m 'Add some AmazingFeature'`)
4. Push a la rama (`git push origin feature/AmazingFeature`)
5. Abre un Pull Request

## 📄 Licencia

Este proyecto está bajo la Licencia MIT - ver el archivo [LICENSE](LICENSE) para detalles.

## 👨‍💻 Autor

**Seiya Japon** - *Desarrollo completo* - [GitHub](https://github.com/SeiyaJapon)

---

⭐ **¡No olvides darle una estrella al proyecto si te gusta!** ⭐
