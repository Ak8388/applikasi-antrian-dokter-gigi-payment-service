package repository

import (
	"database/sql"

	"github.com/Ak8388/applikasi-antrian-dokter-gigi-payment-service/model"
)

type PaymentRepository interface {
	CustomerPayment(req model.TransModel) (model.PaymentResponse, error)
	TrackingTransactionNotications(req model.TrackerTransaction) error
}

type paymentRepository struct {
	db *sql.DB
}

func (p *paymentRepository) CustomerPayment(req model.TransModel) (model.PaymentResponse, error) {
	tx, err := p.db.Begin()
	qry := "Insert Into payment (user_id,amount,order_id,status) Values($1,$2,$3,$4)"

	if err != nil {
		return model.PaymentResponse{}, err
	}

	_, err = tx.Exec(qry)

	if err != nil {
		return model.PaymentResponse{}, err
	}

	return model.PaymentResponse{}, nil
}

func (p *paymentRepository) TrackingTransactionNotications(req model.TrackerTransaction) error {
	return nil
}

func NewRepository(db *sql.DB) PaymentRepository {
	return &paymentRepository{db}
}
