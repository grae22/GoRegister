package domain

import (
	"errors"
	"strings"
	"time"
)

type EventRegisterEntry struct {
	IdempotencyId        string
	PersonName           string
	PersonContactNumber  string
	VehicleRegistration  string
	EntrantCountByType   map[string]int
	AmountPaidInC        int
	IsConditionsAccepted bool
	Timestamp            time.Time
}

func NewEventRegisterEntry(
	idempotencyId string,
	personName string,
	personContactNumber string,
	vehicleRegistration string,
	entrantCountByType map[string]int,
	amountPaidInC int,
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

	if entrantCountByType == nil {
		return nil, errors.New("Nil entrant count data")
	}

	if amountPaidInC < 0 {
		return nil, errors.New("Negative amount paid")
	}

	e := EventRegisterEntry{
		IdempotencyId:        idempotencyId,
		PersonName:           personName,
		PersonContactNumber:  personContactNumber,
		VehicleRegistration:  vehicleRegistration,
		EntrantCountByType:   entrantCountByType,
		AmountPaidInC:        amountPaidInC,
		IsConditionsAccepted: isConditionsAccepted,
		Timestamp:            timestamp,
	}

	return &e, nil
}
