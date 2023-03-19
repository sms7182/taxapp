package pkg

import (
	"context"
	"tax-management/external"
)

type Repository interface {
	LogReqRes(taxRawId *uint, taxProcessId *uint, requestUniqueId string, apiName string, url string, statusCode int, req string, res *string, errorMsg *string) error
	InsertTaxData(ctx context.Context, rawType string, taxData external.RawTransaction) (*string, error)
}
