package model

type TransModel struct {
	TransactionDetails TransactionDetails `json:"transaction_detail"`
	CustomerDetails    CustomerDetails    `json:"customer_detail"`
}

type TransactionDetails struct {
	OrderID     string `json:"order_id"`
	GrossAmount uint   `json:"gross_amount"`
}

type CustomerDetails struct {
	CustomerID string `json:"customerID"`
}

type PaymentResponse struct {
	Token       string `json:"token"`
	RedirectURL string `json:"redirect_url"`
}
