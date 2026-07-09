package webapi

import (
	"errors"
	"goregister/domain"
	"goregister/dto"
	"goregister/services"
	"goregister/utils"
	"net/http"
	"time"
)

type EventsApi struct {
	events *services.EventsService
	users  *services.UsersService
}

func NewEventsApi(
	events *services.EventsService,
	users *services.UsersService,
) (*EventsApi, error) {

	if events == nil {
		return nil, errors.New("Received nil events service")
	}

	if users == nil {
		return nil, errors.New("Received nil users service")
	}

	return &EventsApi{
			events: events,
			users:  users,
		},
		nil
}

func (a *EventsApi) HandleEvents(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		handleAddEvent(w, r, a)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func handleAddEvent(
	w http.ResponseWriter,
	r *http.Request,
	api *EventsApi,
) {
	requestCtx := utils.NewRequestContext(api.users, r)

	if !requestCtx.User.HasPermission(domain.PermissionManageEvents) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	isUpdate := r.FormValue("isUpdate") == "true"

	date, err := time.Parse("2006-01-02 15:04", r.FormValue("date")+" "+r.FormValue("time"))
	if err != nil {
		http.Error(w, "Invalid date", http.StatusBadRequest)
		return
	}

	e := dto.AddEventDto{
		IdempotencyId:        r.FormValue("idempotencyId"),
		Date:                 date,
		Title:                r.FormValue("title"),
		OrganiserId:          r.FormValue("organiserId"),
		AreNewEntriesAllowed: r.FormValue("allowNewEntries") == "allowed",
		IsPaymentCompleted:   r.FormValue("paymentCompleted") == "paid",
	}

	if !isUpdate {
		_, err = api.events.AddEvent(e)
	} else {
		err = api.events.UpdateEvent(e)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	http.Redirect(
		w,
		r,
		"/events",
		http.StatusSeeOther)
}
