package pkg

import "tax-management/utility"

type TaxClient interface {
	GetServerInformation() (*string, error)
	GetToken() (*utility.TokenResponse, error)
}
