package app

import (
	"github.com/SeiyaJapon/iot-sensor-app/internal/iotcontext/application"
	"github.com/SeiyaJapon/iot-sensor-app/internal/iotcontext/domain"
	iot_persistence "github.com/SeiyaJapon/iot-sensor-app/internal/iotcontext/infrastructure/persistence"
	"github.com/SeiyaJapon/iot-sensor-app/internal/metricscontext/infrastructure/events"
	"github.com/SeiyaJapon/iot-sensor-app/internal/metricscontext/infrastructure/persistence"
	"github.com/joho/godotenv"
	"log"
	"os"
)

type AppContainer struct {
	DeviceUC          *application.DeviceUseCase
	SensorUC          *application.SensorUseCase
	ReadingsUC        *application.ReadingsUsecase
	SimulatorUC       *application.SimulatorUseCase
	Metrics           *persistence.PrometheusMetricsImpl
	EventPublisher    domain.EventPublisher
	SensorRepo        domain.SensorRepository
	SensorReadingRepo domain.SensorReadingRepository
	DeviceRepo        domain.DeviceRepository
	SimulatorRepo     domain.SimulatorRepository
}

func NewAppContainer() *AppContainer {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	db := iot_persistence.NewDB()

	sensorRepo := iot_persistence.NewPostgresSensorRepository(db)
	sensorReadingRepo := iot_persistence.NewPostgresSensorReadingRepository(db)
	deviceRepo := iot_persistence.NewPostgresDeviceRepository(db)
	simulatorRepo := iot_persistence.NewSimulatorRepository(sensorRepo, sensorReadingRepo, nil)

	natsURL := os.Getenv("NATS_URL")
	if natsURL == "" {
		natsURL = "nats://localhost:4222"
	}

	eventPub, err := events.NewNatsPublisher(&natsURL)
	if err != nil {
		log.Fatalf("Failed to create NATS publisher: %v", err)
	}

	metics := persistence.NewPrometheusMetrics()

	deviceUC := application.NewDeviceUseCase(deviceRepo)
	sensorUC := application.NewSensorUseCase(sensorRepo, metics, eventPub)
	readingsUC := application.NewReadingsUsecase(sensorReadingRepo)
	simulatorUC := application.NewSimulatorUseCase(sensorRepo, simulatorRepo)

	simulatorRepo = iot_persistence.NewSimulatorRepository(sensorRepo, sensorReadingRepo, simulatorUC)

	return &AppContainer{
		DeviceUC:          deviceUC,
		SensorUC:          sensorUC,
		ReadingsUC:        readingsUC,
		SimulatorUC:       simulatorUC,
		Metrics:           metics,
		EventPublisher:    eventPub,
		SensorRepo:        sensorRepo,
		SensorReadingRepo: sensorReadingRepo,
		DeviceRepo:        deviceRepo,
		SimulatorRepo:     simulatorRepo,
	}
}
