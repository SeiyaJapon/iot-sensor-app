package main

import (
	"github.com/SeiyaJapon/iot-sensor-app/cmd/app"
	"github.com/SeiyaJapon/iot-sensor-app/internal"
	"log"
	"net/http"
)

func main() {
	container := app.NewAppContainer()

	router := internal.NewRouter(container)

	log.Println(http.ListenAndServe(":8080", router))
}
