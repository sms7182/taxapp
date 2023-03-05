package models

import (
	"time"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgtype"
	"gorm.io/gorm"
)

type TaxStatus int

const (
	InProcess TaxStatus = iota
	Retry
	Failed
	Completed
)

func (ts TaxStatus) String() string {
	return []string{"inprocess", "retry", "failed", "completed"}[ts]
}

type RawTaxDomain struct {
	Id        uint         `gorm:"autoIncrement,primaryKey"`
	TraceId   string       `gorm:"column:trace_id"`
	CreatedAt time.Time    `gorm:"column:created_at"`
	TaxData   pgtype.JSONB `gorm:"type:jsonb;default:'[]'"`
	TaxType   string       `gorm:"column:tax_type"`
	Status    string       `gorm:"column:status"`
	UpdatedAt time.Time    `gorm:"column:updated_at"`
}

func (obj *RawTaxDomain) BeforeCreate(_ *gorm.DB) error {
	id, err := uuid.NewV4()
	if err != nil {
		return err
	}
	obj.TraceId = obj.TaxType + "-" + id.String()
	return nil
}

func (RawTaxDomain) TableName() string {
	return "raw_tax"
}
