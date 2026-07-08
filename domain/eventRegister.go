package domain

import (
	"errors"
	"strings"
	"time"
)

type EventRegister struct {
	IdempotencyId      string
	Date               time.Time
	Title              string
	OrganiserId        string
	PaymentOptionsById map[string]PaymentOption
	Entries            []*EventRegisterEntry
	IsPaymentCompleted bool
}

func NewEventRegister(
	idempotencyId string,
	date time.Time,
	title string,
	organiserId string,
	paymentOptionsById map[string]PaymentOption,
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

	if len(paymentOptionsById) == 0 {
		return nil, errors.New("No payment options supplied")
	}

	r := EventRegister{
		IdempotencyId:      idempotencyId,
		Date:               date,
		Title:              title,
		OrganiserId:        organiserId,
		Entries:            []*EventRegisterEntry{},
		PaymentOptionsById: paymentOptionsById,
	}

	return &r, nil
}

func (r *EventRegister) AddEntry(entry *EventRegisterEntry) {
	for _, e := range r.Entries {
		if e.IdempotencyId == entry.IdempotencyId {
			return
		}
	}

	r.Entries = append(r.Entries, entry)
}

func (r *EventRegister) TogglePaymentComplete() {
	r.IsPaymentCompleted = !r.IsPaymentCompleted
}

func (r *EventRegister) CalculateAmountDueInC(paymentOptions map[string]PaymentOption) int {
	totalInC := 0

	for _, e := range r.Entries {
		for id, count := range e.EntrantCountByPaymentTypeId {
			opt, ok := paymentOptions[id]
			if !ok {
				continue
			}

			totalInC += opt.ValueInC * count
		}
	}

	return totalInC
}
