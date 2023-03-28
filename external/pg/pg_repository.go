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

type RepositoryImpl struct {
	DB *gorm.DB
}

func (repository RepositoryImpl) LogReqRes(taxRawId *uint, taxProcessId *uint, requestUniqueId string, apiName string, url string, statusCode int, req string, res *string, errorMsg *string) error {
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
	return repository.DB.Create(&taxOfficeRequest).Error
}

func (repository RepositoryImpl) InsertTaxData(ctx context.Context, rawType string, taxData external.RawTransaction, companyName string) (uint, uint, string, error) {
	tax := models2.TaxRawDomain{
		TaxType:  rawType,
		UniqueId: taxData.After.Trn + "-" + rawType,
	}
	tax.TaxData.Set(taxData)
	taxProcess := toTaxProcess(tax, rawType, companyName)

	err := repository.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if e := tx.Clauses(clause.Returning{}).Create(&tax).Error; e != nil {
			return e
		}
		taxProcess.TaxRawId = tax.Id

		if e := tx.Clauses(clause.Returning{}).Create(&taxProcess).Error; e != nil {
			return e
		}
		taxId := terminal.GenerateTaxID(taxData.After.Username, taxProcess.Id, time.UnixMilli(taxData.After.Indatim))
		taxProcess.TaxId = &taxId

		return tx.Model(&models2.TaxProcess{}).Where("id = ?", taxProcess.Id).Update("tax_id", taxId).Error
	})
	if err != nil {
		return 0, 0, "", err
	}
	return tax.Id, taxProcess.Id, *taxProcess.TaxId, nil
}
func (repository RepositoryImpl) UpdateTaxProcessStandartInvoice(ctx context.Context, taxProcessId uint, invoice types.StandardInvoice) error {
	updTax := models2.TaxProcess{
		Id: taxProcessId,
	}
	updTax.StandardInvoice.Set(invoice)
	return repository.DB.Model(&models2.TaxProcess{}).Where("id = ?", taxProcessId).Updates(updTax).Error
}
func (repository RepositoryImpl) UpdateTaxReferenceId(ctx context.Context, taxProcessId uint, taxOrgReferenceId string, taxOrgInternalTrn *string, taxOrgInquiryUuid *string) error {
	updTax := models2.TaxProcess{
		Id:                taxProcessId,
		Status:            models2.InProgress.String(),
		TaxOrgReferenceId: &taxOrgReferenceId,
		InternalTrn:       taxOrgInternalTrn,
		InquiryUuid:       taxOrgInquiryUuid,
	}
	return repository.DB.Model(&models2.TaxProcess{}).Where("id = ?", taxProcessId).Updates(updTax).Error
}

func (repository RepositoryImpl) UpdateTaxProcessStatus(ctx context.Context, taxProcessId uint, status string, confirmationReferenceId *string) error {
	updTax := models2.TaxProcess{
		Id:                      taxProcessId,
		Status:                  status,
		ConfirmationReferenceId: confirmationReferenceId,
	}
	return repository.DB.Model(&models2.TaxProcess{}).Where("id = ?", taxProcessId).Updates(updTax).Error
}

func toTaxProcess(tax models2.TaxRawDomain, rawType string, companyName string) models2.TaxProcess {
	taxP := models2.TaxProcess{
		TaxType:  rawType,
		TaxRawId: tax.Id,
		CompanyName: &companyName,
	}
	return taxP
}

func (repository RepositoryImpl) GetInprogressTaxProcess(ctx context.Context) (taxProcesses []models.RawProcessTaxData, err error) {

	var rawPTData []models.RawProcessTaxData
	sqlStr := `select tp.id, tp.tax_raw_id, tp.tax_org_reference_id,tp.tax_id 	from tax_process tp where tp.status = 'in-progress' order by tp.created_at limit 256`
	if e := repository.DB.Raw(sqlStr).Scan(&rawPTData).Error; e != nil {
		return nil, e
	}

	return rawPTData, nil

}
