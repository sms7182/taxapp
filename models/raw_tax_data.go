package models

import (
	"time"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgtype"
	"gorm.io/gorm"
)

type TaxStatus int

const (
	InProgress TaxStatus = iota
	Retry
	Failed
	Completed
)

func (ts TaxStatus) String() string {
	return []string{"inprogress", "retry", "failed", "completed"}[ts]
}

type TaxRawDomain struct {
	Id            uint         `gorm:"autoIncrement,primaryKey"`
	TraceId       string       `gorm:"column:trace_id"`
	CreatedAt     time.Time    `gorm:"column:created_at"`
	TaxData       pgtype.JSONB `gorm:"type:jsonb;default:'[]'"`
	TaxType       string       `gorm:"column:tax_type"`
	ProcessStatus string       `gorm:"column:process_status"`
	UpdatedAt     time.Time    `gorm:"column:updated_at"`
}

func (obj *TaxRawDomain) BeforeCreate(_ *gorm.DB) error {
	id, err := uuid.NewV4()
	if err != nil {
		return err
	}
	obj.TraceId = obj.TaxType + "-" + id.String()
	obj.CreatedAt = time.Now()
	obj.UpdatedAt = obj.CreatedAt
	return nil
}

func (TaxRawDomain) TableName() string {
	return "tax_raw_data"
}
