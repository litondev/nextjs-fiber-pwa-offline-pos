package models

import (
	"time"

	"gorm.io/gorm"
)

type Category struct {
	gorm.Model
	ID          uint            `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name        string          `gorm:"column:name;size:50" json:"name"`
	Description *string         `gorm:"column:description;type:text" json:"description"`
	CreatedAt   *time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   *time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt   *gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
}
