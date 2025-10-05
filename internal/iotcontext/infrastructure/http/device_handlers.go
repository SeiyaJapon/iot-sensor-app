package http

import (
	"encoding/json"
	"github.com/SeiyaJapon/iot-sensor-app/internal/iotcontext/application"
	"github.com/SeiyaJapon/iot-sensor-app/internal/iotcontext/domain"
	"github.com/google/uuid"
	"net/http"
)

type DeviceHandlers struct {
	deviceUseCase application.DeviceUseCase
}

func NewDeviceHandlers(deviceUseCase application.DeviceUseCase) *DeviceHandlers {
	return &DeviceHandlers{
		deviceUseCase: deviceUseCase,
	}
}

func (h *DeviceHandlers) All(w http.ResponseWriter, r *http.Request) {
	devices, err := h.deviceUseCase.GetAllDevices()
	if err != nil {
		http.Error(w, "Failed to retrieve devices", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(devices); err != nil {
		http.Error(w, "Failed to encode devices", http.StatusInternalServerError)
		return
	}
}

func (h *DeviceHandlers) GetByID(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Missing device ID", http.StatusBadRequest)
		return
	}

	device, err := h.deviceUseCase.GetDeviceByID(domain.DeviceID(id))
	if err != nil {
		http.Error(w, "Device not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(device); err != nil {
		http.Error(w, "Failed to encode device", http.StatusInternalServerError)
		return
	}
}

func (h *DeviceHandlers) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name string `json:"name"`
		Type string `json:"type"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	id := domain.DeviceID(uuid.New().String())

	device, err := h.deviceUseCase.CreateDevice(id, req.Name, req.Type)
	if err != nil {
		http.Error(w, "Failed to create device: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(device); err != nil {
		http.Error(w, "Failed to encode device", http.StatusInternalServerError)
		return
	}
}
