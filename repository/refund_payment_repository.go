package repository

import (
	"database/sql"
	"strconv"

	"github.com/Ak8388/applikasi-antrian-dokter-gigi-payment-service/model"
	"github.com/Ak8388/applikasi-antrian-dokter-gigi-payment-service/model/dto"
)

type RefundPaymentRepository interface {
	CreateRefundPayment(model.UserRefundPayment) error
	ChangeStatus(id, status string) error
	ChangePayment(request dto.RefundPaymentChangeData) error
	GetDataRefundPaymentForPatient(idUser, status string) ([]model.UserRefundPayment, error)
	GetDataRefundPaymentForAdmin(status string) ([]model.UserRefundPayment, error)
	GetRefundByID(id string) (model.UserRefundPayment, error)
	DeleteDataRefund(id string) error
}

type refundPaymentRepository struct {
	db *sql.DB
}

func (r *refundPaymentRepository) CreateRefundPayment(refReq model.UserRefundPayment) error {
	qry := "Insert Into refaund_payment (user_id, payment_id, bank_name,bank_number, description, status, funds) Values($1, $2, $3 ,$4 ,$5 ,$6, $7)"

	_, err := r.db.Exec(qry, refReq.UserID, refReq.PaymentId, refReq.BankName, refReq.BankNumber, refReq.Description, refReq.Status, refReq.Funds)

	return err
}

func (r *refundPaymentRepository) ChangeStatus(id, status string) error {
	qry := "Update refaund_payment Set status=$1 Where id=$2"

	_, err := r.db.Exec(qry, status, id)

	return err
}

func (r *refundPaymentRepository) GetDataRefundPaymentForPatient(idUser, status string) (data []model.UserRefundPayment, err error) {
	qry := "Select * From refaund_payment Where user_id=$1"
	var value []interface{}
	index := 1
	value = append(value, idUser)

	if status != "" {
		index++
		qry += " AND status=$" + strconv.Itoa(index)
		value = append(value, status)
	}
	qry += " Order By created_at DESC"
	row, err := r.db.Query(qry, value...)

	if err != nil {
		return nil, err
	}

	for row.Next() {
		var userRefund model.UserRefundPayment
		err = row.Scan(&userRefund.ID, &userRefund.UserID, &userRefund.PaymentId, &userRefund.BankName, &userRefund.BankNumber, &userRefund.Description, &userRefund.Status, &userRefund.Funds, &userRefund.CreatedAt, &userRefund.UpdatedAt)

		if err != nil {
			return nil, err
		}

		data = append(data, userRefund)
	}

	return
}

func (r *refundPaymentRepository) GetDataRefundPaymentForAdmin(status string) (data []model.UserRefundPayment, err error) {
	qry := ""
	var value []interface{}

	if status != "" {
		qry = "Select * From refaund_payment Where status=$1"
		value = append(value, status)
	} else {
		qry = "Select * From refaund_payment"
	}
	qry += " Order By created_at DESC"
	row, err := r.db.Query(qry, value...)

	if err != nil {
		return nil, err
	}

	for row.Next() {
		var userRefund model.UserRefundPayment
		err = row.Scan(&userRefund.ID, &userRefund.UserID, &userRefund.PaymentId, &userRefund.BankName, &userRefund.BankNumber, &userRefund.Description, &userRefund.Status, &userRefund.Funds, &userRefund.CreatedAt, &userRefund.UpdatedAt)
		if err != nil {
			return nil, err
		}

		data = append(data, userRefund)
	}

	return
}

func (r *refundPaymentRepository) GetRefundByID(id string) (userRefund model.UserRefundPayment, err error) {
	qry := "Select * from refaund_payment Where id=$1"

	err = r.db.QueryRow(qry, id).Scan(&userRefund.ID, &userRefund.UserID, &userRefund.PaymentId, &userRefund.BankName, &userRefund.BankNumber, &userRefund.Description, &userRefund.Status, &userRefund.Funds, &userRefund.CreatedAt, &userRefund.UpdatedAt)

	return
}

func (r *refundPaymentRepository) DeleteDataRefund(id string) error {
	qry := "Delete from refaund_payment Where id=$1"

	_, err := r.db.Exec(qry, id)

	return err
}

func (r *refundPaymentRepository) ChangePayment(request dto.RefundPaymentChangeData) error {
	qry := "Update refaund_payment Set bank_name=$1, bank_number=$2 Where id=$3"

	_, err := r.db.Exec(qry, request.BankName, request.BankNumber, request.Id)

	return err
}

func NewRefundPayment(db *sql.DB) RefundPaymentRepository {
	return &refundPaymentRepository{db}
}
