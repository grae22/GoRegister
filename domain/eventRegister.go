package domain

import (
	"errors"
	"strings"
	"time"
)

type EventRegister struct {
	IdempotencyId string
	Date          time.Time
	Title         string
	OrganiserId   string
	Entries       []*EventRegisterEntry
}

func NewEventRegister(
	idempotencyId string,
	date time.Time,
	title string,
	organiserId string,
) (*EventRegister, error) {

	idempotencyId = strings.TrimSpace(idempotencyId)
	title = strings.TrimSpace(title)
	organiserId = strings.TrimSpace(organiserId)

	if idempotencyId == "" {
		return nil, errors.New("Invalid idempotency id")
	}

	if title == "" {
		return nil, errors.New("Invalid title")
	}

	if organiserId == "" {
		return nil, errors.New("Invalid title")
	}

	r := EventRegister{
		IdempotencyId: idempotencyId,
		Date:          date,
		Title:         title,
		OrganiserId:   organiserId,
		Entries:       []*EventRegisterEntry{},
	}

	return &r, nil
}

func (er *EventRegister) AddEntry(entry *EventRegisterEntry) {
	for _, e := range er.Entries {
		if e.IdempotencyId == entry.IdempotencyId {
			return
		}
	}

	er.Entries = append(er.Entries, entry)
}
