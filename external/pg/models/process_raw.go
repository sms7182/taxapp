package models

type RawProcessTaxData struct {
	Id             uint   `gorm:"column:id"`
	TaxRawId       uint   `gorm:"column:tax_raw_id"`
	OrgReferenceId string `gorm:"column:tax_org_reference_id"`
	TaxId          string `gorm:"column:tax_id"`
}
