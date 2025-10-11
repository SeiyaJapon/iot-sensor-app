package persistence

import (
	"encoding/json"
	"errors"
	"github.com/SeiyaJapon/iot-sensor-app/internal/iotcontext/domain"
	"gorm.io/gorm"
)

type PostgresSensorRepository struct {
	db *DB
}

func NewPostgresSensorRepository(db *DB) domain.SensorRepository {
	return &PostgresSensorRepository{db: db}
}

func (r *PostgresSensorRepository) FindByID(id domain.SensorID) (*domain.Sensor, error) {
	var model SensorModel
	if err := r.db.conn.First(&model, "id = ?", string(id)).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrSensorNotFound
		}

		return nil, err
	}

	return unmarshalSensor(&model)
}

func (r *PostgresSensorRepository) FindAll() ([]*domain.Sensor, error) {
	var models []SensorModel
	if err := r.db.conn.Find(&models).Error; err != nil {
		return nil, err
	}

	var sensors []*domain.Sensor
	for _, model := range models {
		sensor, err := unmarshalSensor(&model)
		if err != nil {
			return nil, err
		}

		sensors = append(sensors, sensor)
	}

	return sensors, nil
}

func (r *PostgresSensorRepository) Save(sensor *domain.Sensor) error {
	model := SensorModel{
		ID:        string(sensor.ID),
		DeviceID:  string(sensor.DeviceID),
		Name:      sensor.Name,
		Type:      string(sensor.Type),
		Config:    marshalConfig(sensor.Config),
		CreatedAt: sensor.CreatedAt,
		UpdatedAt: sensor.UpdatedAt,
	}

	return r.db.conn.Create(&model).Error
}

func (r *PostgresSensorRepository) Update(sensor *domain.Sensor) error {
	model := SensorModel{
		ID:        string(sensor.ID),
		DeviceID:  string(sensor.DeviceID),
		Name:      sensor.Name,
		Type:      string(sensor.Type),
		Config:    marshalConfig(sensor.Config),
		CreatedAt: sensor.CreatedAt,
		UpdatedAt: sensor.UpdatedAt,
	}

	return r.db.conn.Save(&model).Error
}

func marshalConfig(config domain.SensorConfig) []byte {
	b, _ := json.Marshal(config)
	return b
}

func unmarshalSensor(model *SensorModel) (*domain.Sensor, error) {
	var config domain.SensorConfig
	if err := json.Unmarshal(model.Config, &config); err != nil {
		return nil, err
	}
	return &domain.Sensor{
		ID:        domain.SensorID(model.ID),
		DeviceID:  domain.DeviceID(model.DeviceID),
		Name:      model.Name,
		Type:      domain.SensorType(model.Type),
		Config:    config,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}, nil
}
