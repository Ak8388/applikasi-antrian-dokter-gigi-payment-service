package controller

import (
	"fmt"
	"net/http"

	"github.com/Ak8388/applikasi-antrian-dokter-gigi-payment-service/driver/middleware"
	"github.com/Ak8388/applikasi-antrian-dokter-gigi-payment-service/model"
	"github.com/Ak8388/applikasi-antrian-dokter-gigi-payment-service/usecase"
	"github.com/gin-gonic/gin"
)

type paymentController struct {
	am        middleware.AuthMiddleware
	pyUsecase usecase.PaymentUscase
	rg        *gin.RouterGroup
}

func (p *paymentController) createPayment(c *gin.Context) {
	var payReq model.TransModel

	if err := c.ShouldBindJSON(&payReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	idCust, exist := c.Get("userID")

	if !exist {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "user is empty"})
		return
	}

	payReq.CustomerDetails.CustomerID = idCust.(string)

	res, err := p.pyUsecase.CustomerPayment(payReq)

	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Success create payment",
		"data":    res,
	})
}

func (p *paymentController) trackingPayment(c *gin.Context) {
	var trackReq model.TrackerTransaction

	if err := c.ShouldBind(&trackReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	err := p.pyUsecase.TrackingTransactionNotications(trackReq)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{"Error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Success Tracking Payment"})
}

func (p *paymentController) paymentCanceled(c *gin.Context) {
	var OrderID struct {
		OrderId string `json:"orderId"`
	}

	if err := c.ShouldBindJSON(&OrderID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	err := p.pyUsecase.CancelPaymentUser(OrderID.OrderId)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"message": "success canceled payment"})
}

func (p *paymentController) viewPaymentbyUserId(c *gin.Context) {
	var id any
	var exist bool

	idUser := c.Query("idUser")

	if idUser == "" {
		id, exist = c.Get("userID")

		if !exist {
			fmt.Println("no Exist")
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "failed to acsess id user"})
			return
		}

		idUser = id.(string)
	}

	status := c.Query("status")

	res, err := p.pyUsecase.ViewPaymentbyUserId(idUser, status)

	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success get data payment",
		"data":    res,
	})
}

func (p *paymentController) PaymentRouter() {
	r := p.rg.Group("payment-reservations")

	r.POST("", p.am.JwtVerify("Patient"), p.createPayment)
	r.POST("canceled-payment", p.am.JwtVerify("Patient"), p.paymentCanceled)
	r.POST("tracking-payment", p.trackingPayment)
	r.GET("find-payment", p.am.JwtVerify("Patient", "Admin", "Doctor"), p.viewPaymentbyUserId)
}

func NewControllerPayment(am middleware.AuthMiddleware, pyUsecase usecase.PaymentUscase, rg *gin.RouterGroup) *paymentController {
	return &paymentController{
		am:        am,
		pyUsecase: pyUsecase,
		rg:        rg,
	}
}
