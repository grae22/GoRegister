package web

import (
	"goregister/domain"
	"goregister/utils"
	"html/template"
	"net/http"
	"slices"
	"strings"

	"github.com/google/uuid"
)

type addRegisterEntryPageData struct {
	Layout         Layout
	EventId        string
	IdempotencyId  string
	IsUpdate       bool
	Entry          domain.EventRegisterEntry
	PaymentOptions []domain.PaymentOption
}

func (c *RegistersController) HandleAddRegisterEntry(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	requestCtx := utils.NewRequestContext(c.usersService, r)

	paymentOptionsById := c.settingsService.GetPaymentOptions()
	paymentOptions := []domain.PaymentOption{}

	for _, po := range paymentOptionsById {
		paymentOptions = append(paymentOptions, po)
	}

	slices.SortFunc(
		paymentOptions,
		func(a domain.PaymentOption, b domain.PaymentOption) int {
			return strings.Compare(a.Name, b.Name)
		})

	data := addRegisterEntryPageData{
		Layout:         NewLayout(true, *requestCtx),
		EventId:        r.URL.Query().Get("eventId"),
		IdempotencyId:  uuid.New().String(),
		Entry:          domain.EventRegisterEntry{},
		PaymentOptions: paymentOptions,
	}

	if !requestCtx.User.IsGuest {
		data.Entry.PersonName = requestCtx.User.Name
	}

	tmpl := template.Must(template.ParseFiles("html/layout.html", "html/registerEntry.html"))
	tmpl.ExecuteTemplate(w, "layout", data)
}
