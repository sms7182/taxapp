package pkg

import "tax-management/external/exkafka/messages"

type TaxClient interface {
	GetServerInformation() (*string, error)
	GetToken() (string, error)
	GetFiscalInformation(token string)
	SendInvoice(rawdata messages.RawTransaction)
}
