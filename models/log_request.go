package models

import "time"

type LogRequestHistory struct {
	Id           uint      `gorm:"autoIncrement,primaryKey"`
	TraceId      string    `gorm:"column:trace_id"`
	Signature    string    `gorm:"column:signature"`
	CreatedAt    time.Time `gorm:"column:created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at"`
	PacketType   string    `gorm:"column:packet_type"`
	RequestUrl   string    `gorm:"column:request_url"`
	StatusCode   int       `gorm:"column:status_code"`
	Request      string    `gorm:"column:request"`
	Response     *string   `gorm:"column:response"`
	ErrorMessage *string   `gorm:"column:error_message"`
}

func (LogRequestHistory) TableName() string {
	return "log_request_history"
}
