package domain

import (
	"testing"
)

func TestNewSensorConfig(t *testing.T) {
	tests := []struct {
		name           string
		sensorID       SensorID
		samplingRateMs int
		thresholds     Thresholds
		errorRate      float64
		enabled        bool
		expectError    bool
	}{
		{
			name:           "valid config",
			sensorID:       "sensor-123",
			samplingRateMs: 1000,
			thresholds:     Thresholds{},
			errorRate:      0.1,
			enabled:        true,
			expectError:    false,
		},
		{
			name:           "empty sensor id",
			sensorID:       "",
			samplingRateMs: 1000,
			thresholds:     Thresholds{},
			errorRate:      0.1,
			enabled:        true,
			expectError:    true,
		},
		{
			name:           "invalid sampling rate",
			sensorID:       "sensor-123",
			samplingRateMs: 0,
			thresholds:     Thresholds{},
			errorRate:      0.1,
			enabled:        true,
			expectError:    true,
		},
		{
			name:           "negative sampling rate",
			sensorID:       "sensor-123",
			samplingRateMs: -100,
			thresholds:     Thresholds{},
			errorRate:      0.1,
			enabled:        true,
			expectError:    true,
		},
		{
			name:           "invalid error rate negative",
			sensorID:       "sensor-123",
			samplingRateMs: 1000,
			thresholds:     Thresholds{},
			errorRate:      -0.1,
			enabled:        true,
			expectError:    true,
		},
		{
			name:           "invalid error rate above 1",
			sensorID:       "sensor-123",
			samplingRateMs: 1000,
			thresholds:     Thresholds{},
			errorRate:      1.5,
			enabled:        true,
			expectError:    true,
		},
		{
			name:           "valid error rate 0",
			sensorID:       "sensor-123",
			samplingRateMs: 1000,
			thresholds:     Thresholds{},
			errorRate:      0.0,
			enabled:        true,
			expectError:    false,
		},
		{
			name:           "valid error rate 1",
			sensorID:       "sensor-123",
			samplingRateMs: 1000,
			thresholds:     Thresholds{},
			errorRate:      1.0,
			enabled:        true,
			expectError:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config, err := NewSensorConfig(tt.sensorID, tt.samplingRateMs, tt.thresholds, tt.errorRate, tt.enabled)

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

			if config.SensorID != tt.sensorID {
				t.Errorf("expected SensorID %s, got %s", tt.sensorID, config.SensorID)
			}

			if config.SamplingRateMs != tt.samplingRateMs {
				t.Errorf("expected SamplingRateMs %d, got %d", tt.samplingRateMs, config.SamplingRateMs)
			}

			if config.ErrorRate != tt.errorRate {
				t.Errorf("expected ErrorRate %f, got %f", tt.errorRate, config.ErrorRate)
			}

			if config.Enabled != tt.enabled {
				t.Errorf("expected Enabled %t, got %t", tt.enabled, config.Enabled)
			}

			if config.UpdatedAt.IsZero() {
				t.Error("expected UpdatedAt to be set")
			}

			if config.Meta == nil {
				t.Error("expected Meta to be initialized")
			}
		})
	}
}
