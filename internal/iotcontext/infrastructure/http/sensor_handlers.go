package http

import (
	"encoding/json"
	"fmt"
	"github.com/SeiyaJapon/iot-sensor-app/internal/iotcontext/application"
	"github.com/SeiyaJapon/iot-sensor-app/internal/iotcontext/domain"
	"github.com/google/uuid"
	"net/http"
)

type SensorHandlers struct {
	application.SensorUseCase
}

func NewSensorHandlers(sensorUseCase application.SensorUseCase) *SensorHandlers {
	return &SensorHandlers{
		SensorUseCase: sensorUseCase,
	}
}

type CreateSensorRequest struct {
	Name     string                 `json:"name"`
	Type     string                 `json:"type"`
	DeviceID string                 `json:"device_id"`
	Config   map[string]interface{} `json:"config"`
}

func (h *SensorHandlers) SensorsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		if r.URL.Path != "/sensors?" {
			h.GetAllSensors(w)
		} else {
			id := r.URL.Query().Get("id")
			if id == "" {
				http.Error(w, "Missing sensor ID", http.StatusBadRequest)
				return
			}

			h.GetSensorByID(w, r)
		}
	case http.MethodPost:
		h.CreateSensor(w, r)
	case http.MethodPut:
		h.UpdateSensorConfigById(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *SensorHandlers) CreateSensor(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CreateSensorRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if err := h.validateCreateRequest(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sensorID := h.generateSensorID()

	sensorConfig, err := h.unmarshalConfig(req.Config)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sensorConfig.SensorID = sensorID

	if err := h.SensorUseCase.CreateSensor(
		sensorID,
		domain.DeviceID(req.DeviceID),
		req.Name,
		domain.SensorType(req.Type),
		sensorConfig,
	); err != nil {
		http.Error(w, fmt.Sprintf("Failed to create sensor: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	response := map[string]interface{}{"message": "Sensor created successfully", "sensor_id": sensorID}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		return
	}
}

func (h *SensorHandlers) GetSensorByID(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Missing sensor ID", http.StatusBadRequest)
		return
	}

	sensor, err := h.SensorUseCase.GetSensorByID(domain.SensorID(id))
	if err != nil {
		http.Error(w, "Sensor not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(sensor); err != nil {
		http.Error(w, "Failed to encode sensor", http.StatusInternalServerError)
		return
	}
}

func (h *SensorHandlers) GetAllSensors(w http.ResponseWriter) {
	sensors, err := h.SensorUseCase.GetAllSensors()
	if err != nil {
		http.Error(w, "Failed to retrieve sensors", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(sensors); err != nil {
		http.Error(w, "Failed to encode sensors", http.StatusInternalServerError)
		return
	}
}

func (h *SensorHandlers) UpdateSensorConfigById(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Missing sensor ID", http.StatusBadRequest)
		return
	}

	var req map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	sensorConfig, err := h.unmarshalConfig(req)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid config: %v", err), http.StatusBadRequest)
		return
	}

	if err := h.SensorUseCase.UpdateSensorConfigById(domain.SensorID(id), sensorConfig); err != nil {
		http.Error(w, fmt.Sprintf("Failed to update sensor config: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]interface{}{"message": "Sensor config updated successfully"}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *SensorHandlers) validateCreateRequest(req CreateSensorRequest) error {
	if req.Name == "" || req.Type == "" || req.DeviceID == "" {
		return fmt.Errorf("missing required fields: name, type, device_id")
	}
	return nil
}

func (h *SensorHandlers) generateSensorID() domain.SensorID {
	return domain.SensorID(uuid.New().String())
}

func (h *SensorHandlers) unmarshalConfig(configMap map[string]interface{}) (domain.SensorConfig, error) {
	configBytes, err := json.Marshal(configMap)
	if err != nil {
		return domain.SensorConfig{}, fmt.Errorf("could not marshal config: %w", err)
	}

	var sensorConfig domain.SensorConfig
	if err := json.Unmarshal(configBytes, &sensorConfig); err != nil {
		return domain.SensorConfig{}, fmt.Errorf("could not unmarshal config: %w", err)
	}
	return sensorConfig, nil
}
