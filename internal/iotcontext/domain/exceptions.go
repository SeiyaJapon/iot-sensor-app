package domain

import "errors"

var ErrInvalidPaginationParams = errors.New("invalid pagination parameters")
var ErrSensorNotFound = errors.New("sensor not found")
var ErrInvalidAction = errors.New("invalid action")
var ErrDeviceNotFound = errors.New("device not found")
