package dto

type PaymentDTO struct {
	DoctorID   string `json:"doctorId"`
	CustomerID string `json:"patientId"`
	QueueDate  string `json:"queueDate"`
	QueueTime  string `json:"queueTime"`
	Note       string `json:"note"`
}
