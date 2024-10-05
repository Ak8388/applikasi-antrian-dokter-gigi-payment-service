package usecase

import (
	"errors"

	"github.com/Ak8388/applikasi-antrian-dokter-gigi-payment-service/model"
	"github.com/Ak8388/applikasi-antrian-dokter-gigi-payment-service/model/dto"
	"github.com/Ak8388/applikasi-antrian-dokter-gigi-payment-service/repository"
)

type RefundPaymentUsecase interface {
	CreateRefundPayment(model.UserRefundPayment) error
	ChangeStatus(id, status string) error
	ChangePayment(request dto.RefundPaymentChangeData) error
	GetDataRefundPaymentForPatient(idUser, status string) ([]model.UserRefundPayment, error)
	GetDataRefundPaymentForAdmin(status string) ([]model.UserRefundPayment, error)
	DeleteDataRefund(id string) error
}

type refaundPaymentUsecase struct {
	paymentUC         PaymentUscase
	refundPaymentRepo repository.RefundPaymentRepository
}

func (r *refaundPaymentUsecase) CreateRefundPayment(request model.UserRefundPayment) error {
	var feeDeduction float64
	if len(request.BankNumber) < 10 {
		return errors.New("please enter the account number correctly")
	}

	if request.UserID == "" {
		return errors.New("user id can't be empty")
	}

	res, err := r.paymentUC.FindPaymentByUserAndDate(request.UserID, request.Date)

	if err != nil || res.Id == "" {
		return errors.New("payment not found")
	}

	request.PaymentId = res.Id
	feeDeduction = float64(res.Amount) * (10.0 / 100.0)

	request.Funds = int(res.Amount) - int(feeDeduction)

	if !request.ValidateStatus() {
		return errors.New("status invalid")
	}

	return r.refundPaymentRepo.CreateRefundPayment(request)
}

func (r *refaundPaymentUsecase) ChangeStatus(id, status string) error {
	if id == "" {
		return errors.New("refund payment id can't be empty")
	}

	req := model.UserRefundPayment{
		ID:          id,
		UserID:      id,
		PaymentId:   id,
		BankNumber:  "",
		Description: "",
		Status:      status,
		Funds:       0,
	}

	if !req.ValidateStatus() {
		return errors.New("status invalid")
	}

	return r.refundPaymentRepo.ChangeStatus(id, status)
}

func (r *refaundPaymentUsecase) GetDataRefundPaymentForPatient(idUser, status string) ([]model.UserRefundPayment, error) {
	return r.refundPaymentRepo.GetDataRefundPaymentForPatient(idUser, status)
}

func (r *refaundPaymentUsecase) GetDataRefundPaymentForAdmin(status string) ([]model.UserRefundPayment, error) {
	return r.refundPaymentRepo.GetDataRefundPaymentForAdmin(status)
}

func (r *refaundPaymentUsecase) DeleteDataRefund(id string) error {
	if id == "" {
		return errors.New("refund id can't be empty")
	}

	res, err := r.refundPaymentRepo.GetRefundByID(id)

	if err != nil {
		return err
	}

	if res.Status != "Finish" {
		return errors.New("refund payment not valid")
	}

	return r.refundPaymentRepo.DeleteDataRefund(id)
}

func (r *refaundPaymentUsecase) ChangePayment(request dto.RefundPaymentChangeData) error {
	res, err := r.refundPaymentRepo.GetRefundByID(request.Id)

	if err != nil || res.ID == "" {
		return errors.New("data not found")
	}

	return r.refundPaymentRepo.ChangePayment(request)
}

func NewPaymentrefund(refundPaymentRepo repository.RefundPaymentRepository, paymentUC PaymentUscase) RefundPaymentUsecase {
	return &refaundPaymentUsecase{paymentUC, refundPaymentRepo}
}
