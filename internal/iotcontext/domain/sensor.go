package domain

import (
	"errors"
	"time"
)

type Sensor struct {
	ID        SensorID     `json:"id"`
	DeviceID  DeviceID     `json:"device_id"`
	Name      string       `json:"name"`
	Type      SensorType   `json:"type"`
	Config    SensorConfig `json:"config"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
}

func NewSensor(id SensorID, deviceID DeviceID, name string, typ SensorType, config SensorConfig) (*Sensor, error) {
	if id == "" {
		return nil, errors.New("sensor id empty")
	}

	if deviceID == "" {
		return nil, errors.New("device id empty")
	}

	if name == "" {
		return nil, errors.New("name empty")
	}

	if typ == "" {
		typ = Generic
	}

	now := time.Now().UTC()
	config.SensorID = id
	config.UpdatedAt = now

	return &Sensor{
		ID:        id,
		DeviceID:  deviceID,
		Name:      name,
		Type:      typ,
		Config:    config,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func (s *Sensor) UpdateConfig(cfg SensorConfig) error {
	if cfg.SensorID != s.ID {
		return errors.New("config sensor_id mismatch")
	}

	if cfg.SamplingRateMs <= 0 {
		return errors.New("sampling rate must be positive")
	}

	if cfg.ErrorRate < 0 || cfg.ErrorRate > 1 {
		return errors.New("error rate must be between 0 and 1")
	}

	cfg.UpdatedAt = time.Now().UTC()
	s.Config = cfg
	s.UpdatedAt = cfg.UpdatedAt

	return nil
}
