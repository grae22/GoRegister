package web

import (
	"errors"
	"goregister/domain"
	"goregister/services"
	"goregister/utils"
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
	CurrentUser  domain.User
	Events       []*domain.EventRegister
	NameByUserId map[string]string
	UserCanAdd   bool
	UserCanEdit  bool
}

type eventDetailsPageData struct {
	CurrentUser   domain.User
	IdempotencyId string
	IsUpdate      bool
	Date          string
	Time          string
	Title         string
	Users         []domain.User
	OrganiserId   string
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
	requestCtx := utils.NewRequestContext(c.usersService, r)

	if !requestCtx.User.HasPermission(domain.PermissionViewAllEvents) {
		http.Redirect(
			w,
			r,
			"/login",
			http.StatusSeeOther)
		return
	}

	isGetAllEvents, isGetOneEvent, eventId := interpretGetEventUrl(r)

	if isGetAllEvents {
		handleGetAllEvents(w, c, requestCtx)
	} else if isGetOneEvent {
		handleGetOneEvent(w, c, eventId, requestCtx)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (c *EventsController) HandleAddEvent(w http.ResponseWriter, r *http.Request) {
	requestCtx := utils.NewRequestContext(c.usersService, r)

	if !requestCtx.User.HasPermission(domain.PermissionManageEvents) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	tmpl := template.Must(template.ParseFiles("html/layout.html", "html/eventDetails.html"))

	data := eventDetailsPageData{
		CurrentUser:   requestCtx.User,
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
	ctx *utils.RequestCtx,
) {
	if !ctx.User.HasPermission(domain.PermissionViewAllEvents) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	data := eventsPageData{
		CurrentUser:  ctx.User,
		Events:       c.eventsService.GetEvents(),
		NameByUserId: map[string]string{},
	}

	users := c.usersService.GetUsers()
	for _, u := range users {
		data.NameByUserId[u.Id] = u.Name
	}

	slices.SortFunc(
		data.Events,
		func(a *domain.EventRegister, b *domain.EventRegister) int {
			return b.Date.Compare(a.Date)
		})

	tmpl := template.Must(template.ParseFiles("html/layout.html", "html/events.html"))
	tmpl.ExecuteTemplate(w, "layout", data)
}

func handleGetOneEvent(
	w http.ResponseWriter,
	c *EventsController,
	eventId string,
	ctx *utils.RequestCtx,
) {
	if !ctx.User.HasPermission(domain.PermissionManageEvents) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	e := c.eventsService.GetEvent(eventId)
	if e == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	tmpl := template.Must(template.ParseFiles("html/layout.html", "html/eventDetails.html"))

	data := eventDetailsPageData{
		CurrentUser:   ctx.User,
		IdempotencyId: eventId,
		IsUpdate:      true,
		Date:          e.Date.Format("2006-01-02"),
		Time:          e.Date.Format("15:04"),
		Title:         e.Title,
		OrganiserId:   e.OrganiserId,
		Users:         c.usersService.GetUsers(),
	}

	if len(data.OrganiserId) == 0 {
		data.OrganiserId = ctx.User.Id
	}

	tmpl.ExecuteTemplate(w, "layout", data)
}
