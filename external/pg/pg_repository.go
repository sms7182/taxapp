package pg

import (
	"context"
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

func (r RepositoryImpl) IsNotProcessable(ctx context.Context, trn string) bool {
	var tp models2.TaxProcess
	if err := r.DB.WithContext(ctx).Where("internal_trn = ?", trn).Last(&tp).Error; err == gorm.ErrRecordNotFound {
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
	taxProcess := toTaxProcess(tax, rawType, companyName)

	err := r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if e := tx.Clauses(clause.Returning{}).Create(&tax).Error; e != nil {
			return e
		}
		taxProcess.TaxRawId = tax.Id

		if e := tx.Clauses(clause.Returning{}).Create(&taxProcess).Error; e != nil {
			return e
		}
		taxId := terminal.GenerateTaxID(taxData.After.Username, taxProcess.Id, time.UnixMilli(taxData.After.Indatim))
		taxProcess.TaxId = &taxId

		if e := tx.Model(&models2.TaxProcess{}).Where("id = ?", taxProcess.Id).Update("tax_id", taxId).Error; e != nil {
			return e
		}
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

func toTaxProcess(tax models2.TaxRawDomain, rawType string, companyName string) models2.TaxProcess {
	taxP := models2.TaxProcess{
		TaxType:     rawType,
		TaxRawId:    tax.Id,
		CompanyName: &companyName,
	}
	return taxP
}

func (r RepositoryImpl) GetInProgressTaxProcess(ctx context.Context) (taxProcesses []models.RawProcessTaxData, err error) {

	var rawPTData []models.RawProcessTaxData
	sqlStr := `select tp.id, tp.tax_raw_id, tp.tax_org_reference_id,tp.tax_id 	from tax_process tp where tp.status = 'in-progress' order by tp.created_at limit 512`
	if e := r.DB.Raw(sqlStr).Scan(&rawPTData).Error; e != nil {
		return nil, e
	}

	return rawPTData, nil

}
