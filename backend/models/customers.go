package models

import (
	"time"

	"gorm.io/gorm"
)

type Customer struct {
	gorm.Model
	ID        uint            `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name      string          `gorm:"column:name;size:25" json:"name"`
	Email     *string         `gorm:"column:email;unique;size:25" json:"email"`
	Address   *string         `gorm:"column:address;type:text" json:"address"`
	Phone     *string         `gorm:"column:phone;size:25" json:"phone"`
	CreatedAt *time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt *time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt *gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
}
