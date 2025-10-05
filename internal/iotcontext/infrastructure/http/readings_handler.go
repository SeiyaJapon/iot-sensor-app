package http

import (
	"github.com/SeiyaJapon/iot-sensor-app/internal/iotcontext/application"
	"github.com/SeiyaJapon/iot-sensor-app/internal/iotcontext/domain"
)

type ReadingsHandler struct {
	readingsUsecase application.ReadingsUsecase
}

func NewReadingsHandler(readingsUsecase application.ReadingsUsecase) *ReadingsHandler {
	return &ReadingsHandler{
		readingsUsecase: readingsUsecase,
	}
}

func (h *ReadingsHandler) GetPaginatedReadings(sensorID string, from int, to int, limit int) ([]domain.SensorReading, error) {
	return h.readingsUsecase.GetPaginatedReadings(domain.SensorID(sensorID), from, to, limit)
}
