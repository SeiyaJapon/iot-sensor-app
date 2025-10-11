package domain

import (
	"testing"
)

func TestNewSensor(t *testing.T) {
	tests := []struct {
		name        string
		id          SensorID
		deviceID    DeviceID
		sensorName  string
		sensorType  SensorType
		config      SensorConfig
		expectError bool
	}{
		{
			name:        "valid sensor",
			id:          "sensor-123",
			deviceID:    "device-123",
			sensorName:  "Temperature Sensor",
			sensorType:  Temperature,
			config:      SensorConfig{},
			expectError: false,
		},
		{
			name:        "empty sensor id",
			id:          "",
			deviceID:    "device-123",
			sensorName:  "Temperature Sensor",
			sensorType:  Temperature,
			config:      SensorConfig{},
			expectError: true,
		},
		{
			name:        "empty device id",
			id:          "sensor-123",
			deviceID:    "",
			sensorName:  "Temperature Sensor",
			sensorType:  Temperature,
			config:      SensorConfig{},
			expectError: true,
		},
		{
			name:        "empty name",
			id:          "sensor-123",
			deviceID:    "device-123",
			sensorName:  "",
			sensorType:  Temperature,
			config:      SensorConfig{},
			expectError: true,
		},
		{
			name:        "empty type defaults to generic",
			id:          "sensor-123",
			deviceID:    "device-123",
			sensorName:  "Temperature Sensor",
			sensorType:  "",
			config:      SensorConfig{},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sensor, err := NewSensor(tt.id, tt.deviceID, tt.sensorName, tt.sensorType, tt.config)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if sensor.ID != tt.id {
				t.Errorf("expected ID %s, got %s", tt.id, sensor.ID)
			}

			if sensor.DeviceID != tt.deviceID {
				t.Errorf("expected DeviceID %s, got %s", tt.deviceID, sensor.DeviceID)
			}

			if sensor.Name != tt.sensorName {
				t.Errorf("expected Name %s, got %s", tt.sensorName, sensor.Name)
			}

			expectedType := tt.sensorType
			if tt.sensorType == "" {
				expectedType = Generic
			}
			if sensor.Type != expectedType {
				t.Errorf("expected Type %s, got %s", expectedType, sensor.Type)
			}

			if sensor.CreatedAt.IsZero() {
				t.Error("expected CreatedAt to be set")
			}

			if sensor.UpdatedAt.IsZero() {
				t.Error("expected UpdatedAt to be set")
			}
		})
	}
}

func TestSensor_UpdateConfig(t *testing.T) {
	sensor, _ := NewSensor("sensor-123", "device-123", "Test Sensor", Temperature, SensorConfig{})

	validConfig := SensorConfig{
		SensorID:       "sensor-123",
		SamplingRateMs: 1000,
		ErrorRate:      0.1,
		Enabled:        true,
	}

	err := sensor.UpdateConfig(validConfig)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if sensor.Config.SamplingRateMs != 1000 {
		t.Errorf("expected SamplingRateMs 1000, got %d", sensor.Config.SamplingRateMs)
	}

	if sensor.Config.ErrorRate != 0.1 {
		t.Errorf("expected ErrorRate 0.1, got %f", sensor.Config.ErrorRate)
	}

	if !sensor.Config.Enabled {
		t.Error("expected Enabled to be true")
	}

	invalidConfig := SensorConfig{
		SensorID:       "different-sensor",
		SamplingRateMs: 1000,
		ErrorRate:      0.1,
		Enabled:        true,
	}

	err = sensor.UpdateConfig(invalidConfig)
	if err == nil {
		t.Error("expected error for mismatched sensor ID")
	}

	invalidSamplingRate := SensorConfig{
		SensorID:       "sensor-123",
		SamplingRateMs: 0,
		ErrorRate:      0.1,
		Enabled:        true,
	}

	err = sensor.UpdateConfig(invalidSamplingRate)
	if err == nil {
		t.Error("expected error for invalid sampling rate")
	}

	invalidErrorRate := SensorConfig{
		SensorID:       "sensor-123",
		SamplingRateMs: 1000,
		ErrorRate:      1.5,
		Enabled:        true,
	}

	err = sensor.UpdateConfig(invalidErrorRate)
	if err == nil {
		t.Error("expected error for invalid error rate")
	}
}
