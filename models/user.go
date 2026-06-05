// models/user.go
package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID       uuid.UUID `gorm:"type:uuid;primaryKey"`
	NIK      string    `gorm:"unique;not null"`
	Name     string    `gorm:"not null"`
	Phone    string    `gorm:"not null"`
	Password string    `json:"-"`
	Role     string    `gorm:"default:patient"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	u.ID = uuid.New()
	return nil
}
