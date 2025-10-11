package domain

type SensorRepository interface {
	Save(sensor *Sensor) error
	FindByID(id SensorID) (*Sensor, error)
	FindAll() ([]*Sensor, error)
	Update(sensor *Sensor) error
}

type SensorReadingRepository interface {
	Save(reading *SensorReading) error
	FindBySensorID(sensorID SensorID, limit int) ([]SensorReading, error)
}

type DeviceRepository interface {
	Save(device *Device) error
	FindByID(id DeviceID) (Device, error)
	FindAll() ([]Device, error)
	Update(device *Device) error
}

type SimulatorRepository interface {
	Start(sensorID SensorID) error
	Stop(sensorID SensorID) error
	InjectError(sensorID SensorID) error
}
