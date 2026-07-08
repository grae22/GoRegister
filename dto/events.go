package dto

import "time"

type AddEventDto struct {
	IdempotencyId      string
	Date               time.Time
	Title              string
	OrganiserId        string
	IsPaymentCompleted bool
}

type AddRegisterEntry struct {
	IdempotencyId             string
	EventId                   string
	Name                      string
	ContactNumber             string
	VehicleRegistration       string
	RhinoCard                 string
	EntrantCountByPaymentType map[string]int
	AmountDueInC              int
	IsConditionsAccepted      bool
}
