package web

import (
	"goregister/domain"
	"goregister/utils"
	"html/template"
	"net/http"
	"slices"
	"strings"
)

type registerPageData struct {
	Layout               Layout
	Link                 string
	Event                *domain.EventRegister
	PaymentOptionsById   map[string]domain.PaymentOption
	SortedPaymentOptions []domain.PaymentOption
	NameByUserId         map[string]string
	TotalDueInC          int
}

func (c *RegistersController) HandleRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	requestCtx := utils.NewRequestContext(c.usersService, r)

	eventId := getEventIdFromUrl(r)
	if len(eventId) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	event := c.eventsService.GetEvent(eventId)
	if event == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	allPaymentOptions := c.settingsService.GetPaymentOptions()

	data := registerPageData{
		Layout:               NewLayout(true, *requestCtx),
		Link:                 c.settingsService.BaseUrl + "/registers/" + eventId,
		Event:                event,
		PaymentOptionsById:   map[string]domain.PaymentOption{},
		SortedPaymentOptions: []domain.PaymentOption{},
		NameByUserId:         map[string]string{},
	}

	users := c.usersService.GetUsers()
	for _, u := range users {
		data.NameByUserId[u.Id] = u.Name
	}

	for _, e := range event.Entries {
		for t := range e.EntrantCountByPaymentTypeId {
			if _, ok := data.PaymentOptionsById[t]; ok {
				continue
			}

			data.PaymentOptionsById[t] = allPaymentOptions[t]
			data.SortedPaymentOptions = append(data.SortedPaymentOptions, allPaymentOptions[t])
		}

		data.TotalDueInC += e.AmountDueInC
	}

	slices.SortFunc(
		data.SortedPaymentOptions,
		func(a domain.PaymentOption, b domain.PaymentOption) int {
			return strings.Compare(a.Name, b.Name)
		})

	tmpl := template.
		Must(template.New("register").
			Funcs(template.FuncMap{"centsToRandsStr": utils.CentsToRandsStr}).
			ParseFiles("html/layout.html", "html/register.html"))

	tmpl.ExecuteTemplate(w, "layout", data)
}

func getEventIdFromUrl(r *http.Request) string {
	var registerId string

	parts := strings.Split(strings.ToLower(r.URL.Path), "/")
	l := len(parts)

	if l > 1 && parts[l-2] == "registers" {
		registerId = parts[l-1]
	}

	return registerId
}
