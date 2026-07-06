package webapi

import (
	"errors"
	"goregister/dto"
	"goregister/services"
	"net/http"
	"strconv"
	"strings"
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
		EntrantCountByPaymentType: getEntrantCounts(r),
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

func getEntrantCounts(r *http.Request) map[string]int {
	const paymentPrefix string = "payment_"

	countsByPaymentTypeId := map[string]int{}

	for k, v := range r.Form {
		if !strings.HasPrefix(k, paymentPrefix) {
			continue
		}

		paymentTypeId := k[len(paymentPrefix):]

		if len(v) == 0 {
			continue
		}

		entrantCount, err := strconv.Atoi(v[0])
		if err != nil {
			continue
		}

		countsByPaymentTypeId[paymentTypeId] = entrantCount
	}

	return countsByPaymentTypeId
}
