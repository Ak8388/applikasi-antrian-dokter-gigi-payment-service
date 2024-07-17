package manager

import "github.com/Ak8388/applikasi-antrian-dokter-gigi-payment-service/usecase"

type UsecaseManager interface {
	PaymentUsecaseManager() usecase.PaymentUscase
}

type usecaseManager struct {
	repoM RepositoryManager
}

func (u *usecaseManager) PaymentUsecaseManager() usecase.PaymentUscase {
	return usecase.NewPaymentUsecase(u.repoM.PaymentRepoManager())
}

func NewPaymentUsecase(repoM RepositoryManager) UsecaseManager {
	return &usecaseManager{repoM}
}
