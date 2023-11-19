package pg

import (
	"context"
	"log"
	"tax-management/external"
	"tax-management/external/pg/models"
	models2 "tax-management/external/pg/models"
	terminal "tax-management/taxDep"
	"tax-management/taxDep/types"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const NumberOfFailureLimit = 10

type RepositoryImpl struct {
	DB *gorm.DB
}

func (r RepositoryImpl) LogReqRes(taxRawId *uint, taxProcessId *uint, requestUniqueId string, apiName string, url string, statusCode int, req string, res *string, errorMsg *string) error {
	taxOfficeRequest := models2.TaxOfficeRequestResponseLogModel{
		TaxRawId:       taxRawId,
		TaxProcessId:   taxProcessId,
		RequestUiqueId: requestUniqueId,
		ApiName:        apiName,
		LoggedAt:       time.Now(),
		Url:            url,
		StatusCode:     statusCode,
		Request:        req,
		Response:       res,
		ErrorMessage:   errorMsg,
	}
	return r.DB.Create(&taxOfficeRequest).Error
}

func (r RepositoryImpl) IsNotProcessable(ctx context.Context, topic string, trn string) bool {
	var tp models2.TaxProcess
	if err := r.DB.WithContext(ctx).Where("internal_trn = ? and tax_type= ?", trn, topic).Last(&tp).Error; err == gorm.ErrRecordNotFound {
		return false
	} else if err != nil {
		return true
	}

	return tp.Status != models2.TaxStatusFailed.String()
}

func (r RepositoryImpl) NumberOfFailureExceeded() bool {
	var failedCount int64
	r.DB.Model(&models2.TaxProcess{}).Where("status = ?", models2.TaxStatusFailed.String()).Count(&failedCount)
	return failedCount > NumberOfFailureLimit
}

func (r RepositoryImpl) InsertTaxData(ctx context.Context, rawType string, taxData external.RawTransaction, companyName string) (uint, uint, string, error) {
	tax := models2.TaxRawDomain{
		TaxType:  rawType,
		UniqueId: taxData.After.Trn + "-" + rawType,
	}
	tax.TaxData.Set(taxData)
	log.Printf("befor check toTaxProcess")
	taxProcess := toTaxProcess(tax, rawType, companyName, taxData)
	log.Printf("befor insert to db in insertTaxData method")
	err := r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if e := tx.Clauses(clause.Returning{}).Create(&tax).Error; e != nil {
			log.Printf("insert in db has error %s", e)
			return e
		}
		taxProcess.TaxRawId = tax.Id

		if e := tx.Clauses(clause.Returning{}).Create(&taxProcess).Error; e != nil {
			return e
		}
		taxId := terminal.GenerateTaxID(taxData.After.Username, taxProcess.Id, time.UnixMilli(taxData.After.Indatim))
		taxProcess.TaxId = &taxId
		log.Printf("after generate of taxid")

		if e := tx.Model(&models2.TaxProcess{}).Where("id = ?", taxProcess.Id).Update("tax_id", taxId).Error; e != nil {
			return e
		}
		log.Printf("inser to taxprocess history")
		tph := models2.ToTaxProcessHistory(taxProcess)
		return tx.WithContext(ctx).Create(&tph).Error

	})
	if err != nil {
		return 0, 0, "", err
	}
	return tax.Id, taxProcess.Id, *taxProcess.TaxId, nil
}

func (r RepositoryImpl) UpdateTaxProcessStandardInvoice(ctx context.Context, taxProcessId uint, invoice types.StandardInvoice) error {
	updTax := models2.TaxProcess{
		Id: taxProcessId,
	}
	updTax.StandardInvoice.Set(invoice)
	return r.DB.Model(&models2.TaxProcess{}).Where("id = ?", taxProcessId).Updates(updTax).Error
}

func (r RepositoryImpl) UpdateTaxReferenceId(ctx context.Context, taxProcessId uint, taxOrgReferenceId string, taxOrgInternalTrn *string, taxOrgInquiryUuid *string) error {
	updTax := models2.TaxProcess{
		Id:                taxProcessId,
		Status:            models2.TaxStatusInProgress.String(),
		TaxOrgReferenceId: &taxOrgReferenceId,
		InternalTrn:       taxOrgInternalTrn,
		InquiryUuid:       taxOrgInquiryUuid,
	}
	return r.DB.Transaction(func(tx *gorm.DB) error {
		if e := tx.Model(&models2.TaxProcess{}).Where("id = ?", taxProcessId).Updates(updTax).Error; e != nil {
			return e
		}
		var updatedTaxProcess models2.TaxProcess
		if e := tx.Where("id = ?", taxProcessId).First(&updatedTaxProcess).Error; e != nil {
			return e
		}
		tph := models2.ToTaxProcessHistory(updatedTaxProcess)
		return tx.Create(&tph).Error
	})
}

