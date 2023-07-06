package repository

import (
	"context"

	"github.com/KryukovO/gophermart/internal/gophermart/entities"
)

type UserRepo interface {
	CreateUser(ctx context.Context, user *entities.User) error
	UserByLogin(ctx context.Context, user *entities.User) error
}

type OrderRepo interface {
	Orders() error
	CreateOrder() error
}

type BalanceRepo interface {
	Balance() error
	CreateWithdrawal() error
	Withdrawals() error
}
