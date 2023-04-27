package models

import (
	"time"

	"gorm.io/gorm"
)

type Customer struct {
	Id        uint      `gorm:"autoIncrement,primaryKey"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UserName  string    `gorm:"column:user_name"`
	Token     string    `gorm:"column:token"`
}

func (obj *Customer) BeforeCreate(_ *gorm.DB) error {
	obj.CreatedAt = time.Now()
	return nil

}

func (Customer) TableName() string {
	return "customers"
}
