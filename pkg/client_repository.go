package pkg

import (
	"context"
	"tax-management/external"
	"tax-management/external/pg/models"
	"tax-management/taxDep/types"
)

type Repository interface {
	UpdateTaxProcessStandardInvoice(ctx context.Context, taxProcessId uint, invoice types.StandardInvoice) error
	LogReqRes(taxRawId *uint, taxProcessId *uint, requestUniqueId string, apiName string, url string, statusCode int, req string, res *string, errorMsg *string) error
	IsNotProcessable(ctx context.Context, topic string, trn string) bool
	NumberOfFailureExceeded() bool
	InsertTaxData(ctx context.Context, rawType string, taxData external.RawTransaction, companyName string) (rawDataId uint, taxProcessId uint, taxId string, err error)
	UpdateTaxReferenceId(ctx context.Context, taxProcessId uint, taxOrgReferenceId string, taxOrgInternalTrn *string, taxOrgInquiryUuid *string) error
	GetInProgressTaxProcess(ctx context.Context) (taxProcesses []models.RawProcessTaxData, err error)
	UpdateTaxProcessStatus(ctx context.Context, taxProcessId uint, status string, confirmationReferenceId *string) error
	GetFailedTaxProcess(ctx context.Context) (failedTaxProcess []models.FailedTaxProcess, err error)
	UpdateNotifyFailedOfTaxProcess(ctx context.Context, ids []uint) error
	GetByTaxRawId(ctx context.Context, taxProcessId uint) (*models.TaxRawDomain, error)
	GetReadyTaxToRetry(ctx context.Context) ([]models.TaxRawDomain, error)
	GetTaxProcess(ctx context.Context, id uint) (*models.TaxProcess, error)
	GetUserName(ctx context.Context, usrname string) (*models.Customer, error)
	CreateCustomer(ctx context.Context, customer models.Customer) (*uint, error)
}
