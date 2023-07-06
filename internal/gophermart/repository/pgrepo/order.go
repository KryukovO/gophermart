package pgrepo

import (
	"github.com/KryukovO/gophermart/internal/postgres"
)

type OrderRepo struct {
	db *postgres.Postgres
}

func NewOrderRepo(db *postgres.Postgres) *OrderRepo {
	return &OrderRepo{db: db}
}

func (repo *OrderRepo) Orders() error {
	return nil
}

func (repo *OrderRepo) CreateOrder() error {
	return nil
}
