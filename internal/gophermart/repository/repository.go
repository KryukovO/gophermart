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
	ProcessableOrders(ctx context.Context) ([]entities.Order, error)
	UpdateOrder(ctx context.Context, order *entities.Order) error
}

type BalanceRepo interface {
	Balance(ctx context.Context, userID int64) (entities.Balance, error)
	ChangeBalance(ctx context.Context, change *entities.BalanceChange) error
	Withdrawals(ctx context.Context, userID int64) ([]entities.BalanceChange, error)
}
