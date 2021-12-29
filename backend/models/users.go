package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID        uint   `gorm:"primaryKey;autoIncrement"`
	Name      string `gorm:"size:50"`
	Email     string `gorm:"unique;size:50"`
	Password  string `gorm:"type:text"`
	CreatedAt *time.Time
	UpdatedAt *time.Time
	DeletedAt *gorm.DeletedAt
}
