package dto

type AllQueueResponse struct {
	QueueID     string `json:"queue_id"`
	QueueNumber int    `json:"queue_number"`
	PatientName string `json:"patient_name"`
	Poli        string `json:"poli"`
	Status      string `json:"status"`
	QueueDate   string `json:"queue_date"`
}
