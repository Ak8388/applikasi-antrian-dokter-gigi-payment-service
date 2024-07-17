package repository

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"

	"github.com/Ak8388/applikasi-antrian-dokter-gigi-payment-service/model"
	"github.com/Ak8388/applikasi-antrian-dokter-gigi-payment-service/model/dto"
	"github.com/joho/godotenv"
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

	if err != nil {
		return model.PaymentResponse{}, err
	}

	if err := godotenv.Load(); err != nil {
		return model.PaymentResponse{}, err
	}

	var response model.PaymentResponse
	client := &http.Client{}

	dataJson, err := json.Marshal(req)

	if err != nil {
		return model.PaymentResponse{}, err
	}

	request, errReq := http.NewRequest("POST", "https://app.sandbox.midtrans.com/snap/v1/transactions", bytes.NewBuffer(dataJson))

	if errReq != nil {
		return model.PaymentResponse{}, errReq
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.SetBasicAuth(os.Getenv("SERVER_KEY"), "")

	res, errRes := client.Do(request)

	if errRes != nil {
		return model.PaymentResponse{}, errRes
	}

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	marshalErr := json.Unmarshal(body, &response)

	if marshalErr != nil {
		return model.PaymentResponse{}, marshalErr
	}

	qry := "Insert Into payment (user_id, doctor_id, amount, order_id, status, queue_date, queue_time, note) Values($1,$2,$3,$4,$5,$6,$7,$8)"

	_, err = tx.Exec(qry, req.CustomerDetails.CustomerID, req.ItemDetails[0].Id, req.TransactionDetails.GrossAmount, req.TransactionDetails.OrderID, "Order", req.ItemDetails[0].ReservasiDate, req.ItemDetails[0].ReservasiTime, req.ItemDetails[0].Note)

	if err != nil {
		tx.Rollback()
		return model.PaymentResponse{}, err
	}

	tx.Commit()

	return response, nil
}

func (p *paymentRepository) TrackingTransactionNotications(req model.TrackerTransaction) error {
	if req.TransactionStatus == "capture" || req.TransactionStatus == "settlement" {
		dtoPayment := dto.PaymentDTO{}
		tx, err := p.db.Begin()
		qry1 := "Select user_id, doctor_id, queue_date, queue_time, note from payment where order_id = $1"
		qry2 := "Update payment Set status=$1 where order_id=$2"

		client := &http.Client{}

		if err != nil {
			return err
		}

		err = tx.QueryRow(qry1, req.OrderID).Scan(&dtoPayment.CustomerID, &dtoPayment.DoctorID, &dtoPayment.QueueDate, &dtoPayment.QueueTime, &dtoPayment.Note)

		if err != nil {
			tx.Rollback()
			return err
		}

		_, err = tx.Exec(qry2, "Paid")

		if err != nil {
			return err
		}

		dataJson, errDec := json.Marshal(dtoPayment)

		if errDec != nil {
			return errDec
		}

		request, errReq := http.NewRequest("POST", "http://localhost:8888/api-klinik-gigi-vony-nur-santy/queues", bytes.NewBuffer(dataJson))

		if errReq != nil {
			return errReq
		}

		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Accept", "application/json")

		res, err := client.Do(request)

		if err != nil {
			tx.Rollback()
			return err
		}

		defer res.Body.Close()

		if res.StatusCode != 201 {
			tx.Rollback()
			return errors.New("failed created reservation")
		}

		tx.Commit()
	}

	return nil
}

func NewRepository(db *sql.DB) PaymentRepository {
	return &paymentRepository{db}
}
