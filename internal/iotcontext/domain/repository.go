package domain

type SensorRepository interface {
	Save(sensor *Sensor) error
	FindByID(id SensorID) (*Sensor, error)
	FindAll() ([]*Sensor, error)
	Update(sensor *Sensor) error
}

type SensorReadingRepository interface {
	Save(reading SensorReading) error
	FindBySensor(sensorID SensorID, limit int) ([]SensorReading, error)
}
