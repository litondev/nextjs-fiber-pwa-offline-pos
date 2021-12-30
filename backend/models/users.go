package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID        uint            `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string          `gorm:"column:name;size:50" json:"name"`
	Email     string          `gorm:"column:email;unique;size:50" json:"email"`
	Password  string          `gorm:"column:password;type:text" json:"password"`
	CreatedAt *time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt *time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt *gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
}
