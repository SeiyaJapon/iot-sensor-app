package http

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

type MetricsHandler struct{}

func NewMetricsHandler() *MetricsHandler {
	return &MetricsHandler{}
}

func (mh *MetricsHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.Handle("/metrics", promhttp.Handler())
}
