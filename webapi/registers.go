package webapi

import (
	"errors"
	"fmt"
	"goregister/domain"
	"goregister/dto"
	"goregister/services"
	"net/http"
	"strconv"
	"strings"
)

type RegistersApi struct {
	events   *services.EventsService
	settings *services.SettingsService
	users    *services.UsersService
}

func NewRegistersApi(
	events *services.EventsService,
	settings *services.SettingsService,
	users *services.UsersService,
) (*RegistersApi, error) {

	if events == nil {
		return nil, errors.New("Received nil events service")
	}

	if settings == nil {
		return nil, errors.New("Received nil settings service")
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
		RhinoCard:                 r.FormValue("rhinoCard"),
		EntrantCountByPaymentType: getEntrantCounts(r),
		IsConditionsAccepted:      r.FormValue("conditions") == "yes",
	}

	amountDueInC, err := calculateAmountDueInC(
		api.settings.GetPaymentOptions(),
		e.EntrantCountByPaymentType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	e.AmountDueInC = amountDueInC

	if !isUpdate {
		err = api.events.AddRegisterEntry(e)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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

func calculateAmountDueInC(
	paymentOptionsById map[string]domain.PaymentOption,
	entrantCountsByPaymentId map[string]int,
) (int, error) {
	totalInC := 0

	for id, count := range entrantCountsByPaymentId {
		opt, ok := paymentOptionsById[id]
		if !ok {
			return 0, fmt.Errorf("Unknown payment option id: %s", id)
		}

		totalInC += opt.ValueInC * count
	}

	return totalInC, nil
}
