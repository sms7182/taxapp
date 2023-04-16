package models

import "github.com/jackc/pgtype"

type RawProcessTaxData struct {
	Id             uint   `gorm:"column:id"`
	TaxRawId       uint   `gorm:"column:tax_raw_id"`
	OrgReferenceId string `gorm:"column:tax_org_reference_id"`
	TaxId          string `gorm:"column:tax_id"`
}

type FailedTaxProcess struct {
	Id          uint         `gorm:"column:id"`
	Response    pgtype.JSONB `gorm:"type:jsonb;default:'[]'"`
	InternalTrn string       `gorm:"column:internal_trn"`
}

type FailedTaxData struct {
	Id              uint         `gorm:"column:id"`
	Response        pgtype.JSONB `gorm:"type:jsonb;default:'[]'"`
	InternalTrn     string       `gorm:"column:internal_trn"`
	StandardInvoice pgtype.JSONB `gorm:"type:jsonb;default:'[]'"`
}
