package http

import (
	"net/http"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) CheckHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	write, err := w.Write([]byte(`{"status":"ok"}`))
	if err != nil {
		return
	}

	if write == 0 {
		return
	}
}
