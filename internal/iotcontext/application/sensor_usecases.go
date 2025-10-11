package application

import (
	"github.com/SeiyaJapon/iot-sensor-app/internal/iotcontext/domain"
	domain_metrics "github.com/SeiyaJapon/iot-sensor-app/internal/metricscontext/domain"
)

type SensorUseCase struct {
	sensorRepo     domain.SensorRepository
	metrics        domain_metrics.Metrics
	eventPublisher domain.EventPublisher
}

func NewSensorUseCase(sensorRepo domain.SensorRepository, metrics domain_metrics.Metrics, publisher domain.EventPublisher) *SensorUseCase {
	return &SensorUseCase{
		sensorRepo:     sensorRepo,
		metrics:        metrics,
		eventPublisher: publisher,
	}
}

func (uc *SensorUseCase) CreateSensor(
	id domain.SensorID,
	deviceID domain.DeviceID,
	name string,
	typ domain.SensorType,
	config domain.SensorConfig,
) error {
	sensor, err := domain.NewSensor(id, deviceID, name, typ, config)
	if err != nil {
		uc.metrics.IncSensorError(typ, deviceID)
		return err
	}

	if err := uc.sensorRepo.Save(sensor); err != nil {
		uc.metrics.IncSensorError(typ, deviceID)
		return err
	}

	uc.metrics.IncSensorReading(typ, deviceID)

	event := &domain.SensorCreatedEvent{
		SensorID: id,
		DeviceID: deviceID,
		Type:     typ,
		Name:     name,
	}

	return uc.eventPublisher.Publish(event.ToDomainEvent())
}

func (uc *SensorUseCase) GetSensorByID(id domain.SensorID) (*domain.Sensor, error) {
	sensor, err := uc.sensorRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	return sensor, nil
}

func (uc *SensorUseCase) GetAllSensors() ([]*domain.Sensor, error) {
	return uc.sensorRepo.FindAll()
}

func (uc *SensorUseCase) UpdateSensorConfigById(id domain.SensorID, config domain.SensorConfig) error {
	sensor, err := uc.sensorRepo.FindByID(id)
	if err != nil {
		return err
	}

	if err := sensor.UpdateConfig(config); err != nil {
		uc.metrics.IncSensorError(sensor.Type, sensor.DeviceID)
		return err
	}

	if err := uc.sensorRepo.Update(sensor); err != nil {
		uc.metrics.IncSensorError(sensor.Type, sensor.DeviceID)
		return err
	}

	uc.metrics.IncSensorReading(sensor.Type, sensor.DeviceID)

	event := &domain.SensorConfigUpdatedEvent{
		SensorID: id,
		Config:   config,
	}

	return uc.eventPublisher.Publish(event.ToDomainEvent())
}
