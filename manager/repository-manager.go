package manager

import (
	"github.com/Ak8388/applikasi-antrian-dokter-gigi-payment-service/repository"
)

type RepositoryManager interface {
	PaymentRepoManager() repository.PaymentRepository
	PaymentRefundManager() repository.RefundPaymentRepository
}

type repositoryManager struct {
	infM InfraManager
}

func (r *repositoryManager) PaymentRepoManager() repository.PaymentRepository {
	return repository.NewRepository(r.infM.DBConnection())
}

func (r *repositoryManager) PaymentRefundManager() repository.RefundPaymentRepository {
	return repository.NewRefundPayment(r.infM.DBConnection())
}

func NewRepoManager(infM InfraManager) RepositoryManager {
	return &repositoryManager{infM}
}
