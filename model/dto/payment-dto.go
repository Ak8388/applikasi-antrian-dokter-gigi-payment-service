package dto

type PaymentDTO struct {
	ID         string `json:"id"`
	DoctorID   string `json:"doctorId"`
	CustomerID string `json:"patientId"`
	QueueDate  string `json:"queueDate"`
	QueueTime  string `json:"queueTime"`
	Note       string `json:"note"`
}

type PaymentViewDTO struct {
	Id        string `json:"id"`
	UserId    string `json:"userId"`
	DoctorId  string `json:"doctorId"`
	Amount    uint32 `json:"amount"`
	OrderId   string `json:"orderId"`
	QueueDate string `json:"queueDate"`
	QueueTime string `json:"queueTime"`
	Status    string `json:"status"`
	To        string `json:"to_stts"`
	CreatedAt string `json:"date"`
}
