package domain

type Thresholds struct {
	Min *float64 `json:"min"`
	Max *float64 `json:"max"`
}

func (t Thresholds) Exceeds(value float64) (bool, string) {
	if t.Max != nil && value > *t.Max {
		return true, "above_max"
	}

	if t.Min != nil && value < *t.Min {
		return true, "below_min"
	}

	return false, ""
}
