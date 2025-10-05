package application

import "github.com/SeiyaJapon/iot-sensor-app/internal/iotcontext/domain"

type SensorUseCase struct {
	sensorRepo domain.SensorRepository
}

func NewSensorUseCase(sensorRepo domain.SensorRepository) *SensorUseCase {
	return &SensorUseCase{
		sensorRepo: sensorRepo,
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
		return err
	}

	if err := uc.sensorRepo.Save(sensor); err != nil {
		return err
	}
	return nil
}

func (uc *SensorUseCase) GetSensorByID(id domain.SensorID) (*domain.Sensor, error) {
	return uc.sensorRepo.FindByID(id)
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
		return err
	}

	if err := uc.sensorRepo.Update(sensor); err != nil {
		return err
	}

	return nil
}
