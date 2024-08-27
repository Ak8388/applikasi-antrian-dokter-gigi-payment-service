package controller

import (
	"net/http"

	"github.com/Ak8388/applikasi-antrian-dokter-gigi-payment-service/driver/middleware"
	"github.com/Ak8388/applikasi-antrian-dokter-gigi-payment-service/model/dto"
	"github.com/Ak8388/applikasi-antrian-dokter-gigi-payment-service/usecase"
	"github.com/gin-gonic/gin"
)

type paymentRefund struct {
	am    middleware.AuthMiddleware
	pyRef usecase.RefundPaymentUsecase
	rg    *gin.RouterGroup
}

func (p *paymentRefund) createRefundPayment(c *gin.Context) {
	// var request model.UserRefundPayment

}

func (p *paymentRefund) changeStatus(c *gin.Context) {

}

func (p *paymentRefund) getDataRefundPaymentForPatient(c *gin.Context) {

}

func (p *paymentRefund) getDataRefundPaymentForAdmin(c *gin.Context) {

}

func (p *paymentRefund) deleteDataRefund(c *gin.Context) {

}

func (p *paymentRefund) changePayment(c *gin.Context) {
	var request dto.RefundPaymentChangeData

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := p.pyRef.ChangePayment(request)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
}

func (p *paymentRefund) RefundPaymentRouter() {
	r := p.rg.Group("payment-refund")
	r.POST("", p.am.JwtVerify("Admin", "Patient"), p.createRefundPayment)
	r.PUT("", p.am.JwtVerify("Admin"), p.changeStatus)
	r.GET("view-patient", p.am.JwtVerify("Patient"))
}

func NewPaymentRefund(am middleware.AuthMiddleware, pyRef usecase.RefundPaymentUsecase, rg *gin.RouterGroup) *paymentRefund {
	return &paymentRefund{
		am:    am,
		pyRef: pyRef,
		rg:    rg,
	}
}
