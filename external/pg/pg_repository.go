package pg

import (
	"context"
	"tax-management/external/exkafka/messages"
	"tax-management/models"

	"gorm.io/gorm"
)

type RepositoryImpl struct {
	DB *gorm.DB
}

func (repository RepositoryImpl) LogReqRes(requestTraceId string, signature string, packetType string, url string, statusCode int, req string, res *string, errorMsg *string) error {
	return nil
}

func (repository RepositoryImpl) InsertTaxData(ctx context.Context, taxData messages.RawTransaction) (*string, error) {
	tax := models.TaxRawDomain{
		TaxType:       "Raw",
		ProcessStatus: models.InProgress.String(),
	}
	tax.TaxData.Set(taxData)
	result := repository.DB.Create(&tax)

	return &tax.TraceId, result.Error

}
