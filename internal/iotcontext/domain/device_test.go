package domain

import (
	"testing"
)

func TestNewDevice(t *testing.T) {
	tests := []struct {
		name        string
		id          DeviceID
		deviceName  string
		deviceType  string
		expectError bool
	}{
		{
			name:        "valid device",
			id:          "device-123",
			deviceName:  "IoT Device",
			deviceType:  "sensor_hub",
			expectError: false,
		},
		{
			name:        "empty device id",
			id:          "",
			deviceName:  "IoT Device",
			deviceType:  "sensor_hub",
			expectError: true,
		},
		{
			name:        "empty device name",
			id:          "device-123",
			deviceName:  "",
			deviceType:  "sensor_hub",
			expectError: true,
		},
		{
			name:        "both empty",
			id:          "",
			deviceName:  "",
			deviceType:  "sensor_hub",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			device, err := NewDevice(tt.id, tt.deviceName, tt.deviceType)

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

			if device.ID != tt.id {
				t.Errorf("expected ID %s, got %s", tt.id, device.ID)
			}

			if device.Name != tt.deviceName {
				t.Errorf("expected Name %s, got %s", tt.deviceName, device.Name)
			}

			if device.Type != tt.deviceType {
				t.Errorf("expected Type %s, got %s", tt.deviceType, device.Type)
			}

			if device.CreatedAt.IsZero() {
				t.Error("expected CreatedAt to be set")
			}

			if device.UpdatedAt.IsZero() {
				t.Error("expected UpdatedAt to be set")
			}
		})
	}
}
