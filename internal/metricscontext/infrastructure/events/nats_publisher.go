package events

import (
	"encoding/json"
	"github.com/SeiyaJapon/iot-sensor-app/internal/iotcontext/domain"
	"github.com/nats-io/nats.go"
)

type NatsPublisher struct {
	conn *nats.Conn
}

func NewNatsPublisher(url *string) (*NatsPublisher, error) {
	natsURL := nats.DefaultURL
	if url != nil {
		natsURL = *url
	}

	conn, err := nats.Connect(natsURL)
	if err != nil {
		return nil, err
	}

	return &NatsPublisher{conn: conn}, nil
}

func (np *NatsPublisher) Publish(event domain.IoTEvent) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return err
	}

	return np.conn.Publish(event.Type, payload)
}
