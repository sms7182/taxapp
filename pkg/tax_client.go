package pkg

type TaxClient interface {
	GetServerInformation() (*string, error)
	GetToken() (string, error)
	GetFiscalInformation(token string)
}
