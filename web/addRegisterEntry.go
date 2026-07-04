package web

import (
	"goregister/domain"
	"html/template"
	"net/http"

	"github.com/google/uuid"
)

type addRegisterEntryPageData struct {
	EventId       string
	IdempotencyId string
	IsUpdate      bool
	Entry         domain.EventRegisterEntry
	PaymentTypes  []string
}

func (c *RegistersController) HandleAddRegisterEntry(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	data := addRegisterEntryPageData{
		EventId:       r.URL.Query().Get("eventId"),
		IdempotencyId: uuid.New().String(),
		Entry:         domain.EventRegisterEntry{},
		PaymentTypes:  []string{},
	}

	tmpl := template.Must(template.ParseFiles("html/layout.html", "html/registerEntry.html"))
	tmpl.ExecuteTemplate(w, "layout", data)
}
