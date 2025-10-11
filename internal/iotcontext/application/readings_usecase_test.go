package application

import (
	"errors"
	"github.com/SeiyaJapon/iot-sensor-app/internal/iotcontext/domain"
	"testing"
)

func TestReadingsUsecase_GetPaginatedReadings(t *testing.T) {
	tests := []struct {
		name          string
		sensorID      domain.SensorID
		from          int
		to            int
		limit         int
		repoFindErr   error
		expectError   bool
		expectedCount int
	}{
		{
			name:          "valid pagination",
			sensorID:      "sensor-123",
			from:          0,
			to:            2,
			limit:         10,
			repoFindErr:   nil,
			expectError:   false,
			expectedCount: 2,
		},
		{
			name:          "from greater than to",
			sensorID:      "sensor-123",
			from:          5,
			to:            2,
			limit:         10,
			repoFindErr:   nil,
			expectError:   true,
			expectedCount: 0,
		},
		{
			name:          "negative from",
			sensorID:      "sensor-123",
			from:          -1,
			to:            2,
			limit:         10,
			repoFindErr:   nil,
			expectError:   true,
			expectedCount: 0,
		},
		{
			name:          "negative to",
			sensorID:      "sensor-123",
			from:          0,
			to:            -1,
			limit:         10,
			repoFindErr:   nil,
			expectError:   true,
			expectedCount: 0,
		},
		{
			name:          "zero limit",
			sensorID:      "sensor-123",
			from:          0,
			to:            2,
			limit:         0,
			repoFindErr:   nil,
			expectError:   true,
			expectedCount: 0,
		},
		{
			name:          "from equals to",
			sensorID:      "sensor-123",
			from:          2,
			to:            2,
			limit:         10,
			repoFindErr:   nil,
			expectError:   true,
			expectedCount: 0,
		},
		{
			name:          "repository error",
			sensorID:      "sensor-123",
			from:          0,
			to:            2,
			limit:         10,
			repoFindErr:   errors.New("database error"),
			expectError:   true,
			expectedCount: 0,
		},
		{
			name:          "from greater than available readings",
			sensorID:      "sensor-123",
			from:          10,
			to:            15,
			limit:         10,
			repoFindErr:   nil,
			expectError:   false,
			expectedCount: 0,
		},
		{
			name:          "to greater than available readings",
			sensorID:      "sensor-123",
			from:          0,
			to:            15,
			limit:         10,
			repoFindErr:   nil,
			expectError:   false,
			expectedCount: 10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := NewMockSensorReadingRepository()
			mockRepo.findErr = tt.repoFindErr

			if tt.repoFindErr == nil {
				for i := 0; i < 10; i++ {
					reading := domain.NewSensorReading(
						tt.sensorID,
						"device-123",
						domain.Temperature,
						20.0+float64(i),
						"°C",
						domain.SensorReading{}.Timestamp,
					)
					mockRepo.Save(&reading)
				}
			}

			useCase := NewReadingsUsecase(mockRepo)

			readings, err := useCase.GetPaginatedReadings(tt.sensorID, tt.from, tt.to, tt.limit)

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

			if len(readings) != tt.expectedCount {
				t.Errorf("expected %d readings, got %d", tt.expectedCount, len(readings))
			}
		})
	}
}
