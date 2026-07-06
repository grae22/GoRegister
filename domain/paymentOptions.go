package domain

type PaymentOption struct {
	Id                          string
	Name                        string
	ValueInC                    int
	DisplayValueInR             string
	IsEnabledForNewTransactions bool
}
