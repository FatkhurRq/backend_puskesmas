package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Registration struct {
	ID               uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID           uuid.UUID `gorm:"type:uuid;not null"`
	PoliID           uuid.UUID `gorm:"type:uuid;not null"`
	RegistrationDate time.Time
	Status           string

	User User `gorm:"foreignKey:UserID"`
	Poli Poli `gorm:"foreignKey:PoliID"`
}

func (r *Registration) BeforeCreate(tx *gorm.DB) error {
	r.ID = uuid.New()
	return nil
}
