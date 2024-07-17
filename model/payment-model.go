package model

type TransModel struct {
	TransactionDetails TransactionDetails `json:"transaction_detail"`
	ItemDetails        []ItemDetails      `json:"item_details"`
	CustomerDetails    CustomerDetails    `json:"customer_detail"`
}

type TransactionDetails struct {
	OrderID     string `json:"order_id"`
	GrossAmount uint   `json:"gross_amount"`
}

type ItemDetails struct {
	Id            string `json:"id"`
	Price         string `json:"price"`
	Name          string `json:"name"`
	Category      string `json:"category"`
	ReservasiDate string `json:"reservasi_date"`
	ReservasiTime string `json:"reservasi_time"`
	Note          string `json:"note"`
}

type CustomerDetails struct {
	CustomerID string `json:"phone"`
}

type PaymentResponse struct {
	Token       string `json:"token"`
	RedirectURL string `json:"redirect_url"`
}
