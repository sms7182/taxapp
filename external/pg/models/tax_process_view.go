package models

import (
	"tax-management/taxDep/types"
	"time"
)

type TaxProcessViewModel struct {
	Id              uint                  `json:"id"`
	CreatedAt       time.Time             `json:"createdAt"`
	UpdatedAt       time.Time             `json:"updatedAt"`
	TaxType         string                `json:"taxType"`
	Status          string                `json:"status"`
	TaxId           *string               `json:"taxId"`
	InternalTrn     *string               `json:"internalTrn"`
	CompanyName     *string               `gorm:"company_name"`
	StandardInvoice types.StandardInvoice `json:"standardInvoice"`
}
