package pkg

import "tax-management/taxDep/types"

type TaxClient interface {
	SendInvoices(taxRawId *uint, taxProcessId *uint, invoices []types.StandardInvoice, privateKey string, customerid string) (*types.AsyncResponse, error)
	InquiryByReferences(taxRawId *uint, taxProcessId *uint, refs []string, privateKey string) ([]types.InquiryResult, error)
}
