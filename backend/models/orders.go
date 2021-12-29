package models

import (
	"time"

	"gorm.io/gorm"
)

type Order struct {
	gorm.Model
	ID         uint `gorm:"primaryKey;autoIncrement"`
	UserID     *uint
	CustomerID *uint
	Total      float32 `gorm:"type:decimal(10,2);default:0.00"`
	CreatedAt  *time.Time
	UpdatedAt  *time.Time
	DeletedAt  *gorm.DeletedAt
}
