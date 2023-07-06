package usecases

import (
	"time"

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

func (uc *BalanceUseCase) Balance() error {
	return nil
}

func (uc *BalanceUseCase) Withdraw() error {
	return nil
}

func (uc *BalanceUseCase) Withdrawals() error {
	return nil
}
