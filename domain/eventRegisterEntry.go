package domain

import (
	"errors"
	"strings"
	"time"
)

type EventRegisterEntry struct {
	IdempotencyId               string
	PersonName                  string
	PersonContactNumber         string
	VehicleRegistration         string
	RhinoCard                   string
	EntrantCountByPaymentTypeId map[string]int
	AmountDueInC                int
	IsConditionsAccepted        bool
	Timestamp                   time.Time
}

func NewEventRegisterEntry(
	idempotencyId string,
	personName string,
	personContactNumber string,
	vehicleRegistration string,
	rhinoCard string,
	entrantCountByPaymentTypeId map[string]int,
	amountDueInC int,
	isConditionsAccepted bool,
	timestamp time.Time,
) (*EventRegisterEntry, error) {

	personName = strings.TrimSpace(personName)
	personContactNumber = strings.TrimSpace(personContactNumber)
	vehicleRegistration = strings.TrimSpace(vehicleRegistration)

	if idempotencyId == "" {
		return nil, errors.New("Invalid idempotency id")
	}

	if personName == "" {
		return nil, errors.New("Invalid person name")
	}

	if personContactNumber == "" {
		return nil, errors.New("Invalid person contact number")
	}

	if len(entrantCountByPaymentTypeId) == 0 {
		return nil, errors.New("No entrant count data")
	}

	if amountDueInC < 0 {
		return nil, errors.New("Negative amount paid")
	}

	e := EventRegisterEntry{
		IdempotencyId:               idempotencyId,
		PersonName:                  personName,
		PersonContactNumber:         personContactNumber,
		VehicleRegistration:         vehicleRegistration,
		RhinoCard:                   rhinoCard,
		EntrantCountByPaymentTypeId: entrantCountByPaymentTypeId,
		AmountDueInC:                amountDueInC,
		IsConditionsAccepted:        isConditionsAccepted,
		Timestamp:                   timestamp,
	}

	return &e, nil
}
