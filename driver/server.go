package driver

import (
	"github.com/Ak8388/applikasi-antrian-dokter-gigi-payment-service/config"
	"github.com/Ak8388/applikasi-antrian-dokter-gigi-payment-service/driver/controller"
	"github.com/Ak8388/applikasi-antrian-dokter-gigi-payment-service/driver/middleware"
	"github.com/Ak8388/applikasi-antrian-dokter-gigi-payment-service/manager"
	"github.com/Ak8388/applikasi-antrian-dokter-gigi-payment-service/utils/common"
	"github.com/gin-gonic/gin"
)

type serverRequirment struct {
	cfg     *config.Config
	useM    manager.UsecaseManager
	engine  *gin.Engine
	jwtAuth common.JwtToken
	host    string
}

func (sr *serverRequirment) setUpController() {
	rg := sr.engine.Group("api-klinik-gigi-vony-nur-santy")

	am := middleware.NewAuthMiddleware(sr.jwtAuth)

	// payment controller
	controller.NewControllerPayment(am, sr.useM.PaymentUsecaseManager(), rg).PaymentRouter()
	// refund controller
	controller.NewPaymentRefund(am, sr.useM.RefundPaymentUsecase(), rg).RefundPaymentRouter()
}

func (sr *serverRequirment) Run() {
	sr.setUpController()

	if err := sr.engine.Run(":" + sr.host); err != nil {
		panic(err)
	}
}

func NewServer() *serverRequirment {
	eng := gin.Default()
	cfg := config.Cfg()
	infraM := manager.NewInfraManager(cfg)
	repoM := manager.NewRepoManager(infraM)
	useM := manager.NewPaymentUsecase(repoM)
	jwt := common.NewJwtUtils(cfg)
	middleware.AddCors(eng)

	return &serverRequirment{
		cfg:     cfg,
		useM:    useM,
		engine:  eng,
		jwtAuth: jwt,
		host:    cfg.ApiPort,
	}
}
