package usecases

import (
	"context"

	"github.com/KryukovO/gophermart/internal/gophermart/entities"
)

type User interface {
	Register(ctx context.Context, user *entities.User, secret []byte) error
	Login(ctx context.Context, user *entities.User, secret []byte) error
}

type Order interface {
	AddOrder(ctx context.Context, order *entities.Order) error
	Orders(ctx context.Context, userID int64) ([]entities.Order, error)
}

type Balance interface {
	Balance() error
	Withdraw() error
	Withdrawals() error
}
