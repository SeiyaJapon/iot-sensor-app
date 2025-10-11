package application

import (
	"github.com/SeiyaJapon/iot-sensor-app/internal/iotcontext/domain"
)

type MockSensorRepository struct {
	sensors   map[domain.SensorID]*domain.Sensor
	saveErr   error
	findErr   error
	updateErr error
}

func NewMockSensorRepository() *MockSensorRepository {
	return &MockSensorRepository{
		sensors: make(map[domain.SensorID]*domain.Sensor),
	}
}

func (m *MockSensorRepository) Save(sensor *domain.Sensor) error {
	if m.saveErr != nil {
		return m.saveErr
	}
	m.sensors[sensor.ID] = sensor
	return nil
}

func (m *MockSensorRepository) FindByID(id domain.SensorID) (*domain.Sensor, error) {
	if m.findErr != nil {
		return nil, m.findErr
	}
	sensor, exists := m.sensors[id]
	if !exists {
		return nil, domain.ErrSensorNotFound
	}
	return sensor, nil
}

func (m *MockSensorRepository) FindAll() ([]*domain.Sensor, error) {
	if m.findErr != nil {
		return nil, m.findErr
	}
	var sensors []*domain.Sensor
	for _, sensor := range m.sensors {
		sensors = append(sensors, sensor)
	}
	return sensors, nil
}

func (m *MockSensorRepository) Update(sensor *domain.Sensor) error {
	if m.updateErr != nil {
		return m.updateErr
	}
	m.sensors[sensor.ID] = sensor
	return nil
}

func (m *MockSensorRepository) Delete(id domain.SensorID) error {
	delete(m.sensors, id)
	return nil
}

type MockDeviceRepository struct {
	devices   map[domain.DeviceID]domain.Device
	saveErr   error
	findErr   error
	updateErr error
}

func NewMockDeviceRepository() *MockDeviceRepository {
	return &MockDeviceRepository{
		devices: make(map[domain.DeviceID]domain.Device),
	}
}

func (m *MockDeviceRepository) Save(device *domain.Device) error {
	if m.saveErr != nil {
		return m.saveErr
	}
	m.devices[device.ID] = *device
	return nil
}

func (m *MockDeviceRepository) FindByID(id domain.DeviceID) (domain.Device, error) {
	if m.findErr != nil {
		return domain.Device{}, m.findErr
	}
	device, exists := m.devices[id]
	if !exists {
		return domain.Device{}, domain.ErrDeviceNotFound
	}
	return device, nil
}

func (m *MockDeviceRepository) FindAll() ([]domain.Device, error) {
	if m.findErr != nil {
		return nil, m.findErr
	}
	var devices []domain.Device
	for _, device := range m.devices {
		devices = append(devices, device)
	}
	return devices, nil
}

func (m *MockDeviceRepository) Update(device *domain.Device) error {
	if m.updateErr != nil {
		return m.updateErr
	}
	m.devices[device.ID] = *device
	return nil
}

func (m *MockDeviceRepository) Delete(id domain.DeviceID) error {
	delete(m.devices, id)
	return nil
}

type MockSensorReadingRepository struct {
	readings map[domain.SensorID][]domain.SensorReading
	saveErr  error
	findErr  error
}

func NewMockSensorReadingRepository() *MockSensorReadingRepository {
	return &MockSensorReadingRepository{
		readings: make(map[domain.SensorID][]domain.SensorReading),
	}
}

func (m *MockSensorReadingRepository) Save(reading *domain.SensorReading) error {
	if m.saveErr != nil {
		return m.saveErr
	}
	m.readings[reading.SensorID] = append(m.readings[reading.SensorID], *reading)
	return nil
}

func (m *MockSensorReadingRepository) FindBySensorID(sensorID domain.SensorID, limit int) ([]domain.SensorReading, error) {
	if m.findErr != nil {
		return nil, m.findErr
	}
	readings, exists := m.readings[sensorID]
	if !exists {
		return []domain.SensorReading{}, nil
	}
	if limit > 0 && limit < len(readings) {
		return readings[:limit], nil
	}
	return readings, nil
}

func (m *MockSensorReadingRepository) FindByDeviceID(deviceID domain.DeviceID, limit int) ([]domain.SensorReading, error) {
	if m.findErr != nil {
		return nil, m.findErr
	}
	var allReadings []domain.SensorReading
	for _, readings := range m.readings {
		allReadings = append(allReadings, readings...)
	}
	if limit > 0 && limit < len(allReadings) {
		return allReadings[:limit], nil
	}
	return allReadings, nil
}

type MockEventPublisher struct {
	events     []domain.IoTEvent
	publishErr error
}

func NewMockEventPublisher() *MockEventPublisher {
	return &MockEventPublisher{
		events: make([]domain.IoTEvent, 0),
	}
}

func (m *MockEventPublisher) Publish(event domain.IoTEvent) error {
	if m.publishErr != nil {
		return m.publishErr
	}
	m.events = append(m.events, event)
	return nil
}

func (m *MockEventPublisher) GetEvents() []domain.IoTEvent {
	return m.events
}

func (m *MockEventPublisher) ClearEvents() {
	m.events = make([]domain.IoTEvent, 0)
}

type MockMetrics struct {
	sensorReadings map[string]int
	sensorErrors   map[string]int
}

func NewMockMetrics() *MockMetrics {
	return &MockMetrics{
		sensorReadings: make(map[string]int),
		sensorErrors:   make(map[string]int),
	}
}

func (m *MockMetrics) IncSensorReading(sensorType domain.SensorType, deviceID domain.DeviceID) {
	key := string(sensorType) + "_" + string(deviceID)
	m.sensorReadings[key]++
}

func (m *MockMetrics) IncSensorError(sensorType domain.SensorType, deviceID domain.DeviceID) {
	key := string(sensorType) + "_" + string(deviceID)
	m.sensorErrors[key]++
}

func (m *MockMetrics) GetSensorReadings(sensorType domain.SensorType, deviceID domain.DeviceID) int {
	key := string(sensorType) + "_" + string(deviceID)
	return m.sensorReadings[key]
}

func (m *MockMetrics) GetSensorErrors(sensorType domain.SensorType, deviceID domain.DeviceID) int {
	key := string(sensorType) + "_" + string(deviceID)
	return m.sensorErrors[key]
}
