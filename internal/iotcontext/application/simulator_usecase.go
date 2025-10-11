package application

import (
	"github.com/SeiyaJapon/iot-sensor-app/internal/iotcontext/domain"
	"time"
)

type SimulatorUseCase struct {
	sensorRepository domain.SensorRepository
	simulatorRepo    domain.SimulatorRepository
	eventPublisher   domain.EventPublisher
}

func NewSimulatorUseCase(sensorRepo domain.SensorRepository, simulatorRepo domain.SimulatorRepository) *SimulatorUseCase {
	return &SimulatorUseCase{
		sensorRepository: sensorRepo,
		simulatorRepo:    simulatorRepo,
	}
}

func (uc *SimulatorUseCase) Publish(event domain.IoTEvent) error {
	event.Type = "simulator." + event.Type

	return uc.eventPublisher.Publish(event)
}

func (uc *SimulatorUseCase) ControlSensor(sensorID domain.SensorID, action string) error {
	sensor, err := uc.sensorRepository.FindByID(sensorID)
	if err != nil {
		return err
	}

	if sensor == nil {
		return domain.ErrSensorNotFound
	}

	var eventType string
	switch action {
	case "start":
		err = uc.simulatorRepo.Start(sensorID)
		eventType = "started"
	case "stop":
		err = uc.simulatorRepo.Stop(sensorID)
		eventType = "stopped"
	case "inject_error":
		err = uc.simulatorRepo.InjectError(sensorID)
		eventType = "error_injected"
	default:
		return domain.ErrInvalidAction
	}

	if err != nil {
		return err
	}

	event := domain.IoTEvent{
		Type:      eventType,
		Payload:   map[string]interface{}{"sensor_id": sensorID},
		Timestamp: time.Now().UTC(),
	}

	return uc.Publish(event)
}
