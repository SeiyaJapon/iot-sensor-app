package persistence

import (
	"github.com/SeiyaJapon/iot-sensor-app/internal/iotcontext/domain"
	
)

func TestPrometheusMetricsImpl_IncSensorReading(t *testing.T) {
	metrics := NewPrometheusMetrics()

	metrics.IncSensorReading(domain.Temperature, "device-123")
	metrics.IncSensorReading(domain.Humidity, "device-123")
	metrics.IncSensorReading(domain.Temperature, "device-456")

	if metrics == nil {
		t.Error("expected metrics to be created")
	}
}
