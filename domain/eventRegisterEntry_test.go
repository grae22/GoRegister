package domain

import (
	"testing"
	"time"
)

func TestNewEventRegisterEntry_WhenRequiredFieldBlank_ThenErrorReturned(t *testing.T) {
	// Arrange.
	const idempotencyId string = "123"
	const personName string = "name"
	const personContactNumber string = "+27123456789"
	const vehicleRegistration string = "ND123"
	var entrantCountByType map[string]int = map[string]int{}
	const amountPaidInC int = 70
	const isConditionsAccepted bool = true

	// Act.
	_, errIdempotency := NewEventRegisterEntry(
		"",
		personName,
		personContactNumber,
		vehicleRegistration,
		entrantCountByType,
		amountPaidInC,
		isConditionsAccepted,
		time.Now())

	_, errName := NewEventRegisterEntry(
		idempotencyId,
		" ",
		personContactNumber,
		vehicleRegistration,
		entrantCountByType,
		amountPaidInC,
		isConditionsAccepted,
		time.Now())

	_, errContact := NewEventRegisterEntry(
		idempotencyId,
		personName,
		" ",
		vehicleRegistration,
		entrantCountByType,
		amountPaidInC,
		isConditionsAccepted,
		time.Now())

	_, errVehicleReg := NewEventRegisterEntry(
		idempotencyId,
		personName,
		personContactNumber,
		" ",
		entrantCountByType,
		amountPaidInC,
		isConditionsAccepted,
		time.Now())

	_, errEntrants := NewEventRegisterEntry(
		idempotencyId,
		personName,
		personContactNumber,
		vehicleRegistration,
		nil,
		amountPaidInC,
		isConditionsAccepted,
		time.Now())

	_, errPaid := NewEventRegisterEntry(
		idempotencyId,
		personName,
		personContactNumber,
		vehicleRegistration,
		entrantCountByType,
		-1,
		isConditionsAccepted,
		time.Now())

	_, errConditions := NewEventRegisterEntry(
		idempotencyId,
		personName,
		personContactNumber,
		vehicleRegistration,
		entrantCountByType,
		amountPaidInC,
		false,
		time.Now())

	// Assert.
	if errIdempotency == nil {
		t.Error("No error returned when idempotency id invalid")
	}

	if errName == nil {
		t.Error("No error returned when person name invalid")
	}

	if errContact == nil {
		t.Error("No error returned when person contact invalid")
	}

	if errVehicleReg != nil {
		t.Error("Error returned when vehicle registration blank - this is allowed")
	}

	if errEntrants == nil {
		t.Error("No error returned on nil entrants object")
	}

	if errPaid == nil {
		t.Error("No error returned when amount paid invalid")
	}

	if errConditions != nil {
		t.Error("Error returned when conditions accepted false - this is allowed")
	}
}

func TestNewEventRegisterEntry_WhenParamsOk_ThenObjectReturned(t *testing.T) {
	// Arrange.
	const idempotencyId string = "123"
	const personName string = "name"
	const personContactNumber string = "+27123456789"
	const vehicleRegistration string = "ND123"
	var entrantCountByType map[string]int = map[string]int{}
	const amountPaidInC int = 70
	const isConditionsAccepted bool = true

	// Act.
	e, _ := NewEventRegisterEntry(
		idempotencyId,
		personName,
		personContactNumber,
		vehicleRegistration,
		entrantCountByType,
		amountPaidInC,
		isConditionsAccepted,
		time.Now())

	// Assert.
	if e == nil {
		t.Error("Nil entry returned")
	}
}
