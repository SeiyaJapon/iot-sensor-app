package application

import "github.com/SeiyaJapon/iot-sensor-app/internal/iotcontext/domain"

type SimulatorUseCase struct {
	sensorRepository domain.SensorRepository
	simulatorRepo    domain.SimulatorRepository
}

func NewSimulatorUseCase(sensorRepo domain.SensorRepository, simulatorRepo domain.SimulatorRepository) *SimulatorUseCase {
	return &SimulatorUseCase{
		sensorRepository: sensorRepo,
		simulatorRepo:    simulatorRepo,
	}
}

func (uc *SimulatorUseCase) ControlSensor(sensorID domain.SensorID, action string) error {
	sensor, err := uc.sensorRepository.FindByID(sensorID)
	if err != nil {
		return err
	}

	if sensor == nil {
		return domain.ErrSensorNotFound
	}

	switch action {
	case "start":
		return uc.simulatorRepo.Start(sensorID)
	case "stop":
		return uc.simulatorRepo.Stop(sensorID)
	case "inject_error":
		return uc.simulatorRepo.InjectError(sensorID)
	default:
		return domain.ErrInvalidAction
	}
}
