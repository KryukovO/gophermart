package usecases

import (
	"context"
	"time"

	"github.com/KryukovO/gophermart/internal/gophermart/entities"
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

func (uc *OrderUseCase) AddOrder(ctx context.Context, order *entities.Order) error {
	if err := order.Validate(); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(ctx, uc.timeout)
	defer cancel()

	return uc.repo.AddOrder(ctx, order)
}

func (uc *OrderUseCase) Orders(ctx context.Context, userID int64) ([]entities.Order, error) {
	ctx, cancel := context.WithTimeout(ctx, uc.timeout)
	defer cancel()

	return uc.repo.Orders(ctx, userID)
}
