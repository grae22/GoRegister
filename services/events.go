package services

import (
	"errors"
	"goregister/domain"
	"goregister/dto"
	"time"

	"github.com/google/uuid"
)

type EventsService struct {
	settings *SettingsService
	events   []*domain.EventRegister
}

func NewEventsService(settings *SettingsService) (*EventsService, error) {
	if settings == nil {
		return nil, errors.New("Nil settings service")
	}

	date1, _ := time.Parse("2006-01-02 15:04", "2026-07-11 08:00")
	date2, _ := time.Parse("2006-01-02 15:04", "2026-06-13 08:00")

	r1, _ := domain.NewEventRegister(
		uuid.New().String(),
		date1,
		"The Gutter",
		"neilw",
		settings.GetPaymentOptions())

	r2, _ := domain.NewEventRegister(
		uuid.New().String(),
		date2,
		"Rumdoodle",
		"graemeb",
		settings.GetPaymentOptions())

	e, _ := domain.NewEventRegisterEntry(
		uuid.NewString(),
		"Brad Inggs",
		"+27824659740",
		"CW79RZZN",
		"",
		map[string]int{
			"adult": 1,
		},
		7000,
		true,
		date1)
	r2.AddEntry(e)

	e, _ = domain.NewEventRegisterEntry(
		uuid.NewString(),
		"Warwick Hastie",
		"+27761983282",
		"BR27CXZN",
		"G567241",
		map[string]int{
			"rhino": 1,
		},
		0,
		true,
		date1)
	r2.AddEntry(e)

	e, _ = domain.NewEventRegisterEntry(
		uuid.NewString(),
		"Rory Nielson",
		"+27827874152",
		"CB94MJZN",
		"",
		map[string]int{
			"adult": 2,
			"child": 1,
		},
		17500,
		true,
		date1)
	r2.AddEntry(e)

	e, _ = domain.NewEventRegisterEntry(
		uuid.NewString(),
		"Graeme Bruschi",
		"+27723955929",
		"CV96MYZN",
		"G556769",
		map[string]int{
			"rhino": 1,
		},
		0,
		true,
		date1)
	r2.AddEntry(e)

	r2.TogglePaymentComplete()

	return &EventsService{
			events: []*domain.EventRegister{r1, r2},
		},
		nil
}

func (s *EventsService) GetEvents() []*domain.EventRegister {
	return s.events
}

func (s *EventsService) GetEvent(id string) *domain.EventRegister {
	for _, e := range s.GetEvents() {
		if e.IdempotencyId == id {
			return e
		}
	}

	return nil
}

func (s *EventsService) AddEvent(newEvent dto.AddEventDto) (*domain.EventRegister, error) {
	for _, e := range s.GetEvents() {
		if e.IdempotencyId == newEvent.IdempotencyId {
			return e, nil
		}
	}

	e, err := domain.NewEventRegister(
		newEvent.IdempotencyId,
		newEvent.Date,
		newEvent.Title,
		newEvent.OrganiserId,
		s.settings.GetPaymentOptions())

	if err != nil {
		return nil, err
	}

	s.events = append(s.events, e)

	// TODO: Persist.

	return e, nil
}

func (s *EventsService) UpdateEvent(newEvent dto.AddEventDto) error {
	e := s.GetEvent(newEvent.IdempotencyId)
	if e == nil {
		return errors.New("Event not found")
	}

	e.Date = newEvent.Date
	e.Title = newEvent.Title
	e.OrganiserId = newEvent.OrganiserId

	if e.AreNewEntriesAllowed != newEvent.AreNewEntriesAllowed {
		if newEvent.AreNewEntriesAllowed {
			e.UnblockEntries()
		} else {
			e.BlockEntries()
		}
	}

	if e.IsPaymentCompleted != newEvent.IsPaymentCompleted {
		e.TogglePaymentComplete()
	}

	// TODO: Persist.

	return nil
}

func (s *EventsService) AddRegisterEntry(newEntry dto.AddRegisterEntry) error {
	e := s.GetEvent(newEntry.EventId)
	if e == nil {
		return errors.New("Event not found")
	}

	for _, entry := range e.Entries {
		if entry.IdempotencyId == newEntry.IdempotencyId {
			return nil
		}
	}

	ne, err := domain.NewEventRegisterEntry(
		newEntry.IdempotencyId,
		newEntry.Name,
		newEntry.ContactNumber,
		newEntry.VehicleRegistration,
		newEntry.RhinoCard,
		newEntry.EntrantCountByPaymentType,
		newEntry.AmountDueInC,
		newEntry.IsConditionsAccepted,
		time.Now())

	if err != nil {
		return err
	}

	// TODO: Persist.

	return e.AddEntry(ne)
}
