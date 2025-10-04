package domain

import "time"

type SensorReading struct {
	ID        int64                  `json:"id"`
	SensorID  SensorID               `json:"sensor_id"`
	DeviceID  DeviceID               `json:"device_id"`
	Type      SensorType             `json:"type"`
	Value     float64                `json:"value"`
	Unit      string                 `json:"unit"`
	Timestamp time.Time              `json:"timestamp"`
	Meta      map[string]interface{} `json:"meta"`
}

func NewSensorReading(sensorID SensorID, deviceID DeviceID, typ SensorType, value float64, unit string, ts time.Time) SensorReading {
	return SensorReading{
		SensorID:  sensorID,
		DeviceID:  deviceID,
		Type:      typ,
		Value:     value,
		Unit:      unit,
		Timestamp: ts.UTC(),
		Meta:      map[string]interface{}{},
	}
}
