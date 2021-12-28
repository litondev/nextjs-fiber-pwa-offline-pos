package models

import (
	"time"
)

type DetailOrder struct {
	ID        uint `gorm:"primaryKey;autoIncrement"`
	OrderID   *uint
	ProductID *uint
	Qty       int     `gorm:"default:1"`
	Price     float32 `gorm:type:decimal(10,2);default:0.00`
	CreatedAt *time.Time
	UpdatedAt *time.Time
	DeletedAt *time.Time
}