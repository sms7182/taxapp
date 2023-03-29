package models

import (
	"github.com/jackc/pgtype"
	"time"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type TaxStatus int

const (
	TaxStatusSending TaxStatus = iota
	TaxStatusInProgress
	TaxStatusFailed
	TaxStatusCompleted
	TaxStatusUnnecessary
)

func (ts TaxStatus) String() string {
	return []string{"sending", "in-progress", "failed", "completed", "unnecessary"}[ts]
}

type TaxProcess struct {
	Id                      uint         `gorm:"autoIncrement,primaryKey"`
	CreatedAt               time.Time    `gorm:"column:created_at"`
	UpdatedAt               time.Time    `gorm:"column:updated_at"`
	TaxUniqueId             string       `gorm:"column:tax_unique_id"`
	TaxType                 string       `gorm:"column:tax_type"`
	TaxRawId                uint         `gorm:"column:tax_raw_id"`
	Status                  string       `gorm:"column:status"`
	TaxOrgReferenceId       *string      `gorm:"column:tax_org_reference_id"`
	TaxId                   *string      `gorm:"column:tax_id"`
	InternalTrn             *string      `gorm:"column:internal_trn"`
	InquiryUuid             *string      `gorm:"column:inquiry_uuid"`
	ConfirmationReferenceId *string      `gorm:"confirmation_reference_id"`
	CompanyName             *string      `gorm:"company_name"`
	StandardInvoice         pgtype.JSONB `gorm:"type:jsonb;default:'[]'"`
}

func (obj *TaxProcess) BeforeCreate(_ *gorm.DB) error {
	id, err := uuid.NewV4()
	if err != nil {
		return err
	}
	obj.TaxUniqueId = id.String()
	obj.CreatedAt = time.Now()
	obj.UpdatedAt = obj.CreatedAt
	obj.Status = TaxStatusSending.String()

	return nil
}

func (TaxProcess) TableName() string {
	return "tax_process"
}

type TaxProcessHistory struct {
	Id                      uint         `gorm:"autoIncrement,primaryKey"`
	TaxProcessId            uint         `gorm:"column:tax_process_id"`
	CreatedAt               time.Time    `gorm:"column:created_at"`
	TaxUniqueId             string       `gorm:"column:tax_unique_id"`
	TaxType                 string       `gorm:"column:tax_type"`
	TaxRawId                uint         `gorm:"column:tax_raw_id"`
	Status                  string       `gorm:"column:status"`
	TaxOrgReferenceId       *string      `gorm:"column:tax_org_reference_id"`
	TaxId                   *string      `gorm:"column:tax_id"`
	InternalTrn             *string      `gorm:"column:internal_trn"`
	InquiryUuid             *string      `gorm:"column:inquiry_uuid"`
	ConfirmationReferenceId *string      `gorm:"confirmation_reference_id"`
	CompanyName             *string      `gorm:"company_name"`
	StandardInvoice         pgtype.JSONB `gorm:"type:jsonb;default:'[]'"`
}

func (TaxProcessHistory) TableName() string {
	return "tax_process_history"
}

func ToTaxProcessHistory(tp TaxProcess) TaxProcessHistory {
	return TaxProcessHistory{
		TaxProcessId:            tp.Id,
		CreatedAt:               time.Now(),
		TaxUniqueId:             tp.TaxUniqueId,
		TaxType:                 tp.TaxType,
		TaxRawId:                tp.TaxRawId,
		Status:                  tp.Status,
		TaxOrgReferenceId:       tp.TaxOrgReferenceId,
		TaxId:                   tp.TaxId,
		InternalTrn:             tp.InternalTrn,
		InquiryUuid:             tp.InquiryUuid,
		ConfirmationReferenceId: tp.ConfirmationReferenceId,
		CompanyName:             tp.CompanyName,
		StandardInvoice:         tp.StandardInvoice,
	}
}
