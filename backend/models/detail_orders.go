package models

import (
	"time"

	"gorm.io/gorm"
)

type DetailOrder struct {
	gorm.Model
	ID        uint `gorm:"primaryKey;autoIncrement"`
	OrderID   *uint
	ProductID *uint
	Qty       int     `gorm:"default:1"`
	Price     float32 `gorm:type:decimal(10,2);default:0.00`
	CreatedAt *time.Time
	UpdatedAt *time.Time
	DeletedAt *gorm.DeletedAt
}
