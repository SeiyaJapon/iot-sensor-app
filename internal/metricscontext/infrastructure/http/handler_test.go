package http

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMetricsHandler_RegisterRoutes(t *testing.T) {
	handler := NewMetricsHandler()
	mux := http.NewServeMux()

	handler.RegisterRoutes(mux)

	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestNewMetricsHandler(t *testing.T) {
	handler := NewMetricsHandler()

	if handler == nil {
		t.Error("expected handler to be created")
	}
}
