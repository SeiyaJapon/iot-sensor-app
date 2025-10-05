package http

import (
	"encoding/json"
	"fmt"
	"github.com/SeiyaJapon/iot-sensor-app/internal/iotcontext/application"
	"github.com/SeiyaJapon/iot-sensor-app/internal/iotcontext/domain"
)

type SensorHandlers struct {
	application.SensorUseCase
}

func NewSensorHandlers(sensorUseCase application.SensorUseCase) *SensorHandlers {
	return &SensorHandlers{
		SensorUseCase: sensorUseCase,
	}
}

func (h *SensorHandlers) CreateSensor(id string, deviceID string, name string, typ string, config map[string]interface{}) error {
	jsonBytes, err := json.Marshal(config)
	if err != nil {
		return fmt.Errorf("could not marshal config map to JSON: %w", err)
	}

	var sensorConfig domain.SensorConfig
	if err := json.Unmarshal(jsonBytes, &sensorConfig); err != nil {
		return fmt.Errorf("could not unmarshal JSON to SensorConfig: %w", err)
	}

	sensorConfig.SensorID = domain.SensorID(id)

	if err := h.SensorUseCase.CreateSensor(
		domain.SensorID(id),
		domain.DeviceID(deviceID),
		name,
		domain.SensorType(typ),
		sensorConfig,
	); err != nil {
		return err
	}

	// Lanzar evento de creación de sensor (opcional)

	return nil
}

func (h *SensorHandlers) GetSensorByID(id string) (*domain.Sensor, error) {
	return h.SensorUseCase.GetSensorByID(domain.SensorID(id))
}

func (h *SensorHandlers) GetAllSensors() ([]*domain.Sensor, error) {
	return h.SensorUseCase.GetAllSensors()
}

func (h *SensorHandlers) UpdateSensorConfigById(id string, config map[string]interface{}) error {
	jsonBytes, err := json.Marshal(config)
	if err != nil {
		return fmt.Errorf("could not marshal config map to JSON: %w", err)
	}

	var sensorConfig domain.SensorConfig
	if err := json.Unmarshal(jsonBytes, &sensorConfig); err != nil {
		return fmt.Errorf("could not unmarshal JSON to SensorConfig: %w", err)
	}

	sensorConfig.SensorID = domain.SensorID(id)

	if err := h.SensorUseCase.UpdateSensorConfigById(domain.SensorID(id), sensorConfig); err != nil {
		return err
	}

	// Lanzar evento de actualización de configuración de sensor (opcional)

	return nil
}
