package pkg

import (
	"context"
	"errors"
	"fmt"
	"log"
	"tax-management/external"
	"tax-management/external/pg/models"
	"tax-management/notify"
	"tax-management/taxDep/types"
	"time"
)

type Service struct {
	Repository            Repository
	UsernameToCompanyName map[string]string
	NotificationClient    notify.NotificationClient

	TaxClient TaxClient
}

const layout = "2006-01-02T15:04:05"

func (service Service) StartSendingInvoice(data external.RawTransaction) error {
	ctx := context.Background()
	farvardin1, _ := time.Parse(layout, "2023-03-20T23:59:59")

	if time.UnixMilli(data.After.Indatim).Before(farvardin1) {
		return nil
	}
	var usrName string
	if usrName, ok := service.UsernameToCompanyName[data.After.Username]; !ok {
		validCustomer, err := service.Repository.GetUserName(ctx, data.After.Username)
		if err != nil || validCustomer == nil {
			//log in db for notValidCustomer
			return nil
		}
		service.UsernameToCompanyName[data.After.Username] = validCustomer.UserName
		usrName = validCustomer.UserName
		fmt.Printf("usrname for sending invoice %s", usrName)
	}

	rawDataId, taxProcessId, taxId, e := service.Repository.InsertTaxData(context.Background(), "", data, usrName)
	if e != nil {
		panic(fmt.Sprintf("failed insertData, data: %+v", data))
	}

	invoice := data.ToStandardInvoice(taxId)
	if len(invoice) == 1 {
		service.Repository.UpdateTaxProcessStandardInvoice(ctx, taxProcessId, invoice[0])
	}
	res, err := service.TaxClient.SendInvoices(&rawDataId, &taxProcessId, invoice)
	if err != nil {
		service.Repository.UpdateTaxProcessStatus(ctx, taxProcessId, models.TaxStatusFailed.String(), nil)
		return nil
	}
	if len(res.Result) > 0 && len(res.Errors) == 0 {
		arp := res.Result[0]
		return service.Repository.UpdateTaxReferenceId(ctx, taxProcessId, arp.ReferenceNumber, &data.After.Trn, &arp.UID)
	}

	return errors.New("failed to process kafka message")
}
func (s Service) ProcessKafkaMessage(topicName string, data external.RawTransaction) error {
	ctx := context.Background()

	if s.Repository.IsNotProcessable(ctx, topicName, data.After.Trn) {
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

	invoice := data.ToStandardInvoice(taxId)
	if len(invoice) == 1 {
		s.Repository.UpdateTaxProcessStandardInvoice(ctx, taxProcessId, invoice[0])
	}
	res, err := s.TaxClient.SendInvoices(&rawDataId, &taxProcessId, invoice)
	if err != nil {
		s.Repository.UpdateTaxProcessStatus(ctx, taxProcessId, models.TaxStatusFailed.String(), nil)
		return nil
	}
	if len(res.Result) > 0 && len(res.Errors) == 0 {
		arp := res.Result[0]
		return s.Repository.UpdateTaxReferenceId(ctx, taxProcessId, arp.ReferenceNumber, &data.After.Trn, &arp.UID)
	}

	return errors.New("failed to process kafka message")
}

func (s Service) TaxRequestInquiry() {
	taxProcess, err := s.Repository.GetInProgressTaxProcess(context.Background())
	if err != nil {
		log.Printf("Get Inprogress Taxprocess has error %s", err)
	} else if len(taxProcess) > 0 {
		for i := 0; i < len(taxProcess); i++ {
			//		userName := taxProcess[i].TaxId[0:6]

			inquiryResult, err := s.TaxClient.InquiryByReferences(&taxProcess[i].TaxRawId, &taxProcess[i].Id, []string{taxProcess[i].OrgReferenceId})
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

func (service Service) NotifyFailedTax() {
	taxProcess, err := service.Repository.GetFailedTaxProcess(context.Background())
	var taxProcessIds []uint
	ln := len(taxProcess)
	fmt.Printf("len taxprocess is %v", ln)
	if err == nil && ln > 0 {
		var failedBodys []notify.FailedBody
		for i := 0; i < len(taxProcess); i++ {
			taxProcessIds = append(taxProcessIds, taxProcess[i].Id)
			var result types.AutoGenerated
			taxProcess[i].Response.AssignTo(&result)
			var dataErr string
			if len(result.Error) > 0 {
				for j := 0; j < len(result.Error); j++ {
					dataErr += result.Error[j].Msg + "###"
				}

			}

			failedBodys = append(failedBodys, notify.FailedBody{
				FailedMessage: dataErr,
				Int_Trn:       taxProcess[i].InternalTrn,
			})
		}
		notifyResult, err := service.NotificationClient.FailedNotify(context.Background(), failedBodys, "templates/FailedTax.html")
		if err == nil && notifyResult != nil {

			service.Repository.UpdateNotifyFailedOfTaxProcess(context.Background(), taxProcessIds)
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
