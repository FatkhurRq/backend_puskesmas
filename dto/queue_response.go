package dto

type QueueResponse struct {
	QueueNumber int    `json:"queue_number"`
	PatientName string `json:"patient_name"`
	Poli        string `json:"poli"`
	Status      string `json:"status"`
}
