package services

import "goregister/domain"

type SettingsService struct {
}

func NewSettingsService() *SettingsService {
	return &SettingsService{}
}

func (s *SettingsService) GetPaymentOptions() map[string]domain.PaymentOption {
	staff := domain.PaymentOption{
		Id:                          "staff",
		Name:                        "(Free) Staff",
		ValueInC:                    0,
		DisplayValueInR:             "R0",
		IsEnabledForNewTransactions: true,
	}

	rhino := domain.PaymentOption{
		Id:                          "rhino",
		Name:                        "(Free) Rhino card",
		ValueInC:                    0,
		DisplayValueInR:             "R0",
		IsEnabledForNewTransactions: true,
	}

	youngChild := domain.PaymentOption{
		Id:                          "child<3",
		Name:                        "(Free) Child <3yrs",
		ValueInC:                    0,
		DisplayValueInR:             "R0",
		IsEnabledForNewTransactions: true,
	}

	child := domain.PaymentOption{
		Id:                          "child",
		Name:                        "(Cash) Child >3yrs",
		ValueInC:                    3500,
		DisplayValueInR:             "R35",
		IsEnabledForNewTransactions: true,
	}

	adult := domain.PaymentOption{
		Id:                          "adult",
		Name:                        "(Cash) Adult",
		ValueInC:                    7000,
		DisplayValueInR:             "R70",
		IsEnabledForNewTransactions: true,
	}

	return map[string]domain.PaymentOption{
		staff.Id:      staff,
		rhino.Id:      rhino,
		youngChild.Id: youngChild,
		child.Id:      child,
		adult.Id:      adult,
	}
}
