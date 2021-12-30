package models

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	ID          uint            `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name        string          `gorm:"column:name;size:50" json:"name"`
	Code        string          `gorm:"column:code;size:25" json:"code"`
	Description *string         `gorm:"column:description;type:text" json:"description"`
	Stock       int             `gorm:"column:stock;default:0" json:"stock"`
	Price       float32         `gorm:"column:price;type:decimal(10,2);default:0.00" json:"price"`
	Photo       *string         `gorm:"column:photo;size:25" json:"photo"`
	CategoryID  *uint           `gorm:"column:category_id" json:"category_id"`
	Category    Category        `gorm:"Foreignkey:CategoryID;association_foreignkey:ID;" json:"category"`
	CreatedAt   *time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   *time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt   *gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
}
