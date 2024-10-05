package controller

import (
	"fmt"
	"net/http"

	"github.com/Ak8388/applikasi-antrian-dokter-gigi-payment-service/driver/middleware"
	"github.com/Ak8388/applikasi-antrian-dokter-gigi-payment-service/model"
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
	var request model.UserRefundPayment

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Println(request)
	id, exist := c.Get("userID")

	if !exist && request.UserID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed get id user"})
		return
	}

	if request.UserID == "" {
		request.UserID = id.(string)
	}

	if err := p.pyRef.CreateRefundPayment(request); err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "ok"})
}

func (p *paymentRefund) changeStatus(c *gin.Context) {
	var Request struct {
		Id     string `json:"id"`
		Status string `json:"status"`
	}

	if err := c.ShouldBindJSON(&Request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := p.pyRef.ChangeStatus(Request.Id, Request.Status); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "ok"})
}

func (p *paymentRefund) getDataRefundPaymentForPatient(c *gin.Context) {
	status := c.Query("status")

	id, exist := c.Get("userID")

	if !exist {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed get id user"})
		return
	}

	res, err := p.pyRef.GetDataRefundPaymentForPatient(id.(string), status)

	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success get data refund", "data": res})
}

func (p *paymentRefund) getDataRefundPaymentForAdmin(c *gin.Context) {
	status := c.Query("status")

	res, err := p.pyRef.GetDataRefundPaymentForAdmin(status)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success get data refund", "data": res})
}

func (p *paymentRefund) deleteDataRefund(c *gin.Context) {
	var Request struct {
		Id string `json:"id"`
	}

	if err := c.ShouldBindJSON(&Request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := p.pyRef.DeleteDataRefund(Request.Id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "sucess delete refund"})
}

func (p *paymentRefund) changePayment(c *gin.Context) {
	var request dto.RefundPaymentChangeData

	if err := c.ShouldBindJSON(&request); err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := p.pyRef.ChangePayment(request)

	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}

func (p *paymentRefund) RefundPaymentRouter() {
	r := p.rg.Group("payment-refund")
	r.POST("", p.am.JwtVerify("Admin", "Patient"), p.createRefundPayment)
	r.PUT("", p.am.JwtVerify("Admin"), p.changeStatus)
	r.GET("view-for-patient", p.am.JwtVerify("Patient"), p.getDataRefundPaymentForPatient)
	r.GET("view-for-admin", p.am.JwtVerify("Admin"), p.getDataRefundPaymentForAdmin)
	r.PUT("data-refund", p.am.JwtVerify("Admin", "Patient"), p.changePayment)
	r.DELETE("", p.am.JwtVerify("Admin"), p.deleteDataRefund)
}

func NewPaymentRefund(am middleware.AuthMiddleware, pyRef usecase.RefundPaymentUsecase, rg *gin.RouterGroup) *paymentRefund {
	return &paymentRefund{
		am:    am,
		pyRef: pyRef,
		rg:    rg,
	}
}
