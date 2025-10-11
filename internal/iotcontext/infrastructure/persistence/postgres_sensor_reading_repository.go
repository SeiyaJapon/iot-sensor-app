package persistence

import (
	"encoding/json"
	"github.com/SeiyaJapon/iot-sensor-app/internal/iotcontext/domain"
)

type PostgresSensorReadingRepository struct {
	db *DB
}

func NewPostgresSensorReadingRepository(db *DB) domain.SensorReadingRepository {
	return &PostgresSensorReadingRepository{db: db}
}

func (r *PostgresSensorReadingRepository) Save(reading *domain.SensorReading) error {
	model := &SensorReadingModel{
		ID:        reading.ID,
		SensorID:  string(reading.SensorID),
		DeviceID:  string(reading.DeviceID),
		Type:      string(reading.Type),
		Value:     reading.Value,
		Unit:      reading.Unit,
		Timestamp: reading.Timestamp,
		Meta:      marshalMeta(reading.Meta),
	}

	return r.db.conn.Create(model).Error
}

func (r *PostgresSensorReadingRepository) FindBySensorID(sensorID domain.SensorID, limit int) ([]domain.SensorReading, error) {
	var models []SensorReadingModel
	query := r.db.conn.Where("sensor_id = ?", string(sensorID))

	if limit > 0 {
		query = query.Limit(limit)
	}

	if err := query.Find(&models).Error; err != nil {
		return nil, err
	}

	readings := unmarshalMeta(models)

	return readings, nil
}

func marshalMeta(meta map[string]interface{}) []byte {
	b, _ := json.Marshal(meta)

	return b
}

func unmarshalMeta(models []SensorReadingModel) []domain.SensorReading {
	var readings []domain.SensorReading
	for _, model := range models {
		var meta map[string]interface{}
		_ = json.Unmarshal(model.Meta, &meta)

		reading := domain.SensorReading{
			ID:        model.ID,
			SensorID:  domain.SensorID(model.SensorID),
			DeviceID:  domain.DeviceID(model.DeviceID),
			Type:      domain.SensorType(model.Type),
			Value:     model.Value,
			Unit:      model.Unit,
			Timestamp: model.Timestamp,
			Meta:      meta,
		}

		readings = append(readings, reading)
	}

	return readings
}
