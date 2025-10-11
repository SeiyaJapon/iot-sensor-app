package domain

import (
	"testing"
)

func TestThresholds_Exceeds(t *testing.T) {
	tests := []struct {
		name          string
		thresholds    Thresholds
		value         float64
		expectExceeds bool
		expectReason  string
	}{
		{
			name: "no thresholds set",
			thresholds: Thresholds{
				Min: nil,
				Max: nil,
			},
			value:         50.0,
			expectExceeds: false,
			expectReason:  "",
		},
		{
			name: "value within range",
			thresholds: Thresholds{
				Min: floatPtr(10.0),
				Max: floatPtr(90.0),
			},
			value:         50.0,
			expectExceeds: false,
			expectReason:  "",
		},
		{
			name: "value above max",
			thresholds: Thresholds{
				Min: floatPtr(10.0),
				Max: floatPtr(90.0),
			},
			value:         95.0,
			expectExceeds: true,
			expectReason:  "above_max",
		},
		{
			name: "value below min",
			thresholds: Thresholds{
				Min: floatPtr(10.0),
				Max: floatPtr(90.0),
			},
			value:         5.0,
			expectExceeds: true,
			expectReason:  "below_min",
		},
		{
			name: "value equals min",
			thresholds: Thresholds{
				Min: floatPtr(10.0),
				Max: floatPtr(90.0),
			},
			value:         10.0,
			expectExceeds: false,
			expectReason:  "",
		},
		{
			name: "value equals max",
			thresholds: Thresholds{
				Min: floatPtr(10.0),
				Max: floatPtr(90.0),
			},
			value:         90.0,
			expectExceeds: false,
			expectReason:  "",
		},
		{
			name: "only min threshold set - value above",
			thresholds: Thresholds{
				Min: floatPtr(10.0),
				Max: nil,
			},
			value:         50.0,
			expectExceeds: false,
			expectReason:  "",
		},
		{
			name: "only min threshold set - value below",
			thresholds: Thresholds{
				Min: floatPtr(10.0),
				Max: nil,
			},
			value:         5.0,
			expectExceeds: true,
			expectReason:  "below_min",
		},
		{
			name: "only max threshold set - value below",
			thresholds: Thresholds{
				Min: nil,
				Max: floatPtr(90.0),
			},
			value:         50.0,
			expectExceeds: false,
			expectReason:  "",
		},
		{
			name: "only max threshold set - value above",
			thresholds: Thresholds{
				Min: nil,
				Max: floatPtr(90.0),
			},
			value:         95.0,
			expectExceeds: true,
			expectReason:  "above_max",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exceeds, reason := tt.thresholds.Exceeds(tt.value)

			if exceeds != tt.expectExceeds {
				t.Errorf("expected exceeds %t, got %t", tt.expectExceeds, exceeds)
			}

			if reason != tt.expectReason {
				t.Errorf("expected reason %s, got %s", tt.expectReason, reason)
			}
		})
	}
}

func floatPtr(f float64) *float64 {
	return &f
}
