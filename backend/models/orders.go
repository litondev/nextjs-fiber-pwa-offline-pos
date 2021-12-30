package models

import (
	"time"

	"gorm.io/gorm"
)

type Order struct {
	gorm.Model
	ID           uint            `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserID       *uint           `gorm:"column:user_id" json:"user_id"`
	CustomerID   *uint           `gorm:"column:customer_id" json:"customer_id"`
	Total        float32         `gorm:"column:total;type:decimal(10,2);default:0.00" json:"total"`
	CreatedAt    *time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt    *time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt    *gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
	User         User            `gorm:"Foreignkey:UserID;association_foreignkey:ID;" json:"user"`
	Customer     Customer        `gorm:"Foreignkey:CustomerID;association_foreignkey:ID;" json:"customer"`
	OrderDetails []DetailOrder   `gorm:"Foreignkey:OrderID;association_foreignkey:ID;" json:"order"`
}
