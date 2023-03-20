package pkg

import (
	"context"
	"tax-management/external"
	terminal "tax-management/taxDep"
)

type Service struct {
	Repository Repository
	TaxClient  map[string]*terminal.Terminal
}

func (s Service) ProcessKafkaMessage(topicName string, data external.RawTransaction) error {
	s.Repository.InsertTaxData(context.Background(), topicName, data)
	println("data")

	//taxClient := s.TaxClient[""]

	//taxClient.SendInvoices()

	return nil
}
