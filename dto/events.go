package dto

import "time"

type AddEventDto struct {
	IdempotencyId string
	Date          time.Time
	Title         string
	OrganiserId   string
}

type AddRegisterEntry struct {
	IdempotencyId             string
	EventId                   string
	Name                      string
	ContactNumber             string
	VehicleRegistration       string
	EntrantCountByPaymentType map[string]int
	AmountDueInC              int
	IsConditionsAccepted      bool
}