func (r RepositoryImpl) UpdateTaxProcessStatus(ctx context.Context, taxProcessId uint, status string, confirmationReferenceId *string) error {
	updTax := models2.TaxProcess{
		Id:                      taxProcessId,
		Status:                  status,
		ConfirmationReferenceId: confirmationReferenceId,
	}
	return r.DB.Transaction(func(tx *gorm.DB) error {
		if e := tx.Model(&models2.TaxProcess{}).Where("id = ?", taxProcessId).Updates(updTax).Error; e != nil {
			return e
		}
		var updatedTaxProcess models2.TaxProcess
		if e := tx.Where("id = ?", taxProcessId).First(&updatedTaxProcess).Error; e != nil {
			return e
		}
		tph := models2.ToTaxProcessHistory(updatedTaxProcess)
		return tx.Create(&tph).Error
	})

}

func toTaxProcess(tax models2.TaxRawDomain, rawType string, companyName string, rawTransaction external.RawTransaction) models2.TaxProcess {
	taxP := models2.TaxProcess{
		TaxType:     rawType,
		TaxRawId:    tax.Id,
		CompanyName: &companyName,
		InternalTrn: &rawTransaction.After.Trn,
	}
	return taxP
}

func (r RepositoryImpl) GetInProgressTaxProcess(ctx context.Context) (taxProcesses []models.RawProcessTaxData, err error) {

	var rawPTData []models.RawProcessTaxData
	sqlStr := `select tp.id, tp.tax_raw_id, tp.tax_org_reference_id,tp.tax_id 	from tax_process tp where (tp.status = 'in-progress' or tp.status='pending') order by tp.created_at limit 512`
	if e := r.DB.Raw(sqlStr).Scan(&rawPTData).Error; e != nil {
		return nil, e
	}

	return rawPTData, nil

}

func (r RepositoryImpl) GetFailedTaxProcess(ctx context.Context) (faileds []models.FailedTaxProcess, err error) {
	var failedTaxProcess []models.FailedTaxProcess
	sqlStr := `select tp.internal_trn, tp.id, ((torrl.response::json -> 'result')::json -> 'data' -> 0)::json -> 'data' as response
from tax_process tp
         join tax_office_request_response_log torrl on tp.id = torrl.tax_process_id
where tp.status = 'failed'
  and tp.failed_notify = false
  and torrl.api_name = 'INQUIRY_BY_REFERENCE_NUMBER'
  and not exists(select * from tax_process where status != 'failed' and tp.internal_trn = internal_trn and tp.tax_type=tax_type)`
	if e := r.DB.Raw(sqlStr).Scan(&failedTaxProcess).Error; e != nil {
		return nil, e
	}
	return failedTaxProcess, err
}

func (r RepositoryImpl) GetTaxProcess(ctx context.Context, id uint) (*models.TaxProcess, error) {
	var taxProcess models.TaxProcess

	if e := r.DB.Model(models.TaxProcess{}).Where("id = ?", id).Scan(&taxProcess).Error; e != nil {
		return nil, e
	}
	return &taxProcess, nil
}

func (r RepositoryImpl) UpdateNotifyFailedOfTaxProcess(ctx context.Context, ids []uint) error {
	if e := r.DB.
		Model(models.TaxProcess{}).
		Where("id in ?", ids).
		Update("failed_notify", true).
		Error; e != nil {
		return e
	}

	return nil
}

func (r RepositoryImpl) GetByTaxRawId(ctx context.Context, taxRawId uint) (*models2.TaxRawDomain, error) {
	var t models2.TaxRawDomain
	if e := r.DB.WithContext(ctx).Where("id = ?", taxRawId).First(&t).Error; e != nil {
		return nil, e
	}
	return &t, nil
}

func (r RepositoryImpl) GetReadyTaxToRetry(ctx context.Context) ([]models.TaxRawDomain, error) {
	var readyToRetry []models.TaxRawDomain
	sqlStr := `select trd.* from tax_process tp 
	join  tax_office_request_response_log torrl on tp.id = torrl.tax_process_id
	join tax_raw_data trd on trd.id=tp.tax_raw_id
	where tp.status='failed'  
	
	  and torrl.api_name = 'SendInvoice'
	  and torrl.status_code=403
	  and tp.internal_trn is not null
	 and not exists(select * from tax_process where status != 'failed' and tp.internal_trn = internal_trn and tp.tax_type=tax_type) `
	if e := r.DB.Raw(sqlStr).Scan(&readyToRetry).Error; e != nil {
		return nil, e
	}
	return readyToRetry, nil
}

func (rep RepositoryImpl) GetUserName(ctx context.Context, token string) (*models.Customer, error) {
	var customer models.Customer
	if e := rep.DB.WithContext(ctx).Where(" finance_id = ?", token).First(&customer).Error; e != nil {
		return nil, e
	}
	return &customer, nil
}

func (rep RepositoryImpl) CreateCustomer(ctx context.Context, customer models.Customer) (*uint, error) {
	err := rep.DB.WithContext(ctx).Create(&customer).Error
	if err == nil {
		return &customer.Id, nil
	}
	return nil, err
}
