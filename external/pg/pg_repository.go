package pg

import (
	"context"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"tax-management/external"
	models2 "tax-management/external/pg/models"
	"time"
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

func (repository RepositoryImpl) InsertTaxData(ctx context.Context, rawType string, taxData external.RawTransaction) (*string, error) {
	tax := models2.TaxRawDomain{
		TaxType:  rawType,
		UniqueId: taxData.After.Trn + "-" + rawType,
	}
	tax.TaxData.Set(taxData)

	var taxUniqId *string
	err := repository.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if e := tx.Clauses(clause.Returning{}).Create(&tax).Error; e != nil {
			return e
		}
		taxProcess := toTaxProcess(tax, rawType)

		err := tx.Create(&taxProcess).Error
		if err == nil {
			taxUniqId = &taxProcess.TaxUniqueId
		}
		return err
	})
	if err != nil {
		return nil, err
	}

	return taxUniqId, nil

}
func toTaxProcess(tax models2.TaxRawDomain, rawType string) models2.TaxProcess {
	taxP := models2.TaxProcess{
		TaxType:  rawType,
		TaxRawId: tax.Id,
	}
	return taxP
}
