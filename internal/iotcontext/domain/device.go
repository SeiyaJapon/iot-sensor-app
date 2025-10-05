package domain

import (
	"errors"
	"time"
)

type Device struct {
	ID        DeviceID  `json:"id"`
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewDevice(id DeviceID, name string, typ string) (*Device, error) {
	if id == "" || name == "" {
		return nil, errors.New("device id or name empty")
	}

	now := time.Now().UTC()

	return &Device{
		ID:        id,
		Name:      name,
		Type:      typ,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}
