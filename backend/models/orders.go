package models

import (
	"time"
)

type Order struct {
	ID         uint `gorm:"primaryKey;autoIncrement"`
	UserID     *uint
	CustomerID *uint
	Total      float32 `gorm:"type:decimal(10,2);default:0.00"`
	CreatedAt  *time.Time
	UpdatedAt  *time.Time
	DeletedAt  *time.Time
}
