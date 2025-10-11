package http

import (
	"github.com/SeiyaJapon/iot-sensor-app/internal/iotcontext/application"
	"github.com/SeiyaJapon/iot-sensor-app/internal/iotcontext/domain"
	"net/http"
)

type SimulatorHandler struct {
	simulatorUsecase application.SimulatorUseCase
}

func NewSimulatorHandler(simulatorUsecase application.SimulatorUseCase) *SimulatorHandler {
	return &SimulatorHandler{
		simulatorUsecase: simulatorUsecase,
	}
}

func (h *SimulatorHandler) ControlSensor(r *http.Request, w http.ResponseWriter) error {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return nil
	}

	sensorID := r.URL.Query().Get("sensor_id")
	action := r.URL.Query().Get("action")
	if sensorID == "" || action == "" {
		http.Error(w, "Missing sensor_id or action parameter", http.StatusBadRequest)
		return nil
	}

	if err := h.simulatorUsecase.ControlSensor(domain.SensorID(sensorID), action); err != nil {
		switch err {
		case domain.ErrSensorNotFound:
			http.Error(w, "Sensor not found", http.StatusNotFound)
		case domain.ErrInvalidAction:
			http.Error(w, "Invalid action", http.StatusBadRequest)
		default:
			http.Error(w, "Failed to control sensor", http.StatusInternalServerError)
		}
		return err
	}

	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("Sensor " + action + " command executed successfully"))
	if err != nil {
		return err
	}
	return nil
}

func (h *SimulatorHandler) SimulatorsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		err := h.ControlSensor(r, w)
		if err != nil {
			return
		}
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
