package web

import (
	"errors"
	"goregister/services"
)

type RegistersController struct {
	eventsService   *services.EventsService
	settingsService *services.SettingsService
	usersService    *services.UsersService
}

func NewRegistersController(
	es *services.EventsService,
	ss *services.SettingsService,
	us *services.UsersService,
) (*RegistersController, error) {

	if es == nil {
		return nil, errors.New("Received nil events service")
	}

	if ss == nil {
		return nil, errors.New("Received nil settings service")
	}

	if us == nil {
		return nil, errors.New("Received nil users service")
	}

	return &RegistersController{
			eventsService:   es,
			settingsService: ss,
			usersService:    us,
		},
		nil
}
