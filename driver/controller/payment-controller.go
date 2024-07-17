package controller

import (
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

	if err := c.ShouldBind(&payReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	res, err := p.pyUsecase.CustomerPayment(payReq)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"Message": "Success create payment",
		"Data":    res,
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

	c.JSON(http.StatusOK, gin.H{"Message": "Success Tracking Payment"})
}

func (p *paymentController) PaymentRouter() {
	r := p.rg.Group("payment-reservations")

	r.POST("", p.am.JwtVerify("Patient"), p.createPayment)
	r.POST("tracking-payment", p.trackingPayment)
}

func NewControllerPayment(am middleware.AuthMiddleware, pyUsecase usecase.PaymentUscase, rg *gin.RouterGroup) *paymentController {
	return &paymentController{
		am:        am,
		pyUsecase: pyUsecase,
		rg:        rg,
	}
}
