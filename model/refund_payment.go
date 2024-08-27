package model

type UserRefundPayment struct {
	ID          string `json:"id"`
	UserID      string `json:"userId"`
	PaymentId   string `json:"paymentId"`
	BankName    string `json:"bankName"`
	BankNumber  string `json:"bankNumber"`
	Description string `json:"description"`
	Status      string `json:"status"`
	Date        string `json:"date"`
	Funds       int    `json:"fund"`
}

func (uRePay UserRefundPayment) ValidateStatus() bool {
	return uRePay.Status == "Created" || uRePay.Status == "Approve" || uRePay.Status == "Finish"
}
