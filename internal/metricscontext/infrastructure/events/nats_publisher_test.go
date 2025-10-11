package events

import (
	"github.com/SeiyaJapon/iot-sensor-app/internal/iotcontext/domain"
	
)

func TestNewNatsPublisher(t *testing.T) {
	tests := []struct {
		name    string
		url     *string
		wantErr bool
	}{
		{
			name:    "valid URL",
			url:     stringPtr("nats://localhost:4222"),
			wantErr: false,
		},
		{
			name:    "nil URL uses default",
			url:     nil,
			wantErr: false,
		},
		{
			name:    "invalid URL",
			url:     stringPtr("invalid://url"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			publisher, err := NewNatsPublisher(tt.url)

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if publisher == nil {
				t.Error("expected publisher to be created")
			}

			if publisher.conn == nil {
				t.Error("expected connection to be initialized")
			}
		})
	}
}

func TestNatsPublisher_Publish(t *testing.T) {
	publisher, err := NewNatsPublisher(stringPtr("nats://localhost:4222"))
	if err != nil {
		t.Skipf("skipping test due to NATS connection error: %v", err)
	}

	event := domain.IoTEvent{
		Type:      "test.event",
		Payload:   map[string]interface{}{"test": "data"},
		Timestamp: domain.IoTEvent{}.Timestamp,
	}

	err = publisher.Publish(event)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func stringPtr(s string) *string {
	return &s
}
