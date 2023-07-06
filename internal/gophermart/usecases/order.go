package usecases

import (
	"time"

	"github.com/KryukovO/gophermart/internal/gophermart/repository"
)

type OrderUseCase struct {
	repo    repository.OrderRepo
	timeout time.Duration
}

func NewOrderUseCase(repo repository.OrderRepo, timeout time.Duration) *OrderUseCase {
	return &OrderUseCase{
		repo:    repo,
		timeout: timeout,
	}
}

func (uc *OrderUseCase) Orders() error {
	return nil
}

func (uc *OrderUseCase) AddOrder() error {
	return nil
}
