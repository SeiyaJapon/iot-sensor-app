package persistence

import "github.com/SeiyaJapon/iot-sensor-app/internal/iotcontext/domain"

type PostgresDeviceRepository struct {
	db *DB
}

func NewPostgresDeviceRepository(db *DB) domain.DeviceRepository {
	return &PostgresDeviceRepository{db: db}
}

func (r *PostgresDeviceRepository) Save(device *domain.Device) error {
	model := DeviceModel{
		ID:        string(device.ID),
		Name:      device.Name,
		Type:      device.Type,
		CreatedAt: device.CreatedAt,
		UpdatedAt: device.UpdatedAt,
	}

	return r.db.conn.Create(&model).Error
}

func (r *PostgresDeviceRepository) FindByID(id domain.DeviceID) (domain.Device, error) {
	var model DeviceModel
	if err := r.db.conn.First(&model, "id = ?", string(id)).Error; err != nil {
		return domain.Device{}, err
	}

	return domain.Device{
		ID:        domain.DeviceID(model.ID),
		Name:      model.Name,
		Type:      model.Type,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}, nil
}

func (r *PostgresDeviceRepository) FindAll() ([]domain.Device, error) {
	var models []DeviceModel
	if err := r.db.conn.Find(&models).Error; err != nil {
		return nil, err
	}

	var devices []domain.Device
	for _, model := range models {
		device := domain.Device{
			ID:        domain.DeviceID(model.ID),
			Name:      model.Name,
			Type:      model.Type,
			CreatedAt: model.CreatedAt,
			UpdatedAt: model.UpdatedAt,
		}
		devices = append(devices, device)
	}

	return devices, nil
}

func (r *PostgresDeviceRepository) Update(device *domain.Device) error {
	model := DeviceModel{
		ID:        string(device.ID),
		Name:      device.Name,
		Type:      device.Type,
		CreatedAt: device.CreatedAt,
		UpdatedAt: device.UpdatedAt,
	}

	return r.db.conn.Save(&model).Error
}
