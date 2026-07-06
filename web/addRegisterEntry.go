package web

import (
	"goregister/domain"
	"html/template"
	"maps"
	"net/http"
	"slices"
	"strings"

	"github.com/google/uuid"
)

type addRegisterEntryPageData struct {
	EventId       string
	IdempotencyId string
	IsUpdate      bool
	Entry         domain.EventRegisterEntry
	PaymentTypes  []domain.PaymentOption
}

func (c *RegistersController) HandleAddRegisterEntry(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	paymentOptionsById := c.settingsService.GetPaymentOptions()
	paymentOptions := []domain.PaymentOption{}

	for po := range maps.Values(paymentOptionsById) {
		paymentOptions = append(paymentOptions, po)
	}

	slices.SortFunc(
		paymentOptions,
		func(a domain.PaymentOption, b domain.PaymentOption) int {
			return strings.Compare(a.Name, b.Name)
		})

	data := addRegisterEntryPageData{
		EventId:       r.URL.Query().Get("eventId"),
		IdempotencyId: uuid.New().String(),
		Entry:         domain.EventRegisterEntry{},
		PaymentTypes:  paymentOptions,
	}

	tmpl := template.Must(template.ParseFiles("html/layout.html", "html/registerEntry.html"))
	tmpl.ExecuteTemplate(w, "layout", data)
}
