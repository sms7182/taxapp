package pkg

import (
	"context"
	"tax-management/external"
	"tax-management/external/pg/models"
)

type Repository interface {
	LogReqRes(taxRawId *uint, taxProcessId *uint, requestUniqueId string, apiName string, url string, statusCode int, req string, res *string, errorMsg *string) error
	InsertTaxData(ctx context.Context, rawType string, taxData external.RawTransaction) (taxId uint, taxProcessId uint, err error)
	UpdateTaxReferenceId(ctx context.Context, taxProcessId uint, taxOrgReferenceId string) error
	GetInprogressTaxProcess(ctx context.Context) (taxProcesses []models.RawProcessTaxData, err error)
	UpdateTaxProcessStatus(ctx context.Context, taxProcessId uint, status string) error
}
