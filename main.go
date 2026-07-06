package main

import (
	"goregister/services"
	"goregister/web"
	"goregister/webapi"
	"log"
	"net/http"
)

func main() {
	settingsService := services.NewSettingsService()

	eventsService, err := services.NewEventsService(settingsService)
	if err != nil {
		log.Fatalf("Failed to create events service: %v", err)
		return
	}

	usersService := services.NewUsersService()

	eventsController, err := web.NewEventsController(eventsService, usersService)
	if err != nil {
		log.Fatalf("Failed to create events controller: %v", err)
		return
	}

	eventsApi, err := webapi.NewEventsApi(eventsService, usersService)
	if err != nil {
		log.Fatalf("Failed to create events api: %v", err)
		return
	}

	registersController, err := web.NewRegistersController(
		eventsService,
		settingsService,
		usersService)
	if err != nil {
		log.Fatalf("Failed to create registers controller: %v", err)
		return
	}

	registersApi, err := webapi.NewRegistersApi(
		eventsService,
		settingsService,
		usersService)
	if err != nil {
		log.Fatalf("Failed to create registers api: %v", err)
		return
	}

	http.HandleFunc("/events/", eventsController.HandleEvents)
	http.HandleFunc("/addEvent", eventsController.HandleAddEvent)
	http.HandleFunc("/api/events", eventsApi.HandleEvents)

	http.HandleFunc("/registers/", registersController.HandleRegister)
	http.HandleFunc("/addRegisterEntry", registersController.HandleAddRegisterEntry)
	http.HandleFunc("/api/registers", registersApi.Handle)

	http.ListenAndServe(":10222", nil)
}
