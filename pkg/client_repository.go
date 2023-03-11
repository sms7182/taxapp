package pkg

import (
	"context"
	"tax-management/external/exkafka/messages"
)

type ClientRepository interface {
	LogReqRes(requestTraceId string, signature string, packetType string, url string, statusCode int, req string, res *string, errorMsg *string) error
	InsertTaxData(ctx context.Context, taxData messages.RawTransaction) (*string, error)
}
