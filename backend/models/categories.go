package models

import (
	"time"

	"gorm.io/gorm"
)

type Category struct {
	gorm.Model
	ID          uint    `gorm:"primaryKey;autoIncrement"`
	Name        string  `gorm:"size:50"`
	Description *string `gorm:"type:text"`
	CreatedAt   *time.Time
	UpdatedAt   *time.Time
	DeletedAt   *gorm.DeletedAt
}
