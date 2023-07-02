package usecases

import "time"

type OrderUseCase struct {
	repo    OrderRepo
	timeout time.Duration
}

func NewOrderUseCase(repo OrderRepo, timeout time.Duration) *OrderUseCase {
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
