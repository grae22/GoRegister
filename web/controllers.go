package web

import (
	"errors"
	"goregister/services"
)

type RegistersController struct {
	eventsService *services.EventsService
	usersService  *services.UsersService
}

func NewRegistersController(
	es *services.EventsService,
	us *services.UsersService,
) (*RegistersController, error) {

	if es == nil {
		return nil, errors.New("Received nil events service")
	}

	if us == nil {
		return nil, errors.New("Received nil users service")
	}

	return &RegistersController{
			eventsService: es,
			usersService:  us,
		},
		nil
}
