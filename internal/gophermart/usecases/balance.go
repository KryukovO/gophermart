package usecases

import (
	"context"
	"time"

	"github.com/KryukovO/gophermart/internal/gophermart/entities"
	"github.com/KryukovO/gophermart/internal/gophermart/repository"
)

type BalanceUseCase struct {
	repo    repository.BalanceRepo
	timeout time.Duration
}

func NewBalanceUseCase(repo repository.BalanceRepo, timeout time.Duration) *BalanceUseCase {
	return &BalanceUseCase{
		repo:    repo,
		timeout: timeout,
	}
}

func (uc *BalanceUseCase) Balance(ctx context.Context, userID int64) (entities.Balance, error) {
	ctx, cancel := context.WithTimeout(ctx, uc.timeout)
	defer cancel()

	return uc.repo.Balance(ctx, userID)
}

func (uc *BalanceUseCase) ChangeBalance(ctx context.Context, change *entities.BalanceChange) error {
	ctx, cancel := context.WithTimeout(ctx, uc.timeout)
	defer cancel()

	if err := change.Validate(); err != nil {
		return err
	}

	return uc.repo.ChangeBalance(ctx, change)
}

func (uc *BalanceUseCase) Withdrawals(ctx context.Context, userID int64) ([]entities.BalanceChange, error) {
	ctx, cancel := context.WithTimeout(ctx, uc.timeout)
	defer cancel()

	return uc.repo.Withdrawals(ctx, userID)
}
