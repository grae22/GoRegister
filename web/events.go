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
	settingsService *services.SettingsService
	eventsService   *services.EventsService
	usersService    *services.UsersService
}

type eventsPageData struct {
	Layout                Layout
	Events                []*domain.EventRegister
	AmountDueInCByEventId map[string]int
	NameByUserId          map[string]string
	UserCanAdd            bool
	UserCanEdit           bool
	TotalOwedInC          int
}

type eventDetailsPageData struct {
	Layout               Layout
	IdempotencyId        string
	IsUpdate             bool
	Date                 string
	Time                 string
	Title                string
	Users                []domain.User
	OrganiserId          string
	AreNewEntriesAllowed bool
	IsPaymentCompleted   bool
}

func NewEventsController(
	ss *services.SettingsService,
	es *services.EventsService,
	us *services.UsersService,
) (*EventsController, error) {

	if ss == nil {
		return nil, errors.New("Received nil settings service")
	}

	if es == nil {
		return nil, errors.New("Received nil events service")
	}

	if us == nil {
		return nil, errors.New("Received nil users service")
	}

	return &EventsController{
			settingsService: ss,
			eventsService:   es,
			usersService:    us,
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
		handleGetAllEvents(w, c, *requestCtx)
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
		Layout:               NewLayout(true, *requestCtx),
		IdempotencyId:        uuid.New().String(),
		Users:                c.usersService.GetUsers(),
		AreNewEntriesAllowed: true,
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
	ctx utils.RequestCtx,
) {
	if !ctx.User.HasPermission(domain.PermissionViewAllEvents) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	data := eventsPageData{
		Layout:                NewLayout(true, ctx),
		Events:                c.eventsService.GetEvents(),
		AmountDueInCByEventId: map[string]int{},
		NameByUserId:          map[string]string{},
		UserCanAdd:            ctx.User.HasPermission(domain.PermissionManageEvents),
		UserCanEdit:           ctx.User.HasPermission(domain.PermissionManageEvents),
	}

	users := c.usersService.GetUsers()
	for _, u := range users {
		data.NameByUserId[u.Id] = u.Name
	}

	paymentOptions := c.settingsService.GetPaymentOptions()

	for _, e := range data.Events {
		dueInC := e.CalculateAmountDueInC(paymentOptions)

		data.AmountDueInCByEventId[e.IdempotencyId] = dueInC

		if !e.IsPaymentCompleted {
			data.TotalOwedInC += dueInC
		}
	}

	slices.SortFunc(
		data.Events,
		func(a *domain.EventRegister, b *domain.EventRegister) int {
			return b.Date.Compare(a.Date)
		})

	tmpl := template.
		Must(template.New("events").
			Funcs(template.FuncMap{"centsToRandsStr": utils.CentsToRandsStr}).
			ParseFiles("html/layout.html", "html/events.html"))

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
		Layout:               NewLayout(true, *ctx),
		IdempotencyId:        eventId,
		IsUpdate:             true,
		Date:                 e.Date.Format("2006-01-02"),
		Time:                 e.Date.Format("15:04"),
		Title:                e.Title,
		OrganiserId:          e.OrganiserId,
		AreNewEntriesAllowed: e.AreNewEntriesAllowed,
		IsPaymentCompleted:   e.IsPaymentCompleted,
		Users:                c.usersService.GetUsers(),
	}

	if len(data.OrganiserId) == 0 {
		data.OrganiserId = ctx.User.Id
	}

	tmpl.ExecuteTemplate(w, "layout", data)
}
