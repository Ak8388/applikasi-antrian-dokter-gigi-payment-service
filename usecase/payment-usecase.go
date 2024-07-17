package usecase

import (
	"errors"

	"github.com/Ak8388/applikasi-antrian-dokter-gigi-payment-service/model"
	"github.com/Ak8388/applikasi-antrian-dokter-gigi-payment-service/repository"
	"github.com/google/uuid"
)

type PaymentUscase interface {
	CustomerPayment(req model.TransModel) (model.PaymentResponse, error)
	TrackingTransactionNotications(req model.TrackerTransaction) error
}

type paymentUsecase struct {
	repoPayment repository.PaymentRepository
}

func (p *paymentUsecase) CustomerPayment(req model.TransModel) (model.PaymentResponse, error) {
	if req.CustomerDetails.CustomerID == "" || req.ItemDetails[0].Id == "" || req.ItemDetails[0].Name == "" || req.ItemDetails[0].Category == "" || req.ItemDetails[0].ReservasiDate == "" || req.ItemDetails[0].ReservasiTime == "" || req.TransactionDetails.GrossAmount < 20000 {
		return model.PaymentResponse{}, errors.New("please fill in the data correctly")
	}

	rendomUid, err := uuid.NewRandom()

	if err != nil {
		return model.PaymentResponse{}, err
	}

	req.TransactionDetails.OrderID = rendomUid.String()

	return p.repoPayment.CustomerPayment(req)
}

func (p *paymentUsecase) TrackingTransactionNotications(req model.TrackerTransaction) error {
	return p.repoPayment.TrackingTransactionNotications(req)
}

func NewPaymentUsecase(repoPayment repository.PaymentRepository) PaymentUscase {
	return &paymentUsecase{repoPayment}
}
