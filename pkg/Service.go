package pkg

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"tax-management/external"
	"tax-management/external/pg/models"
	"tax-management/taxDep/types"
	"time"
)

type Service struct {
	Repository            Repository
	UsernameToCompanyName map[string]string

	TaxClient TaxClient
}

const layout = "2006-01-02T15:04:05"

func (service Service) InitialCustomer(dto *external.CustomerDto) (*uint, error) {
	ctx := context.Background()
	id, err := service.Repository.CreateCustomer(ctx, models.Customer{
		FinanceId:  dto.UserName,
		Token:      dto.Token,
		PublicKey:  dto.PublicKey,
		PrivateKey: dto.PrivateKey,
		ExpireTime: time.Now().AddDate(1, 0, 0),
	})
	if err != nil {
		return nil, err
	}
	return id, nil
}
func (service Service) StartSendingInvoice(ctx context.Context, data external.RawTransaction) error {
	farvardin1, _ := time.Parse(layout, "2023-03-20T23:59:59")
	log.Printf("befor check  farvardin")
	if time.UnixMilli(data.After.Indatim).Before(farvardin1) {
		return nil
	}
	log.Printf("after check farvardin")
	validCustomer, err := service.Repository.GetUserName(ctx, data.After.Username)
	if err != nil || validCustomer == nil {
		if err != nil {
			fmt.Sprintf("founding user has error %s", err)
			return err
		}
		fmt.Sprintf("customer not found")
		return nil
	}

	rawDataId, taxProcessId, taxId, e := service.Repository.InsertTaxData(ctx, "", data, validCustomer.FinanceId)
	if e != nil {
		panic(fmt.Sprintf("failed insertData, data: %+v", data))
	}

	invoice := data.ToStandardInvoice(taxId)
	if len(invoice) == 1 {
		service.Repository.UpdateTaxProcessStandardInvoice(ctx, taxProcessId, invoice[0])
	}
	res, err := service.TaxClient.SendInvoices(&rawDataId, &taxProcessId, invoice, validCustomer.PrivateKey, validCustomer.FinanceId)
	if err != nil {
		log.Printf("sending invoices has error %s", err)
		service.Repository.UpdateTaxProcessStatus(ctx, taxProcessId, models.TaxStatusFailed.String(), nil)
		return err
	}
	if len(res.Result) > 0 && len(res.Errors) == 0 {
		arp := res.Result[0]
		return service.Repository.UpdateTaxReferenceId(ctx, taxProcessId, arp.ReferenceNumber, &data.After.Trn, &arp.UID)
	}

	return errors.New("failed to process kafka message")
}

// func (s Service) ProcessKafkaMessage(topicName string, data external.RawTransaction) error {
// 	ctx := context.Background()

// 	if s.Repository.IsNotProcessable(ctx, topicName, data.After.Trn) {
// 		return nil
// 	}

// 	farvardin1, _ := time.Parse(layout, "2023-03-20T23:59:59")

// 	if time.UnixMilli(data.After.Indatim).Before(farvardin1) {
// 		return nil
// 	}
// 	rawDataId, taxProcessId, taxId, e := s.Repository.InsertTaxData(context.Background(), topicName, data, s.UsernameToCompanyName[data.After.Username])
// 	if e != nil {
// 		panic(fmt.Sprintf("failed, topic: %s, data: %+v", topicName, data))
// 	}

// 	invoice := data.ToStandardInvoice(taxId)
// 	if len(invoice) == 1 {
// 		s.Repository.UpdateTaxProcessStandardInvoice(ctx, taxProcessId, invoice[0])
// 	}
// 	res, err := s.TaxClient.SendInvoices(&rawDataId, &taxProcessId, invoice)
// 	if err != nil {
// 		s.Repository.UpdateTaxProcessStatus(ctx, taxProcessId, models.TaxStatusFailed.String(), nil)
// 		return nil
// 	}
// 	if len(res.Result) > 0 && len(res.Errors) == 0 {
// 		arp := res.Result[0]
// 		return s.Repository.UpdateTaxReferenceId(ctx, taxProcessId, arp.ReferenceNumber, &data.After.Trn, &arp.UID)
// 	}

// 	return errors.New("failed to process kafka message")
// }

func (s Service) TaxRequestInquiry(userName string) {
	taxProcess, err := s.Repository.GetInProgressTaxProcess(context.Background())
	user, err := s.Repository.GetUserName(context.Background(), userName)

	if err != nil {
		log.Printf("Get Inprogress Taxprocess has error %s", err)
	} else if len(taxProcess) > 0 {
		for i := 0; i < len(taxProcess); i++ {
			//		userName := taxProcess[i].TaxId[0:6]

			inquiryResult, err := s.TaxClient.InquiryByReferences(&taxProcess[i].TaxRawId, &taxProcess[i].Id, []string{taxProcess[i].OrgReferenceId}, user.PrivateKey, user.FinanceId)
			if err == nil && len(inquiryResult) > 0 {
				if inquiryResult[0].Data.Success {
					s.Repository.UpdateTaxProcessStatus(context.Background(), taxProcess[i].Id, models.TaxStatusCompleted.String(), &inquiryResult[0].Data.ConfirmationReferenceID)
				} else if strings.ToLower(inquiryResult[0].Status) == models.TaxStatusFailed.String() {
					s.Repository.UpdateTaxProcessStatus(context.Background(), taxProcess[i].Id, models.TaxStatusFailed.String(), nil)
				} else if strings.ToLower(inquiryResult[0].Status) == models.TaxStatusPending.String() {
					s.Repository.UpdateTaxProcessStatus(context.Background(), taxProcess[i].Id, models.TaxStatusPending.String(), nil)
				}
			} else {
				s.Repository.UpdateTaxProcessStatus(context.Background(), taxProcess[0].Id, models.TaxStatusFailed.String(), nil)
			}

		}
	}
}

func (s Service) AutoRetry(ctx context.Context) {
	taxRaws, err := s.Repository.GetReadyTaxToRetry(ctx)
	if err != nil {
		fmt.Printf("retry data failed")
	}

	for i := 0; i < len(taxRaws); i++ {

	}

}

func (s Service) GetTaxProcess(ctx context.Context, id uint) (*models.TaxProcessViewModel, error) {
	tp, err := s.Repository.GetTaxProcess(ctx, id)
	if err != nil {
		return nil, err
	}
	tpvm := models.TaxProcessViewModel{
		Id:          tp.Id,
		CreatedAt:   tp.CreatedAt,
		UpdatedAt:   tp.UpdatedAt,
		TaxType:     tp.TaxType,
		Status:      tp.Status,
		TaxId:       tp.TaxId,
		InternalTrn: tp.InternalTrn,
		CompanyName: tp.CompanyName,
	}

	var result types.StandardInvoice
	tp.StandardInvoice.AssignTo(&result)
	tpvm.StandardInvoice = result
	return &tpvm, err

}
