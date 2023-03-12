package models

import (
	"time"

	"github.com/gofrs/uuid"
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

type TaxProcess struct {
	Id          uint      `gorm:"autoIncrement,primaryKey"`
	CreatedAt   time.Time `gorm:"column:created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at"`
	TaxUniqueId string    `gorm:"column:tax_unique_id"`
	TaxType     string    `gorm:"column:tax_type"`
	TaxRawId    uint      `gorm:"column:tax_raw_id"`
	Status      string    `gorm:"column:status"`
}

func (obj *TaxProcess) BeforeCreate(_ *gorm.DB) error {
	id, err := uuid.NewV4()
	if err != nil {
		return err
	}
	obj.TaxUniqueId = id.String()
	obj.CreatedAt = time.Now()
	obj.UpdatedAt = obj.CreatedAt
	obj.Status = InProgress.String()

	return nil
}

func (TaxProcess) TableName() string {
	return "tax_process"
}
