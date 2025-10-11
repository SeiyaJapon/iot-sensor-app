package http

import (
	"github.com/SeiyaJapon/iot-sensor-app/internal/iotcontext/application"
	"github.com/SeiyaJapon/iot-sensor-app/internal/iotcontext/domain"
	"net/http"
	"strconv"
)

type ReadingsHandler struct {
	readingsUsecase application.ReadingsUsecase
}

func NewReadingsHandler(readingsUsecase application.ReadingsUsecase) *ReadingsHandler {
	return &ReadingsHandler{
		readingsUsecase: readingsUsecase,
	}
}

func (h *ReadingsHandler) GetPaginatedReadings(r *http.Request, w http.ResponseWriter) ([]domain.SensorReading, error) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return nil, nil
	}

	sensorID, from, to, limit, _ := h.parsePaginationParams(r, w)

	readings, err := h.readingsUsecase.GetPaginatedReadings(domain.SensorID(sensorID), from, to, limit)
	if err != nil {
		http.Error(w, "Failed to retrieve readings", http.StatusInternalServerError)
		return nil, err
	}

	return readings, nil
}

func (h *ReadingsHandler) parsePaginationParams(r *http.Request, w http.ResponseWriter) (string, int, int, int, error) {
	sensorID := r.URL.Query().Get("sensor_id")
	from, err := strconv.Atoi(r.URL.Query().Get("from"))
	if err != nil {
		http.Error(w, "Invalid 'from' parameter", http.StatusBadRequest)
		return "", 0, 0, 0, err
	}
	to, err := strconv.Atoi(r.URL.Query().Get("to"))
	if err != nil {
		http.Error(w, "Invalid 'to' parameter", http.StatusBadRequest)
		return "", 0, 0, 0, err
	}
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		http.Error(w, "Invalid 'limit' parameter", http.StatusBadRequest)
		return "", 0, 0, 0, err
	}

	return sensorID, from, to, limit, nil
}

func (h *ReadingsHandler) SensorReadingsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		_, err := h.GetPaginatedReadings(r, w)
		if err != nil {
			return
		}
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
