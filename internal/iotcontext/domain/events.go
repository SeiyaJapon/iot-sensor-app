package domain

import "time"

type EventPublisher interface {
	Publish(event IoTEvent) error
}

type IoTEvent struct {
	Type      string    `json:"type"`
	Timestamp time.Time `json:"timestamp"`
	Payload   any       `json:"payload"`
}

type SensorCreatedEvent struct {
	SensorID SensorID   `json:"sensor_id"`
	DeviceID DeviceID   `json:"device_id"`
	Type     SensorType `json:"type"`
	Name     string     `json:"name"`
}

type SensorConfigUpdatedEvent struct {
	SensorID SensorID     `json:"sensor_id"`
	Config   SensorConfig `json:"config"`
}

type SensorReadingPublishedEvent struct {
	SensorID SensorID `json:"sensor_id"`
	Reading  string   `json:"reading"`
}

type SensorReadingErrorEvent struct {
	SensorID SensorID `json:"sensor_id"`
	Type     string   `json:"type"`
}

func (e *SensorCreatedEvent) ToDomainEvent() IoTEvent {
	return IoTEvent{
		Type:      "sensor.created",
		Timestamp: time.Now().UTC(),
		Payload:   e,
	}
}

func (e *SensorConfigUpdatedEvent) ToDomainEvent() IoTEvent {
	return IoTEvent{
		Type:      "sensor.config.updated",
		Timestamp: time.Now().UTC(),
		Payload:   e,
	}
}

func (e *SensorReadingPublishedEvent) ToDomainEvent() IoTEvent {
	return IoTEvent{
		Type:      "sensor.reading.published",
		Timestamp: time.Now().UTC(),
		Payload:   e,
	}
}

func (e *SensorReadingErrorEvent) ToDomainEvent() IoTEvent {
	return IoTEvent{
		Type:      "sensor.reading.error",
		Timestamp: time.Now().UTC(),
		Payload:   e,
	}
}
