package pgrepo

import (
	"github.com/KryukovO/gophermart/internal/postgres"
)

type BalanceRepo struct {
	db *postgres.Postgres
}

func NewBalanceRepo(db *postgres.Postgres) *BalanceRepo {
	return &BalanceRepo{db: db}
}

func (repo *BalanceRepo) Balance() error {
	return nil
}

func (repo *BalanceRepo) AddWithdrawal() error {
	return nil
}

func (repo *BalanceRepo) Withdrawals() error {
	return nil
}
