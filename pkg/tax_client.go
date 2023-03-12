package pkg

type TaxClient interface {
	GetServerInformation() (*string, error)
}
