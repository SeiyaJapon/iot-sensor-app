package domain

import (
	"testing"
	"time"
)

func TestNewSensorReading(t *testing.T) {
	now := time.Now()
	sensorID := SensorID("sensor-123")
	deviceID := DeviceID("device-123")
	sensorType := Temperature
	value := 25.5
	unit := "°C"

	reading := NewSensorReading(sensorID, deviceID, sensorType, value, unit, now)

	if reading.SensorID != sensorID {
		t.Errorf("expected SensorID %s, got %s", sensorID, reading.SensorID)
	}

	if reading.DeviceID != deviceID {
		t.Errorf("expected DeviceID %s, got %s", deviceID, reading.DeviceID)
	}

	if reading.Type != sensorType {
		t.Errorf("expected Type %s, got %s", sensorType, reading.Type)
	}

	if reading.Value != value {
		t.Errorf("expected Value %f, got %f", value, reading.Value)
	}

	if reading.Unit != string(sensorType) {
		t.Errorf("expected Unit %s, got %s", string(sensorType), reading.Unit)
	}

	if reading.Timestamp != now.UTC() {
		t.Errorf("expected Timestamp %v, got %v", now.UTC(), reading.Timestamp)
	}

	if reading.Meta == nil {
		t.Error("expected Meta to be initialized")
	}
}

func TestNewSensorReadingWithDifferentTypes(t *testing.T) {
	now := time.Now()
	sensorID := SensorID("sensor-123")
	deviceID := DeviceID("device-123")

	tests := []struct {
		name       string
		sensorType SensorType
		value      float64
		unit       string
	}{
		{
			name:       "temperature",
			sensorType: Temperature,
			value:      25.5,
			unit:       "°C",
		},
		{
			name:       "humidity",
			sensorType: Humidity,
			value:      60.0,
			unit:       "%",
		},
		{
			name:       "pressure",
			sensorType: Pressure,
			value:      1013.25,
			unit:       "hPa",
		},
		{
			name:       "generic",
			sensorType: Generic,
			value:      100.0,
			unit:       "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reading := NewSensorReading(sensorID, deviceID, tt.sensorType, tt.value, tt.unit, now)

			if reading.Type != tt.sensorType {
				t.Errorf("expected Type %s, got %s", tt.sensorType, reading.Type)
			}

			if reading.Value != tt.value {
				t.Errorf("expected Value %f, got %f", tt.value, reading.Value)
			}

			if reading.Unit != string(tt.sensorType) {
				t.Errorf("expected Unit %s, got %s", string(tt.sensorType), reading.Unit)
			}
		})
	}
}
