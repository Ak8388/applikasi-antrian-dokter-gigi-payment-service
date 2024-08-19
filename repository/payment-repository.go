package repository

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/Ak8388/applikasi-antrian-dokter-gigi-payment-service/model"
	"github.com/Ak8388/applikasi-antrian-dokter-gigi-payment-service/model/dto"
	"github.com/joho/godotenv"
)

type PaymentRepository interface {
	CustomerPayment(req model.TransModel) (model.PaymentResponse, error)
	TrackingTransactionNotications(req model.TrackerTransaction) error
	FindPaymentById(userId, to string) bool
	FindPaymentPaidByDate(date, userId string) bool
	CancelPaymentUser(orderId string) error
	ViewPaymentbyUserId(user, status string) (data []dto.PaymentViewDTO, err error)
}

type paymentRepository struct {
	db *sql.DB
}

func (p *paymentRepository) CustomerPayment(req model.TransModel) (model.PaymentResponse, error) {
	var dtoValidate dto.ValidatePayment
	client := &http.Client{}

	dtoValidate = dto.ValidatePayment{
		DoctorId: req.ItemDetails[0].Id,
		OpenTime: req.ItemDetails[0].ReservasiTime,
		Date:     req.ItemDetails[0].ReservasiDate,
		ToStts:   req.ItemDetails[0].To,
	}

	dataJsonVal, err := json.Marshal(dtoValidate)

	if err != nil {
		return model.PaymentResponse{}, err
	}

	requestVal, errReq1 := http.NewRequest("POST", "http://localhost:8888/api-klinik-gigi-vony-nur-santy/queues/validate", bytes.NewBuffer(dataJsonVal))

	if errReq1 != nil {
		return model.PaymentResponse{}, errReq1
	}

	resVal, errRes := client.Do(requestVal)

	if errRes != nil {
		return model.PaymentResponse{}, errRes
	}

	defer resVal.Body.Close()
	body, _ := io.ReadAll(resVal.Body)

	if resVal.StatusCode > 299 {
		return model.PaymentResponse{}, errors.New(string(body))
	}

	var qry string
	tx, err := p.db.Begin()

	if err != nil {
		return model.PaymentResponse{}, err
	}

	if err := godotenv.Load(); err != nil {
		return model.PaymentResponse{}, err
	}

	var response = model.PaymentResponse{
		Token:       "",
		RedirectURL: "",
		TransactionDetails: model.TransactionDetails{
			OrderID: req.TransactionDetails.OrderID,
		},
	}

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
	body, _ = io.ReadAll(res.Body)

	if res.StatusCode > 299 {
		return model.PaymentResponse{}, errors.New(string(body))
	}

	marshalErr := json.Unmarshal(body, &response)

	if marshalErr != nil {
		return model.PaymentResponse{}, marshalErr
	}

	if req.ItemDetails[0].To == "Created" {
		qry = "Insert Into payment (user_id, doctor_id, amount, order_id, status, queue_date, queue_time, note, to_stts) Values($1,$2,$3,$4,$5,$6,$7,$8,$9)"
		_, err = tx.Exec(qry, req.CustomerDetails.CustomerID, req.ItemDetails[0].Id, req.TransactionDetails.GrossAmount, req.TransactionDetails.OrderID, "Order", req.ItemDetails[0].ReservasiDate, req.ItemDetails[0].ReservasiTime, req.ItemDetails[0].Note, req.ItemDetails[0].To)
	} else {
		qry = "Insert Into payment (user_id, doctor_id, amount, order_id, status, queue_date, queue_time, note, to_stts, id_reserv) Values($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)"
		_, err = tx.Exec(qry, req.CustomerDetails.CustomerID, req.ItemDetails[0].Id, req.TransactionDetails.GrossAmount, req.TransactionDetails.OrderID, "Order", req.ItemDetails[0].ReservasiDate, req.ItemDetails[0].ReservasiTime, req.ItemDetails[0].Note, req.ItemDetails[0].To, req.CustomerDetails.ReservID)
	}

	if err != nil {
		tx.Rollback()
		return model.PaymentResponse{}, err
	}

	tx.Commit()

	return response, nil
}

