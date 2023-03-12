package pkg

import (
	"context"
	"tax-management/external/exkafka/messages"
)

type ClientRepository interface {
	LogReqRes(taxRawId *uint, taxProcessId *uint, requestUniqueId string, apiName string, url string, statusCode int, req string, res *string, errorMsg *string) error
	InsertTaxData(ctx context.Context, taxData messages.RawTransaction) (*string, error)
}
