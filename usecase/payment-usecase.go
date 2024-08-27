package usecase

import (
	"errors"

	"github.com/Ak8388/applikasi-antrian-dokter-gigi-payment-service/model"
	"github.com/Ak8388/applikasi-antrian-dokter-gigi-payment-service/model/dto"
	"github.com/Ak8388/applikasi-antrian-dokter-gigi-payment-service/repository"
	"github.com/google/uuid"
)

type PaymentUscase interface {
	CustomerPayment(req model.TransModel) (model.PaymentResponse, error)
	TrackingTransactionNotications(req model.TrackerTransaction) error
	CancelPaymentUser(orderId string) error
	ViewPaymentbyUserId(user, status string) (data []dto.PaymentViewDTO, err error)
	FindPaymentByUserAndDate(userID, date string) (dto.PaymentViewDTO, error)
}

type paymentUsecase struct {
	repoPayment repository.PaymentRepository
}

func (p *paymentUsecase) CustomerPayment(req model.TransModel) (model.PaymentResponse, error) {
	if req.CustomerDetails.CustomerID == "" || req.ItemDetails[0].Id == "" || req.ItemDetails[0].Name == "" || req.ItemDetails[0].Category == "" || req.ItemDetails[0].ReservasiDate == "" || req.ItemDetails[0].ReservasiTime == "" || req.TransactionDetails.GrossAmount < 10000 {
		return model.PaymentResponse{}, errors.New("please fill in the data correctly")
	}

	if p.repoPayment.FindPaymentById(req.CustomerDetails.CustomerID, req.ItemDetails[0].To) {
		return model.PaymentResponse{}, errors.New("there seems to be a payment that you have not paid")
	}

	if p.repoPayment.FindPaymentPaidByDate(req.ItemDetails[0].ReservasiDate, req.CustomerDetails.CustomerID) {
		return model.PaymentResponse{}, errors.New("looks like you've already taken the queue on the day you chose")
	}

	randomUid, err := uuid.NewRandom()

	if err != nil {
		return model.PaymentResponse{}, err
	}

	req.TransactionDetails.OrderID = randomUid.String()
	req.ItemDetails[0].Price = uint64(req.TransactionDetails.GrossAmount)
	req.ItemDetails[0].Qty = 1

	return p.repoPayment.CustomerPayment(req)
}

func (p *paymentUsecase) TrackingTransactionNotications(req model.TrackerTransaction) error {
	return p.repoPayment.TrackingTransactionNotications(req)
}

func (p *paymentUsecase) CancelPaymentUser(orderId string) error {
	if orderId == "" {
		return errors.New("order id cannot be empty")
	}

	return p.repoPayment.CancelPaymentUser(orderId)
}

func (p *paymentUsecase) ViewPaymentbyUserId(user, status string) (data []dto.PaymentViewDTO, err error) {

	if user == "" {
		return nil, errors.New("user id can't be empty")
	}

	return p.repoPayment.ViewPaymentbyUserId(user, status)
}

func (p *paymentUsecase) FindPaymentByUserAndDate(userID, date string) (dto.PaymentViewDTO, error) {
	return p.repoPayment.FindPaymentByUserAndDate(userID, date)
}

func NewPaymentUsecase(repoPayment repository.PaymentRepository) PaymentUscase {
	return &paymentUsecase{repoPayment}
}
