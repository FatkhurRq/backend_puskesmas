package dto

type CurrentQueueResponse struct {
	CurrentQueue int    `json:"current_queue"`
	PatientName  string `json:"patient_name"`
	Poli         string `json:"poli"`
	Status       string `json:"status"`
}
