package domain

type SensorType string

const (
	Temperature SensorType = "temperature"
	Humidity    SensorType = "humidity"
	Pressure    SensorType = "pressure"
	Generic     SensorType = "generic"
)

type SensorID string
type DeviceID string
