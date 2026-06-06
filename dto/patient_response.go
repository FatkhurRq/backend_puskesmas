package dto

import "github.com/google/uuid"

type PatientResponse struct {
	ID    uuid.UUID `json:"id"`
	NIK   string    `json:"nik"`
	Name  string    `json:"name"`
	Phone string    `json:"phone"`
	Role  string    `json:"role"`
}
