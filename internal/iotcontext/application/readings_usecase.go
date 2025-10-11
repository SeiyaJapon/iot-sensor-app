package application

import "github.com/SeiyaJapon/iot-sensor-app/internal/iotcontext/domain"

type ReadingsUsecase struct {
	readingsRepo domain.SensorReadingRepository
}

func NewReadingsUsecase(readingsRepo domain.SensorReadingRepository) *ReadingsUsecase {
	return &ReadingsUsecase{
		readingsRepo: readingsRepo,
	}
}

func (uc *ReadingsUsecase) GetPaginatedReadings(id domain.SensorID, from int, to int, limit int) ([]domain.SensorReading, error) {
	if from < 0 || to < 0 || limit <= 0 || from >= to {
		return nil, domain.ErrInvalidPaginationParams
	}

	readings, err := uc.readingsRepo.FindBySensorID(id, limit)
	if err != nil {
		return nil, err
	}

	if from >= len(readings) {
		return []domain.SensorReading{}, nil
	}

	if to > len(readings) {
		to = len(readings)
	}

	return readings[from:to], nil
}
