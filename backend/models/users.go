package models

import (
	"time"
)

type User struct {
	ID        uint   `gorm:"primaryKey;autoIncrement"`
	Name      string `gorm:"size:50"`
	Email     string `gorm:"unique;size:50"`
	Password  string `gorm:"type:text"`
	CreatedAt *time.Time
	UpdatedAt *time.Time
	DeletedAt *time.Time
}
