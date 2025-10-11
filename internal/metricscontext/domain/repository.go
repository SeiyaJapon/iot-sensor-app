package domain

import "github.com/SeiyaJapon/iot-sensor-app/internal/iotcontext/domain"

type Metrics interface {
	IncSensorReading(sensorType domain.SensorType, id domain.DeviceID)
	IncSensorError(sensorType domain.SensorType, id domain.DeviceID)
}
