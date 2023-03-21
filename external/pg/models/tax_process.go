package models

import (
	"time"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type TaxStatus int

const (
	Sending TaxStatus = iota
	InProgress
	Retry
	Failed
	Completed
)

func (ts TaxStatus) String() string {
	return []string{"sending", "in-progress", "retry", "failed", "completed"}[ts]
}

type TaxProcess struct {
	Id                uint      `gorm:"autoIncrement,primaryKey"`
	CreatedAt         time.Time `gorm:"column:created_at"`
	UpdatedAt         time.Time `gorm:"column:updated_at"`
	TaxUniqueId       string    `gorm:"column:tax_unique_id"`
	TaxType           string    `gorm:"column:tax_type"`
	TaxRawId          uint      `gorm:"column:tax_raw_id"`
	Status            string    `gorm:"column:status"`
	TaxOrgReferenceId *string   `gorm:"column:tax_org_reference_id"`
}

func (obj *TaxProcess) BeforeCreate(_ *gorm.DB) error {
	id, err := uuid.NewV4()
	if err != nil {
		return err
	}
	obj.TaxUniqueId = id.String()
	obj.CreatedAt = time.Now()
	obj.UpdatedAt = obj.CreatedAt
	obj.Status = Sending.String()

	return nil
}

func (TaxProcess) TableName() string {
	return "tax_process"
}