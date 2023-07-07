package repository

import (
	"context"

	"github.com/KryukovO/gophermart/internal/gophermart/entities"
)

type UserRepo interface {
	AddUser(ctx context.Context, user *entities.User) error
	User(ctx context.Context, user *entities.User) error
}

type OrderRepo interface {
	AddOrder(ctx context.Context, order *entities.Order) error
	Orders(ctx context.Context, userID int64) ([]entities.Order, error)
}

type BalanceRepo interface {
	Balance() error
	AddWithdrawal() error
	Withdrawals() error
}