func (p *paymentRepository) TrackingTransactionNotications(req model.TrackerTransaction) error {
	fmt.Println("This Req From Midtrans := ", req)
	if req.TransactionStatus == "capture" || req.TransactionStatus == "settlement" {
		var amount int64
		var errReq error
		var request *http.Request
		dtoPayment := dto.PaymentDTO{}
		var dataAny any
		tx, err := p.db.Begin()

		qry1 := "Select user_id, doctor_id, queue_date, queue_time, note, amount, id_reserv from payment where order_id = $1"
		qry2 := "Update payment Set status=$1 where order_id=$2"

		client := &http.Client{}

		if err != nil {
			return err
		}

		err = tx.QueryRow(qry1, req.OrderID).Scan(&dtoPayment.CustomerID, &dtoPayment.DoctorID, &dtoPayment.QueueDate, &dtoPayment.QueueTime, &dtoPayment.Note, &amount, &dataAny)

		if err != nil {
			tx.Rollback()
			return err
		}

		qdateFix := strings.Replace(dtoPayment.QueueDate, "T", " ", -1)
		qdateFix2 := strings.Replace(qdateFix, "Z", "", -1)
		qdatSplit := strings.Split(qdateFix2, " ")

		qTimeFix := strings.Replace(dtoPayment.QueueTime, "T", " ", -1)
		qTimeFix2 := strings.Replace(qTimeFix, "Z", "", -1)
		qTimeFix3 := strings.Replace(qTimeFix2, "0000-01-01", qdatSplit[0], -1)

		dtoPayment.QueueDate = qdateFix2
		dtoPayment.QueueTime = qTimeFix3

		_, err = tx.Exec(qry2, "Paid", req.OrderID)

		if err != nil {
			return err
		}

		if amount <= 10000 {
			if dataAny != nil {
				var idUint = dataAny.([]uint8)
				dtoPayment.ID = string(idUint)
			}

			dataJson, errDec := json.Marshal(dtoPayment)

			if errDec != nil {
				return errDec
			}

			request, errReq = http.NewRequest("PUT", "http://localhost:8888/api-klinik-gigi-vony-nur-santy/queues/reschedules", bytes.NewBuffer(dataJson))
		} else {
			dataJson, errDec := json.Marshal(dtoPayment)

			if errDec != nil {
				return errDec
			}

			request, errReq = http.NewRequest("POST", "http://localhost:8888/api-klinik-gigi-vony-nur-santy/queues", bytes.NewBuffer(dataJson))
		}

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
		body, _ := io.ReadAll(res.Body)

		if res.StatusCode > 299 {
			tx.Rollback()
			return errors.New("failed created reservation :" + string(body))
		}

		tx.Commit()
		return nil
	}

	if req.TransactionStatus == "expire" {
		qry := "Update payment Set status=$1 where order_id=$2"

		_, err := p.db.Exec(qry, "Expired", req.OrderID)

		if err != nil {
			return err
		}

		return nil
	}

	return errors.New("still panding")
}

func (p *paymentRepository) FindPaymentById(userId, to_stts string) bool {
	var exist int32
	qry := "Select Count(id) from payment Where user_id=$1 AND to_stts=$2 AND status=$3"

	err := p.db.QueryRow(qry, userId, to_stts, "Order").Scan(&exist)

	if err != nil {
		if sql.ErrNoRows == err {
			return false
		} else {
			return true
		}
	}

	if exist > 0 {
		return true
	}

	return false
}

func (p *paymentRepository) FindPaymentPaidByDate(date, userId string) bool {
	var exist int32
	qry := "Select Count(id) from payment Where user_id=$1 AND queue_date=$2 AND status=$3"

	err := p.db.QueryRow(qry, userId, date, "Paid").Scan(&exist)

	if err != nil {
		if sql.ErrNoRows == err {
			return false
		} else {
			return true
		}
	}

	if exist > 0 {
		return true
	}

	return false
}

func (p *paymentRepository) CancelPaymentUser(orderId string) error {
	client := http.Client{}
	var status string
	tx, err := p.db.Begin()

	qry1 := "Select status from payment Where order_id=$1"
	qry2 := "Update payment Set status=$1 where order_id=$2"

	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.QueryRow(qry1, orderId).Scan(&status)

	if err != nil {
		tx.Rollback()
		return err
	}

	if status != "Order" {
		tx.Rollback()
		return errors.New("payment not valid")
	}

	url := fmt.Sprintf("https://api.sandbox.midtrans.com/v2/%s/cancel", orderId)

	request, errReq := http.NewRequest("POST", url, nil)

	if errReq != nil {
		tx.Rollback()
		return errReq
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.SetBasicAuth(os.Getenv("SERVER_KEY"), "")

	res, errRes := client.Do(request)

	if errRes != nil {
		tx.Rollback()
		return errRes
	}

	if res.StatusCode != 200 {
		tx.Rollback()
		return errors.New("failed cancel payment")
	}

	_, err = tx.Exec(qry2, "Cancel", orderId)

	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil
}

func (p *paymentRepository) ViewPaymentbyUserId(userId, status string) (data []dto.PaymentViewDTO, err error) {
	var args []interface{}
	qry := "Select id, user_id, doctor_id, amount, order_id, status, to_stts, queue_date, queue_time, created_at From payment Where user_id=$1"
	args = append(args, userId)

	if status != "" {
		qry += " AND status=$2"
		args = append(args, status)
	}

	qry += " Order By created_at DESC"

	rows, err := p.db.Query(qry, args...)
	fmt.Println(qry)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		paymentDto := dto.PaymentViewDTO{}

		err = rows.Scan(&paymentDto.Id, &paymentDto.UserId, &paymentDto.DoctorId, &paymentDto.Amount, &paymentDto.OrderId, &paymentDto.Status, &paymentDto.To, &paymentDto.QueueDate, &paymentDto.QueueTime, &paymentDto.CreatedAt)

		if err != nil {
			return nil, err
		}

		data = append(data, paymentDto)
	}

	return
}

func NewRepository(db *sql.DB) PaymentRepository {
	return &paymentRepository{db}
}
