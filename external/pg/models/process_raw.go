package models

import "github.com/jackc/pgtype"

type RawProcessTaxData struct {
	Id             uint         `gorm:"column:id"`
	TaxRawId       uint         `gorm:"column:tax_raw_id"`
	OrgReferenceId string       `gorm:"column:tax_org_reference_id"`
	TaxData        pgtype.JSONB `gorm:"column:tax_data"`
}
