package pkg

import (
	"context"
	"errors"
	"fmt"
	"log"
	"tax-management/external"
	"tax-management/external/pg/models"
	"time"
)

type Service struct {
	Repository            Repository
	TaxClient             map[string]TaxClient
	UsernameToCompanyName map[string]string
}

const layout = "2006-01-02T15:04:05"

func (s Service) ProcessKafkaMessage(topicName string, data external.RawTransaction) error {
	ctx := context.Background()

	//if s.Repository.NumberOfFailureExceeded() {
	//	panic("number of failure exceeded")
	//}

	if s.Repository.IsNotProcessable(ctx, data.After.Trn) {
		return nil
	}

	farvardin1, _ := time.Parse(layout, "2023-03-20T23:59:59")

	if time.UnixMilli(data.After.Indatim).Before(farvardin1) {
		return nil
	}
	rawDataId, taxProcessId, taxId, e := s.Repository.InsertTaxData(context.Background(), topicName, data, s.UsernameToCompanyName[data.After.Username])
	if e != nil {
		panic(fmt.Sprintf("failed, topic: %s, data: %+v", topicName, data))
	}

	if client, ok := s.TaxClient[data.After.Username]; ok {
		invoice := data.ToStandardInvoice(taxId)
		if len(invoice) == 1 {
			s.Repository.UpdateTaxProcessStandardInvoice(ctx, taxProcessId, invoice[0])
		}
		res, err := client.SendInvoices(&rawDataId, &taxProcessId, invoice)
		if err != nil {
			s.Repository.UpdateTaxProcessStatus(ctx, taxProcessId, models.TextStatusUnknown.String(), nil)
			return nil
		}
		if len(res.Result) > 0 && len(res.Errors) == 0 {
			arp := res.Result[0]
			return s.Repository.UpdateTaxReferenceId(ctx, taxProcessId, arp.ReferenceNumber, &data.After.Trn, &arp.UID)
		}
	}
	return errors.New("failed to process kafka message")
}

func (s Service) TaxRequestInquiry() {
	taxProcess, err := s.Repository.GetInProgressTaxProcess(context.Background())
	if err != nil {
		log.Printf("Get Inprogress Taxprocess has error %s", err)
	} else if len(taxProcess) > 0 {
		for i := 0; i < len(taxProcess); i++ {
			userName := taxProcess[i].TaxId[0:6]
			if client, ok := s.TaxClient[userName]; ok {
				inquiryResult, err := client.InquiryByReferences(&taxProcess[i].TaxRawId, &taxProcess[i].Id, []string{taxProcess[i].OrgReferenceId})
				if err == nil && len(inquiryResult) > 0 {
					if inquiryResult[0].Data.Success {
						s.Repository.UpdateTaxProcessStatus(context.Background(), taxProcess[i].Id, models.TaxStatusCompleted.String(), &inquiryResult[0].Data.ConfirmationReferenceID)
					} else {
						s.Repository.UpdateTaxProcessStatus(context.Background(), taxProcess[i].Id, models.TaxStatusFailed.String(), nil)
					}
				} else {
					s.Repository.UpdateTaxProcessStatus(context.Background(), taxProcess[i].Id, models.TaxStatusFailed.String(), nil)
				}
			}
		}
	}
}
