package pkg

import (
	"context"
	"tax-management/external"
)

type Service struct {
	Repository Repository
	TaxClient  map[string]TaxClient
}

func (s Service) ProcessKafkaMessage(topicName string, data external.RawTransaction) error {
	taxId, taxProcessId, e := s.Repository.InsertTaxData(context.Background(), topicName, data)
	if e != nil {
		panic("")
	}

	println("data")

	if client, ok := s.TaxClient[data.After.Taxid]; ok {
		invoice := data.ToStandardInvoice()
		res, err := client.SendInvoices(&taxId, &taxProcessId, invoice)
		if err != nil {
			panic("")
		}
		println(res)
	}

	//taxClient := s.TaxClient[""]

	//taxClient.SendInvoices()

	return nil
}
