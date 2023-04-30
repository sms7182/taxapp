package models

import (
	"time"

	"gorm.io/gorm"
)

type Customer struct {
	Id         uint      `gorm:"autoIncrement,primaryKey"`
	CreatedAt  time.Time `gorm:"column:created_at"`
	FinanceId  string    `gorm:"column:finance_id"`
	Token      string    `gorm:"column:token"`
	PublicKey  string    `gorm:"column:"public_key"`
	PrivateKey string    `gorm:"column:"private_key"`
	ExpireTime time.Time `gorm:"column:"expire_time"`
}

func (obj *Customer) BeforeCreate(_ *gorm.DB) error {
	obj.CreatedAt = time.Now()
	return nil

}

func (Customer) TableName() string {
	return "customers"
}
