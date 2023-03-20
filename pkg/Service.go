package pkg

import (
	"context"
	"errors"
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

		if len(res.Result) > 0 && len(res.Errors) == 0 {
			arp := res.Result[0]
			return s.Repository.UpdateTaxReferenceId(context.Background(), taxProcessId, arp.ReferenceNumber)
		}

	}
	return errors.New("failed to process kafka message")
}
