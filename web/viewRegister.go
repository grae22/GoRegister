package web

import (
	"goregister/domain"
	"html/template"
	"net/http"
	"slices"
	"strings"
)

type registerPageData struct {
	Event        *domain.EventRegister
	PaymentTypes []string
}

func (c *RegistersController) HandleRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

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

	data := registerPageData{
		Event:        event,
		PaymentTypes: []string{},
	}

	for _, e := range event.Entries {
		for t := range e.EntrantCountByPaymentTypeId {
			if slices.Contains(data.PaymentTypes, t) {
				continue
			}

			data.PaymentTypes = append(data.PaymentTypes, t)
		}
	}

	slices.Sort(data.PaymentTypes)

	tmpl := template.Must(template.ParseFiles("html/layout.html", "html/register.html"))
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
