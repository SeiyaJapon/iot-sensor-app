package application

import (
	"errors"
	"github.com/SeiyaJapon/iot-sensor-app/internal/iotcontext/domain"
	"testing"
)

func TestSensorUseCase_CreateSensor(t *testing.T) {
	tests := []struct {
		name        string
		id          domain.SensorID
		deviceID    domain.DeviceID
		sensorName  string
		sensorType  domain.SensorType
		config      domain.SensorConfig
		repoSaveErr error
		expectError bool
		expectEvent bool
	}{
		{
			name:       "successful creation",
			id:         "sensor-123",
			deviceID:   "device-123",
			sensorName: "Temperature Sensor",
			sensorType: domain.Temperature,
			config: domain.SensorConfig{
				SamplingRateMs: 1000,
				ErrorRate:      0.1,
				Enabled:        true,
			},
			repoSaveErr: nil,
			expectError: false,
			expectEvent: true,
		},
		{
			name:       "repository save error",
			id:         "sensor-123",
			deviceID:   "device-123",
			sensorName: "Temperature Sensor",
			sensorType: domain.Temperature,
			config: domain.SensorConfig{
				SamplingRateMs: 1000,
				ErrorRate:      0.1,
				Enabled:        true,
			},
			repoSaveErr: errors.New("database error"),
			expectError: true,
			expectEvent: false,
		},
		{
			name:       "invalid sensor data",
			id:         "",
			deviceID:   "device-123",
			sensorName: "Temperature Sensor",
			sensorType: domain.Temperature,
			config: domain.SensorConfig{
				SamplingRateMs: 1000,
				ErrorRate:      0.1,
				Enabled:        true,
			},
			repoSaveErr: nil,
			expectError: true,
			expectEvent: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := NewMockSensorRepository()
			mockRepo.saveErr = tt.repoSaveErr
			mockPublisher := NewMockEventPublisher()
			mockMetrics := NewMockMetrics()

			useCase := NewSensorUseCase(mockRepo, mockMetrics, mockPublisher)

			err := useCase.CreateSensor(tt.id, tt.deviceID, tt.sensorName, tt.sensorType, tt.config)

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

			sensor, err := mockRepo.FindByID(tt.id)
			if err != nil {
				t.Errorf("sensor not found in repository: %v", err)
			}

			if sensor.Name != tt.sensorName {
				t.Errorf("expected name %s, got %s", tt.sensorName, sensor.Name)
			}

			if sensor.Type != tt.sensorType {
				t.Errorf("expected type %s, got %s", tt.sensorType, sensor.Type)
			}

			if tt.expectEvent {
				events := mockPublisher.GetEvents()
				if len(events) != 1 {
					t.Errorf("expected 1 event, got %d", len(events))
				}

				if events[0].Type != "sensor.created" {
					t.Errorf("expected event type 'sensor.created', got %s", events[0].Type)
				}
			}

			readings := mockMetrics.GetSensorReadings(tt.sensorType, tt.deviceID)
			if readings != 1 {
				t.Errorf("expected 1 sensor reading metric, got %d", readings)
			}
		})
	}
}

func TestSensorUseCase_GetSensorByID(t *testing.T) {
	mockRepo := NewMockSensorRepository()
	mockPublisher := NewMockEventPublisher()
	mockMetrics := NewMockMetrics()

	useCase := NewSensorUseCase(mockRepo, mockMetrics, mockPublisher)

	sensor, err := domain.NewSensor("sensor-123", "device-123", "Test Sensor", domain.Temperature, domain.SensorConfig{})
	if err != nil {
		t.Fatalf("failed to create test sensor: %v", err)
	}

	mockRepo.Save(sensor)

	retrievedSensor, err := useCase.GetSensorByID("sensor-123")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if retrievedSensor.ID != "sensor-123" {
		t.Errorf("expected ID sensor-123, got %s", retrievedSensor.ID)
	}

	_, err = useCase.GetSensorByID("nonexistent")
	if err == nil {
		t.Error("expected error for nonexistent sensor")
	}
}

func TestSensorUseCase_GetAllSensors(t *testing.T) {
	mockRepo := NewMockSensorRepository()
	mockPublisher := NewMockEventPublisher()
	mockMetrics := NewMockMetrics()

	useCase := NewSensorUseCase(mockRepo, mockMetrics, mockPublisher)

	sensor1, _ := domain.NewSensor("sensor-1", "device-1", "Sensor 1", domain.Temperature, domain.SensorConfig{})
	sensor2, _ := domain.NewSensor("sensor-2", "device-2", "Sensor 2", domain.Humidity, domain.SensorConfig{})

	mockRepo.Save(sensor1)
	mockRepo.Save(sensor2)

	sensors, err := useCase.GetAllSensors()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(sensors) != 2 {
		t.Errorf("expected 2 sensors, got %d", len(sensors))
	}
}

func TestSensorUseCase_UpdateSensorConfigById(t *testing.T) {
	tests := []struct {
		name          string
		sensorID      domain.SensorID
		config        domain.SensorConfig
		repoFindErr   error
		repoUpdateErr error
		expectError   bool
		expectEvent   bool
	}{
		{
			name:     "successful update",
			sensorID: "sensor-123",
			config: domain.SensorConfig{
				SensorID:       "sensor-123",
				SamplingRateMs: 2000,
				ErrorRate:      0.2,
				Enabled:        true,
			},
			repoFindErr:   nil,
			repoUpdateErr: nil,
			expectError:   false,
			expectEvent:   true,
		},
		{
			name:     "sensor not found",
			sensorID: "nonexistent",
			config: domain.SensorConfig{
				SensorID:       "nonexistent",
				SamplingRateMs: 2000,
				ErrorRate:      0.2,
				Enabled:        true,
			},
			repoFindErr:   domain.ErrSensorNotFound,
			repoUpdateErr: nil,
			expectError:   true,
			expectEvent:   false,
		},
		{
			name:     "repository update error",
			sensorID: "sensor-123",
			config: domain.SensorConfig{
				SensorID:       "sensor-123",
				SamplingRateMs: 2000,
				ErrorRate:      0.2,
				Enabled:        true,
			},
			repoFindErr:   nil,
			repoUpdateErr: errors.New("update error"),
			expectError:   true,
			expectEvent:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := NewMockSensorRepository()
			mockRepo.findErr = tt.repoFindErr
			mockRepo.updateErr = tt.repoUpdateErr
			mockPublisher := NewMockEventPublisher()
			mockMetrics := NewMockMetrics()

			if tt.sensorID == "sensor-123" && tt.repoFindErr == nil {
				sensor, _ := domain.NewSensor("sensor-123", "device-123", "Test Sensor", domain.Temperature, domain.SensorConfig{})
				mockRepo.Save(sensor)
			}

			useCase := NewSensorUseCase(mockRepo, mockMetrics, mockPublisher)

			err := useCase.UpdateSensorConfigById(tt.sensorID, tt.config)

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

			if tt.expectEvent {
				events := mockPublisher.GetEvents()
				if len(events) != 1 {
					t.Errorf("expected 1 event, got %d", len(events))
				}

				if events[0].Type != "sensor.config.updated" {
					t.Errorf("expected event type 'sensor.config.updated', got %s", events[0].Type)
				}
			}
		})
	}
}
