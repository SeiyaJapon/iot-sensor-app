package persistence

import (
	"errors"
	"math/rand"
	"sync"
	"time"

	"github.com/SeiyaJapon/iot-sensor-app/internal/iotcontext/domain"
)

type SimulatorRepositoryImpl struct {
	sensorRepo        domain.SensorRepository
	sensorReadingRepo domain.SensorReadingRepository
	eventPublisher    domain.EventPublisher
	activeSensors     map[domain.SensorID]*simulatorState
	mu                sync.RWMutex
}

func NewSimulatorRepository(sensorRepo domain.SensorRepository, sensorReadingRepo domain.SensorReadingRepository, eventPublisher domain.EventPublisher) domain.SimulatorRepository {
	return &SimulatorRepositoryImpl{
		sensorRepo:        sensorRepo,
		sensorReadingRepo: sensorReadingRepo,
		eventPublisher:    eventPublisher,
		activeSensors:     make(map[domain.SensorID]*simulatorState),
	}
}

func (s *SimulatorRepositoryImpl) Start(sensorID domain.SensorID) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.activeSensors[sensorID]; ok {
		return errors.New("sensor already active")
	}

	sensor, err := s.sensorRepo.FindByID(sensorID)
	if err != nil {
		return err
	}
	if !sensor.Config.Enabled {
		return errors.New("sensor is disabled")
	}

	state := &simulatorState{
		sensor:      sensor,
		stopCh:      make(chan struct{}),
		ticker:      time.NewTicker(time.Duration(sensor.Config.SamplingRateMs) * time.Millisecond),
		injectError: false,
	}
	s.activeSensors[sensorID] = state

	go s.simulateReadings(sensorID, state)

	return nil
}

func (s *SimulatorRepositoryImpl) Stop(sensorID domain.SensorID) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	state, ok := s.activeSensors[sensorID]
	if !ok {
		return errors.New("sensor not active")
	}

	close(state.stopCh)
	state.ticker.Stop()
	delete(s.activeSensors, sensorID)

	return nil
}

func (s *SimulatorRepositoryImpl) InjectError(sensorID domain.SensorID) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	state, ok := s.activeSensors[sensorID]
	if !ok {
		return errors.New("sensor not active")
	}

	state.injectError = true

	return nil
}

func (s *SimulatorRepositoryImpl) simulateReadings(sensorID domain.SensorID, state *simulatorState) {
	defer state.ticker.Stop()

	for {
		select {
		case <-state.stopCh:
			return
		case <-state.ticker.C:
			if state.injectError {
				errorEvent := &domain.SensorReadingErrorEvent{
					SensorID: sensorID,
					Type:     "injection",
				}
				_ = s.eventPublisher.Publish(errorEvent.ToDomainEvent())
				state.injectError = false
				continue
			}

			if rand.Float64() < state.sensor.Config.ErrorRate {
				continue
			}

			value := s.generateValue(state.sensor.Type)
			unit := s.getUnit(state.sensor.Type)

			reading := domain.NewSensorReading(
				sensorID,
				state.sensor.DeviceID,
				state.sensor.Type,
				value,
				unit,
				time.Now().UTC(),
			)

			if err := s.sensorReadingRepo.Save(&reading); err != nil {
				continue
			}

			readingEvent := &domain.SensorReadingPublishedEvent{
				SensorID: sensorID,
				Reading:  reading.ID,
			}
			if err := s.eventPublisher.Publish(readingEvent.ToDomainEvent()); err != nil {
				continue
			}
		}
	}
}

type simulatorState struct {
	ticker      *time.Ticker
	stopCh      chan struct{}
	sensor      *domain.Sensor
	injectError bool
}

func (s *SimulatorRepositoryImpl) generateValue(typ domain.SensorType) float64 {
	switch typ {
	case domain.Temperature:
		return 20 + rand.Float64()*60
	case domain.Humidity:
		return rand.Float64() * 100
	case domain.Pressure:
		return 900 + rand.Float64()*200
	default:
		return rand.Float64() * 100
	}
}

func (s *SimulatorRepositoryImpl) getUnit(typ domain.SensorType) string {
	switch typ {
	case domain.Temperature:
		return "°C"
	case domain.Humidity:
		return "%"
	case domain.Pressure:
		return "hPa"
	default:
		return ""
	}
}
