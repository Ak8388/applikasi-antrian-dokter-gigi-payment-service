package model

type TransModel struct {
	TransactionDetails TransactionDetails `json:"transaction_details"`
	ItemDetails        []ItemDetails      `json:"item_details"`
	CustomerDetails    CustomerDetails    `json:"customer_detail"`
}

type TransactionDetails struct {
	OrderID     string `json:"order_id"`
	GrossAmount uint64 `json:"gross_amount"`
}

type ItemDetails struct {
	Id            string `json:"id"`
	Name          string `json:"name"`
	Category      string `json:"category"`
	Price         uint64 `json:"price"`
	Qty           uint64 `json:"quantity"`
	ReservasiDate string `json:"reservasi_date"`
	ReservasiTime string `json:"reservasi_time"`
	Note          string `json:"note"`
	To            string `json:"to"`
}

type CustomerDetails struct {
	CustomerID string `json:"phone"`
	ReservID   string `json:"idResev"`
}

type PaymentResponse struct {
	Token              string             `json:"token"`
	RedirectURL        string             `json:"redirect_url"`
	TransactionDetails TransactionDetails `json:"transaction_details"`
}
