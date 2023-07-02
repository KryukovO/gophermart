package usecases

import "time"

type BalanceUseCase struct {
	repo    BalanceRepo
	timeout time.Duration
}

func NewBalanceUseCase(repo BalanceRepo, timeout time.Duration) *BalanceUseCase {
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
