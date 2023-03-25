package pg

import (
	"context"
	"tax-management/external"
	"tax-management/external/pg/models"
	models2 "tax-management/external/pg/models"
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

func (repository RepositoryImpl) InsertTaxData(ctx context.Context, rawType string, taxData external.RawTransaction) (uint, uint, error) {
	tax := models2.TaxRawDomain{
		TaxType:  rawType,
		UniqueId: taxData.After.Trn + "-" + rawType,
	}
	tax.TaxData.Set(taxData)
	taxProcess := toTaxProcess(tax, rawType)

	err := repository.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if e := tx.Clauses(clause.Returning{}).Create(&tax).Error; e != nil {
			return e
		}
		taxProcess.TaxRawId = tax.Id
		return tx.Create(&taxProcess).Error
	})
	if err != nil {
		return 0, 0, err
	}
	return tax.Id, taxProcess.Id, nil
}

func (repository RepositoryImpl) UpdateTaxReferenceId(ctx context.Context, taxProcessId uint, taxOrgReferenceId string) error {
	updTax := models2.TaxProcess{
		Id:                taxProcessId,
		Status:            models2.InProgress.String(),
		TaxOrgReferenceId: &taxOrgReferenceId,
	}
	return repository.DB.Model(&models2.TaxProcess{}).Where("id = ?", taxProcessId).Updates(updTax).Error
}

func (repository RepositoryImpl) UpdateTaxProcessStatus(ctx context.Context, taxProcessId uint, status string) error {
	updTax := models2.TaxProcess{
		Id:     taxProcessId,
		Status: status,
	}
	return repository.DB.Model(&models2.TaxProcess{}).Where("id = ?", taxProcessId).Updates(updTax).Error
}

func toTaxProcess(tax models2.TaxRawDomain, rawType string) models2.TaxProcess {
	taxP := models2.TaxProcess{
		TaxType:  rawType,
		TaxRawId: tax.Id,
	}
	return taxP
}

func (repository RepositoryImpl) GetInprogressTaxProcess(ctx context.Context) (taxProcesses []models.RawProcessTaxData, err error) {
	var rawPTData []models.RawProcessTaxData
	sqlStr := "select tp.id,tp.tax_raw_id,tp.tax_org_reference_id,trd.tax_data 	from tax_process tp join tax_raw_data trd on tp.tax_raw_id=trd.id	where tp.status='in-progress'"
	if e := repository.DB.Raw(sqlStr).Scan(&rawPTData).Error; e != nil {
		return nil, e
	}

	return rawPTData, nil
	// var tps []models.TaxProcess
	// if e := repository.DB.WithContext(ctx).Model(&models.TaxProcess{}).Where("status = ?", models.InProgress.String()).Scan(&tps).Error; e != nil {
	// 	return nil, e
	// }
	// return tps, nil
}
