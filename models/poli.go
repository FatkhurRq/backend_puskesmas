package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Poli struct {
	ID   uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name string    `gorm:"not null;unique"`
}

func (p *Poli) BeforeCreate(tx *gorm.DB) error {
	p.ID = uuid.New()
	return nil
}
