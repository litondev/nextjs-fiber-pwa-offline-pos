package models

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	ID          uint    `gorm:"primaryKey;autoIncrement"`
	Name        string  `gorm:"size:50"`
	Code        string  `gorm:"size:25"`
	Description *string `gorm:"type:text"`
	Stock       int     `gorm:"default:0"`
	Price       float32 `gorm:"type:decimal(10,2);default:0.00`
	Photo       *string `gorm:"size:25"`
	CategoryID  *uint
	CreatedAt   *time.Time
	UpdatedAt   *time.Time
	DeletedAt   *gorm.DeletedAt
}
