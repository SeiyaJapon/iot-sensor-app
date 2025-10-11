package internal

import (
	"github.com/SeiyaJapon/iot-sensor-app/cmd/app"
	iot_http "github.com/SeiyaJapon/iot-sensor-app/internal/iotcontext/infrastructure/http"
	metrics_http "github.com/SeiyaJapon/iot-sensor-app/internal/metricscontext/infrastructure/http"
	"log"
	"net/http"
)

type Router struct {
	mux *http.ServeMux
}

func NewRouter(container *app.AppContainer) *Router {
	r := &Router{
		mux: http.NewServeMux(),
	}

	logMW := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Printf("%s %s", r.Method, r.URL.Path)
			next.ServeHTTP(w, r)
		})
	}

	deviceHandlers := iot_http.NewDeviceHandlers(*container.DeviceUC)
	r.mux.Handle("/devices", logMW(http.HandlerFunc(deviceHandlers.DevicesHandler)))

	sensorHandlers := iot_http.NewSensorHandlers(*container.SensorUC)
	r.mux.Handle("/sensors", logMW(http.HandlerFunc(sensorHandlers.SensorsHandler)))

	readingsHandlers := iot_http.NewReadingsHandler(*container.ReadingsUC)
	r.mux.HandleFunc("/readings", readingsHandlers.SensorReadingsHandler)

	simulatorHandlers := iot_http.NewSimulatorHandler(*container.SimulatorUC)
	r.mux.HandleFunc("/simulator/", simulatorHandlers.SimulatorsHandler)

	r.mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("OK"))
		if err != nil {
			return
		}
	})

	metricsHandler := metrics_http.NewMetricsHandler()
	metricsHandler.RegisterRoutes(r.mux)

	return r
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
}
