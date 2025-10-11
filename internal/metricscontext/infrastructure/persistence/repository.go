package persistence

import (
	"github.com/SeiyaJapon/iot-sensor-app/internal/iotcontext/domain"
	"github.com/prometheus/client_golang/prometheus"
)

type PrometheusMetricsImpl struct {
	readingsTotal *prometheus.CounterVec
	errorsTotal   *prometheus.CounterVec
}

func NewPrometheusMetrics() *PrometheusMetricsImpl {
	readings := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "sensor_readings_total",
			Help: "Total number of sensor readings",
		},
		[]string{"sensor_type", "device_id"},
	)

	errors := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "sensor_errors_total",
			Help: "Total number of sensor errors",
		},
		[]string{"sensor_type", "device_id"},
	)

	prometheus.MustRegister(readings, errors)

	return &PrometheusMetricsImpl{
		readingsTotal: readings,
		errorsTotal:   errors,
	}
}

func (pm *PrometheusMetricsImpl) IncSensorReading(sensorType domain.SensorType, deviceID domain.DeviceID) {
	pm.readingsTotal.WithLabelValues(string(sensorType), string(deviceID)).Inc()
}

func (pm *PrometheusMetricsImpl) IncSensorError(sensorType domain.SensorType, deviceID domain.DeviceID) {
	pm.errorsTotal.WithLabelValues(string(sensorType), string(deviceID)).Inc()
}
