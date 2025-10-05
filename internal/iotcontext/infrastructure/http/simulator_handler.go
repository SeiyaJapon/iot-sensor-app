package http

import (
	"github.com/SeiyaJapon/iot-sensor-app/internal/iotcontext/application"
	"github.com/SeiyaJapon/iot-sensor-app/internal/iotcontext/domain"
)

type SimulatorHandler struct {
	simulatorUsecase application.SimulatorUseCase
}

func NewSimulatorHandler(simulatorUsecase application.SimulatorUseCase) *SimulatorHandler {
	return &SimulatorHandler{
		simulatorUsecase: simulatorUsecase,
	}
}

func (h *SimulatorHandler) ControlSensor(sensorID string, action string) error {
	return h.simulatorUsecase.ControlSensor(domain.SensorID(sensorID), action)
}
