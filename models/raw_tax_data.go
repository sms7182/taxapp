package models

import (
	"time"

	"github.com/jackc/pgtype"
	"gorm.io/gorm"
)

type TaxRawDomain struct {
	Id        uint         `gorm:"autoIncrement,primaryKey"`
	CreatedAt time.Time    `gorm:"column:created_at"`
	TaxData   pgtype.JSONB `gorm:"type:jsonb;default:'[]'"`
}

func (obj *TaxRawDomain) BeforeCreate(_ *gorm.DB) error {

	obj.CreatedAt = time.Now()
	return nil
}

func (TaxRawDomain) TableName() string {
	return "tax_raw_data"
}
