package application

import (
	"github.com/SeiyaJapon/iot-sensor-app/internal/iotcontext/domain"
	"time"
)

type DeviceUseCase struct {
	deviceRepo domain.DeviceRepository
}

func NewDeviceUseCase(deviceRepo domain.DeviceRepository) *DeviceUseCase {
	return &DeviceUseCase{
		deviceRepo: deviceRepo,
	}
}

func (uc *DeviceUseCase) CreateDevice(id domain.DeviceID, name string, typ string) (*domain.Device, error) {
	device, err := domain.NewDevice(id, name, typ)
	if err != nil {
		return nil, err
	}

	if err := uc.deviceRepo.Save(device); err != nil {
		return nil, err
	}

	return device, nil
}

func (uc *DeviceUseCase) GetDeviceByID(id domain.DeviceID) (*domain.Device, error) {
	return uc.deviceRepo.FindByID(id)
}

func (uc *DeviceUseCase) GetAllDevices() ([]*domain.Device, error) {
	return uc.deviceRepo.FindAll()
}

func (uc *DeviceUseCase) UpdateDevice(device *domain.Device) error {
	device.UpdatedAt = time.Now().UTC()

	return uc.deviceRepo.Update(device)
}
