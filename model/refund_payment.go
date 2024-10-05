package model

import "time"

type UserRefundPayment struct {
	ID          string    `json:"id"`
	UserID      string    `json:"userId"`
	PaymentId   string    `json:"paymentId"`
	BankName    string    `json:"bankName"`
	BankNumber  string    `json:"bankNumber"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	Date        string    `json:"date"`
	Funds       int       `json:"fund"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func (uRePay UserRefundPayment) ValidateStatus() bool {
	return uRePay.Status == "Created" || uRePay.Status == "Approve" || uRePay.Status == "Finish"
}
