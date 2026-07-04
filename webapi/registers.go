package webapi

import (
	"errors"
	"goregister/dto"
	"goregister/services"
	"net/http"
)

type RegistersApi struct {
	events *services.EventsService
	users  *services.UsersService
}

func NewRegistersApi(
	events *services.EventsService,
	users *services.UsersService,
) (*RegistersApi, error) {

	if events == nil {
		return nil, errors.New("Received nil events service")
	}

	if users == nil {
		return nil, errors.New("Received nil users service")
	}

	return &RegistersApi{
			events: events,
			users:  users,
		},
		nil
}

func (a *RegistersApi) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		handleAddRegisterEntry(w, r, a)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func handleAddRegisterEntry(
	w http.ResponseWriter,
	r *http.Request,
	api *RegistersApi,
) {
	isUpdate := r.FormValue("isUpdate") == "true"

	e := dto.AddRegisterEntry{
		IdempotencyId:             r.FormValue("idempotencyId"),
		EventId:                   r.FormValue("eventId"),
		Name:                      r.FormValue("name"),
		ContactNumber:             r.FormValue("contact"),
		VehicleRegistration:       r.FormValue("vehicleReg"),
		EntrantCountByPaymentType: map[string]int{},
	}

	var err error

	if !isUpdate {
		err = api.events.AddRegisterEntry(e)
	} else {
		//err = api.events.UpdateEvent(e)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	http.Redirect(
		w,
		r,
		"/registers/"+e.EventId,
		http.StatusSeeOther)
}
