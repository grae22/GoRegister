package web

import (
	"errors"
	"goregister/domain"
	"goregister/services"
	"html/template"
	"net/http"
	"slices"
	"strings"

	"github.com/google/uuid"
)

type EventsController struct {
	eventsService *services.EventsService
	usersService  *services.UsersService
}

type eventsPageData struct {
	Events []*domain.EventRegister
}

type eventDetailsPageData struct {
	IdempotencyId string
	IsUpdate      bool
	Date          string
	Time          string
	Title         string
	Users         []domain.User
	OrganiserId   string
	CurrentUserId string
}

func NewEventsController(
	es *services.EventsService,
	us *services.UsersService,
) (*EventsController, error) {

	if es == nil {
		return nil, errors.New("Received nil events service")
	}

	if us == nil {
		return nil, errors.New("Received nil users service")
	}

	return &EventsController{
			eventsService: es,
			usersService:  us,
		},
		nil
}

func (c *EventsController) HandleEvents(w http.ResponseWriter, r *http.Request) {
	isGetAllEvents, isGetOneEvent, eventId := interpretGetEventUrl(r)

	if isGetAllEvents {
		handleGetAllEvents(w, c)
	} else if isGetOneEvent {
		handleGetOneEvent(w, c, eventId)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (c *EventsController) HandleAddEvent(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("html/layout.html", "html/eventDetails.html"))

	data := eventDetailsPageData{
		IdempotencyId: uuid.New().String(),
		Users:         c.usersService.GetUsers(),
	}

	tmpl.ExecuteTemplate(w, "layout", data)
}

func interpretGetEventUrl(r *http.Request) (isGetAllEvents bool, isGetOneEvent bool, eventId string) {
	parts := strings.Split(strings.ToLower(r.URL.Path), "/")
	l := len(parts)

	if l > 1 && parts[l-2] == "events" {
		eventId = parts[l-1]
		isGetOneEvent = len(eventId) > 0
		isGetAllEvents = !isGetOneEvent
	}

	return
}

func handleGetAllEvents(
	w http.ResponseWriter,
	c *EventsController,
) {
	tmpl := template.Must(template.ParseFiles("html/layout.html", "html/events.html"))

	data := eventsPageData{
		Events: c.eventsService.GetEvents(),
	}

	slices.SortFunc(
		data.Events,
		func(a *domain.EventRegister, b *domain.EventRegister) int {
			return b.Date.Compare(a.Date)
		})

	tmpl.ExecuteTemplate(w, "layout", data)
}

func handleGetOneEvent(
	w http.ResponseWriter,
	c *EventsController,
	eventId string,
) {
	e := c.eventsService.GetEvent(eventId)
	if e == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	tmpl := template.Must(template.ParseFiles("html/layout.html", "html/eventDetails.html"))

	data := eventDetailsPageData{
		IdempotencyId: eventId,
		IsUpdate:      true,
		Date:          e.Date.Format("2006-01-02"),
		Time:          e.Date.Format("15:04"),
		Title:         e.Title,
		OrganiserId:   e.OrganiserId,
		CurrentUserId: e.OrganiserId,
		Users:         c.usersService.GetUsers(),
	}

	tmpl.ExecuteTemplate(w, "layout", data)
}
