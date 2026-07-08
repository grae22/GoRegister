package services

import (
	"errors"
	"goregister/domain"
	"goregister/dto"
	"time"
)

type EventsService struct {
	settings *SettingsService
	events   []*domain.EventRegister
}

func NewEventsService(settings *SettingsService) (*EventsService, error) {
	if settings == nil {
		return nil, errors.New("Nil settings service")
	}

	r, _ := domain.NewEventRegister(
		"123",
		time.Now(),
		"Test",
		"graemeb",
		map[string]domain.PaymentOption{
			"rhino": {
				Id:                          "rhino",
				Name:                        "Rhino card",
				ValueInC:                    0,
				DisplayValueInR:             "0",
				IsEnabledForNewTransactions: true,
			},
			"adult": {
				Id:                          "adult",
				Name:                        "Adult",
				ValueInC:                    7000,
				DisplayValueInR:             "70",
				IsEnabledForNewTransactions: true,
			},
		})

	e1, _ := domain.NewEventRegisterEntry(
		"123",
		"Person 1",
		"+27...",
		"ND...",
		"G...",
		map[string]int{
			"rhino": 2,
			"adult": 1,
		},
		7000,
		true,
		time.Now())

	e2, _ := domain.NewEventRegisterEntry(
		"456",
		"Person 2",
		"+27...",
		"NP...",
		"G...",
		map[string]int{
			"rhino": 1,
			"adult": 2,
			"child": 1,
		},
		14000,
		true,
		time.Now())

	r.AddEntry(e1)
	r.AddEntry(e2)

	return &EventsService{
			events: []*domain.EventRegister{r},
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

	e.Entries = append(e.Entries, ne)

	return nil
}
