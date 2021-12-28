package models

import (
	"time"
)

type Customer struct {
	ID        uint    `gorm:"primaryKey;autoIncrement"`
	Name      string  `gorm:"size:25"`
	Email     *string `gorm:"unique;size:25"`
	Address   *string `gorm:"type:text"`
	Phone     *string `gorm:"size:25"`
	CreatedAt *time.Time
	UpdatedAt *time.Time
	DeletedAt *time.Time
}
