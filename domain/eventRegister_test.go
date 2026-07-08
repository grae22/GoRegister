package domain

import (
	"strings"
	"testing"
	"time"
)

func TestNewEventRegister_WhenBlankIdempotencyId_ThenReturnsError(t *testing.T) {
	// Arrange.
	// Act.
	_, err := NewEventRegister(
		" ",
		time.Now(),
		"title",
		"oid",
		map[string]PaymentOption{})

	// Assert.
	if err == nil {
		t.Error("Blank idempotency id should produce error")
	}
}
func TestNewEventRegister_WhenNoTitle_ThenReturnsError(t *testing.T) {
	// Arrange.
	// Act.
	_, err := NewEventRegister(
		"123",
		time.Now(),
		"",
		"oid",
		map[string]PaymentOption{})

	// Assert.
	if err == nil {
		t.Error("Empty title should produce error")
	}
}

func TestNewEventRegister_WhenBlankTitle_ThenReturnsError(t *testing.T) {
	// Arrange.
	// Act.
	_, err := NewEventRegister(
		"123",
		time.Now(),
		" ",
		"oid",
		map[string]PaymentOption{})

	// Assert.
	if err == nil {
		t.Error("Blank title should produce error")
	}
}

func TestNewEventRegister_WhenNoOrganiser_ThenReturnsError(t *testing.T) {
	// Arrange.
	// Act.
	_, err := NewEventRegister(
		"123",
		time.Now(),
		"title",
		"",
		map[string]PaymentOption{})

	// Assert.
	if err == nil {
		t.Error("Empty title should produce error")
	}
}

func TestNewEventRegister_WhenBlankOrganiser_ThenReturnsError(t *testing.T) {
	// Arrange.
	// Act.
	_, err := NewEventRegister(
		"123",
		time.Now(),
		"title",
		" ",
		map[string]PaymentOption{})

	// Assert.
	if err == nil {
		t.Error("Blank title should produce error")
	}
}

func TestNewEventRegister_WhenParamsOk_ThenReturnsEventRegister(t *testing.T) {
	// Arrange.
	// Act.
	er, _ := NewEventRegister(
		"123",
		time.Now(),
		"title",
		"oid",
		map[string]PaymentOption{
			"cash": {},
		})

	// Assert.
	if er == nil {
		t.Error("Return object should not be null")
	}
}

func TestNewEventRegister_WhenParamsOk_ThenReturnsEventRegisterWithCorrectValues(t *testing.T) {
	// Arrange.
	idempotencyId := "123"
	date := time.Now()
	title := " title "
	organiserId := " oid "
	paymentOptions := map[string]PaymentOption{
		"o1": {},
	}

	// Act.
	er, _ := NewEventRegister(
		idempotencyId,
		date,
		title,
		organiserId,
		paymentOptions)

	// Assert.
	if er.IdempotencyId != idempotencyId {
		t.Error("Incorrect idempotency id")
	}

	if er.Date != date {
		t.Error("Incorrect date")
	}

	if er.Title != strings.TrimSpace(title) {
		t.Error("Incorrect title")
	}

	if er.OrganiserId != strings.TrimSpace(organiserId) {
		t.Error("Incorrect organiser id")
	}

	if er.Entries == nil {
		t.Error("Entries is nil")
	}

	if _, ok := er.PaymentOptionsById["o1"]; !ok {
		t.Error("Payment option missing")
	}
}

func TestAddEntry_WhenDuplicateIdempotencyReceived_ThenEntryIsNotAdded(t *testing.T) {
	// Arrange.
	const idempotencyId string = "abc"

	r, _ := NewEventRegister(
		"123",
		time.Now(),
		"title",
		"oid",
		map[string]PaymentOption{
			"o1": {},
		})

	e, _ := NewEventRegisterEntry(
		idempotencyId,
		"name",
		"+27...",
		"ND123",
		"G...",
		map[string]int{
			"cash": 0,
		},
		0,
		true,
		time.Now())

	r.AddEntry(e)

	// Act.
	r.AddEntry(e)

	// Assert.
	if len(r.Entries) > 1 {
		t.Error("Only one entry should exist")
	}
}
