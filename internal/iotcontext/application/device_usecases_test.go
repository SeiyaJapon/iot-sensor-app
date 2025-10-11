package application

import (
	"errors"
	"github.com/SeiyaJapon/iot-sensor-app/internal/iotcontext/domain"
	"testing"
)

func TestDeviceUseCase_CreateDevice(t *testing.T) {
	tests := []struct {
		name        string
		id          domain.DeviceID
		deviceName  string
		deviceType  string
		repoSaveErr error
		expectError bool
	}{
		{
			name:        "successful creation",
			id:          "device-123",
			deviceName:  "IoT Device",
			deviceType:  "sensor_hub",
			repoSaveErr: nil,
			expectError: false,
		},
		{
			name:        "repository save error",
			id:          "device-123",
			deviceName:  "IoT Device",
			deviceType:  "sensor_hub",
			repoSaveErr: errors.New("database error"),
			expectError: true,
		},
		{
			name:        "invalid device data",
			id:          "",
			deviceName:  "IoT Device",
			deviceType:  "sensor_hub",
			repoSaveErr: nil,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := NewMockDeviceRepository()
			mockRepo.saveErr = tt.repoSaveErr

			useCase := NewDeviceUseCase(mockRepo)

			device, err := useCase.CreateDevice(tt.id, tt.deviceName, tt.deviceType)

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
		})
	}
}

func TestDeviceUseCase_GetDeviceByID(t *testing.T) {
	tests := []struct {
		name        string
		deviceID    domain.DeviceID
		repoFindErr error
		expectError bool
	}{
		{
			name:        "device found",
			deviceID:    "device-123",
			repoFindErr: nil,
			expectError: false,
		},
		{
			name:        "device not found",
			deviceID:    "nonexistent",
			repoFindErr: domain.ErrDeviceNotFound,
			expectError: true,
		},
		{
			name:        "repository error",
			deviceID:    "device-123",
			repoFindErr: errors.New("database error"),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := NewMockDeviceRepository()
			mockRepo.findErr = tt.repoFindErr

			if tt.deviceID == "device-123" && tt.repoFindErr == nil {
				device, _ := domain.NewDevice("device-123", "Test Device", "sensor_hub")
				mockRepo.Save(device)
			}

			useCase := NewDeviceUseCase(mockRepo)

			device, err := useCase.GetDeviceByID(tt.deviceID)

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

			if device.ID != tt.deviceID {
				t.Errorf("expected ID %s, got %s", tt.deviceID, device.ID)
			}
		})
	}
}

func TestDeviceUseCase_GetAllDevices(t *testing.T) {
	tests := []struct {
		name          string
		repoFindErr   error
		expectError   bool
		expectedCount int
	}{
		{
			name:          "successful retrieval",
			repoFindErr:   nil,
			expectError:   false,
			expectedCount: 2,
		},
		{
			name:          "repository error",
			repoFindErr:   errors.New("database error"),
			expectError:   true,
			expectedCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := NewMockDeviceRepository()
			mockRepo.findErr = tt.repoFindErr

			if tt.repoFindErr == nil {
				device1, _ := domain.NewDevice("device-1", "Device 1", "sensor_hub")
				device2, _ := domain.NewDevice("device-2", "Device 2", "gateway")
				mockRepo.Save(device1)
				mockRepo.Save(device2)
			}

			useCase := NewDeviceUseCase(mockRepo)

			devices, err := useCase.GetAllDevices()

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

			if len(devices) != tt.expectedCount {
				t.Errorf("expected %d devices, got %d", tt.expectedCount, len(devices))
			}
		})
	}
}

func TestDeviceUseCase_UpdateDevice(t *testing.T) {
	tests := []struct {
		name          string
		device        *domain.Device
		repoUpdateErr error
		expectError   bool
	}{
		{
			name: "successful update",
			device: &domain.Device{
				ID:   "device-123",
				Name: "Updated Device",
				Type: "sensor_hub",
			},
			repoUpdateErr: nil,
			expectError:   false,
		},
		{
			name: "repository update error",
			device: &domain.Device{
				ID:   "device-123",
				Name: "Updated Device",
				Type: "sensor_hub",
			},
			repoUpdateErr: errors.New("update error"),
			expectError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := NewMockDeviceRepository()
			mockRepo.updateErr = tt.repoUpdateErr

			useCase := NewDeviceUseCase(mockRepo)

			err := useCase.UpdateDevice(tt.device)

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

			if tt.device.UpdatedAt.IsZero() {
				t.Error("expected UpdatedAt to be set")
			}
		})
	}
}
