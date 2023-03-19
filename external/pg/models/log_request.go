package models

import "time"

type TaxOfficeRequestResponseLogModel struct {
	Id             uint      `gorm:"autoIncrement,primaryKey"`
	TaxRawId       *uint     `gorm:"column:tax_raw_id"`
	TaxProcessId   *uint     `gorm:"column:tax_process_id"`
	RequestUiqueId string    `gorm:"column:request_unique_id"`
	ApiName        string    `gorm:"column:api_name"`
	LoggedAt       time.Time `gorm:"column:logged_at"`
	Url            string    `gorm:"column:url"`
	StatusCode     int       `gorm:"column:status_code"`
	Request        string    `gorm:"column:request"`
	Response       *string   `gorm:"column:response"`
	ErrorMessage   *string   `gorm:"column:error_message"`
}

func (TaxOfficeRequestResponseLogModel) TableName() string {
	return "tax_office_request_response_log"
}
