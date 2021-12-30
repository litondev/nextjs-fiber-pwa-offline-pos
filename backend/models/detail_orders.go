package models

import (
	"time"

	"gorm.io/gorm"
)

type DetailOrder struct {
	gorm.Model
	ID        uint            `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	OrderID   *uint           `gorm:"column:order_id" json:"order_id"`
	ProductID *uint           `gorm:"column:product_id" json:"product_id"`
	Qty       int             `gorm:"column:qty;default:1" json:"qty"`
	Price     float32         `gorm:"column:price;type:decimal(10,2);default:0.00" json:"price"`
	CreatedAt *time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt *time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt *gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
	Order     Order           `gorm:"Foreignkey:OrderID;association_foreignkey:ID;" json:"order"`
	Product   Product         `gorm:"Foreignkey:ProductID;association_foreignkey:ID;" json:"product"`
}
