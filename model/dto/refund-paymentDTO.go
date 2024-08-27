package dto

type RefundPaymentChangeData struct {
	Id          string `json:"id"`
	BankName    string `json:"bankName"`
	BankNumber  string `json:"bankNumber"`
	Description string `json:"description"`
}
