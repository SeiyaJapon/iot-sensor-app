package domain

import (
	"errors"
	"time"
)

type SensorConfig struct {
	SensorID       SensorID               `json:"sensor_id"`
	SamplingRateMs int                    `json:"sampling_rate_ms"`
	Thresholds     Thresholds             `json:"thresholds"`
	ErrorRate      float64                `json:"error_rate"`
	Enabled        bool                   `json:"enabled"`
	UpdatedAt      time.Time              `json:"updated_at"`
	Meta           map[string]interface{} `json:"meta"`
}

func NewSensorConfig(sensorID SensorID, samplingRateMs int, thresholds Thresholds, errorRate float64, enabled bool) (SensorConfig, error) {
	if sensorID == "" {
		return SensorConfig{}, errors.New("sensor id empty")
	}

	if samplingRateMs <= 0 {
		return SensorConfig{}, errors.New("sampling rate must be positive")
	}

	if errorRate < 0 || errorRate > 1 {
		return SensorConfig{}, errors.New("error rate between 0 and 1")
	}

	return SensorConfig{
		SensorID:       sensorID,
		SamplingRateMs: samplingRateMs,
		Thresholds:     thresholds,
		ErrorRate:      errorRate,
		Enabled:        enabled,
		UpdatedAt:      time.Now().UTC(),
		Meta:           map[string]interface{}{},
	}, nil
}
