package pkg

import (
	"context"
	"errors"
	"log"
	"tax-management/external"
	"time"
)

type Service struct {
	Repository Repository
	TaxClient  map[string]TaxClient
}

const layout = "2006-01-02T15:04:05"

func (s Service) ProcessKafkaMessage(topicName string, data external.RawTransaction) error {
	taxId, taxProcessId, e := s.Repository.InsertTaxData(context.Background(), topicName, data)
	if e != nil {
		panic("")
	}

	farvardin1, _ := time.Parse(layout, "2023-03-20T23:59:59")

	if time.UnixMicro(data.After.Indatim).Before(farvardin1) {
		return nil
	}

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

func (s Service) TaxRequestInquiry() {
	taxProcess, err := s.Repository.GetInprogressTaxProcess(context.Background())
	if err != nil {
		log.Printf("Get Inprogress Taxprocess has error %s", err)
	} else if len(taxProcess) > 0 {
		for i := 0; i < len(taxProcess); i++ {
			//	taxProcess[i].
		}
	}
}
