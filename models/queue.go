package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Queue struct {
	ID             uuid.UUID `gorm:"type:uuid;primaryKey"`
	RegistrationID uuid.UUID `gorm:"type:uuid;not null"`
	QueueNumber    int
	QueueDate      time.Time
	Status         string

	Registration Registration `gorm:"foreignKey:RegistrationID"`
}

func (q *Queue) BeforeCreate(tx *gorm.DB) error {
	q.ID = uuid.New()
	return nil
}
