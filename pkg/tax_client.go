package pkg

import "tax-management/taxDep/types"

type TaxClient interface {
	SendInvoices(taxRawId *uint, taxProcessId *uint, invoices []types.StandardInvoice) (*types.AsyncResponse, error)
	InquiryByReferences(taxRawId *uint, taxProcessId *uint, refs []string) ([]types.InquiryResult, error)
}
