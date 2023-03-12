package pg

import (
	"context"
	"tax-management/external/exkafka/messages"
	"tax-management/models"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type RepositoryImpl struct {
	DB *gorm.DB
}

func (repository RepositoryImpl) LogReqRes(taxRawId *uint, taxProcessId *uint, requestUniqueId string, apiName string, url string, statusCode int, req string, res *string, errorMsg *string) error {
	taxOfficeRequest := models.TaxOfficeRequestResponseLogModel{
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

func (repository RepositoryImpl) InsertTaxData(ctx context.Context, taxData messages.RawTransaction) (*string, error) {
	tax := models.TaxRawDomain{}
	tax.TaxData.Set(taxData)
	var taxUniqId *string
	err := repository.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if e := tx.Clauses(clause.Returning{}).Create(&tax).Error; e != nil {
			return e
		}
		taxProcess := toTaxProcess(tax)

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
func toTaxProcess(tax models.TaxRawDomain) models.TaxProcess {
	taxP := models.TaxProcess{
		TaxType:  "",
		TaxRawId: tax.Id,
	}
	return taxP
}
