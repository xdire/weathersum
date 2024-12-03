package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/xdire/weathersum/handlers"
	"github.com/xdire/weathersum/service"
)

func main() {

	var serverOnPort int
	flag.IntVar(&serverOnPort, "port", 8000, "start server on port")
	flag.Parse()

	// Set up route
	router := mux.NewRouter()
	router.HandleFunc("/", handlers.APIHome)
	router.HandleFunc("/v1/weather", handlers.SimplifiedWeather)

	// Start server
	err := service.StartWeatherService(router, serverOnPort)
	if err != nil {
		fmt.Println(err)
	}
}
